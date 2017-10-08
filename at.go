package at

import (
	"log"
	"runtime"
	"sort"
	"time"

	"github.com/pborman/uuid"
)

type At struct {
	Log *log.Logger

	entries  []*Entry
	add      chan *Entry
	remove   chan string
	stop     chan struct{}
	snapshot chan []*Entry
	running  bool
	location *time.Location
}

type Entry struct {
	// The schedule on which this job should be run.
	Schedule Schedule

	// The time the job will run.
	At time.Time

	// The job to run
	Job Job

	// Indicates whether the job has been executed
	Ran bool

	// Job ID
	ID string
}

type Schedule interface {
	At(t time.Time) time.Time
}

type Job interface {
	Run()
}

// New returns a new At job runner, in the local time zone.
func New() *At {
	return NewWithLocation(time.Now().Location())
}

// NewWithLocation returns a new At job runner.
func NewWithLocation(locaton *time.Location) *At {
	return &At{
		entries:  nil,
		add:      make(chan *Entry),
		remove:   make(chan string),
		stop:     make(chan struct{}),
		snapshot: make(chan []*Entry),
		running:  false,
		Log:      nil,
		location: locaton,
	}
}

// A wrapper that turns a func() into a at.Job
type FuncJob func()

func (f FuncJob) Run() {
	f()
}

func (a *At) AddFuncWithID(id, spec string, cmd func()) error {
	return a.AddJobWithID(id, spec, FuncJob(cmd))
}

func (a *At) AddJobWithID(id, spec string, cmd Job) error {
	schedule, err := Parse(spec)
	if err != nil {
		return err
	}
	a.ScheduleWithID(id, schedule, cmd)
	return nil
}

func (a *At) ScheduleWithID(id string, schedule Schedule, cmd Job) {
	entry := &Entry{
		Schedule: schedule,
		Job:      cmd,
		Ran:      false,
		ID:       id,
	}
	if !a.running {
		a.entries = append(a.entries, entry)
		return
	}

	a.add <- entry
}

// AddFunc adds a func to the At to be run on the given schedule.
func (a *At) AddFunc(spec string, cmd func()) error {
	return a.AddJob(spec, FuncJob(cmd))
}

// AddJob adds a Job to the At to be run on the given schedule.
func (a *At) AddJob(spec string, cmd Job) error {
	schedule, err := Parse(spec)
	if err != nil {
		return err
	}
	a.Schedule(schedule, cmd)
	return nil
}

// Schedule adds a Job to the At to be run on the given schedule.
func (a *At) Schedule(schedule Schedule, cmd Job) {
	entry := &Entry{
		Schedule: schedule,
		Job:      cmd,
		Ran:      false,
		ID:       uuid.NewUUID().String(),
	}
	if !a.running {
		a.entries = append(a.entries, entry)
		return
	}

	a.add <- entry
}

// Start the at scheduler in its own go-routine, or no-op if already started.
func (a *At) Start() {
	if a.running {
		return
	}
	a.running = true
	go a.run()
}

// Run the at scheduler, or no-op if already running.
func (a *At) Run() {
	if a.running {
		return
	}
	a.running = true
	a.run()
}

func (a *At) Remove(id string) {
	if !a.running {
		es := make([]*Entry, 0)
		for _, e := range a.entries {
			if e.ID == id {
				continue
			}
			es = append(es, e)
		}
		a.entries = es
		return
	}
	a.remove <- id
}

// Stop stops the at scheduler if it is running; otherwise it does nothing.
func (a *At) Stop() {
	if !a.running {
		return
	}

	a.stop <- struct{}{}
	a.running = false
}

// Entries returns a snapshot of the at entries.
func (a *At) Entries() []*Entry {
	if a.running {
		a.snapshot <- nil
		x := <-a.snapshot
		return x
	}

	return a.entrySnapshot()
}

// Location gets the time zone location
func (a *At) Location() *time.Location {
	return a.location
}

func (a *At) run() {
	now := a.now()
	for _, entry := range a.entries {
		entry.At = entry.Schedule.At(now)
	}

	for {
		sort.Sort(byTime(a.entries))

		var timer *time.Timer
		if len(a.entries) == 0 || a.entries[0].Ran {
			// If there are no entries yet, just sleep - it still handles new entries
			// and stop requests.
			timer = time.NewTimer(100000 * time.Hour)
		} else {
			timer = time.NewTimer(a.entries[0].At.Sub(now))
		}

		for {
			select {
			case now = <-timer.C:
				now = now.In(a.location)
				for _, e := range a.entries {
					if e.At.After(now) || e.Ran {
						break
					}

					go a.runWithRecovery(e.Job)
					e.Ran = true
				}

			case newEntry := <-a.add:
				timer.Stop()
				now = a.now()
				newEntry.At = newEntry.Schedule.At(now)

				es := make([]*Entry, 0, len(a.entries)+1)
				for _, e := range a.entries {
					if e.Ran {
						continue
					}
					es = append(es, e)
				}
				es = append(es, newEntry)
				a.entries = es

			case id := <-a.remove:
				timer.Stop()
				now = a.now()

				es := make([]*Entry, 0)
				for _, e := range a.entries {
					if e.ID == id || e.Ran {
						continue
					}
					es = append(es, e)
				}
				a.entries = es

			case <-a.snapshot:
				a.snapshot <- a.entrySnapshot()
				continue

			case <-a.stop:
				timer.Stop()
				return
			}

			break
		}
	}
}

// entrySnapshot returns a copy of the current at entry list.
func (a *At) entrySnapshot() []*Entry {
	entries := []*Entry{}
	for _, e := range a.entries {
		entries = append(entries, &Entry{
			Schedule: e.Schedule,
			At:       e.At,
			Job:      e.Job,
			Ran:      e.Ran,
		})
	}

	return entries
}

// Logs an error to stderr or to the configured error log
func (a *At) logf(format string, args ...interface{}) {
	if a.Log != nil {
		a.Log.Printf(format, args...)
	} else {
		log.Printf(format, args...)
	}
}

func (a *At) runWithRecovery(j Job) {
	defer func() {
		if r := recover(); r != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			a.logf("at: panic running job: %v\n%s", r, buf)
		}
	}()

	j.Run()
}

// now returns current time in location
func (a *At) now() time.Time {
	return time.Now().In(a.location)
}

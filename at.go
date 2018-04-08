package at

import (
	"log"
	"runtime"
	"time"

	"github.com/gotoxu/at/queue"
)

type At struct {
	Log *log.Logger

	entries  *queue.PriorityQueue
	add      chan *entry
	stop     chan struct{}
	running  bool
	location *time.Location
}

type entry struct {
	// The schedule on which this job should be run.
	Schedule Schedule

	// The time the job will run.
	At time.Time

	// The job to run
	Job Job
}

func (e entry) Compare(other queue.Item) int {
	oe := other.(entry)
	if e.At.Before(oe.At) {
		return -1
	} else if e.At.After(oe.At) {
		return 1
	}

	return 0
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
		entries:  queue.NewPriorityQueue(1),
		add:      make(chan *entry),
		stop:     make(chan struct{}),
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
	entry := &entry{
		Schedule: schedule,
		Job:      cmd,
	}
	if !a.running {
		a.entries.Push(entry)
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

// Stop stops the at scheduler if it is running; otherwise it does nothing.
func (a *At) Stop() {
	if !a.running {
		return
	}

	a.stop <- struct{}{}
	a.running = false
}

// Location gets the time zone location
func (a *At) Location() *time.Location {
	return a.location
}

func (a *At) run() {
	now := a.now()
	// for _, entry := range a.entries {
	// 	entry.At = entry.Schedule.At(now)
	// }

	for {
		var timer *time.Timer
		e := a.entries.Peek()
		if e == nil {
			// If there are no entries yet, just sleep - it still handles new entries
			// and stop requests.
			timer = time.NewTimer(100000 * time.Hour)
		} else {
			entry := e.(entry)
			entry.At = entry.Schedule.At(now)
			timer = time.NewTimer(entry.At.Sub(now))
		}

		for {
			select {
			case now = <-timer.C:
				now = now.In(a.location)

				e, err := a.entries.Pop()
				if err == nil && e != nil {
					entry := e.(entry)
					go a.runWithRecovery(entry.Job)
				}

			case newEntry := <-a.add:
				timer.Stop()
				now = a.now()
				newEntry.At = newEntry.Schedule.At(now)

				a.entries.Push(newEntry)

			case <-a.stop:
				timer.Stop()
				a.entries.Dispose()
				return
			}

			break
		}
	}
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

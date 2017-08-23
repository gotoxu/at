package at

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// DateTimeSchedule specifies a time for job executing
// Time definition format is 'HH:mm:ss yyyy-mm-dd'
type DateTimeSchedule struct {
	Hour   int
	Minute int
	Second int
	Year   int
	Month  time.Month
	Date   int
}

func (s *DateTimeSchedule) At(t time.Time) time.Time {
	at := time.Date(s.Year, s.Month, s.Date, s.Hour, s.Minute, s.Second, 0, t.Location())
	return at
}

func NewDateTimeSchedule(clock string, date string) (*DateTimeSchedule, error) {
	sched := &DateTimeSchedule{}
	cgs := strings.FieldsFunc(clock, func(r rune) bool { return r == ':' })
	if len(cgs) < 2 || len(cgs) > 3 {
		return nil, fmt.Errorf("error time schedule spec format: %s", clock)
	}

	dgs := strings.FieldsFunc(date, func(r rune) bool { return r == '-' })
	if len(dgs) != 3 {
		return nil, fmt.Errorf("error date schedule spec format: %s", date)
	}

	hour, err := getField(cgs[0], hours)
	if err != nil {
		return nil, err
	}
	sched.Hour = hour

	minute, err := getField(cgs[1], minutes)
	if err != nil {
		return nil, err
	}
	sched.Minute = minute

	if len(cgs) == 3 {
		sec, err := getField(cgs[2], seconds)
		if err != nil {
			return nil, err
		}
		sched.Second = sec
	}

	year, err := strconv.Atoi(dgs[0])
	if err != nil {
		return nil, err
	}
	sched.Year = year

	mon, err := getField(dgs[1], months)
	if err != nil {
		return nil, err
	}
	sched.Month = time.Month(mon)

	day, err := strconv.Atoi(dgs[2])
	if err != nil {
		return nil, err
	}
	sched.Date = day

	return sched, nil
}

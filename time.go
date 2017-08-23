package at

import (
	"fmt"
	"strings"
	"time"
)

// TimeSchedule specifies a time for job executing
// Time definition format is 'HH:mm:ss'
type TimeSchedule struct {
	Hour   int
	Minute int
	Second int
}

func (s *TimeSchedule) At(t time.Time) time.Time {
	year, mon, day := t.Date()
	at := time.Date(year, mon, day, s.Hour, s.Minute, s.Second, 0, t.Location())
	if at.Before(t) {
		at = at.AddDate(0, 0, 1)
	}
	return at
}

func NewTimeSchedule(spec string) (*TimeSchedule, error) {
	parts := strings.FieldsFunc(spec, func(r rune) bool { return r == ':' })
	if len(parts) < 2 || len(parts) > 3 {
		return nil, fmt.Errorf("error time schedule spec format: %s", spec)
	}

	hour, err := getField(parts[0], hours)
	if err != nil {
		return nil, err
	}
	min, err := getField(parts[1], minutes)
	if err != nil {
		return nil, err
	}

	if len(parts) == 2 {
		return &TimeSchedule{
			Hour:   hour,
			Minute: min,
			Second: 0,
		}, nil
	}

	sec, err := getField(parts[2], seconds)
	if err != nil {
		return nil, err
	}

	return &TimeSchedule{
		Hour:   hour,
		Minute: min,
		Second: sec,
	}, nil
}

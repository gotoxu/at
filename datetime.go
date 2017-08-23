package at

import (
	"fmt"
	"strings"
	"time"
)

// DateTimeSchedule specifies a time for job executing
// Time definition format is 'HH:mm:ss yyyy-mm-dd'
type DateTimeSchedule struct {
	layout string
	value  string
}

func (s *DateTimeSchedule) At(t time.Time) time.Time {
	at, _ := time.ParseInLocation(s.layout, s.value, t.Location())
	return at
}

func NewDateTimeSchedule(clock string, date string) (*DateTimeSchedule, error) {
	cgs := strings.FieldsFunc(clock, func(r rune) bool { return r == ':' })
	if len(cgs) < 2 || len(cgs) > 3 {
		return nil, fmt.Errorf("error time schedule spec format: %s", clock)
	}

	dgs := strings.FieldsFunc(date, func(r rune) bool { return r == '-' })
	if len(dgs) != 3 {
		return nil, fmt.Errorf("error date schedule spec format: %s", date)
	}

	value := fmt.Sprintf("%s %s", date, clock)
	if len(cgs) == 2 {
		layout := "2006-01-02 15:04"
		_, err := time.Parse(layout, value)
		if err != nil {
			return nil, err
		}

		return &DateTimeSchedule{
			layout: layout,
			value:  value,
		}, nil
	}

	layout := "2006-01-02 15:04:05"
	_, err := time.Parse(layout, value)
	if err != nil {
		return nil, err
	}

	return &DateTimeSchedule{
		layout: layout,
		value:  value,
	}, nil
}

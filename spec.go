package at

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// SpecSchedule specifies a time for job executing
// Time definition format is 'HH:MM[am|pm] [Month] [Date]'
type SpecSchedule struct {
	Month  int
	Date   int
	Hour   int
	Minute int
	Second int
}

func (s *SpecSchedule) At(t time.Time) time.Time {
	at := time.Date(t.Year(), time.Month(s.Month), s.Date, s.Hour, s.Minute, s.Second, 0, t.Location())
	return at
}

func NewSpecSchedule(clock, month, day string) (*SpecSchedule, error) {
	retval := &SpecSchedule{}
	meridiem := strings.ToLower(clock[len(clock)-2:])
	if meridiem != "am" && meridiem != "pm" {
		return nil, fmt.Errorf("unknow clock spec string(%s), use like this: HH:MM[am|pm] or HH[am|pm]", clock)
	}

	clock = clock[:len(clock)-2]
	parts := strings.FieldsFunc(clock, func(r rune) bool { return r == ':' })

	hour, err := getField(parts[0], hours)
	if err != nil {
		return nil, err
	}
	if meridiem == "pm" {
		retval.Hour = hour + 12
	} else {
		retval.Hour = hour
	}

	if len(parts) >= 2 {
		min, err := getField(parts[1], minutes)
		if err != nil {
			return nil, err
		}
		retval.Minute = min
	}

	if len(parts) == 3 {
		sec, err := getField(parts[2], seconds)
		if err != nil {
			return nil, err
		}
		retval.Second = sec
	}

	mon, ok := months.ranges[strings.ToLower(month)]
	if !ok {
		return nil, fmt.Errorf("invalid month(%s), valid month are: jan, feb, mar, apr, ..., dec", month)
	}
	retval.Month = mon

	date, err := strconv.Atoi(day)
	if err != nil {
		return nil, err
	}
	retval.Date = date

	value := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", time.Now().Year(), retval.Month, retval.Date, retval.Hour, retval.Minute, retval.Second)
	_, err = time.Parse("2006-01-02 15:04:05", value)
	if err != nil {
		return nil, err
	}

	return retval, nil
}

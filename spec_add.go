package at

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

var validUnit = []string{
	"seconds",
	"minutes",
	"hours",
	"days",
	"weeks",
	"months",
}

func ValidateUnit(unit string, now bool) (string, error) {
	lower := strings.ToLower(unit)
	if now {
		for _, u := range validUnit {
			if lower == u {
				return lower, nil
			}
		}
		return lower, fmt.Errorf("Invalid unit: %s, Valid unit are: %v", lower, validUnit)
	}

	for _, u := range validUnit {
		if u == "seconds" || u == "minutes" || u == "hours" {
			continue
		}
		if lower == u {
			return lower, nil
		}
	}

	return lower, fmt.Errorf("Invalid unit: %s, Valid unit are: %v", lower, []string{"days", "weeks", "months"})
}

// SpecAddSchedule specifies a time for job executing
// Time definition format is 'HH:MM[am|pm] + number [seconds|minutes|hours|days|weeks|months]'
type SpecAddSchedule struct {
	Now    bool
	Hour   int
	Minute int
	Second int
	Number int
	Unit   string
}

func (s *SpecAddSchedule) At(t time.Time) time.Time {
	var st time.Time
	if s.Now {
		st = t
	} else {
		year, mon, day := t.Date()
		st = time.Date(year, mon, day, s.Hour, s.Minute, s.Second, 0, t.Location())
	}

	var at time.Time
	if s.Unit == "seconds" || s.Unit == "minutes" || s.Unit == "hours" {
		var dur time.Duration
		if s.Unit == "seconds" {
			dur = time.Duration(s.Number) * time.Second
		} else if s.Unit == "minutes" {
			dur = time.Duration(s.Number) * time.Minute
		} else if s.Unit == "hours" {
			dur = time.Duration(s.Number) * time.Hour
		}

		at = st.Add(dur)
	} else {
		if s.Unit == "days" {
			at = st.AddDate(0, 0, s.Number)
		} else if s.Unit == "weeks" {
			at = st.AddDate(0, 0, s.Number*7)
		} else {
			at = st.AddDate(0, s.Number, 0)
		}
	}

	return at
}

func NewSpecAddSchedule(clock string, number string, unit string) (*SpecAddSchedule, error) {
	sched := &SpecAddSchedule{}
	now := clock == "now"

	u, err := ValidateUnit(unit, now)
	if err != nil {
		return nil, err
	}

	sched.Now = now
	sched.Unit = u

	if !now {
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
			sched.Hour = hour + 12
		} else {
			sched.Hour = hour
		}

		if len(parts) >= 2 {
			min, err := getField(parts[1], minutes)
			if err != nil {
				return nil, err
			}
			sched.Minute = min
		}

		if len(parts) == 3 {
			sec, err := getField(parts[2], seconds)
			if err != nil {
				return nil, err
			}
			sched.Second = sec
		}
	}

	num, err := strconv.Atoi(number)
	if err != nil {
		return nil, err
	}

	sched.Number = num
	return sched, nil
}

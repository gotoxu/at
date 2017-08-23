package at

import (
	"fmt"
	"strconv"
	"strings"
)

// bounds provides a range of acceptable values
type bounds struct {
	min, max int
	ranges   map[string]int
}

// The bounds for each field.
var (
	seconds = bounds{0, 59, nil}
	minutes = bounds{0, 59, nil}
	hours   = bounds{0, 23, nil}
	months  = bounds{1, 12, map[string]int{
		"jan": 1,
		"feb": 2,
		"mar": 3,
		"apr": 4,
		"may": 5,
		"jun": 6,
		"jul": 7,
		"aug": 8,
		"sep": 9,
		"oct": 10,
		"nov": 11,
		"dec": 12,
	}}
)

func Parse(spec string) (Schedule, error) {
	if len(spec) == 0 {
		return nil, fmt.Errorf("Empty spec string")
	}
	spec = strings.TrimSpace(spec)

	fields := strings.Fields(spec)
	if len(fields) == 1 {
		return NewTimeSchedule(fields[0])
	}

	if len(fields) == 2 {
		return NewDateTimeSchedule(fields[0], fields[1])
	}

	if len(fields) == 3 {
		return NewSpecSchedule(fields[0], fields[1], fields[2])
	}

	if len(fields) == 4 {
		plus := strings.TrimSpace(fields[1])
		if plus != "+" {
			return nil, fmt.Errorf("unknow spec string, use like: 'HH:MM[am|pm] + number [seconds|minutes|hours|days|weeks]'")
		}

		return NewSpecAddSchedule(fields[0], fields[2], fields[3])
	}

	return nil, fmt.Errorf("unknow spec format")
}

func getField(s string, r bounds) (int, error) {
	is, err := strconv.Atoi(s)
	if err != nil {
		return is, err
	}

	if is < r.min {
		return 0, fmt.Errorf("value (%d) below minimum (%d): %s", is, r.min, s)
	}
	if is > r.max {
		return 0, fmt.Errorf("value (%d) above maximum (%d): %s", is, r.max, s)
	}

	return is, nil
}

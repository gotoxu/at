package at

import (
	"testing"
	"time"
)

func TestSpecAddAt(t *testing.T) {
	runs := []struct {
		time, spec string
		expected   string
	}{
		{"Wed Aug 9 14:45 2017", "04:30pm + 5 days", "Mon Aug 14 16:30 2017"},
		{"Wed Aug 9 14:45 2017", "now + 5 hours", "Wed Aug 9 19:45 2017"},
		{"Wed Aug 9 14:45 2017", "04:30am + 10 days", "Sat Aug 19 4:30 2017"},
		{"Wed Aug 9 14:45 2017", "04am + 2 months", "Mon Oct 9 4:00 2017"},
		{"Wed Aug 9 14:45 2017", "04:20:31am + 15 days", "Thu Aug 24 4:20:31 2017"},
	}

	for _, c := range runs {
		sched, err := Parse(c.spec)
		if err != nil {
			t.Error(err)
			continue
		}
		actual := sched.At(getTime(c.time))
		expected := getTime(c.expected)
		if !actual.Equal(expected) {
			t.Errorf("%s, \"%s\": (expected) %v != %v (actual)", c.time, c.spec, expected, actual)
		}
	}
}

func TestSpecAt(t *testing.T) {
	runs := []struct {
		time, spec string
		expected   string
	}{
		{"Wed Aug 9 14:45 2017", "04:30pm Aug 15", "Tue Aug 15 16:30 2017"},
		{"Wed Aug 9 14:45 2017", "04pm Aug 15", "Tue Aug 15 16:00 2017"},
		{"Wed Aug 9 14:45 2017", "04:30am Aug 15", "Tue Aug 15 4:30 2017"},
		{"Wed Aug 9 14:45 2017", "04am Aug 15", "Tue Aug 15 4:00 2017"},
		{"Wed Aug 9 14:45 2017", "04:20:31am Aug 15", "Tue Aug 15 4:20:31 2017"},
	}

	for _, c := range runs {
		sched, err := Parse(c.spec)
		if err != nil {
			t.Error(err)
			continue
		}
		actual := sched.At(getTime(c.time))
		expected := getTime(c.expected)
		if !actual.Equal(expected) {
			t.Errorf("%s, \"%s\": (expected) %v != %v (actual)", c.time, c.spec, expected, actual)
		}
	}
}

func getTime(value string) time.Time {
	if value == "" {
		return time.Time{}
	}
	t, err := time.Parse("Mon Jan 2 15:04 2006", value)
	if err != nil {
		t, err = time.Parse("Mon Jan 2 15:04:05 2006", value)
		if err != nil {
			t, err = time.Parse("2006-01-02T15:04:05-0700", value)
			if err != nil {
				panic(err)
			}
			// Daylight savings time tests require location
			if ny, err := time.LoadLocation("China/Shanghai"); err == nil {
				t = t.In(ny)
			}
		}
	}

	return t
}

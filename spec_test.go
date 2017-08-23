package at

import (
	"testing"
	"time"
)

func TestSpecAt(t *testing.T) {
	runs := []struct {
		time, spec string
		expected   string
	}{
		{"Mon Aug 9 14:45 2017", "04:30pm Aug 15", "Mon Aug 15 16:30 2017"},
		{"Mon Aug 9 14:45 2017", "04pm Aug 15", "Mon Aug 15 16:00 2017"},
		{"Mon Aug 9 14:45 2017", "04:30am Aug 15", "Mon Aug 15 4:30 2017"},
		{"Mon Aug 9 14:45 2017", "04am Aug 15", "Mon Aug 15 4:00 2017"},
		{"Mon Aug 9 14:45 2017", "04:20:31am Aug 15", "Mon Aug 15 4:20:31 2017"},
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
			if ny, err := time.LoadLocation("America/New_York"); err == nil {
				t = t.In(ny)
			}
		}
	}

	return t
}

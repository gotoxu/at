package at

import (
	"time"
)

// DateTimeSchedule specifies a time for job executing
// Time definition format is 'HH:mm:ss yyyy-mm-dd'
type DateTimeSchedule struct {
	Year   int
	Month  int
	Day    int
	Hour   int
	Minute int
	Second int
}

func (s *DateTimeSchedule) Next(t time.Time) time.Time {
	return time.Now()
}

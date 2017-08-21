package at

// byTime is a wrapper for sorting the entry array by time
type byTime []*Entry

func (s byTime) Len() int {
	return len(s)
}

func (s byTime) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byTime) Less(i, j int) bool {
	// The job that has been ran should return false.
	// To sort it at end of the list.
	if s[i].Ran {
		return false
	}
	if s[j].Ran {
		return true
	}

	return s[i].At.Before(s[j].At)
}

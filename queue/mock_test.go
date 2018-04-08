package queue

type mockItem int

func (m mockItem) Compare(other Item) int {
	om := other.(mockItem)
	if m > om {
		return 1
	} else if m == om {
		return 0
	}

	return -1
}

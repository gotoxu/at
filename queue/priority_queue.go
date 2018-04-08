package queue

import (
	"sync"
)

// Item 代表可以被添加到优先队列中的项
type Item interface {
	// Compare 用来决定优先队列中项的顺序。
	// 返回1表示当前项大于other，返回0表示相等，返回-1表示小于
	Compare(other Item) int
}

type priorityItems []Item

func (items *priorityItems) swap(i, j int) {
	(*items)[i], (*items)[j] = (*items)[j], (*items)[i]
}

func (items *priorityItems) pop() Item {
	size := len(*items)

	items.swap(size-1, 0)
	item := (*items)[size-1]
	(*items)[size-1], *items = nil, (*items)[:size-1]

	index := 0
	childL, childR := 2*index+1, 2*index+2
	for len(*items) > childL {
		child := childL
		if len(*items) > childR && (*items)[childR].Compare((*items)[childL]) < 0 {
			child = childR
		}

		if (*items)[child].Compare((*items)[index]) < 0 {
			items.swap(index, child)

			index = child
			childL, childR = 2*index+1, 2*index+2
		} else {
			break
		}
	}

	return item
}

func (items *priorityItems) push(item Item) {
	*items = append(*items, item)

	index := len(*items) - 1
	parent := int((index - 1) / 2)
	for parent >= 0 && (*items)[parent].Compare(item) > 0 {
		items.swap(index, parent)

		index = parent
		parent = int((index - 1/2))
	}
}

// PriorityQueue 是一个优先队列
type PriorityQueue struct {
	items       priorityItems
	lock        sync.Mutex
	disposeLock sync.Mutex
	disposed    bool
}

// Push 将item添加到优先队列中
func (pq *PriorityQueue) Push(item Item) error {
	pq.lock.Lock()
	defer pq.lock.Unlock()

	if pq.disposed {
		return ErrDisposed
	}

	pq.items.push(item)
	return nil
}

// Pop 弹出优先队列中的队首项
func (pq *PriorityQueue) Pop() (Item, error) {
	pq.lock.Lock()
	defer pq.lock.Unlock()

	if pq.disposed {
		return nil, ErrDisposed
	}

	if len(pq.items) == 0 {
		return nil, nil
	}

	item := pq.items.pop()
	return item, nil
}

// Peek 返回优先队列的队首项，但是不会删除它
func (pq *PriorityQueue) Peek() Item {
	pq.lock.Lock()
	defer pq.lock.Unlock()

	if len(pq.items) > 0 {
		return pq.items[0]
	}

	return nil
}

// Empty 表明队列中是否包含任何项
func (pq *PriorityQueue) Empty() bool {
	pq.lock.Lock()
	defer pq.lock.Unlock()

	return len(pq.items) == 0
}

// Len 返回优先队列的长度
func (pq *PriorityQueue) Len() int {
	pq.lock.Lock()
	defer pq.lock.Unlock()

	return len(pq.items)
}

// Disposed 表明优先队列是否已经释放
func (pq *PriorityQueue) Disposed() bool {
	pq.disposeLock.Lock()
	defer pq.disposeLock.Unlock()

	return pq.disposed
}

// Dispose 释放当前队列
func (pq *PriorityQueue) Dispose() {
	pq.lock.Lock()
	defer pq.lock.Unlock()

	pq.disposeLock.Lock()
	defer pq.disposeLock.Unlock()

	pq.disposed = true
	pq.items = nil
}

// NewPriorityQueue 创建一个新的优先队列
func NewPriorityQueue(hint int) *PriorityQueue {
	return &PriorityQueue{
		items: make(priorityItems, 0, hint),
	}
}

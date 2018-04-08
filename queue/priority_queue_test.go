package queue

import (
	"sync"
	"testing"

	"github.com/gotoxu/assert"
)

func TestPriorityPush(t *testing.T) {
	q := NewPriorityQueue(1)
	q.Push(mockItem(2))
	assert.Len(t, q.items, 1)
	assert.DeepEqual(t, q.items[0], mockItem(2))

	q.Push(mockItem(1))
	assert.Len(t, q.items, 2)
	assert.DeepEqual(t, q.items[0], mockItem(1))
	assert.DeepEqual(t, q.items[1], mockItem(2))
}

func TestPriorityPop(t *testing.T) {
	q := NewPriorityQueue(1)

	q.Push(mockItem(2))
	result, err := q.Pop()
	assert.Nil(t, err)
	assert.Len(t, q.items, 0)

	assert.DeepEqual(t, result, mockItem(2))

	result, err = q.Pop()
	assert.Nil(t, err)
	assert.Nil(t, result)
}

func TestPriorityEmpty(t *testing.T) {
	q := NewPriorityQueue(1)
	assert.True(t, q.Empty())

	q.Push(mockItem(1))
	assert.False(t, q.Empty())
}

func TestPriorityLen(t *testing.T) {
	q := NewPriorityQueue(1)
	assert.DeepEqual(t, q.Len(), 0)

	q.Push(mockItem(1))
	assert.DeepEqual(t, q.Len(), 1)

	q.Push(mockItem(2))
	assert.DeepEqual(t, q.Len(), 2)
}

func TestPriorityPeek(t *testing.T) {
	q := NewPriorityQueue(1)
	q.Push(mockItem(1))

	assert.DeepEqual(t, q.Peek(), mockItem(1))
	assert.DeepEqual(t, q.Len(), 1)
}

func BenchmarkPriority(b *testing.B) {
	q := NewPriorityQueue(b.N)
	var wg sync.WaitGroup
	wg.Add(1)
	i := 0

	go func() {
		for {
			q.Pop()
			i++
			if i == b.N {
				wg.Done()
				break
			}
		}
	}()

	for i := 0; i < b.N; i++ {
		q.Push(mockItem(i))
	}

	wg.Wait()
}

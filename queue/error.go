package queue

import "errors"

var (
	// ErrDisposed 当在一个已经释放的队列上进行操作时返回该错误
	ErrDisposed = errors.New("queue: disposed")
)

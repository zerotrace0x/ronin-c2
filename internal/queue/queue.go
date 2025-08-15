package queue

import "sync"

type Queue[T any] struct {
	mu    sync.Mutex
	items []T
}

func New[T any]() *Queue[T] { return &Queue[T]{} }

func (q *Queue[T]) Enqueue(v T) {
	q.mu.Lock()
	q.items = append(q.items, v)
	q.mu.Unlock()
}
func (q *Queue[T]) Dequeue() (T, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.items) == 0 {
		var zero T
		return zero, false
	}
	v := q.items[0]
	q.items = q.items[1:]
	return v, true
}
func (q *Queue[T]) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.items)
}

package history

import (
	"sync"
)

// RingBuffer is a thread-safe generic fixed-capacity circular buffer.
type RingBuffer[T any] struct {
	mu       sync.RWMutex
	data     []T
	capacity int
	head     int
	size     int
}

// NewRingBuffer allocates a RingBuffer with the given maximum capacity.
func NewRingBuffer[T any](capacity int) *RingBuffer[T] {
	if capacity <= 0 {
		capacity = 100
	}
	return &RingBuffer[T]{
		data:     make([]T, capacity),
		capacity: capacity,
	}
}

// Push adds an item to the buffer. If full, the oldest element is overwritten.
func (r *RingBuffer[T]) Push(item T) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.data[r.head] = item
	r.head = (r.head + 1) % r.capacity
	if r.size < r.capacity {
		r.size++
	}
}

// Slice returns all elements in chronological order (oldest to newest).
func (r *RingBuffer[T]) Slice() []T {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.size == 0 {
		return nil
	}

	result := make([]T, r.size)
	if r.size < r.capacity {
		copy(result, r.data[:r.size])
		return result
	}

	n := copy(result, r.data[r.head:])
	copy(result[n:], r.data[:r.head])
	return result
}

// Latest returns the most recently pushed item.
func (r *RingBuffer[T]) Latest() (T, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var zero T
	if r.size == 0 {
		return zero, false
	}
	idx := (r.head - 1 + r.capacity) % r.capacity
	return r.data[idx], true
}

// Size returns current element count.
func (r *RingBuffer[T]) Size() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.size
}

// Capacity returns max element capacity.
func (r *RingBuffer[T]) Capacity() int {
	return r.capacity
}

/*
Package queue provides a fast, ring-buffer queue based on the version suggested by Dariusz Górecki.
Using this instead of other, simpler, queue implementations (slice+append or linked list) provides
substantial memory and time benefits, and fewer GC pauses.
The queue implemented here is as fast as it is for an additional reason: it is *not* thread-safe.
*/
package queue

import "github.com/snple/types"

// minQueueLen is smallest capacity that queue may have.
// Must be power of 2 for bitwise modulus: x % n == x & (n - 1).
const minQueueLen = 16

// Queue represents a single instance of the queue data structure.
type Queue[T any] struct {
	buf               []types.Option[T]
	head, tail, count int
}

// New constructs and returns a new Queue.
func New[T any]() Queue[T] {
	return Queue[T]{
		buf: make([]types.Option[T], minQueueLen),
	}
}

// Length returns the number of elements currently stored in the queue.
func (q *Queue[T]) Length() int {
	return q.count
}

// resizes the queue to fit exactly twice its current contents
// this can result in shrinking if the queue is less than half-full
func (q *Queue[T]) resize() {
	newBuf := make([]types.Option[T], q.count<<1)

	if q.tail > q.head {
		copy(newBuf, q.buf[q.head:q.tail])
	} else {
		n := copy(newBuf, q.buf[q.head:])
		copy(newBuf[n:], q.buf[:q.tail])
	}

	q.head = 0
	q.tail = q.count
	q.buf = newBuf
}

// Push puts an element on the end of the queue.
func (q *Queue[T]) Push(elem T) {
	if q.count == len(q.buf) {
		q.resize()
	}

	q.buf[q.tail] = types.Some(elem)
	// bitwise modulus
	q.tail = (q.tail + 1) & (len(q.buf) - 1)
	q.count++
}

// Peek returns the element at the head of the queue.
func (q *Queue[T]) Peek() types.Option[T] {
	if q.count <= 0 {
		return types.None[T]()
	}
	return q.buf[q.head]
}

// Get returns the element at index i in the queue.
// This method accepts both positive and negative index values.
// Index 0 refers to the first element, and index -1 refers
// to the last.
func (q *Queue[T]) Get(i int) types.Option[T] {
	// If indexing backwards, convert to positive index.
	if i < 0 {
		i += q.count
	}
	if i < 0 || i >= q.count {
		return types.None[T]()
	}
	// bitwise modulus
	return q.buf[(q.head+i)&(len(q.buf)-1)]
}

// Pop removes and returns the element from the front of the queue.
func (q *Queue[T]) Pop() types.Option[T] {
	if q.count <= 0 {
		return types.None[T]()
	}
	ret := q.buf[q.head]
	q.buf[q.head] = types.None[T]()
	// bitwise modulus
	q.head = (q.head + 1) & (len(q.buf) - 1)
	q.count--
	// Resize down if buffer 1/4 full.
	if len(q.buf) > minQueueLen && (q.count<<2) == len(q.buf) {
		q.resize()
	}
	return ret
}

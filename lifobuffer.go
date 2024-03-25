package robin

// LIFOBuffer is a stack-like buffer with a fixed capacity.
// If the buffer is full, pushing a new value will overwrite
// the oldest value in the buffer.
//
// All operations are O(1). Like [Robin], it is backed by a
// map to keep track of the values in the buffer. When used
// by [Robin], the buffer will only receive unique values,
// although it is not enforced by the buffer itself.
type LIFOBuffer[T comparable] struct {
	buf   []T
	i     int
	n     int
	count map[T]int

	capacity int
}

// NewLIFOBuffer creates a new [LIFOBuffer] with the given capacity.
func NewLIFOBuffer[T comparable](capacity int) *LIFOBuffer[T] {
	return &LIFOBuffer[T]{
		capacity: capacity,
		buf:      make([]T, capacity),
		count:    make(map[T]int, capacity),
	}
}

// keeps track of the number of total values as well as the
// number of occurrences of each value in the buffer
func (b *LIFOBuffer[T]) incCount(v T) {
	b.n++
	b.count[v]++
}

// decrements counts and removes the value from the map if
// the count reaches zero
func (b *LIFOBuffer[T]) decCount(v T) {
	b.n--
	b.count[v]--
	if b.count[v] == 0 {
		delete(b.count, v)
	}
}

// Push a value to the buffer. If the buffer is full, the oldest
// value will be overwritten.
func (b *LIFOBuffer[T]) Push(v T) {
	if b.n == b.capacity {
		b.decCount(b.buf[b.i])
	}
	b.incCount(v)
	b.buf[b.i] = v
	b.i = (b.i + 1) % b.capacity
}

// Pop a value from the buffer. If the buffer is empty, the
// second return value is false.
func (b *LIFOBuffer[T]) Pop() (T, bool) {
	if b.n == 0 {
		return *new(T), false
	}
	b.i = (b.i - 1 + b.capacity) % b.capacity
	v := b.buf[b.i]
	b.decCount(v)
	return v, true
}

// Contains returns true if the value is in the buffer.
func (b *LIFOBuffer[T]) Contains(v T) bool {
	_, ok := b.count[v]
	return ok
}

// Len returns the number of values in the buffer.
func (b *LIFOBuffer[T]) Len() int {
	return b.n
}

// Reset the buffer.
func (b *LIFOBuffer[T]) Reset() {
	b.i = 0
	b.n = 0
	b.count = make(map[T]int, b.capacity)
}

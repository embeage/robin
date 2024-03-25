package robin

type Buffer[T comparable] interface {
	Push(v T)
	Pop() (T, bool)
	Contains(v T) bool
	Len() int
	Reset()
}

type node[T comparable] struct {
	v    T
	prev *node[T]
	next *node[T]
}

// Robin is a round-robin data structure for comparable types that
// supports addition and removal of values. It can grow indefinitely,
// see [NewUnbounded], or be bounded by a maximum length, see
// [NewBounded], in which case an optional buffer can be provided.
// The purpose of the buffer is to automatically replace values that
// are removed from the robin, see the [Buffer] interface and the
// [WithBuffer] option.
//
// The values in the robin must be unique. Duplicate values are ignored.
// All operations are O(1), or O(n) for variadic operations where n is
// the number of arguments. A buffer implementation, [LIFOBuffer], is
// provided in the package. If a custom buffer is used, the time
// complexity of the operations may be affected.
//
// Robin uses a map internally so if [T] is a complex type with poor
// hashing and comparison performance, the Robin performance will
// be affected. Therefore it is recommended to use simple types.
//
// Robin is not thread-safe by default. A mutex or some other form of
// synchronization should be used for concurrent access.
type Robin[T comparable] struct {
	next  *node[T]
	nodes map[T]*node[T]

	maxLen int
	buffer Buffer[T]
}

// Create a new unbounded [Robin].
func NewUnbounded[T comparable]() *Robin[T] {
	return &Robin[T]{nodes: make(map[T]*node[T])}
}

type BoundedOption[T comparable] func(*Robin[T])

// WithBuffer sets the buffer for a bounded [Robin]. When the [Robin]
// is full, added values will be pushed to the buffer. When a value is
// removed, it will be replaced by popping a value from the buffer if
// one is available.
func WithBuffer[T comparable](buffer Buffer[T]) BoundedOption[T] {
	return func(r *Robin[T]) {
		r.buffer = buffer
	}
}

// Create a new bounded [Robin] with a maximum length. An optional
// buffer can be provided with the [WithBuffer] option. If the length
// is negative or zero, an unbounded [Robin] will be returned and
// any option will be ignored.
func NewBounded[T comparable](len int, options ...BoundedOption[T]) *Robin[T] {
	if len <= 0 {
		return NewUnbounded[T]()
	}
	r := &Robin[T]{nodes: make(map[T]*node[T], len), maxLen: len}
	for _, option := range options {
		option(r)
	}
	return r
}

// attach added nodes to the circular doubly linked list between the
// next node and its predecessor and update next node to the new head
func (r *Robin[T]) attach(head, tail *node[T]) {
	if head == nil {
		return
	}

	if r.next == nil {
		head.prev = tail
		tail.next = head
		r.next = head
		return
	}

	prev := r.next.prev
	next := r.next
	head.prev = prev
	tail.next = next
	prev.next = head
	next.prev = tail
	r.next = head
}

// Add values to the robin between current position. A
// subsequent call to [Next] will return the first added value.
// If the robin is bounded and full and a buffer is provided, the
// values are pushed to the buffer, otherwise they are ignored.
// Values already in the robin or in the buffer are ignored.
func (r *Robin[T]) Add(vs ...T) {
	if r.maxLen > 0 && len(r.nodes) == r.maxLen && r.buffer == nil {
		return
	}

	var (
		head *node[T]
		tail *node[T]
	)

	for _, v := range vs {
		if _, ok := r.nodes[v]; ok {
			continue
		}
		if r.maxLen > 0 && len(r.nodes) == r.maxLen {
			if r.buffer == nil {
				break
			}
			if !r.buffer.Contains(v) {
				r.buffer.Push(v)
			}
			continue
		}
		node := &node[T]{v: v}
		r.nodes[v] = node
		if head == nil {
			head = node
			tail = head
			continue
		}
		node.prev = tail
		tail.next = node
		tail = node
	}

	r.attach(head, tail)
}

// removes a node from the circular doubly linked list
func (r *Robin[T]) unlink(node *node[T]) {
	// reset if removed value was the last
	if node == node.next {
		r.next = nil
		return
	}

	node.prev.next = node.next
	node.next.prev = node.prev

	// advance robin if removed value belonged to next node
	if node == r.next {
		r.next = node.next
	}
}

// replaces a removed value with a value from the buffer if possible;
// if not, the node has to be unlinked from the robin
func (r *Robin[T]) replaceValue(node *node[T]) bool {
	if r.buffer != nil {
		if v, ok := r.buffer.Pop(); ok {
			node.v = v
			r.nodes[v] = node
			return true
		}
	}
	return false
}

// Remove values from the robin. If the robin is bounded and there is a
// non-empty buffer, each removed value will be replaced by popping a
// value from the buffer. Values not in the robin, including values in
// the buffer, are ignored.
func (r *Robin[T]) Remove(vs ...T) {
	for _, v := range vs {
		if node, ok := r.nodes[v]; ok {
			delete(r.nodes, v)
			if !r.replaceValue(node) {
				r.unlink(node)
			}
		}
	}
}

// Next returns the next value in the robin. If the robin is empty, the
// second return value is false.
func (r *Robin[T]) Next() (T, bool) {
	if r.next == nil {
		return *new(T), false
	}
	v := r.next.v
	r.next = r.next.next
	return v, true
}

// Contains returns true if the value is in the robin.
func (r *Robin[T]) Contains(v T) bool {
	_, ok := r.nodes[v]
	return ok
}

// BufferContains returns true if the value is in the buffer.
// If there is no buffer, false is returned.
func (r *Robin[T]) BufferContains(v T) bool {
	if r.buffer == nil {
		return false
	}
	return r.buffer.Contains(v)
}

// Len returns the number of values in the robin.
func (r *Robin[T]) Len() int {
	return len(r.nodes)
}

// BufferLen returns the number of values in the buffer.
// If there is no buffer, 0 is returned.
func (r *Robin[T]) BufferLen() int {
	if r.buffer == nil {
		return 0
	}
	return r.buffer.Len()
}

// Reset the robin. If there is a buffer, it is reset as well.
func (r *Robin[T]) Reset() {
	r.next = nil
	if r.buffer == nil {
		r.nodes = make(map[T]*node[T])
		return
	} 
	r.buffer.Reset()
	r.nodes = make(map[T]*node[T], r.maxLen)
}

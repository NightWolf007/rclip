package store

// Element represents a single element of the list.
type Element struct {
	Value []byte

	next *Element
	prev *Element
}

// Next returns the next element from the list.
func (e *Element) Next() *Element {
	return e.next
}

// Prev returns the previous element from the list.
func (e *Element) Prev() *Element {
	return e.prev
}

// List represents list struct.
type List struct {
	head *Element
	tail *Element
	size uint
}

// NewList builds new list.
func NewList() *List {
	return &List{}
}

// First returns the first element from the list.
func (l List) First() *Element {
	return l.head
}

// Last returns the last element from the list.
func (l List) Last() *Element {
	return l.tail
}

// Len returns length of the list.
func (l List) Len() uint {
	return l.size
}

// Push pushes new element at the back of the list.
func (l *List) Push(val []byte) *Element {
	elem := &Element{
		Value: val,
		prev:  l.tail,
	}

	if l.tail != nil {
		l.tail.next = elem
	}

	if l.head == nil {
		l.head = elem
	}

	l.tail = elem
	l.size++

	return elem
}

// Pop pops element from the front of the list.
func (l *List) Pop() *Element {
	if l.head == nil {
		return nil
	}

	elem := l.head
	l.head = l.head.next
	l.size--

	if l.head != nil {
		l.head.prev = nil
	}

	return elem
}

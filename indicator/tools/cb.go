package tools

import (
	"fmt"
	"strings"
)

// Assert Queue implementation

// Queue holds values in a slice.
type Queue struct {
	Vals    []interface{}
	Start   int
	End     int
	full    bool
	MaxSize int
	size    int
}

// New instantiates a new empty queue with the specified size of maximum number of elements that it can hold.
// This max size of the buffer cannot be changed.
func New(maxSize int) *Queue {
	if maxSize < 1 {
		panic("Invalid maxSize, should be at least 1")
	}
	queue := &Queue{MaxSize: maxSize}
	queue.Clear()
	return queue
}

// Enqueue adds a value to the end of the queue
func (queue *Queue) Enqueue(value interface{}) {
	if queue.Full() {
		queue.Dequeue()
	}
	queue.Vals[queue.End] = value //138.8 ns/op-17.64 ns/op
	//copy(queue.values[queue.end:], value)
	queue.End = queue.End + 1
	if queue.End >= queue.MaxSize {
		queue.End = 0
	}
	if queue.End == queue.Start {
		queue.full = true
	}

	queue.size = queue.calculateSize()
}

// Dequeue removes first element of the queue and returns it, or nil if queue is empty.
// Second return parameter is true, unless the queue was empty and there was nothing to dequeue.
func (queue *Queue) Dequeue() (value interface{}, ok bool) {
	if queue.Empty() {
		return nil, false
	}

	value, ok = queue.Vals[queue.Start], true

	if value != nil {
		queue.Vals[queue.Start] = nil
		queue.Start = queue.Start + 1
		if queue.Start >= queue.MaxSize {
			queue.Start = 0
		}
		queue.full = false
	}

	queue.size = queue.size - 1

	return
}

// Peek returns first element of the queue without removing it, or nil if queue is empty.
// Second return parameter is true, unless the queue was empty and there was nothing to peek.
func (queue *Queue) Peek() (value interface{}, ok bool) {
	if queue.Empty() {
		return nil, false
	}
	return queue.Vals[queue.Start], true
}

// Empty returns true if queue does not contain any elements.
func (queue *Queue) Empty() bool {
	return queue.Size() == 0
}

// Full returns true if the queue is full, i.e. has reached the maximum number of elements that it can hold.
func (queue *Queue) Full() bool {
	return queue.Size() == queue.MaxSize
}

// Size returns number of elements within the queue.
func (queue *Queue) Size() int {
	return queue.size
}

// Clear removes all elements from the queue.
func (queue *Queue) Clear() {
	queue.Vals = make([]interface{}, queue.MaxSize, queue.MaxSize)
	queue.Start = 0
	queue.End = 0
	queue.full = false
	queue.size = 0
}

// Values returns all elements in the queue (FIFO order).
func (queue *Queue) Values() []interface{} {
	values := make([]interface{}, queue.Size(), queue.Size())
	for i := 0; i < queue.Size(); i++ {
		values[i] = queue.Vals[(queue.Start+i)%queue.MaxSize]
	}
	return values
}

// String returns a string representation of container
func (queue *Queue) String() string {
	str := "CircularBuffer\n"
	var values []string
	for _, value := range queue.Values() {
		values = append(values, fmt.Sprintf("%v", value))
	}
	str += strings.Join(values, ", ")
	return str
}

// Check that the index is within bounds of the list
func (queue *Queue) withinRange(index int) bool {
	return index >= 0 && index < queue.size
}

func (queue *Queue) calculateSize() int {
	if queue.End < queue.Start {
		return queue.MaxSize - queue.Start + queue.End
	} else if queue.End == queue.Start {
		if queue.full {
			return queue.MaxSize
		}
		return 0
	}
	return queue.End - queue.Start
}

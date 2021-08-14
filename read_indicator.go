package lrc

import "sync/atomic"

// readIndicator uses an int32 as the counter to count the readers
type readIndicator struct {
	count *int32
}

// newReadIndicator creates a readIndicator
func newReadIndicator() *readIndicator {
	r := &readIndicator{
		count: new(int32),
	}
	*r.count = 0
	return r
}

// arrive should be called by the reader goroutine when start reading, increments the counter
func (r *readIndicator) arrive() {
	atomic.AddInt32(r.count, 1)
}

// depart should be called by the reader goroutine when finish reading, decrements the counter
func (r *readIndicator) depart() {
	atomic.AddInt32(r.count, -1)
}

// isEmpty returns true if no readers
func (r *readIndicator) isEmpty() bool {
	return atomic.LoadInt32(r.count) == 0
}

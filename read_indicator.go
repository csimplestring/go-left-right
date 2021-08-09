package lrc

import "sync/atomic"

type ReadIndicator struct {
	count *int32
}

func newReadIndicator() *ReadIndicator {
	r := &ReadIndicator{
		count: new(int32),
	}
	*r.count = 0
	return r
}

func (r *ReadIndicator) arrive() {
	atomic.AddInt32(r.count, 1)
}

func (r *ReadIndicator) depart() {
	atomic.AddInt32(r.count, -1)
}

func (r *ReadIndicator) isEmpty() bool {
	return atomic.LoadInt32(r.count) == 0
}

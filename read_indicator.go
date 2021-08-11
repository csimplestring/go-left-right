package lrc

import "sync/atomic"

type readIndicator struct {
	count *int32
}

func newReadIndicator() *readIndicator {
	r := &readIndicator{
		count: new(int32),
	}
	*r.count = 0
	return r
}

func (r *readIndicator) arrive() {
	atomic.AddInt32(r.count, 1)
}

func (r *readIndicator) depart() {
	atomic.AddInt32(r.count, -1)
}

func (r *readIndicator) isEmpty() bool {
	return atomic.LoadInt32(r.count) == 0
}

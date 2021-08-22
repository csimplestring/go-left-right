package primitive

import (
	"sync/atomic"

	g "github.com/linxGnu/go-adder"
)

type ReadIndicator interface {
	arrive()
	depart()
	isEmpty() bool
}

type distributedAtomicReadIndicator struct {
	ingress g.LongAdder
	egress  g.LongAdder
}

func newDistributedAtomicReadIndicator() *distributedAtomicReadIndicator {
	return &distributedAtomicReadIndicator{
		ingress: g.NewJDKAdder(),
		egress:  g.NewJDKAdder(),
	}
}

func (d *distributedAtomicReadIndicator) arrive() {
	d.ingress.Inc()
}

func (d *distributedAtomicReadIndicator) depart() {
	d.egress.Inc()
}

func (d *distributedAtomicReadIndicator) isEmpty() bool {
	return d.egress.Sum() == d.ingress.Sum()
}

// atomicReadIndicator uses an int32 as the counter to count the readers
type atomicReadIndicator struct {
	count *int32
}

// newAtomicReadIndicator creates a readIndicator
func newAtomicReadIndicator() *atomicReadIndicator {
	r := &atomicReadIndicator{
		count: new(int32),
	}
	*r.count = 0
	return r
}

// arrive should be called by the reader goroutine when start reading, increments the counter
func (r *atomicReadIndicator) arrive() {
	atomic.AddInt32(r.count, 1)
}

// depart should be called by the reader goroutine when finish reading, decrements the counter
func (r *atomicReadIndicator) depart() {
	atomic.AddInt32(r.count, -1)
}

// isEmpty returns true if no readers
func (r *atomicReadIndicator) isEmpty() bool {
	return atomic.LoadInt32(r.count) == 0
}

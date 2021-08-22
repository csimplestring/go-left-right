package primitive

import (
	"sync/atomic"

	g "github.com/linxGnu/go-adder"
)

// ReadIndicator defines the top level interface to record how many readers.
type ReadIndicator interface {
	arrive()
	depart()
	isEmpty() bool
}

// ingressEgressReaderIndicator uses the LongAdder to keep track of how many readers arrive and depart.
// An ingress/egress technique is used.
// A java implementation: https://github.com/pramalhe/ConcurrencyFreaks/blob/master/Java/com/concurrencyfreaks/readindicators/RIIngressEgressLongAdder.java
type ingressEgressReaderIndicator struct {
	ingress g.LongAdder
	egress  g.LongAdder
}

// newDistributedAtomicReadIndicator creates a new distributedAtomicReadIndicator.
func newDistributedAtomicReadIndicator() *ingressEgressReaderIndicator {
	return &ingressEgressReaderIndicator{
		ingress: g.NewJDKAdder(),
		egress:  g.NewJDKAdder(),
	}
}

// arrive indicates a new reader enters.
func (d *ingressEgressReaderIndicator) arrive() {
	d.ingress.Inc()
}

// depart indicates a reader left.
func (d *ingressEgressReaderIndicator) depart() {
	d.egress.Inc()
}

// isEmpty shows if all the readers depart.
func (d *ingressEgressReaderIndicator) isEmpty() bool {
	// the order is very important.
	// the LongAdder is sequentially consitent only if you use Inc and Sum, see: http://concurrencyfreaks.blogspot.com/2013/09/longadder-is-not-sequentially-consistent.html
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

package primitive

import (
	"runtime"
	"sync/atomic"
)

const ReadOnLeft int32 = -1
const ReadOnRight int32 = 1

// LeftRightPrimitive provides the basic core of the leftt-right pattern.
type LeftRightPrimitive struct {
	// readIndicators is an array of 2 read-indicators, counting the reader numbers on the left/right instance
	readIndicators [2]ReadIndicator
	// versionIndex is the index for readIndicators, 0 means reading on left, 1 means reading on right
	versionIndex *int32
	// sideToRead represents which instance to read
	sideToRead *int32
}

// New creates a LeftRightPrimitive
func New() *LeftRightPrimitive {

	m := &LeftRightPrimitive{
		readIndicators: [2]ReadIndicator{
			newDistributedAtomicReadIndicator(),
			newDistributedAtomicReadIndicator(),
		},
		versionIndex: new(int32),
		sideToRead:   new(int32),
	}

	// starts reading on the left side
	*m.versionIndex = 0
	*m.sideToRead = ReadOnLeft
	return m
}

// readerArrive shall be called by the reader goroutine before start reading
func (lr *LeftRightPrimitive) readerArrive() int {
	idx := atomic.LoadInt32(lr.versionIndex)
	lr.readIndicators[idx].arrive()
	return int(idx)
}

// readerDepart shall be called by the reader goroutine after finish reading
func (lr *LeftRightPrimitive) readerDepart(localVI int) {
	lr.readIndicators[localVI].depart()
}

// writerToggleVersionAndWait shall be called by a single writer goroutine when applying the modification
func (lr *LeftRightPrimitive) writerToggleVersionAndWait() {

	localVI := atomic.LoadInt32(lr.versionIndex)
	prevVI := int(localVI % 2)
	nextVI := int((localVI + 1) % 2)

	// waiting for all the readers on the current side (the same side where the writer is) to complete
	for !lr.readIndicators[nextVI].isEmpty() {
		runtime.Gosched()
	}

	// toggle the version index, so all the following readers start reading on the opposite side
	atomic.StoreInt32(lr.versionIndex, int32(nextVI))

	// waiting for all the reader on the previous side (the opposite side where the writer was) to complete
	for !lr.readIndicators[prevVI].isEmpty() {
		runtime.Gosched()
	}
}

// ApplyReadFn applies read operation on the chosen instance, oh, I really need generics, interface type is ugly
func (lr *LeftRightPrimitive) ApplyReadFn(l interface{}, r interface{}, fn func(interface{})) {

	lvi := lr.readerArrive()

	which := atomic.LoadInt32(lr.sideToRead)
	if which == ReadOnLeft {
		fn(l)
	} else {
		fn(r)
	}

	lr.readerDepart(lvi)
	return
}

// ApplyWriteFn applies write operation on the chosen instance, write operation is done twice, on the left and right
// instance respectively, this might make writing longer, but the readers are wait-free.
func (lr *LeftRightPrimitive) ApplyWriteFn(l interface{}, r interface{}, fn func(interface{})) {

	side := atomic.LoadInt32(lr.sideToRead)
	if side == ReadOnLeft {
		// write on right
		fn(r)
		atomic.StoreInt32(lr.sideToRead, ReadOnRight)
		lr.writerToggleVersionAndWait()
		fn(l)
	} else if side == ReadOnRight {
		// write on left
		fn(l)
		atomic.StoreInt32(lr.sideToRead, ReadOnLeft)
		lr.writerToggleVersionAndWait()
		fn(r)
	} else {
		panic("illegal state: you can only read on LEFT or RIGHT")
	}
}

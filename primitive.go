package lrc

import (
	"runtime"
	"sync/atomic"
)

const ReadOnLeft int32 = -1
const ReadOnRight int32 = 1

type LeftRightPrimitive struct {
	readIndicators [2]*readIndicator
	versionIndex   *int32
	sideToRead     *int32
}

func New() *LeftRightPrimitive {

	m := &LeftRightPrimitive{
		readIndicators: [2]*readIndicator{
			newReadIndicator(),
			newReadIndicator(),
		},
		versionIndex: new(int32),
		sideToRead:   new(int32),
	}

	*m.versionIndex = 0
	*m.sideToRead = ReadOnLeft
	return m
}

func (lr *LeftRightPrimitive) ReaderArrive() int {
	idx := atomic.LoadInt32(lr.versionIndex)
	lr.readIndicators[idx].arrive()
	return int(idx)
}

func (lr *LeftRightPrimitive) ReaderDepart(localVI int) {
	lr.readIndicators[localVI].depart()
}

func (lr *LeftRightPrimitive) WriterToggleVersionAndWait() {

	localVI := atomic.LoadInt32(lr.versionIndex)
	prevVI := int(localVI % 2)
	nextVI := int((localVI + 1) % 2)

	for !lr.readIndicators[nextVI].isEmpty() {
		runtime.Gosched()
	}

	atomic.StoreInt32(lr.versionIndex, int32(nextVI))

	for !lr.readIndicators[prevVI].isEmpty() {
		runtime.Gosched()
	}
}

func (lr *LeftRightPrimitive) ApplyReadFn(l interface{}, r interface{}, fn func(interface{})) {

	lvi := lr.ReaderArrive()

	which := atomic.LoadInt32(lr.sideToRead)
	if which == ReadOnLeft {
		fn(l)
	} else {
		fn(r)
	}

	lr.ReaderDepart(lvi)
	return
}

func (lr *LRMap) ApplyWriteFn(l interface{}, r interface{}, fn func(interface{})) {

	side := atomic.LoadInt32(lr.sideToRead)
	if side == ReadOnLeft {
		// write on right
		fn(r)
		atomic.StoreInt32(lr.sideToRead, ReadOnRight)
		lr.WriterToggleVersionAndWait()
		fn(l)
	} else if side == ReadOnRight {
		// write on left
		fn(l)
		atomic.StoreInt32(lr.sideToRead, ReadOnLeft)
		lr.WriterToggleVersionAndWait()
		fn(r)
	} else {
		panic("illegal state: you can only read on LEFT or RIGHT")
	}
}

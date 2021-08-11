package lrc

import (
	"runtime"
	"sync/atomic"
)

const ReadOnLeft int32 = -1
const ReadOnRight int32 = 1

type LRMap struct {
	readIndicators [2]*readIndicator
	versionIndex   *int32
	sideToRead     *int32

	left  map[int]int
	right map[int]int
}

func New() *LRMap {

	m := &LRMap{
		readIndicators: [2]*readIndicator{
			newReadIndicator(),
			newReadIndicator(),
		},
		versionIndex: new(int32),
		sideToRead:   new(int32),

		left:  make(map[int]int),
		right: make(map[int]int),
	}

	*m.versionIndex = 0
	*m.sideToRead = ReadOnLeft
	return m
}

func (lr *LRMap) arrive() int {
	idx := atomic.LoadInt32(lr.versionIndex)
	lr.readIndicators[idx].arrive()
	return int(idx)
}

func (lr *LRMap) depart(localVI int) {
	lr.readIndicators[localVI].depart()
}

func (lr *LRMap) toggleVersionAndWait() {

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

func (lr *LRMap) Get(k int) (val int, exist bool) {

	lvi := lr.arrive()

	which := atomic.LoadInt32(lr.sideToRead)
	if which == ReadOnLeft {
		val, exist = lr.left[k]
	} else {
		val, exist = lr.right[k]
	}

	lr.depart(lvi)
	return
}

func (lr *LRMap) Put(key, val int) {

	side := atomic.LoadInt32(lr.sideToRead)
	if side == ReadOnLeft {
		// write on right
		lr.right[key] = val
		atomic.StoreInt32(lr.sideToRead, ReadOnRight)
		lr.toggleVersionAndWait()
		lr.left[key] = val
	} else if side == ReadOnRight {
		// write on left
		lr.left[key] = val
		atomic.StoreInt32(lr.sideToRead, ReadOnLeft)
		lr.toggleVersionAndWait()
		lr.right[key] = val
	} else {
		panic("illegal state: you can only read on LEFT or RIGHT")
	}
}

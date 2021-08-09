package lrc

import (
	"fmt"
	"sync"
	"sync/atomic"
)

const read_left int32 = -1
const read_right int32 = 1

type LRMap struct {
	readIndicators [2]*ReadIndicator
	versionIndex   *int32
	leftRight      *int32

	left  map[int]int
	right map[int]int

	wm sync.Mutex
}

func New() *LRMap {

	m := &LRMap{
		readIndicators: [2]*ReadIndicator{
			newReadIndicator(),
			newReadIndicator(),
		},
		versionIndex: new(int32),
		leftRight:    new(int32),

		left:  make(map[int]int),
		right: make(map[int]int),

		wm: sync.Mutex{},
	}

	*m.versionIndex = 0
	*m.leftRight = read_left
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

	}

	atomic.StoreInt32(lr.versionIndex, int32(nextVI))

	for !lr.readIndicators[prevVI].isEmpty() {

	}
}

func (lr *LRMap) Get(k int) (val int, exist bool) {

	lvi := lr.arrive()

	which := atomic.LoadInt32(lr.leftRight)
	if which == read_left {
		val, exist = lr.left[k]
	} else {
		val, exist = lr.right[k]
	}

	lr.depart(lvi)
	return
}

func (lr *LRMap) Put(key, val int) {

	which := atomic.LoadInt32(lr.leftRight)
	if which == read_left {
		// write on right
		lr.right[key] = val
		atomic.StoreInt32(lr.leftRight, read_right)
		lr.toggleVersionAndWait()
		lr.left[key] = val
	} else if which == read_right {
		// write on left
		lr.left[key] = val
		atomic.StoreInt32(lr.leftRight, read_left)
		lr.toggleVersionAndWait()
		lr.right[key] = val
	} else {
		fmt.Println("fuuuu")
	}
}

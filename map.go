package lrc

import "github.com/csimplestring/go-left-right/primitive"

// LRMap utilises the left-right pattern to handle concurrent read-write.
type LRMap struct {
	*primitive.LeftRightPrimitive[map[int]int]

	left  map[int]int
	right map[int]int
}

func newIntMap() *LRMap {

	m := &LRMap{
		left:  make(map[int]int),
		right: make(map[int]int),
	}

	m.LeftRightPrimitive = primitive.New[map[int]int]()

	return m
}

func (lr *LRMap) Get(k int) (val int, exist bool) {

	lr.ApplyReadFn(lr.left, lr.right, func(instance map[int]int) {
		val, exist = instance[k]
	})

	return
}

func (lr *LRMap) Put(key, val int) {
	lr.ApplyWriteFn(lr.left, lr.right, func(instance map[int]int) {
		m := instance
		m[key] = val
	})
}

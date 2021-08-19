package lrmap

import "github.com/csimplestring/go-left-right/primitive"

// LRMap utilises the left-right pattern to handle concurrent read-write.
type LRMap struct {
	*primitive.LeftRightPrimitive

	left  map[int]int
	right map[int]int
}

// NewIntMap creates a default LRMap
func NewIntMap() *LRMap {

	m := &LRMap{
		left:  make(map[int]int),
		right: make(map[int]int),
	}

	m.LeftRightPrimitive = primitive.New()

	return m
}

func (lr *LRMap) Get(k int) (val int, exist bool) {

	lr.ApplyReadFn(lr.left, lr.right, func(instance interface{}) {
		m := instance.(map[int]int)
		val, exist = m[k]
	})

	return
}

func (lr *LRMap) Put(key, val int) {
	lr.ApplyWriteFn(lr.left, lr.right, func(instance interface{}) {
		m := instance.(map[int]int)
		m[key] = val
	})
}

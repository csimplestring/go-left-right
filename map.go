package lrc

// LRMap utilises the left-right pattern to handle concurrent read-write.
type LRMap struct {
	*LeftRightPrimitive

	left  map[int]int
	right map[int]int
}

func newIntMap() *LRMap {

	m := &LRMap{
		left:  make(map[int]int),
		right: make(map[int]int),
	}

	m.LeftRightPrimitive = New()

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

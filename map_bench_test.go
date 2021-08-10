package lrc

import (
	"math/rand"
	"sync"
	"testing"
)

type testMap interface {
	Put(k, v int)
	Get(k int) (int, bool)
}

type LockMap struct {
	m    map[int]int
	lock sync.RWMutex
}

func (l *LockMap) Put(k, v int) {
	l.lock.Lock()
	l.m[k] = v
	l.lock.Unlock()
}

func (l *LockMap) Get(k int) (int, bool) {
	l.lock.RLock()
	v, ok := l.m[k]
	l.lock.RUnlock()
	return v, ok
}

func InitLockMap(num int) *LockMap {
	lockmap := &LockMap{
		m:    make(map[int]int),
		lock: sync.RWMutex{},
	}

	for i := 0; i < num; i++ {
		lockmap.Put(i, i)
	}

	return lockmap
}

func InitLRMap(num int) *LRMap {
	lrmap := New()

	for i := 0; i < num; i++ {
		lrmap.Put(i, i)
	}

	return lrmap
}

func BenchmarkLRMap_Write(b *testing.B) {
	lrmap := InitLRMap(0)

	for i := 0; i < b.N; i++ {
		k := rand.Intn(10000)
		lrmap.Put(k, k)
	}
}

func BenchmarkLockMap_Write(b *testing.B) {
	lockmap := InitLockMap(0)

	for i := 0; i < b.N; i++ {
		k := rand.Intn(10000)
		lockmap.Put(k, k)
	}
}

func BenchmarkLRMap_Read(b *testing.B) {
	lrmap := InitLRMap(1000000)

	for i := 0; i < b.N; i++ {
		k := rand.Intn(1000000)
		lrmap.Get(k)
	}
}

func BenchmarkLockMap_Read(b *testing.B) {
	lockmap := InitLockMap(1000000)

	for i := 0; i < b.N; i++ {
		k := rand.Intn(1000000)
		lockmap.Get(k)
	}
}

func run(m testMap, reader int) {
	wg := sync.WaitGroup{}

	for k := 0; k < reader; k++ {
		wg.Add(1)
		go func() {
			for i := 0; i < 1000000; i++ {
				m.Get(rand.Intn(10000))
			}
			wg.Done()
		}()
	}

	wg.Add(1)
	go func() {
		for i := 0; i < 1000000; i++ {
			k := rand.Intn(10000)
			m.Put(k, k)
		}
		wg.Done()
	}()

	wg.Wait()
}

func BenchmarkLRMap_Read_Write_5_1(b *testing.B) {
	m := InitLRMap(1000000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		run(m, 5)
	}
}

func BenchmarkLockMap_Read_Write_5_1(b *testing.B) {
	m := InitLockMap(1000000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		run(m, 5)
	}
}

func BenchmarkLRMap_Read_Write_10_1(b *testing.B) {
	m := InitLRMap(1000000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		run(m, 10)
	}
}

func BenchmarkLockMap_Read_Write_10_1(b *testing.B) {
	m := InitLockMap(1000000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		run(m, 10)
	}
}

func BenchmarkLRMap_Read_Write_50_1(b *testing.B) {
	m := InitLRMap(1000000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		run(m, 50)
	}
}

func BenchmarkLockMap_Read_Write_50_1(b *testing.B) {
	m := InitLockMap(1000000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		run(m, 50)
	}
}

func BenchmarkLRMap_Read_Write_100_1(b *testing.B) {
	m := InitLRMap(1000000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		run(m, 100)
	}
}

func BenchmarkLockMap_Read_Write_100_1(b *testing.B) {
	m := InitLockMap(1000000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		run(m, 100)
	}
}

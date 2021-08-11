package lrc

import (
	"math/rand"
	"sync"
	"testing"
)

func TestLRMap(t *testing.T) {
	lrmap := newIntMap()

	wg := sync.WaitGroup{}

	for i := 0; i < 100000; i++ {
		wg.Add(1)
		go func() {
			lrmap.Get(rand.Intn(10000))
			wg.Done()
		}()
	}

	wg.Add(1)
	go func() {
		k := rand.Intn(10000)
		lrmap.Put(k, k)
		wg.Done()
	}()

	wg.Wait()
}

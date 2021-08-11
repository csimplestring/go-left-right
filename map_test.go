package lrc

import (
	"testing"
)

func TestLRMap(t *testing.T) {
	lrmap := newIntMap()

	_, exist := lrmap.Get(1)
	if exist {
		t.Error("should not exist")
	}

	lrmap.Put(1, 1)
	v, exist := lrmap.Get(1)
	if v != 1 {
		t.Error("not equal")
	}
	if !exist {
		t.Error("should exist")
	}
}

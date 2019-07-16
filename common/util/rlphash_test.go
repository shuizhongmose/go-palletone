package util

import "testing"

type A struct {
	Data uint
}

func TestRlpHash(t *testing.T) {
	for i := 0; i < 10; i++ {
		a := &A{Data: uint(i)}
		hash := RlpHash(a)
		t.Logf("Number:%d,Hash:%x,Hash2:%x", i, hash, RlpHash(i))
	}
}

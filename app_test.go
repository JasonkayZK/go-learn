package main

import (
	"go.uber.org/goleak"
	"testing"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
	//m.Run()
}

func TestA(t *testing.T) {
	leak()
}

func leak() {
	ch := make(chan struct{})
	go func() {
		ch <- struct{}{}
	}()
}

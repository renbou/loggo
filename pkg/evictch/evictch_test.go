package evictch

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const testCapacity = 2

func Test_Chan_NoWait(t *testing.T) {
	ch := NewChan[int](testCapacity)

	// Write 0..capacity, which should throw out 0 in the end
	for i := 0; i < testCapacity+1; i++ {
		ch.Write(i)
	}

	for i := 0; i < testCapacity; i++ {
		got, ok := ch.Read()
		assert.Equal(t, i+1, got)
		assert.True(t, ok)
	}
}

func Test_Chan_Wait(t *testing.T) {
	ch := NewChan[int](testCapacity)

	expect := 0
	var got int
	var ok bool

	// Launch a Read which should wait for a Write
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		got, ok = ch.Read()
	}()

	time.Sleep(time.Millisecond * 20)
	ch.Write(expect)

	wg.Wait()
	assert.Equal(t, expect, got)
	assert.True(t, ok)
}

func Test_Chan_Close(t *testing.T) {
	ch := NewChan[int](testCapacity)

	var ok bool

	// Launch a Read which should wait for a Close
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, ok = ch.Read()
	}()

	time.Sleep(time.Millisecond * 20)
	ch.Close()

	wg.Wait()
	assert.False(t, ok)
}

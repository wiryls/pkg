package flow

import (
	"context"
	"runtime"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFlowAppend(t *testing.T) {
	assert := assert.New(t)

	{
		count := runtime.NumGoroutine()
		waste := func() { time.Sleep(1 * time.Microsecond) }

		f := &Flow{}
		for i := 0; i < 100; i++ {
			f.Push(waste)
		}

		delta := int(f.limit)
		if delta > 100 {
			delta = 100
		}

		assert.Equal(count+delta, runtime.NumGoroutine())

		f.Wait()
		assert.EqualValues(0, f.count)
	}

	{
		ctx, cancel := context.WithCancel(context.Background())
		waste := func() {
			select {
			case <-ctx.Done():
			case <-time.After(1 * time.Microsecond):
			}
		}

		count := runtime.NumGoroutine()

		f := &Flow{}
		for i := 0; i < 1000; i++ {
			f.Push(waste)
		}

		delta := int(f.limit)
		assert.Equal(count+delta, runtime.NumGoroutine())

		cancel()
		f.Wait()
		assert.EqualValues(0, f.count)
	}

	{
		total := 100
		count := uint32(0)
		adder := func() { atomic.AddUint32(&count, 1) }

		f := Flow{}
		for i := 0; i < total; i++ {
			f.Push(adder)
		}

		f.Wait()
		assert.EqualValues(total, count)
	}
}

type Slime struct {
	B []bool
	F func(func())
}

func (s *Slime) Run() {
	if len(s.B) > 33 {
		s.F((&Slime{B: s.B[33:], F: s.F}).Run)
		s.B = s.B[:33]
	}
	for i := range s.B {
		s.B[i] = true
	}
}

func TestSampleTask(t *testing.T) {
	assert := assert.New(t)
	{
		bools := [1000]bool{}
		f := Flow{}
		f.Push((&Slime{B: bools[:], F: f.Push}).Run)
		f.Wait()

		for i := range bools {
			assert.True(bools[i], i)
		}
	}
}

// note: please test this with -race
func TestRace(t *testing.T) {
	assert := assert.New(t)
	{
		f := Flow{}
		for i := 0; i < 10; i++ {
			go func() {
				for x, n := 0, 1000; x < n; x++ {
					f.Push(func() {})
				}
			}()
		}
		for i := 0; i < 10; i++ {
			go f.Wait()
		}
		f.Wait()
		assert.True(true)
	}
}

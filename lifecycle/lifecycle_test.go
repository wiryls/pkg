package lifecycle_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/wiryls/pkg/lifecycle"
)

// SUGGEST using data race detector to run this test.
// `go test -race <THIS_FILE_PATH>`

/////////////////////////////////////////////////////////////////////////////

func TestNoCallback(t *testing.T) {
	assert := assert.New(t)

	srv := lifecycle.LifeCycle{}
	assert.EqualValues(srv.State(), lifecycle.StateStopped)
	assert.Error(srv.Close())
	assert.Error(srv.WhileRunning(nil))
	assert.NoError(srv.Run())
}

/////////////////////////////////////////////////////////////////////////////

type Dummy interface {
	lifecycle.RunnerCloser

	State() lifecycle.State
	HanldePing() error
	HanldeClose() error
}

func NewDummy() Dummy {
	d := &dummy{}
	d.Bind(d)
	return d
}

type dummy struct {
	// lifecycle
	lifecycle.LifeCycle

	// data
	input chan chan<- bool
}

func (d *dummy) HanldePing() error {
	return d.LifeCycle.WhileRunningChan(func(done <-chan struct{}) error {
		something := make(chan bool)

		select {
		case d.input <- something:
		case <-done:
		}

		select {
		case <-something:
		case <-done:
		}

		return nil
	})
}

func (d *dummy) HanldeClose() error {
	return d.LifeCycle.WhileRunning(func() error {
		return d.LifeCycle.CloseAsync()
	})
}

func (d *dummy) BeforeRunning() error {
	d.input = make(chan chan<- bool)
	time.Sleep(10 * time.Millisecond)
	return nil
}

func (d *dummy) AfterRunning() error {
	time.Sleep(10 * time.Millisecond)
	close(d.input)
	d.input = nil
	return nil
}

func (d *dummy) Running(cancel <-chan struct{}) error {

loop:
	for {
		select {
		case in := <-d.input:

			select {
			case in <- true:
				close(in)
			case <-time.After(time.Second):
			}

		case <-cancel:
			break loop
		}
	}

	return nil
}

func TestDummyService(t *testing.T) {
	assert := assert.New(t)
	{ // Not Running
		srv := NewDummy()
		assert.EqualValues(srv.State(), lifecycle.StateStopped)
		assert.Error(srv.Close())
		assert.Error(srv.HanldePing())
		assert.Error(srv.HanldeClose())
	}

	{ // CloseAsync
		srv := NewDummy()
		assert.EqualValues(srv.State(), lifecycle.StateStopped)

		c := make(chan error, 100)
		go func() { c <- srv.Run() }()

		time.Sleep(5 * time.Millisecond)
		assert.EqualValues(srv.State(), lifecycle.StateBooting)

		time.Sleep(10 * time.Millisecond)
		assert.EqualValues(srv.State(), lifecycle.StateRunning)

		for i := 0; i < 98; i++ {
			go func() { c <- srv.HanldePing() }()
		}
		time.Sleep(1 * time.Millisecond)
		go func() { c <- srv.HanldeClose() }()

		time.Sleep(14 * time.Millisecond)
		assert.EqualValues(srv.State(), lifecycle.StateStopped)

	out:
		for {
			select {
			case err := <-c:
				assert.NoError(err)
			default:
				break out
			}
		}

		assert.Error(srv.Close())
		assert.Error(srv.HanldePing())
		assert.Error(srv.HanldeClose())
	}

	{ // Close
		srv := NewDummy()
		assert.EqualValues(srv.State(), lifecycle.StateStopped)

		c := make(chan error)
		go func() { c <- srv.Run() }()

		time.Sleep(5 * time.Millisecond)
		assert.EqualValues(srv.State(), lifecycle.StateBooting)

		time.Sleep(10 * time.Millisecond)
		assert.EqualValues(srv.State(), lifecycle.StateRunning)

		assert.NoError(srv.Close())
		assert.NoError(<-c)

		time.Sleep(15 * time.Millisecond)
		assert.EqualValues(srv.State(), lifecycle.StateStopped)
		assert.Error(srv.Close())
		assert.Error(srv.HanldePing())
	}
}

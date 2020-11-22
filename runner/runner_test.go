package runner_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/wiryls/pkg/runner"
)

// SUGGEST using data race detector to run this test.
// `go test -race <THIS_FILE_PATH>`

/////////////////////////////////////////////////////////////////////////////

func TestNoCallback(t *testing.T) {
	assert := assert.New(t)

	srv := runner.Determination{}
	assert.EqualValues(srv.State(), runner.StateStopped)
	assert.Error(srv.Close())
	assert.Error(srv.WhileRunning(nil))
	assert.NoError(srv.Run())
}

/////////////////////////////////////////////////////////////////////////////

type Dummy interface {
	runner.Runner

	State() runner.State
	HanldePing() error
	HanldeClose() error
}

func NewDummy() Dummy {
	d := &dummy{}
	d.Bind(d)
	return d
}

type dummy struct {
	// runner
	runner.Determination

	// data
	input chan chan<- bool
}

func (d *dummy) HanldePing() error {
	return d.Determination.WhileRunning(func(done <-chan struct{}) error {
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
	return d.Determination.WhileRunning(func(<-chan struct{}) error {
		return d.Determination.CloseAsync()
	})
}

func (d *dummy) BeforeRunning(exit <-chan struct{}) error {
	d.input = make(chan chan<- bool)
	select {
	case <-time.After(10 * time.Millisecond):
	case <-exit:
	}
	return nil
}

func (d *dummy) AfterRunning() error {
	time.Sleep(10 * time.Millisecond)
	close(d.input)
	d.input = nil
	return nil
}

func (d *dummy) Running(exit <-chan struct{}) error {

loop:
	for {
		select {
		case in := <-d.input:

			select {
			case in <- true:
				close(in)
			case <-time.After(time.Second):
			case <-exit:
			}

		case <-exit:
			break loop
		}
	}

	return nil
}

func TestDummyService(t *testing.T) {
	assert := assert.New(t)
	{ // Not Running
		srv := NewDummy()
		assert.EqualValues(srv.State(), runner.StateStopped)
		assert.Error(srv.Close())
		assert.Error(srv.HanldePing())
		assert.Error(srv.HanldeClose())
	}

	{ // CloseAsync
		srv := NewDummy()
		assert.EqualValues(srv.State(), runner.StateStopped)

		c := make(chan error, 100)
		go func() { c <- srv.Run() }()

		time.Sleep(5 * time.Millisecond)
		assert.EqualValues(srv.State(), runner.StateBooting)

		time.Sleep(10 * time.Millisecond)
		assert.EqualValues(srv.State(), runner.StateRunning)

		for i := 0; i < 98; i++ {
			go func() { c <- srv.HanldePing() }()
		}
		time.Sleep(1 * time.Millisecond)
		go func() { c <- srv.HanldeClose() }()

		time.Sleep(14 * time.Millisecond)
		assert.EqualValues(srv.State(), runner.StateStopped)

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
		assert.EqualValues(srv.State(), runner.StateStopped)

		c := make(chan error)
		go func() { c <- srv.Run() }()

		time.Sleep(5 * time.Millisecond)
		assert.EqualValues(srv.State(), runner.StateBooting)

		time.Sleep(10 * time.Millisecond)
		assert.EqualValues(srv.State(), runner.StateRunning)

		assert.NoError(srv.Close())
		assert.NoError(<-c)

		time.Sleep(15 * time.Millisecond)
		assert.EqualValues(srv.State(), runner.StateStopped)
		assert.Error(srv.Close())
		assert.Error(srv.HanldePing())
	}
}

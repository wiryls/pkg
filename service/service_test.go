package service_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wiryls/pkg/service"
)

// SUGGEST using data race detector to run this test.
// `go test -race <THIS_FILE_PATH>`

/////////////////////////////////////////////////////////////////////////////

func TestNoCallback(t *testing.T) {
	assert := assert.New(t)

	srv := service.Service{}
	assert.EqualValues(srv.State(), service.StateStopped)
	assert.Error(srv.Close())
	assert.Error(srv.WhileRunning(nil))
	assert.NoError(srv.Run())
}

/////////////////////////////////////////////////////////////////////////////

type Dummy dummy
type dummy struct {
	service.Service
}

func NewDummy() *Dummy {
	d := &dummy{}
	d.Bind(d)
	return (*Dummy)(d)
}

func (d *Dummy) Run() error {
	return d.Service.Run()
}

func (d *Dummy) HanldePing() error {
	return d.Service.WhileRunning(func() error {
		time.Sleep(time.Millisecond)
		return nil
	})
}

func (d *Dummy) HanldeClose() error {
	return d.Service.WhileRunning(func() error {
		return d.Service.CloseAsync()
	})
}

func (d *Dummy) State() service.State { return d.Service.State() }

func (d *Dummy) Close() error { return d.Service.Close() }

func (d *dummy) BeforeRunning() error {
	time.Sleep(10 * time.Millisecond)
	return nil
}

func (d *dummy) Running(cancel <-chan struct{}) error {
	<-cancel
	return nil
}

func (d *dummy) AfterRunning() error {
	time.Sleep(10 * time.Millisecond)
	return nil
}

func TestDummyService(t *testing.T) {
	assert := assert.New(t)
	{ // Not Running
		srv := NewDummy()
		assert.EqualValues(srv.State(), service.StateStopped)
		assert.Error(srv.Close())
		assert.Error(srv.HanldePing())
		assert.Error(srv.HanldeClose())
	}

	{ // CloseAsync
		srv := NewDummy()
		assert.EqualValues(srv.State(), service.StateStopped)

		c := make(chan error, 100)
		go func() { c <- srv.Run() }()

		time.Sleep(5 * time.Millisecond)
		assert.EqualValues(srv.State(), service.StateBooting)

		time.Sleep(10 * time.Millisecond)
		assert.EqualValues(srv.State(), service.StateRunning)

		for i := 0; i < 98; i++ {
			go func() { c <- srv.HanldePing() }()
		}
		time.Sleep(1 * time.Millisecond)
		go func() { c <- srv.HanldeClose() }()

		time.Sleep(14 * time.Millisecond)
		assert.EqualValues(srv.State(), service.StateStopped)

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
		assert.EqualValues(srv.State(), service.StateStopped)

		c := make(chan error)
		go func() { c <- srv.Run() }()

		time.Sleep(5 * time.Millisecond)
		assert.EqualValues(srv.State(), service.StateBooting)

		time.Sleep(10 * time.Millisecond)
		assert.EqualValues(srv.State(), service.StateRunning)

		assert.NoError(srv.Close())
		assert.NoError(<-c)

		time.Sleep(15 * time.Millisecond)
		assert.EqualValues(srv.State(), service.StateStopped)
		assert.Error(srv.Close())
		assert.Error(srv.HanldePing())
	}
}

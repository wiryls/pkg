# lifecycle

The `lifecycle.LifeCycle` in this package is responsible for lifecycle management.

Please see [lifecycle_test](./lifecycle_test.go).

## Sample

```golang

// Service is a service like dummy object. It uses `lifecycle.LifeCycle` to run and close.
type Service interface {
    lifecycle.RunnerCloser

    State() lifecycle.State
    HanldePing() error
    HanldeClose() error
}

// New Service.
func New() Service {
    s := &service{}
    s.Bind(s)
    return s
}

type service struct {
    // lifecycle
    lifecycle.LifeCycle

    // data
    input chan chan<- bool
}

func (s *service) HanldePing() error {
    return s.LifeCycle.WhileRunningChan(func(done <-chan struct{}) error {
        something := make(chan bool)

        select {
        case s.input <- something:
        case <-done:
        }

        select {
        case <-something:
        case <-done:
        }

        return nil
    })
}

func (s *service) HanldeClose() error {
    return s.LifeCycle.WhileRunning(func() error {
        return s.LifeCycle.CloseAsync()
    })
}

func (s *service) BeforeRunning() error {
    s.input = make(chan chan<- bool)
    return nil
}

func (s *service) Running(cancel <-chan struct{}) error {

loop:
    for {
        select {
        case in := <-s.input:

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

func (s *service) AfterRunning() error {
    close(s.input)
    s.input = nil
    return nil
}

```

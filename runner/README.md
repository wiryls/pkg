# runner

The `runner.Determination` in this package is responsible for runner management.

Please see [runner_test](./runner_test.go).

## Sample

```golang
// Service is a service like dummy object. It uses `runner.Determination` to run and close.
type Service interface {
    runner.Runner

    State() runner.State
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
    // runner
    runner.Determination

    // data
    input chan chan<- bool
}

func (s *service) HanldePing() error {
    return s.Determination.WhileRunning(func(exit <-chan struct{}) error {
        something := make(chan bool)

        select {
        case s.input <- something:
        case <-exit:
        }

        select {
        case <-something:
        case <-exit:
        }

        return nil
    })
}

func (s *service) HanldeClose() error {
    return s.Determination.WhileRunning(func(<-chan struct{}) error {
        return s.Determination.CloseAsync()
    })
}

func (s *service) BeforeRunning() error {
    s.input = make(chan chan<- bool)
    return nil
}

func (s *service) Running(exit <-chan struct{}) error {

loop:
    for {
        select {
        case in := <-s.input:

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

func (s *service) AfterRunning() error {
    close(s.input)
    s.input = nil
    return nil
}

```

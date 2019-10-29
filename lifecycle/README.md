# lifecycle

The `lifecycle.LifeCycle` in this package is responsible for lifecycle management.

Please see [lifecycle_test](./lifecycle_test.go).

## Sample

```golang
// Service is a service like dummy object. It uses `lifecycle.LifeCycle` to run and close.
type Service service
type service struct {
    lifecycle.LifeCycle

    input chan chan<- bool
}

func NewService() *Service {
    s := &service{}
    s.Bind(s)
    return (*Service)(s)
}

func (s *Service) Run() error { return s.LifeCycle.Run() }

func (s *Service) State() lifecycle.State { return s.LifeCycle.State() }

func (s *Service) Close() error { return s.LifeCycle.Close() }

func (s *Service) HanldePing() error {
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

func (s *Service) HanldeClose() error {
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
    close(d.input)
    d.input = nil
    return nil
}
```

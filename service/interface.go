package service

// Runnable is the operation of a Service.
type Runnable interface {

	// BeforeRunning, do some initialization.
	BeforeRunning() error

	// Running start some tasks and blocks.
	Running(cancel <-chan struct{}) error

	// AfterRunning, do some cleaning.
	AfterRunning() error
}

// Runner has a method `Run`.
type Runner interface {
	Run() error
}

// Closer has a method `Close`.
type Closer interface {
	Close() error
}

// RunnerCloser has both `Run` and `Close`
type RunnerCloser interface {
	Runner
	Closer
}

package runner

import (
	"context"
)

// Runner has both `Run` and `Close` methods.
type Runner interface {
	Run() error
	Close() error
}

// Runnable contains three operations: BeforeRunning, Running, AfterRunning.
type Runnable interface {

	// BeforeRunning, do some initialization.
	BeforeRunning(exit <-chan struct{}) error

	// Running blocks and starts some tasks.
	Running(exit <-chan struct{}) error

	// AfterRunning, do some cleaning.
	AfterRunning() error
}

// RunnableWithContext is the context version of Runnable.
type RunnableWithContext interface {

	// BeforeRunning, do some initialization.
	BeforeRunning(ctx context.Context) error

	// Running blocks and starts some tasks.
	Running(ctx context.Context) error

	// AfterRunning, do some cleaning.
	AfterRunning() error
}

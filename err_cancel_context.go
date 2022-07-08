package main

import (
	"context"
	"sync"
)

// CancelWithErrFunc is a function that is called to cancel a context with a given error
type CancelWithErrFunc func(error)

// ErrCancelContext provides a similar type of context as you get with context.WithCancel, however this context allows
// an error to be specified in the cancel call.  Contexts that are cancelled with a specified error return that error
// on calls to Err()
type ErrCancelContext struct {
	context.Context
	err error
	mu  *sync.Mutex
}

// Err returns the error that the context was cancelled with
func (ec *ErrCancelContext) Err() error {
	ec.mu.Lock()
	defer ec.mu.Unlock()

	if ec.err != nil {
		return ec.err
	}

	return ec.Context.Err()
}

// NewErrContextWithCancel returns a new context, and a function which can be used to cancel the context
func NewErrContextWithCancel(parent context.Context) (context.Context, CancelWithErrFunc) {
	newCtx, cancelFunc := context.WithCancel(parent)
	errCtx := &ErrCancelContext{
		Context: newCtx,
		err:     nil,
		mu:      &sync.Mutex{},
	}

	return errCtx, func(err error) {
		errCtx.mu.Lock()
		defer errCtx.mu.Unlock()

		errCtx.err = err
		cancelFunc()
	}
}

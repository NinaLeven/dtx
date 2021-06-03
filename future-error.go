package dtx

import (
	"context"
	"dtx/future"
	"errors"
)

var ErrFutureNotError = errors.New("future value is not of type error")

type FutureErrorGetter interface {
	GetError(ctx context.Context) (error, error)
	Chan() chan error
}

type FutureErrorSetter interface {
	SetError(error) error
}

type FutureError interface {
	FutureErrorSetter
	FutureErrorGetter
}

type futureError struct {
	f future.Future
}

func NewFutureError() FutureError {
	return &futureError{
		f: future.NewFuture(),
	}
}

func (e *futureError) SetError(err error) error {
	return e.f.Set(err)
}

func (e *futureError) GetError(ctx context.Context) (error, error) {
	val, err := e.f.Get(ctx)
	if err != nil {
		return nil, err
	}
	errVal, ok := val.(error)
	if !ok {
		return nil, ErrFutureNotError
	}
	return errVal, nil
}

func (e *futureError) Chan() chan error {
	c := make(chan error, 1)
	go func() {
		defer close(c)
		val, err := e.GetError(context.Background())
		if err != nil {
			return
		}
		c <- val
	}()
	return c
}
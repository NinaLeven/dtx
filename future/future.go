package future

import (
	"context"
	"errors"
	"sync/atomic"
)

var ErrAlreadySet = errors.New("value is already set")

type FutureGetter interface {
	Get(ctx context.Context) (interface{}, error)
	Chan() chan interface{}
}

type FutureSetter interface {
	Set(interface{}) error
}

type Future interface {
	FutureSetter
	FutureGetter
}

type future struct {
	isSet, isGotten int32
	c               chan interface{}
	val             interface{}
}

func NewFuture() Future {
	return &future{
		c: make(chan interface{}, 1),
	}
}

func (f *future) Set(val interface{}) error {
	if !atomic.CompareAndSwapInt32(&f.isSet, 0, 1) {
		return ErrAlreadySet
	}
	f.c <- val
	close(f.c)
	return nil
}

func (f *future) Get(ctx context.Context) (interface{}, error) {
	if atomic.CompareAndSwapInt32(&f.isGotten, 1, 1) {
		return f.val, nil
	}
	select {
	case <-ctx.Done():
		return nil, context.Canceled
	case val, ok := <-f.c:
		if !ok {
			return f.Get(ctx)
		}
		if atomic.CompareAndSwapInt32(&f.isGotten, 0, 1) {
			f.val = val
		}
		return f.val, nil
	}
}

func (f *future) Chan() chan interface{} {
	c := make(chan interface{}, 1)
	go func() {
		defer close(c)
		val, err := f.Get(context.Background())
		if err != nil {
			return
		}
		c <- val
	}()
	return c
}

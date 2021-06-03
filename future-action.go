package dtx

import (
	"context"
	"dtx/future"
	"errors"
)

var ErrFutureNotCommitAction = errors.New("future value is not of type commitAction")

type futureCommitAction interface {
	SetCommitAction(action commitAction) error
	Chan() chan commitAction
}

type futureAction struct {
	f future.Future
}

func NewFutureAction() futureCommitAction {
	return &futureAction{
		f: future.NewFuture(),
	}
}

func (e *futureAction) SetCommitAction(err commitAction) error {
	return e.f.Set(err)
}

func (e *futureAction) getCommitAction(ctx context.Context) (commitAction, error) {
	val, err := e.f.Get(ctx)
	if err != nil {
		return commitActionNone, err
	}
	action, ok := val.(commitAction)
	if !ok {
		return commitActionNone, ErrFutureNotError
	}
	return action, nil
}

func (e *futureAction) Chan() chan commitAction {
	c := make(chan commitAction, 1)
	go func() {
		defer close(c)
		val, err := e.getCommitAction(context.Background())
		if err != nil {
			return
		}
		c <- val
	}()
	return c
}
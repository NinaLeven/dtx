package dtx

import (
	"context"
	"errors"
)

var (
	ErrTxNotPresent     = errors.New("tx is not present in context")
	ErrTxAlreadyPresent = errors.New("tx already present in context")
)

type stringContextKeyType string

const (
	txIdType = stringContextKeyType("tx_id")
	parentTxIdType = stringContextKeyType("parent_tx_id")
)

type stringContextKey struct {
	Type     stringContextKeyType
	Incoming bool
}

func stringFromContext(ctx context.Context, key stringContextKey) (string, error) {
	val, ok := ctx.Value(key).(string)
	if !ok {
		return "", ErrTxNotPresent
	}
	return val, nil
}

func hasStringInContext(ctx context.Context, key stringContextKey) bool {
	_, ok := ctx.Value(key).(string)
	if !ok {
		return false
	}
	return true
}

func stringToContext(ctx context.Context, key stringContextKey, txId string) (context.Context, error) {
	has := hasStringInContext(ctx, key)
	if has {
		return nil, ErrTxAlreadyPresent
	}
	return context.WithValue(ctx, key, txId), nil
}

var contextValues = contextValuesEntrypoint{}

type contextValuesEntrypoint struct {}

func (c contextValuesEntrypoint) ParentTxId() contextValueType {
	return parentTxIdType
}

func (c contextValuesEntrypoint)  TxId() contextValueType {
	return txIdType
}

type contextValueType interface {
	Incoming() contextValueAction
	Outgoing() contextValueAction
}

func (c stringContextKeyType) Incoming() contextValueAction {
	return stringContextKey{
		Type:     c,
		Incoming: true,
	}
}

func (c stringContextKeyType) Outgoing() contextValueAction {
	return stringContextKey{
		Type:     c,
		Incoming: false,
	}
}

type contextValueAction interface {
	Exists(ctx context.Context) bool
	Get(ctx context.Context) (string, error)
	Set(ctx context.Context, val string) (context.Context, error)
}

func (c stringContextKey) Exists(ctx context.Context) bool {
	return hasStringInContext(ctx, c)
}

func (c stringContextKey) Get(ctx context.Context) (string, error) {
	return stringFromContext(ctx, c)
}

func (c stringContextKey) Set(ctx context.Context, val string) (context.Context, error) {
	return stringToContext(ctx, c, val)
}

package dtx

import (
	"context"
	"time"
)

type CreateRemoteTxOptions struct {
	TxId string
}

type UpdateRemoteTxOptions struct {
	TxId   string
	Status TxStatus
}

type GetRemoteTxOptions struct {
	TxId string
}

type RemoteTx struct {
	TxId      string
	Status    TxStatus
	UpdatedAt time.Time
}

type RemoteTxEvent struct {
	RemoteTx
}

type RemoteStorage interface {
	CreateTx(ctx context.Context, opts CreateRemoteTxOptions) error
	UpdateTx(ctx context.Context, opts UpdateRemoteTxOptions) error
	GetTx(ctx context.Context, opts GetRemoteTxOptions) (RemoteTx, error)
	Subscribe(txIds ...string) (chan RemoteTxEvent, error)
	Unsubscribe(txIds ...string) error
}
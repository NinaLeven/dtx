package dtx

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

type TxStatus int8

const (
	TxStatusCreated = TxStatus(iota + 1)
	TxStatusExecuting
	TxStatusReadyToCommit
	TxStatusRollback
	TxStatusCommit
)

var txStatusToString = map[TxStatus]string{
	TxStatusCreated:            "created",
	TxStatusExecuting:          "executing",
	TxStatusReadyToCommit:      "ready_to_commit",
	TxStatusRollback:           "rollback",
	TxStatusCommit:             "commit",
}

func (r TxStatus) String() string {
	return txStatusToString[r]
}

func TxStatusFromString(s string) TxStatus {
	for status, str := range txStatusToString {
		if str == s {
			return status
		}
	}
	return TxStatus(0)
}

type commitAction int8

const (
	commitActionRollback = commitAction(iota)
	commitActionCommit
)

type DistributedTxFactory interface {
	JoinDistributedTx(ctx context.Context) (DistributedTx, error)
	NewDistributedTx(ctx context.Context) (context.Context, DistributedTx, error)
}

type distributedTxFactory struct {
	storage RemoteStorage
}

func (d *distributedTxFactory) JoinDistributedTx(ctx context.Context) (DistributedTx, error) {
	txId, err := contextValues.TxId().Incoming().Get(ctx)
	if err != nil {
		return nil, ErrTxNotPresent
	}

}

func (d *distributedTxFactory) NewDistributedTx(ctx context.Context) (context.Context, DistributedTx, error) {
	ok := contextValues.TxId().Incoming().Exists(ctx)
	if ok {
		return nil, nil, ErrTxAlreadyPresent
	}

	txId := uuid.New().String()
	err := d.storage.CreateTx(ctx, CreateRemoteTxOptions{
		TxId: txId,
	})
	if err != nil {
		return nil, nil, err
	}


}

func (d *distributedTxFactory) constructDistributedTx(ctx context.Context) (DistributedTx, error) {

}

type DistributedTx interface {
	CallWithDistributedTx(ctx context.Context) (context.Context, error)
	Commit(commit, rollback func() error) FutureErrorGetter
	Rollback(rollback func() error) error
}

type distributedTx struct {
	storage RemoteStorage

	parentTx, rootTx *RemoteTx
	children         map[string]RemoteTx
	sub              chan RemoteTxEvent

	currentTx RemoteTx

	commitTimeout time.Duration

	commitAction futureCommitAction
	commitOnce   sync.Once
	commitErr    FutureError

	m sync.Mutex

	doneCh chan struct{}
}

func (tx *distributedTx) done() {
	_, ok := <-tx.doneCh
	if ok {
		close(tx.doneCh)
	}
}

func (tx *distributedTx) watchParent() error {

}

func (tx *distributedTx) watchChild(childTxId string) error {

}

func (tx *distributedTx) CallWithDistributedTx(ctx context.Context) (context.Context, error) {

}

func (tx *distributedTx) applyCommitAction(commit, rollback func() error, action commitAction) {
	tx.commitOnce.Do(func() {
		var err error
		switch action {
		case commitActionCommit:
			err = commit()
		case commitActionRollback:
			err = rollback()
		default:
			err = rollback()
			err = fmt.Errorf("invalid commit action; rollback err: %w", err)
		}
		_ = tx.commitErr.SetError(err)
	})
}

func (tx *distributedTx) waitForCommit(commit, rollback func() error) {
	defer tx.done()
	select {
	case action, ok := <-tx.commitAction.Chan():
		if !ok {
			return
		}
		tx.applyCommitAction(commit, rollback, action)
	case <-time.After(tx.commitTimeout):
		tx.applyCommitAction(commit, rollback, commitActionRollback)
	}
}

var ErrIncorrectStatus = errors.New("incorrect status for status change")

func (tx *distributedTx) Commit(commit, rollback func() error) FutureErrorGetter {
	defer tx.m.Unlock()
	tx.m.Lock()

	if tx.currentTx.Status != TxStatusExecuting {
		return tx.commitErr
	}

	err := tx.storage.UpdateTx(context.Background(), UpdateRemoteTxOptions{
		TxId:   tx.currentTx.TxId,
		Status: TxStatusReadyToCommit,
	})
	if err != nil {
		_ = tx.commitErr.SetError(err)
		return tx.commitErr
	}
	tx.currentTx = rtx

	go tx.waitForCommit(commit, rollback)

	return tx.commitErr
}

func (tx *distributedTx) Rollback(rollback func() error) error {
	panic("implement me")
}

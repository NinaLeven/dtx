package dtx

import (
	"context"
	"database/sql"
	"dtx/future"
	"errors"
)

type SqlDBWrapper struct {
	*sql.DB
}

func (d *SqlDBWrapper) BeginTxDistributed(ctx context.Context, opts *sql.TxOptions) (context.Context, *TxWrapper, error) {
	//TODO: check ctx, retrieve current global transaction id; check remote server; check remote server
	//TODO:
}

func (d *SqlDBWrapper) NewTxDistributed(ctx context.Context, opts *sql.TxOptions) (context.Context, *TxWrapper, error) {

}

var ErrMustNotUseCommit = errors.New("usage of local Commit in a distributed transaction; use CommitDistributed")
var ErrMustNotUseRollback = errors.New("usage of local Rollback in a distributed transaction; use RollbackDistributed")

type TxWrapper struct {
	*sql.Tx
	txErr future.Future
	DistributedTx
}

func (tx *TxWrapper) Commit() error {
	return ErrMustNotUseCommit
}

func (tx *TxWrapper) Rollback() error {
	return ErrMustNotUseRollback
}

func (tx *TxWrapper) CommitDistributed(ctx context.Context) future.Future {
	tx.
}

func (tx *TxWrapper) RollbackDistributed(ctx context.Context) future.Future {

}




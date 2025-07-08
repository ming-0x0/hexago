package transaction

import (
	"context"

	"gorm.io/gorm"
)

type TxKey string

const (
	Tx TxKey = "tx"
)

//go:generate go tool mockgen -destination mock/transaction.go -package mock github.com/ming-0x0/hexago/shared/transaction TransactionInterface
type TransactionInterface interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

type Transaction struct {
	db *gorm.DB
}

func NewTransaction(db *gorm.DB) *Transaction {
	return &Transaction{
		db: db,
	}
}

func (t *Transaction) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	return t.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, Tx, tx)
		return fn(ctx)
	})
}

func TransactionFromContext(ctx context.Context) (*gorm.DB, bool) {
	v := ctx.Value(Tx)
	if v != nil {
		if tx, ok := v.(*gorm.DB); ok {
			return tx, true
		}
	}
	return nil, false
}

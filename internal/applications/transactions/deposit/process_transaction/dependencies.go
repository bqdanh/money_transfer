package process_transaction

import (
	"context"
	"time"

	"github.com/bqdanh/money_transfer/internal/entities/transaction"
)

//go:generate mockgen -destination=mocks.go -package=process_transaction -source=dependencies.go .

type distributeLock interface {
	AcquireLockForProcessDepositTransaction(ctx context.Context, transactionID int64, lockDuration time.Duration) (releaseLock func(), err error)
}

type transactionRepository interface {
	GetTransactionByID(ctx context.Context, transID int64) (transaction.Transaction, error)
	UpdateTransaction(ctx context.Context, t transaction.Transaction) error
}

type sofProvider interface {
	MakeDepositTransaction(ctx context.Context, trans transaction.Transaction) (transaction.MakeDepositResult, error)
}

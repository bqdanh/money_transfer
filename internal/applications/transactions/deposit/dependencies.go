package deposit

import (
	"context"
	"time"

	"github.com/bqdanh/money_transfer/internal/entities/account"
	"github.com/bqdanh/money_transfer/internal/entities/transaction"
)

//go:generate mockgen -destination=mocks.go -package=deposit -source=dependencies.go .

type distributeLock interface {
	AcquireLockForCreateDepositTransaction(ctx context.Context, requestID string, lockDuration time.Duration) (releaseLock func(), err error)
}

type accountRepository interface {
	GetAccountsByID(ctx context.Context, accountID int64) (account.Account, error)
}

type transactionRepository interface {
	//CreateTransaction create transaction and generate transaction id
	CreateTransaction(ctx context.Context, t transaction.Transaction) (transaction.Transaction, error)
	//GetTransactionByRequestID get transaction in 7day by request id,
	// if notfound return exceptions.PreconditionError SubjectTransaction PreconditionReasonTransactionNotFound
	GetTransactionByRequestID(ctx context.Context, requestID string) (transaction.Transaction, error)
}

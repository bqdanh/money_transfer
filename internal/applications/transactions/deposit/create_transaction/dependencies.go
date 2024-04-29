package create_transaction

import (
	"context"
	"time"

	"github.com/bqdanh/money_transfer/internal/entities/account"
	"github.com/bqdanh/money_transfer/internal/entities/exceptions"
	"github.com/bqdanh/money_transfer/internal/entities/transaction"
)

//go:generate mockgen -destination=mocks.go -package=create_transaction -source=dependencies.go .

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
	GetTransactionByRequestID(ctx context.Context, account account.Account, requestID string) (transaction.Transaction, error)
}

var ErrNotFoundTransaction = exceptions.NewPreconditionError(exceptions.PreconditionReasonTransactionNotFound, exceptions.SubjectTransaction, "transaction not found", nil)

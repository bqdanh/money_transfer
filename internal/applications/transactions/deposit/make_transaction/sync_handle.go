package make_transaction

import (
	"context"
	"fmt"

	"github.com/bqdanh/money_transfer/internal/applications/transactions/deposit/create_transaction"
	"github.com/bqdanh/money_transfer/internal/applications/transactions/deposit/process_transaction"
	"github.com/bqdanh/money_transfer/internal/entities/transaction"
)

type MakeTransactionSync struct {
	create_transaction.CreateDepositTransaction
	process_transaction.ProcessDepositTransaction
}

var (
	ErrFailedToCreateDepositTransaction  = fmt.Errorf("failed to create deposit transaction")
	ErrFailedToProcessDepositTransaction = fmt.Errorf("failed to process deposit transaction")
)

func NewMakeTransactionSync(
	create create_transaction.CreateDepositTransaction,
	process process_transaction.ProcessDepositTransaction,
) MakeTransactionSync {
	return MakeTransactionSync{
		CreateDepositTransaction:  create,
		ProcessDepositTransaction: process,
	}
}

// Handle is a synchronous function to handle the deposit transaction
// It will create a deposit transaction and then process it
// If any error occurs, it will return the error for each step:
// - ErrFailedToCreateDepositTransaction: failed to create deposit transaction
// - ErrFailedToProcessDepositTransaction: failed to process deposit transaction
func (m MakeTransactionSync) Handle(ctx context.Context, p MakeDepositTransactionParams) (transaction.Transaction, error) {
	trans, err := m.CreateDepositTransaction.Handle(ctx, create_transaction.CreateDepositTransactionParams{
		RequestID:    p.RequestID,
		UserID:       p.UserID,
		AccountID:    p.AccountID,
		Amount:       p.Amount,
		Descriptions: p.Descriptions,
		Source:       p.Source,
	})
	if err != nil {
		return transaction.Transaction{}, fmt.Errorf("%w: %w", ErrFailedToCreateDepositTransaction, err)
	}

	trans2, err2 := m.ProcessDepositTransaction.Handle(ctx, trans.ID)
	if err2 != nil {
		return trans, fmt.Errorf("%w: %w", ErrFailedToProcessDepositTransaction, err2)
	}
	return trans2, nil
}

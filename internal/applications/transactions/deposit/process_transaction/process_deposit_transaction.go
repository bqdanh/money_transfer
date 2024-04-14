package process_transaction

import (
	"context"
	"fmt"
	"time"

	"github.com/bqdanh/money_transfer/internal/entities/exceptions"
	"github.com/bqdanh/money_transfer/internal/entities/transaction"
)

type ProcessDepositTransaction struct {
	cfg            Config
	trepo          transactionRepository
	distributeLock distributeLock
	sofProvider    sofProvider
}

type Config struct {
	LockDuration time.Duration `json:"lock_duration" mapstructure:"lock_duration"`
}

type sofProvider interface {
	MakeDepositTransaction(ctx context.Context, trans transaction.Transaction) (transaction.Transaction, error)
}

func (p ProcessDepositTransaction) ProcessDepositTransaction(ctx context.Context, transactionID int64) (transaction.Transaction, error) {
	if transactionID <= 0 {
		return transaction.Transaction{}, exceptions.NewInvalidArgumentError(
			"TransactionID",
			"transaction must greater than 0",
			map[string]interface{}{
				"transaction_id": transactionID,
			},
		)
	}

	releaseLock, err := p.distributeLock.AcquireLockForProcessDepositTransaction(ctx, transactionID, p.cfg.LockDuration)
	if err != nil {
		return transaction.Transaction{}, fmt.Errorf("acquire lock for process_transaction deposit transaction: %w", err)
	}
	defer releaseLock()
	trans, err := p.trepo.GetTransactionByID(ctx, transactionID)
	if err != nil {
		return transaction.Transaction{}, fmt.Errorf("get transaction by id: %w", err)
	}
	if err = trans.ReadyForProcessDeposit(); err != nil {
		return transaction.Transaction{}, fmt.Errorf("transaction is not ready for process_transaction deposit: %w", err)
	}
	trans, err = trans.MakeTransactionDepositProcessing()
	if err != nil {
		return transaction.Transaction{}, fmt.Errorf("make transaction deposit processing: %w", err)
	}
	if err = p.trepo.UpdateTransaction(ctx, trans); err != nil {
		return transaction.Transaction{}, fmt.Errorf("update transaction: %w", err)
	}
	trans, err = p.sofProvider.MakeDepositTransaction(ctx, trans)
	if err != nil {
		return transaction.Transaction{}, fmt.Errorf("make deposit transaction: %w", err)
	}
	if err = p.trepo.UpdateTransaction(ctx, trans); err != nil {
		return transaction.Transaction{}, fmt.Errorf("update transaction: %w", err)
	}
	return trans, nil
}

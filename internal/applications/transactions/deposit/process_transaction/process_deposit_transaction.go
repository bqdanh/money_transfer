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

func NewProcessDepositTransaction(cfg Config, tr transactionRepository, dl distributeLock, sof sofProvider) (ProcessDepositTransaction, error) {
	if tr == nil {
		return ProcessDepositTransaction{}, fmt.Errorf("transaction repository is nil")
	}
	if dl == nil {
		return ProcessDepositTransaction{}, fmt.Errorf("distribute lock is nil")
	}
	if sof == nil {
		return ProcessDepositTransaction{}, fmt.Errorf("sof provider is nil")
	}
	return ProcessDepositTransaction{
		cfg:            cfg,
		trepo:          tr,
		distributeLock: dl,
		sofProvider:    sof,
	}, nil
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
	var evt transaction.Event
	trans, evt, err = trans.MakeTransactionDepositProcessing()
	if err != nil {
		return transaction.Transaction{}, fmt.Errorf("make transaction deposit processing: %w", err)
	}
	if err = p.trepo.UpdateTransaction(ctx, trans, evt); err != nil {
		return transaction.Transaction{}, fmt.Errorf("update transaction: %w", err)
	}
	transDataResult, err := p.sofProvider.MakeDepositTransaction(ctx, trans)
	if err != nil {
		return transaction.Transaction{}, fmt.Errorf("make deposit transaction: %w", err)
	}
	trans, evt, err = trans.WithTransactionResult(transDataResult)
	if err != nil {
		return transaction.Transaction{}, fmt.Errorf("make transaction deposit: %w", err)
	}
	if err = p.trepo.UpdateTransaction(ctx, trans, evt); err != nil {
		return transaction.Transaction{}, fmt.Errorf("update transaction: %w", err)
	}
	return trans, nil
}

package create_transaction

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/bqdanh/money_transfer/internal/entities/currency"
	"github.com/bqdanh/money_transfer/internal/entities/exceptions"
	"github.com/bqdanh/money_transfer/internal/entities/transaction"
	"github.com/bqdanh/money_transfer/internal/entities/transaction/deposit"
)

type CreateDepositTransaction struct {
	cfg                   Config
	distributeLock        distributeLock
	accountRepository     accountRepository
	transactionRepository transactionRepository
}

type Config struct {
	LockDuration       time.Duration `json:"lock_duration" mapstructure:"lock_duration"`
	MinimumAmount      float64       `json:"minimum_amount" mapstructure:"minimum_amount"`
	MaximumAmount      float64       `json:"maximum_amount" mapstructure:"maximum_amount"`
	CurrenciesAccepted []string      `json:"currencies_accepted" mapstructure:"currencies_accepted"`
}

var DefaultConfig = Config{
	LockDuration:       30 * time.Second,
	MinimumAmount:      1,
	MaximumAmount:      1_000_000_000,
	CurrenciesAccepted: []string{string(currency.VND)},
}

type CreateDepositTransactionParams struct {
	//RequestID for detect duplicate request, so request id just need check unique in 7 days, that is trace off for performance
	RequestID    string
	UserID       int64
	AccountID    int64
	Amount       currency.Amount
	Descriptions string
	Source       string //optional
}

func NewCreateDepositTransaction(cfg Config, dl distributeLock, ar accountRepository, tr transactionRepository) (CreateDepositTransaction, error) {
	if dl == nil {
		return CreateDepositTransaction{}, fmt.Errorf("distribute lock is nil")
	}
	if ar == nil {
		return CreateDepositTransaction{}, fmt.Errorf("account repository is nil")
	}
	if tr == nil {
		return CreateDepositTransaction{}, fmt.Errorf("transaction repository is nil")
	}

	return CreateDepositTransaction{
		cfg:                   cfg,
		distributeLock:        dl,
		accountRepository:     ar,
		transactionRepository: tr,
	}, nil
}

func (c CreateDepositTransaction) validateCreateDepositTransactionParams(p CreateDepositTransactionParams) error {
	if p.UserID <= 0 {
		return exceptions.NewInvalidArgumentError(
			"UserID",
			"user must greater than 0",
			map[string]interface{}{
				"user_id": p.UserID,
			},
		)
	}
	if p.AccountID <= 0 {
		return exceptions.NewInvalidArgumentError(
			"AccountID",
			"account must greater than 0",
			map[string]interface{}{
				"account_id": p.AccountID,
			},
		)
	}
	if ok, err := p.Amount.IsLt(c.cfg.MinimumAmount); err != nil || ok {
		return exceptions.NewInvalidArgumentError(
			"Amount",
			fmt.Sprintf("amount must greater than %f", c.cfg.MinimumAmount),
			map[string]interface{}{
				"error":          err,
				"amount":         p.Amount,
				"minimum_amount": c.cfg.MinimumAmount,
			},
		)
	}
	if ok, err := p.Amount.IsGt(c.cfg.MaximumAmount); err != nil || ok {
		return exceptions.NewInvalidArgumentError(
			"Amount",
			fmt.Sprintf("amount must less than %f", c.cfg.MaximumAmount),
			map[string]interface{}{
				"error":  err,
				"amount": p.Amount,
			},
		)
	}
	if !slices.Contains(c.cfg.CurrenciesAccepted, string(p.Amount.Currency)) {
		return exceptions.NewInvalidArgumentError(
			"Amount",
			fmt.Sprintf("currency %s is not accepted", p.Amount.Currency),
			map[string]interface{}{
				"currency":            p.Amount.Currency,
				"currencies_accepted": c.cfg.CurrenciesAccepted,
			},
		)
	}
	return nil
}

func (c CreateDepositTransaction) Handle(ctx context.Context, p CreateDepositTransactionParams) (transaction.Transaction, error) {
	if err := c.validateCreateDepositTransactionParams(p); err != nil {
		return transaction.Transaction{}, fmt.Errorf("validate create deposit transaction params: %w", err)
	}

	ac, err := c.accountRepository.GetAccountsByID(ctx, p.AccountID)
	if err != nil {
		return transaction.Transaction{}, fmt.Errorf("get account by id: %w", err)
	}
	if ac.UserID != p.UserID {
		return transaction.Transaction{}, exceptions.NewPreconditionError(exceptions.PreconditionReasonPermissionDenied, exceptions.SubjectAccount, "use account of another user", map[string]interface{}{
			"account_id": p.AccountID,
			"user_id":    p.UserID,
		})
	}
	if err = ac.IsAvailableForDeposit(); err != nil {
		return transaction.Transaction{}, fmt.Errorf("account is not available for deposit: %w", err)
	}

	releaseLock, err := c.distributeLock.AcquireLockForCreateDepositTransaction(ctx, p.RequestID, c.cfg.LockDuration)
	if err != nil {
		return transaction.Transaction{}, fmt.Errorf("acquire lock for create deposit transaction: %w", err)
	}
	defer releaseLock()

	var trn transaction.Transaction
	trn, err = c.transactionRepository.GetTransactionByRequestID(ctx, p.RequestID)
	if err == nil {
		//TODO: we can return the transaction to the client, but we need to check the status of the transaction and created time to make sure this transaction is valid
		return transaction.Transaction{}, exceptions.NewPreconditionError(exceptions.PreconditionReasonTransactionIsAvailable, exceptions.SubjectTransaction, "request id is duplicated", map[string]interface{}{
			"request_id":     p.RequestID,
			"transaction_id": trn.ID,
		})
	}
	if !errors.Is(err, ErrNotFoundTransaction) {
		return transaction.Transaction{}, fmt.Errorf("get transaction by request id: %w", err)
	}
	depositData, err := deposit.CreateDeposit(ac, p.Source)
	if err != nil {
		return transaction.Transaction{}, fmt.Errorf("create deposit data: %w", err)
	}

	depositTransaction := transaction.CreateTransaction(ac, p.Amount, p.Descriptions, transaction.TypeDeposit, transaction.Data{
		IsTransactionDataItr: depositData,
	})

	depositTransaction, err = c.transactionRepository.CreateTransaction(ctx, depositTransaction)
	if err != nil {
		return transaction.Transaction{}, fmt.Errorf("create deposit transaction: %w", err)
	}
	return depositTransaction, nil
}

package transaction

import (
	"slices"

	"github.com/bqdanh/money_transfer/internal/entities/account"
	"github.com/bqdanh/money_transfer/internal/entities/currency"
	"github.com/bqdanh/money_transfer/internal/entities/exceptions"
)

type Transaction struct {
	ID          int64           `json:"id"`
	Account     account.Account `json:"account"`
	Amount      currency.Amount `json:"amount"`
	Description string          `json:"description"`
	Status      Status          `json:"status"`
	Type        Type            `json:"type"`
	Data        Data            `json:"data"`
}

func CreateTransaction(account account.Account, amount currency.Amount, description string, t Type, data Data) Transaction {
	return Transaction{
		ID:          0,
		Account:     account,
		Amount:      amount,
		Description: description,
		Type:        t,
		Data:        data,
		Status:      StatusInit,
	}
}

type Type string

const (
	TypeWithdraw = Type("withdraw")
	TypeDeposit  = Type("deposit")
)

type Status string

const (
	StatusInit       = Status("init")
	StatusProcessing = Status("processing")
	StatusSuccess    = Status("success")
	StatusFailed     = Status("failed")
)

type Data struct {
	IsTransactionDataItr
}

type IsTransactionDataItr interface {
	isTransactionData()
	GetType() Type
}

type IsTransactionDataImplementMustImport struct {
}

func (b IsTransactionDataImplementMustImport) isTransactionData() {}

// statusReadyForProcessDeposit is a list of status that transaction can be processed for deposit
// NOTE: for at most one execution of process_transaction deposit transaction, plz remove StatusProcessing
var statusReadyForProcessDeposit = []Status{StatusInit, StatusProcessing}

func (t Transaction) ReadyForProcessDeposit() error {
	if t.Type != TypeDeposit {
		return exceptions.NewPreconditionError(
			exceptions.PreconditionReasonTransactionTypeNotMatch,
			exceptions.SubjectTransaction,
			"transaction type not match",
			map[string]interface{}{
				"transaction_type": t.Type,
				"expected_type":    TypeDeposit,
			},
		)
	}

	if slices.Contains(statusReadyForProcessDeposit, t.Status) {
		return exceptions.NewPreconditionError(
			exceptions.PreconditionReasonTransactionStatusNotMatch,
			exceptions.SubjectTransaction,
			"transaction status not match",
			map[string]interface{}{
				"transaction_status": t.Status,
				"expected_status":    statusReadyForProcessDeposit,
			},
		)
	}

	return nil
}

func (t Transaction) MakeTransactionDepositProcessing() (Transaction, error) {
	if slices.Contains(statusReadyForProcessDeposit, t.Status) {
		return t, exceptions.NewPreconditionError(
			exceptions.PreconditionReasonTransactionStatusNotMatch,
			exceptions.SubjectTransaction,
			"transaction status not match",
			map[string]interface{}{
				"transaction_status": t.Status,
				"expected_status":    statusReadyForProcessDeposit,
			},
		)
	}

	t.Status = StatusProcessing
	return t, nil
}

type MakeDepositResult struct {
	Status Status
	Data   Data
}

func (t Transaction) MakeTransactionDeposit(result MakeDepositResult) (Transaction, error) {
	if t.Status != StatusProcessing {
		return t, exceptions.NewPreconditionError(
			exceptions.PreconditionReasonTransactionStatusNotMatch,
			exceptions.SubjectTransaction,
			"transaction status not match",
			map[string]interface{}{
				"transaction_status": t.Status,
				"expected_status":    StatusProcessing,
			},
		)
	}

	t.Status = result.Status
	t.Data = result.Data
	return t, nil
}

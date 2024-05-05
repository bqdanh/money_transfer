package transaction

import (
	"fmt"
	"slices"

	"github.com/bqdanh/money_transfer/internal/entities/account"
	"github.com/bqdanh/money_transfer/internal/entities/currency"
	"github.com/bqdanh/money_transfer/internal/entities/exceptions"
)

type Transaction struct {
	ID          int64           `json:"id"`
	Account     account.Account `json:"account"`
	Amount      currency.Amount `json:"amount"`
	Version     int             `json:"version"`
	RequestID   string          `json:"request_id"`
	Description string          `json:"description"`
	Status      Status          `json:"status"`
	Type        Type            `json:"type"`
	Data        Data            `json:"data"`
}

func CreateTransaction(requestID string, account account.Account, amount currency.Amount, description string, t Type, data Data) Transaction {
	return Transaction{
		ID:          0,
		Account:     account,
		Amount:      amount,
		Version:     0,
		RequestID:   requestID,
		Description: description,
		Status:      StatusInit,
		Type:        t,
		Data:        data,
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
	GetTransactionStatus() Status
	GetPartnerRefTransactionID() string
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

	if !slices.Contains(statusReadyForProcessDeposit, t.Status) {
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

func (t Transaction) IsInitTransaction() bool {
	return t.Status == StatusInit
}

func (t Transaction) IsProcessing() bool {
	return t.Status == StatusInit || t.Status == StatusProcessing
}

func (t Transaction) IsSuccess() bool {
	return t.Status == StatusSuccess
}

func (t Transaction) IsFailed() bool {
	return t.Status == StatusFailed
}

func (t Transaction) MakeTransactionDepositProcessing() (Transaction, Event, error) {
	if err := t.ReadyForProcessDeposit(); err != nil {
		return t, Event{}, fmt.Errorf("ready for process deposit failed: %w", err)
	}

	t.Version += 1

	oldStatus := t.Status
	t.Status = StatusProcessing

	evtData := MakeTransactionProcessing{
		FromStatus: oldStatus,
		ToStatus:   t.Status,
	}
	evt := NewEvent(t.ID, t.Type, t.Version, evtData)

	return t, evt, nil
}

func (t Transaction) GetPartnerRefTransactionID() string {
	return t.Data.GetPartnerRefTransactionID()
}

func (t Transaction) WithTransactionResult(tranData Data) (Transaction, Event, error) {
	if t.Status != StatusProcessing {
		return t, Event{}, exceptions.NewPreconditionError(
			exceptions.PreconditionReasonTransactionStatusNotMatch,
			exceptions.SubjectTransaction,
			"transaction status not match",
			map[string]interface{}{
				"transaction_status": t.Status,
				"expected_status":    StatusProcessing,
			},
		)
	}
	oldStatus := t.Status
	t.Data = tranData
	t.Status = tranData.GetTransactionStatus()
	t.Version += 1

	var evtData EventData
	switch tranData.GetTransactionStatus() {
	case StatusSuccess:
		evtData = EventTransactionSucceed{
			Transaction: t,
			FromStatus:  oldStatus,
			ToStatus:    t.Status,
		}
	case StatusFailed:
		evtData = EventTransactionFailed{
			Transaction: t,
			FromStatus:  oldStatus,
			ToStatus:    t.Status,
		}
	case StatusProcessing:
		evtData = EventTransactionKeepProcessing{
			TransData: tranData,
		}
	default:
		return t, Event{}, exceptions.NewPreconditionError(
			exceptions.PreconditionReasonTransactionInvalidStatus,
			exceptions.SubjectTransaction,
			"invalid transaction status",
			map[string]interface{}{
				"transaction_status": tranData.GetTransactionStatus(),
				"trans_data":         tranData,
			},
		)
	}

	evt := NewEvent(t.ID, t.Type, t.Version, evtData)

	return t, evt, nil
}

func (t Transaction) CreateTransaction(transID int64) (Transaction, Event, error) {
	if t.Status != StatusInit {
		return t, Event{}, exceptions.NewPreconditionError(
			exceptions.PreconditionReasonTransactionStatusNotMatch,
			exceptions.SubjectTransaction,
			"transaction status not match",
			map[string]interface{}{
				"transaction_status": t.Status,
				"expected_status":    StatusInit,
			},
		)
	}
	t.ID = transID
	t.Version = t.Version + 1
	evtData := EventTransactionCreated{
		Transaction: t,
	}
	evt := NewEvent(t.ID, t.Type, t.Version, evtData)
	return t, evt, nil
}

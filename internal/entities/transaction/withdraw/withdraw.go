package withdraw

import "github.com/bqdanh/money_transfer/internal/entities/transaction"

type Withdraw struct {
	transaction.IsTransactionDataImplementMustImport
	//TODO: define Destination entity
	Destination       string `json:"destination"`
	BankTransactionID string `json:"bank_transaction_id"` //id at bank for identify transaction, used for audit and reconciliation
}

func (d Withdraw) GetType() transaction.Type {
	return transaction.TypeWithdraw
}

func (d Withdraw) GetTransactionStatus() transaction.Status {
	//TODO: implement GetTransactionStatus
	return transaction.StatusFailed
}

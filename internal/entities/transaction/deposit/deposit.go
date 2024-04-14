package deposit

import "github.com/bqdanh/money_transfer/internal/entities/transaction"

type Deposit struct {
	transaction.IsTransactionDataImplementMustImport
	//TODO: define Deposit Source entity
	Source            string `json:"source"`
	BankTransactionID string `json:"bank_transaction_id"` //id at bank for identify transaction, used for audit and reconciliation
}

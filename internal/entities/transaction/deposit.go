package transaction

type Deposit struct {
	IsTransactionDataImplementMustImport
	//TODO: define Deposit Source entity
	Source            string
	BankTransactionID string //id at bank for identify transaction, used for audit and reconciliation
}

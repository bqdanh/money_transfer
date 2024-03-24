package transaction

type Withdraw struct {
	IsTransactionDataImplementMustImport
	//TODO: define Destination entity
	Destination       string
	BankTransactionID string //id at bank for identify transaction, used for audit and reconciliation
}

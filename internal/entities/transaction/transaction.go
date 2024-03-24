package transaction

import "github.com/bqdanh/money_transfer/internal/entities/account"

type Transaction struct {
	ID      int64
	Account account.Account

	Amount      int64
	Description string
	Type        Type
	Data        Data
}

type Type string

const (
	TypeWithdraw = Type("withdraw")
	TypeDeposit  = Type("deposit")
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

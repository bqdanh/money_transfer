package transaction

import (
	"github.com/bqdanh/money_transfer/internal/entities/account"
	"github.com/bqdanh/money_transfer/internal/entities/currency"
)

type Transaction struct {
	ID          int64           `json:"id"`
	Account     account.Account `json:"account"`
	Amount      currency.Amount `json:"amount"`
	Description string          `json:"description"`
	Type        Type            `json:"type"`
	Data        Data            `json:"data"`
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

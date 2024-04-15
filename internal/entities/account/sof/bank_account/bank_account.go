package bank_account

import (
	"github.com/bqdanh/money_transfer/internal/entities/account"
)

// BankAccount
type BankAccount struct {
	AccountNumber string
	AccountName   string
}

func (BankAccount) GetSourceOfFundType() account.SourceOfFundType {
	return account.SofTypeBankAccount
}

package sof

import "github.com/bqdanh/money_transfer/internal/entities/account"

const (
	SourceOfFundCodeACB = account.SourceOfFundCode("ACB")
)

// ACBAccountStatus the status of account at ACB bank
type ACBAccountStatus string

const (
	ACBAccountStatusActive   = ACBAccountStatus("active")
	ACBAccountStatusInactive = ACBAccountStatus("inactive")
	ACBAccountStatusDormant  = ACBAccountStatus("dormant")
)

type ACBAccount struct {
	account.IsSourceOfFundImplementMustImport
	BankAccount
	//define ACB specific fields here: example status
	Status ACBAccountStatus
}

func (ACBAccount) GetSourceOfFundCode() account.SourceOfFundCode {
	return SourceOfFundCodeACB
}

package sof

import "github.com/bqdanh/money_transfer/internal/entities/account"

const (
	SourceOfFundCodeVIB = account.SourceOfFundCode("VIB")
)

// VIBAccountStatus the status of account at VIB bank
type VIBAccountStatus string

const (
	VIBAccountStatusActive   = VIBAccountStatus("active")
	VIBAccountStatusInactive = VIBAccountStatus("inactive")
)

type VibAccount struct {
	account.IsSourceOfFundImplementMustImport
	BankAccount
	//define VIB specific fields here: example status
	Status VIBAccountStatus
}

func (VibAccount) GetSourceOfFundCode() account.SourceOfFundCode {
	return SourceOfFundCodeVIB
}

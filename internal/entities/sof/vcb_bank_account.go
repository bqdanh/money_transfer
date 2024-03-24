package sof

import "github.com/bqdanh/money_transfer/internal/entities/account"

const (
	SourceOfFundCodeVCB = account.SourceOfFundCode("VCB")
)

// VCBAccountStatus the status of account at VCB bank
type VCBAccountStatus string

const (
	VCBAccountStatusActive   = VCBAccountStatus("active")
	VCBAccountStatusInactive = VCBAccountStatus("inactive")
	VCBAccountStatusLocked   = VCBAccountStatus("locked")
)

type VcbAccount struct {
	account.IsSourceOfFundImplementMustImport
	BankAccount
	//define VCB specific fields here: example status
	Status VCBAccountStatus
}

func (VcbAccount) GetSourceOfFundCode() account.SourceOfFundCode {
	return SourceOfFundCodeVCB
}

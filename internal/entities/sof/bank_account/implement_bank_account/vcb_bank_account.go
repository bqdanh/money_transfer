package implement_bank_account

import (
	"github.com/bqdanh/money_transfer/internal/entities/account"
	"github.com/bqdanh/money_transfer/internal/entities/sof/bank_account"
)

const (
	SourceOfFundCodeVCB = account.SourceOfFundCode("VCB")
)

func init() {
	bank_account.RegisterSourceOfFundBankAccount(SourceOfFundCodeVCB)
	bank_account.RegisterSourceOfFundBankAccountConstructor(SourceOfFundCodeVCB, func(a bank_account.BankAccount) (account.SourceOfFundData, error) {
		acbAc := VcbAccount{
			IsSourceOfFundImplementMustImport: account.IsSourceOfFundImplementMustImport{},
			BankAccount:                       a,
			Status:                            "",
		}

		return account.SourceOfFundData{
			IsSourceOfFundItr: acbAc,
		}, nil
	})
}

// VCBAccountStatus the status of account at VCB bank
type VCBAccountStatus string

const (
	VCBAccountStatusActive   = VCBAccountStatus("active")
	VCBAccountStatusInactive = VCBAccountStatus("inactive")
	VCBAccountStatusLocked   = VCBAccountStatus("locked")
)

type VcbAccount struct {
	account.IsSourceOfFundImplementMustImport
	bank_account.BankAccount
	//define VCB specific fields here: example status
	Status VCBAccountStatus
}

func (VcbAccount) GetSourceOfFundCode() account.SourceOfFundCode {
	return SourceOfFundCodeVCB
}

func (a VcbAccount) IsTheSameSof(other account.IsSourceOfFundItr) bool {
	if v, ok := other.(account.SourceOfFundData); ok {
		return a.IsTheSameSof(v.IsSourceOfFundItr)
	}

	v, ok := other.(VcbAccount)
	if !ok {
		return false
	}
	return v.BankAccount == a.BankAccount && v.AccountName == a.AccountName
}

package implement_bank_account

import (
	"github.com/bqdanh/money_transfer/internal/entities/account"
	"github.com/bqdanh/money_transfer/internal/entities/sof/bank_account"
)

const (
	SourceOfFundCodeVIB = account.SourceOfFundCode("VIB")
)

func init() {
	bank_account.RegisterSourceOfFundBankAccount(SourceOfFundCodeVIB)
	bank_account.RegisterSourceOfFundBankAccountConstructor(SourceOfFundCodeVIB, func(a bank_account.BankAccount) (account.SourceOfFundData, error) {
		acbAc := VibAccount{
			IsSourceOfFundImplementMustImport: account.IsSourceOfFundImplementMustImport{},
			BankAccount:                       a,
			Status:                            VIBAccountStatusActive,
		}

		return account.SourceOfFundData{
			IsSourceOfFundItr: acbAc,
		}, nil
	})
}

// VIBAccountStatus the status of account at VIB bank
type VIBAccountStatus string

const (
	VIBAccountStatusActive   = VIBAccountStatus("active")
	VIBAccountStatusInactive = VIBAccountStatus("inactive")
)

type VibAccount struct {
	account.IsSourceOfFundImplementMustImport
	bank_account.BankAccount
	//define VIB specific fields here: example status
	Status VIBAccountStatus
}

func (VibAccount) GetSourceOfFundCode() account.SourceOfFundCode {
	return SourceOfFundCodeVIB
}

func (a VibAccount) IsTheSameSof(other account.IsSourceOfFundItr) bool {
	if v, ok := other.(account.SourceOfFundData); ok {
		return a.IsTheSameSof(v.IsSourceOfFundItr)
	}

	v, ok := other.(VibAccount)
	if !ok {
		return false
	}
	return v.BankAccount == a.BankAccount && v.AccountName == a.AccountName
}

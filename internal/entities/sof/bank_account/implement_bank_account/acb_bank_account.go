package implement_bank_account

import (
	"github.com/bqdanh/money_transfer/internal/entities/account"
	"github.com/bqdanh/money_transfer/internal/entities/sof/bank_account"
)

const (
	SourceOfFundCodeACB = account.SourceOfFundCode("ACB")
)

func init() {
	bank_account.RegisterSourceOfFundBankAccount(SourceOfFundCodeACB)
	bank_account.RegisterSourceOfFundBankAccountConstructor(SourceOfFundCodeACB, func(a bank_account.BankAccount) (account.SourceOfFundData, error) {
		acbAc := ACBAccount{
			IsSourceOfFundImplementMustImport: account.IsSourceOfFundImplementMustImport{},
			BankAccount:                       a,
			Status:                            "",
		}

		return account.SourceOfFundData{
			IsSourceOfFundItr: acbAc,
		}, nil
	})
}

// ACBAccountStatus the status of account at ACB bank
type ACBAccountStatus string

const (
	ACBAccountStatusActive   = ACBAccountStatus("active")
	ACBAccountStatusInactive = ACBAccountStatus("inactive")
	ACBAccountStatusDormant  = ACBAccountStatus("dormant")
)

type ACBAccount struct {
	account.IsSourceOfFundImplementMustImport
	bank_account.BankAccount
	//define ACB specific fields here: example status
	Status ACBAccountStatus
}

func (ACBAccount) GetSourceOfFundCode() account.SourceOfFundCode {
	return SourceOfFundCodeACB
}

func (a ACBAccount) IsTheSameSof(other account.IsSourceOfFundItr) bool {
	if v, ok := other.(account.SourceOfFundData); ok {
		return a.IsTheSameSof(v.IsSourceOfFundItr)
	}

	v, ok := other.(ACBAccount)
	if !ok {
		return false
	}
	return v.BankAccount == a.BankAccount && v.AccountName == a.AccountName
}

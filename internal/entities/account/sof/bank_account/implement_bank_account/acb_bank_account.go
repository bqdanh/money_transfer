package implement_bank_account

import (
	"encoding/json"
	"fmt"

	"github.com/bqdanh/money_transfer/internal/entities/account"
	"github.com/bqdanh/money_transfer/internal/entities/account/sof/bank_account"
	"github.com/bqdanh/money_transfer/internal/entities/exceptions"
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

	bank_account.RegisterBankAccountDecoder(SourceOfFundCodeACB, decodeACB)
	bank_account.RegisterBankAccountEncoder(SourceOfFundCodeACB, encodeACB)
}

func decodeACB(bs []byte) (account.IsSourceOfFundItr, error) {
	var acbAc ACBAccount
	err := json.Unmarshal(bs, &acbAc)
	if err != nil {
		return ACBAccount{}, fmt.Errorf("failed to unmarshal ACB account: %w", err)
	}

	return acbAc, nil
}

func encodeACB(sof account.IsSourceOfFundItr) ([]byte, error) {
	bs, err := json.Marshal(sof)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal ACB account: %w", err)
	}
	return bs, nil
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

func (a ACBAccount) IsAvailableForDeposit() error {
	if a.Status != ACBAccountStatusActive {
		return exceptions.NewPreconditionError(
			exceptions.PreconditionReasonSOFBankStatusNotReadyForDeposit,
			exceptions.SubjectSofBank,
			"sof status not ready for deposit",
			map[string]interface{}{
				"sof_code": a.GetSourceOfFundCode(),
				"status":   a.Status,
				"expected": ACBAccountStatusActive,
			},
		)
	}
	return nil
}

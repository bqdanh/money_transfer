package implement_bank_account

import (
	"encoding/json"
	"fmt"

	"github.com/bqdanh/money_transfer/internal/entities/account"
	"github.com/bqdanh/money_transfer/internal/entities/exceptions"
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

	bank_account.RegisterBankAccountDecoder(SourceOfFundCodeVIB, decodeVIB)
	bank_account.RegisterBankAccountEncoder(SourceOfFundCodeVIB, encodeVIB)
}

func decodeVIB(data []byte) (account.IsSourceOfFundItr, error) {
	var ac VibAccount
	err := json.Unmarshal(data, &ac)
	if err != nil {
		return VibAccount{}, fmt.Errorf("failed to decode VIB account: %w", err)
	}
	return ac, nil
}

func encodeVIB(ac account.IsSourceOfFundItr) ([]byte, error) {
	bs, err := json.Marshal(ac)
	if err != nil {
		return nil, fmt.Errorf("failed to encode VIB account: %w", err)
	}
	return bs, nil
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

func (a VibAccount) IsAvailableForDeposit() error {
	if a.Status != VIBAccountStatusActive {
		return exceptions.NewPreconditionError(
			exceptions.PreconditionReasonSOFBankStatusNotReadyForDeposit,
			exceptions.SubjectSofBank,
			"sof status not ready for deposit",
			map[string]interface{}{
				"sof_code": a.GetSourceOfFundCode(),
				"status":   a.Status,
				"expected": VIBAccountStatusActive,
			},
		)
	}
	return nil
}

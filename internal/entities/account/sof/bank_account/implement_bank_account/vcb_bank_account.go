package implement_bank_account

import (
	"encoding/json"
	"fmt"

	"github.com/bqdanh/money_transfer/internal/entities/account"
	"github.com/bqdanh/money_transfer/internal/entities/account/sof/bank_account"
	"github.com/bqdanh/money_transfer/internal/entities/exceptions"
)

const (
	SourceOfFundCodeVCB = account.SourceOfFundCode("VCB")
)

func init() {
	bank_account.RegisterSourceOfFundBankAccount(SourceOfFundCodeVCB)
	bank_account.RegisterSourceOfFundBankAccountConstructor(SourceOfFundCodeVCB, func(a bank_account.BankAccount) (account.SourceOfFundData, error) {
		vcbAc := VCBAccount{
			IsSourceOfFundImplementMustImport: account.IsSourceOfFundImplementMustImport{},
			BankAccount:                       a,
			Status:                            "",
		}

		return account.SourceOfFundData{
			IsSourceOfFundItr: vcbAc,
		}, nil
	})
	bank_account.RegisterBankAccountDecoder(SourceOfFundCodeVCB, decodeVCB)
	bank_account.RegisterBankAccountEncoder(SourceOfFundCodeVCB, encodeVCB)
}

func decodeVCB(bs []byte) (account.IsSourceOfFundItr, error) {
	var vcbAc VCBAccount
	err := json.Unmarshal(bs, &vcbAc)
	if err != nil {
		return VCBAccount{}, fmt.Errorf("failed to unmarshal VCB account: %w", err)
	}

	return vcbAc, nil
}

func encodeVCB(sof account.IsSourceOfFundItr) ([]byte, error) {
	bs, err := json.Marshal(sof)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal VCB account: %w", err)
	}
	return bs, nil
}

// VCBAccountStatus the status of account at VCB bank
type VCBAccountStatus string

const (
	VCBAccountStatusActive   = VCBAccountStatus("active")
	VCBAccountStatusInactive = VCBAccountStatus("inactive")
	VCBAccountStatusLocked   = VCBAccountStatus("locked")
)

type VCBAccount struct {
	account.IsSourceOfFundImplementMustImport
	bank_account.BankAccount
	//define VCB specific fields here: example status
	Status VCBAccountStatus
}

func (VCBAccount) GetSourceOfFundCode() account.SourceOfFundCode {
	return SourceOfFundCodeVCB
}

func (a VCBAccount) IsTheSameSof(other account.IsSourceOfFundItr) bool {
	if v, ok := other.(account.SourceOfFundData); ok {
		return a.IsTheSameSof(v.IsSourceOfFundItr)
	}

	v, ok := other.(VCBAccount)
	if !ok {
		return false
	}
	return v.BankAccount == a.BankAccount && v.AccountName == a.AccountName
}

func (a VCBAccount) IsAvailableForDeposit() error {
	if a.Status != VCBAccountStatusActive {
		return exceptions.NewPreconditionError(
			exceptions.PreconditionReasonSOFBankStatusNotReadyForDeposit,
			exceptions.SubjectSofBank,
			"sof status not ready for deposit",
			map[string]interface{}{
				"sof_code": a.GetSourceOfFundCode(),
				"status":   a.Status,
				"expected": VCBAccountStatusActive,
			},
		)
	}
	return nil
}

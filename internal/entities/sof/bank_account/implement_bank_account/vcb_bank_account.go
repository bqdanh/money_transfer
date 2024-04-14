package implement_bank_account

import (
	"encoding/json"
	"fmt"

	"github.com/bqdanh/money_transfer/internal/entities/account"
	"github.com/bqdanh/money_transfer/internal/entities/sof/bank_account"
)

const (
	SourceOfFundCodeVCB = account.SourceOfFundCode("VCB")
)

func init() {
	bank_account.RegisterSourceOfFundBankAccount(SourceOfFundCodeVCB)
	bank_account.RegisterSourceOfFundBankAccountConstructor(SourceOfFundCodeVCB, func(a bank_account.BankAccount) (account.SourceOfFundData, error) {
		vcbAc := VcbAccount{
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

func decodeVCB(bs []byte) (account.IsSourceOfFundItr, error) {
	var vcbAc VcbAccount
	err := json.Unmarshal(bs, &vcbAc)
	if err != nil {
		return VcbAccount{}, fmt.Errorf("failed to unmarshal VCB account: %w", err)
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

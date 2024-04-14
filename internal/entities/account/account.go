package account

import (
	"fmt"

	"github.com/bqdanh/money_transfer/internal/entities/exceptions"
)

type Account struct {
	ID               int64            `json:"id"`
	UserID           int64            `json:"user_id"`
	Status           Status           `json:"status"`
	SourceOfFundData SourceOfFundData `json:"source_of_fund_data"`
}

type Status string

const (
	StatusNormal   = Status("normal")
	StatusLocked   = Status("locked")
	StatusUnlinked = Status("unlinked")
)

type SourceOfFundType string

const (
	SofTypeBankAccount = SourceOfFundType("bank_account")

	//example for another type of source of fund
	SofTypeEWallet   = SourceOfFundType("ewallet")
	SofTypeBankToken = SourceOfFundType("bank_token")
)

// SourceOfFundCode is the code of source of fund: banks: VIB, Vietcombank, Techcombank, etc, finance institutions: LFVN, Momo, ZaloPay, etc
type SourceOfFundCode string

type SourceOfFundData struct {
	IsSourceOfFundItr
}

func (s SourceOfFundData) GetSourceOfFundCode() SourceOfFundCode {
	if s.IsSourceOfFundItr == nil {
		return ""
	}
	return s.IsSourceOfFundItr.GetSourceOfFundCode()
}

func (s SourceOfFundData) GetSourceOfFundType() SourceOfFundType {
	if s.IsSourceOfFundItr == nil {
		return ""
	}
	return s.IsSourceOfFundItr.GetSourceOfFundType()
}

type IsSourceOfFundItr interface {
	isSourceOfFund()
	GetSourceOfFundCode() SourceOfFundCode
	GetSourceOfFundType() SourceOfFundType
	IsTheSameSof(other IsSourceOfFundItr) bool
	IsAvailableForDeposit() error
}

type IsSourceOfFundImplementMustImport struct {
}

func (b IsSourceOfFundImplementMustImport) isSourceOfFund() {}

func (a Account) IsAvailableForDeposit() error {
	if a.Status != StatusNormal {
		return exceptions.NewPreconditionError(
			exceptions.PreconditionReasonAccountStatusNotReadyForDeposit,
			exceptions.SubjectAccount,
			"account is not ready for deposit",
			map[string]interface{}{
				"account_id": a.ID,
				"status":     a.Status,
				"expected":   StatusNormal,
			},
		)
	}
	if err := a.SourceOfFundData.IsAvailableForDeposit(); err != nil {
		return fmt.Errorf("source of fund is not available for deposit: %w", err)
	}
	return nil
}

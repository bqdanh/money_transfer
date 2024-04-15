package deposit

import (
	"fmt"

	"github.com/bqdanh/money_transfer/internal/entities/account"
	"github.com/bqdanh/money_transfer/internal/entities/transaction"
)

type Deposit struct {
	transaction.IsTransactionDataImplementMustImport
	Source            string `json:"source"`
	BankTransactionID string `json:"bank_transaction_id"` //id at bank for identify transaction, used for audit and reconciliation
	Data              Data   `json:"data"`
}

func CreateDeposit(ac account.Account, source string) (Deposit, error) {
	builder, err := GetDepositDataBuilder(ac.SourceOfFundData.GetSourceOfFundType(), ac.SourceOfFundData.GetSourceOfFundCode())
	if err != nil {
		return Deposit{}, fmt.Errorf("get deposit data builder error: %w", err)
	}
	return Deposit{
		IsTransactionDataImplementMustImport: transaction.IsTransactionDataImplementMustImport{},
		Source:                               source,
		BankTransactionID:                    "",
		Data: Data{
			PartnerData: builder(ac),
		},
	}, nil
}

type Data struct {
	PartnerData IsPartnerDepositData `json:"partner_data"`
}

type IsPartnerDepositData interface {
	isPartnerDepositData()
	GetTransactionStatus() transaction.Status
	GetSOFType() account.SourceOfFundType
	GetSOFCode() account.SourceOfFundCode
	GetPartnerTransactionID() string
}

type ImplementIsPartnerDepositData struct {
}

func (ImplementIsPartnerDepositData) isPartnerDepositData() {}

func (d Deposit) GetType() transaction.Type {
	return transaction.TypeDeposit
}

func (d Deposit) GetTransactionStatus() transaction.Status {
	return d.Data.PartnerData.GetTransactionStatus()
}

func (d Deposit) GetPartnerRefTransactionID() string {
	return d.Data.PartnerData.GetPartnerTransactionID()
}

func (d Deposit) UpdateData(data Data) Deposit {
	d.Data = data
	d.BankTransactionID = data.PartnerData.GetPartnerTransactionID()
	return d
}

package bank_account_vib

import (
	"encoding/json"

	"github.com/bqdanh/money_transfer/internal/entities/account"
	"github.com/bqdanh/money_transfer/internal/entities/account/sof/bank_account/implement_bank_account"
	"github.com/bqdanh/money_transfer/internal/entities/transaction"
	"github.com/bqdanh/money_transfer/internal/entities/transaction/deposit"
)

func init() {
	deposit.RegisterDepositDataBuilder(
		account.SofTypeBankAccount,
		implement_bank_account.SourceOfFundCodeVIB,
		func(ac account.Account) deposit.IsPartnerDepositData {
			return DepositData{
				Status:           DepositStatusInit,
				TransactionRefID: "",
			}
		},
	)

	deposit.RegisterDepositDataEncoder(
		account.SofTypeBankAccount,
		implement_bank_account.SourceOfFundCodeVIB,
		func(data deposit.IsPartnerDepositData) ([]byte, error) {
			return json.Marshal(data)
		},
	)

	deposit.RegisterDepositDataDecoder(
		account.SofTypeBankAccount,
		implement_bank_account.SourceOfFundCodeVIB,
		func(data []byte) (deposit.IsPartnerDepositData, error) {
			var d DepositData
			if err := json.Unmarshal(data, &d); err != nil {
				return nil, err
			}
			return d, nil
		},
	)
}

// DepositStatus represents the status of a deposit transaction at VIB bank.
type DepositStatus string

const (
	DepositStatusInit       = DepositStatus("init")
	DepositStatusProcessing = DepositStatus("processing")
	DepositStatusSuccess    = DepositStatus("success")
	DepositStatusFailed     = DepositStatus("failed")
)

type DepositData struct {
	deposit.ImplementIsPartnerDepositData
	Status           DepositStatus `json:"status"`
	TransactionRefID string        `json:"transaction_ref_id"`
}

func (d DepositData) GetTransactionStatus() transaction.Status {
	//TODO: Implement this method
	return transaction.StatusFailed
}

func (d DepositData) GetSOFType() account.SourceOfFundType {
	return account.SofTypeBankAccount
}

func (d DepositData) GetSOFCode() account.SourceOfFundCode {
	return implement_bank_account.SourceOfFundCodeVIB
}

func (d DepositData) GetPartnerTransactionID() string {
	return d.TransactionRefID
}

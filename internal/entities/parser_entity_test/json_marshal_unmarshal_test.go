package parser_entity_test

import (
	"encoding/json"
	"testing"

	"github.com/bqdanh/money_transfer/internal/entities/account"
	"github.com/bqdanh/money_transfer/internal/entities/account/sof/bank_account"
	"github.com/bqdanh/money_transfer/internal/entities/account/sof/bank_account/implement_bank_account"
	"github.com/bqdanh/money_transfer/internal/entities/currency"
	"github.com/bqdanh/money_transfer/internal/entities/transaction"
	"github.com/bqdanh/money_transfer/internal/entities/transaction/deposit"
	"github.com/bqdanh/money_transfer/internal/entities/transaction/deposit/bank_account_acb"
	"github.com/bqdanh/money_transfer/internal/entities/transaction/deposit/bank_account_vcb"
	"github.com/bqdanh/money_transfer/internal/entities/transaction/deposit/bank_account_vib"
	"github.com/stretchr/testify/assert"
)

var _ = bank_account_vib.DepositData{}
var _ = bank_account_acb.DepositData{}
var _ = bank_account_vcb.DepositData{}

func TestAccountJsonMarshalAndUnMarshal(t *testing.T) {
	testcases := []struct {
		name    string
		account account.Account
		err     error
	}{
		{
			name: "VIB Account",
			account: account.Account{
				ID:     1,
				UserID: 1,
				Status: account.StatusNormal,
				SourceOfFundData: account.SourceOfFundData{
					IsSourceOfFundItr: implement_bank_account.VIBAccount{
						BankAccount: bank_account.BankAccount{
							AccountNumber: "account_number_xyz",
							AccountName:   "account name xyz",
						},
						Status: implement_bank_account.VIBAccountStatusActive,
					},
				},
			},
			err: nil,
		},
		{
			name: "VIB Account",
			account: account.Account{
				ID:     1,
				UserID: 1,
				Status: account.StatusNormal,
				SourceOfFundData: account.SourceOfFundData{
					IsSourceOfFundItr: implement_bank_account.VIBAccount{
						BankAccount: bank_account.BankAccount{
							AccountNumber: "account_number_xyz",
							AccountName:   "account name xyz",
						},
						Status: implement_bank_account.VIBAccountStatusInactive,
					},
				},
			},
			err: nil,
		},
		{
			name: "ACB Account",
			account: account.Account{
				ID:     1,
				UserID: 1,
				Status: account.StatusNormal,
				SourceOfFundData: account.SourceOfFundData{
					IsSourceOfFundItr: implement_bank_account.ACBAccount{
						BankAccount: bank_account.BankAccount{
							AccountNumber: "account_number_xyz",
							AccountName:   "account name xyz",
						},
						Status: implement_bank_account.ACBAccountStatusActive,
					},
				},
			},
			err: nil,
		},
		{
			name: "VIB Account",
			account: account.Account{
				ID:     1,
				UserID: 1,
				Status: account.StatusNormal,
				SourceOfFundData: account.SourceOfFundData{
					IsSourceOfFundItr: implement_bank_account.VIBAccount{
						BankAccount: bank_account.BankAccount{
							AccountNumber: "account_number_xyz",
							AccountName:   "account name xyz",
						},
						Status: implement_bank_account.VIBAccountStatusActive,
					},
				},
			},
			err: nil,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			bs, err := json.Marshal(tc.account)
			assert.ErrorIs(t, err, tc.err)
			t.Log(string(bs))
			if err != nil {
				return
			}
			ac := account.Account{}
			err = json.Unmarshal(bs, &ac)
			assert.NoError(t, err)
			assert.Equal(t, tc.account, ac)
		})
	}

}

func TestJsonMarshalVIBDepositTransaction(t *testing.T) {
	ac := account.Account{
		ID:     1,
		UserID: 11,
		Status: account.StatusNormal,
		SourceOfFundData: account.SourceOfFundData{
			IsSourceOfFundItr: implement_bank_account.VIBAccount{
				BankAccount: bank_account.BankAccount{
					AccountNumber: "AccountNumber_XYZ",
					AccountName:   "AccountName_XYZ",
				},
				Status: implement_bank_account.VIBAccountStatusActive,
			},
		},
	}
	depositData, err := deposit.CreateDeposit(ac, "testing deposit data")
	assert.Nil(t, err)
	if err != nil {
		t.Fatal(err)
	}
	amount := currency.Amount{
		Currency: currency.VND,
		Amount:   10_000,
	}

	depositTransaction := transaction.CreateTransaction(ac, amount, "ut description", transaction.TypeDeposit, transaction.Data{
		IsTransactionDataItr: depositData,
	})
	bs, err := json.Marshal(depositTransaction)
	assert.Nil(t, err)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("json: %s", string(bs))

	depositTransaction2 := transaction.Transaction{}
	err = json.Unmarshal(bs, &depositTransaction2)
	assert.Nil(t, err)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, depositTransaction, depositTransaction2)
}

func TestJsonMarshalACBDepositTransaction(t *testing.T) {
	ac := account.Account{
		ID:     1,
		UserID: 11,
		Status: account.StatusNormal,
		SourceOfFundData: account.SourceOfFundData{
			IsSourceOfFundItr: implement_bank_account.ACBAccount{
				BankAccount: bank_account.BankAccount{
					AccountNumber: "AccountNumber_XYZ",
					AccountName:   "AccountName_XYZ",
				},
				Status: implement_bank_account.ACBAccountStatusActive,
			},
		},
	}
	depositData, err := deposit.CreateDeposit(ac, "testing deposit data")
	assert.Nil(t, err)
	if err != nil {
		t.Fatal(err)
	}
	amount := currency.Amount{
		Currency: currency.VND,
		Amount:   10_000,
	}

	depositTransaction := transaction.CreateTransaction(ac, amount, "ut description", transaction.TypeDeposit, transaction.Data{
		IsTransactionDataItr: depositData,
	})
	bs, err := json.Marshal(depositTransaction)
	assert.Nil(t, err)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("json: %s", string(bs))

	depositTransaction2 := transaction.Transaction{}
	err = json.Unmarshal(bs, &depositTransaction2)
	assert.Nil(t, err)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, depositTransaction, depositTransaction2)
}

func TestJsonMarshalVCBDepositTransaction(t *testing.T) {
	ac := account.Account{
		ID:     1,
		UserID: 11,
		Status: account.StatusNormal,
		SourceOfFundData: account.SourceOfFundData{
			IsSourceOfFundItr: implement_bank_account.VCBAccount{
				BankAccount: bank_account.BankAccount{
					AccountNumber: "AccountNumber_XYZ",
					AccountName:   "AccountName_XYZ",
				},
				Status: implement_bank_account.VCBAccountStatusActive,
			},
		},
	}
	depositData, err := deposit.CreateDeposit(ac, "testing deposit data")
	assert.Nil(t, err)
	if err != nil {
		t.Fatal(err)
	}
	amount := currency.Amount{
		Currency: currency.VND,
		Amount:   10_000,
	}

	depositTransaction := transaction.CreateTransaction(ac, amount, "ut description", transaction.TypeDeposit, transaction.Data{
		IsTransactionDataItr: depositData,
	})
	bs, err := json.Marshal(depositTransaction)
	assert.Nil(t, err)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("json: %s", string(bs))

	depositTransaction2 := transaction.Transaction{}
	err = json.Unmarshal(bs, &depositTransaction2)
	assert.Nil(t, err)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, depositTransaction, depositTransaction2)
}

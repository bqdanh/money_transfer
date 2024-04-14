package parser_entity_test

import (
	"encoding/json"
	"testing"

	"github.com/bqdanh/money_transfer/internal/entities/account"
	"github.com/bqdanh/money_transfer/internal/entities/sof/bank_account"
	"github.com/bqdanh/money_transfer/internal/entities/sof/bank_account/implement_bank_account"
	"github.com/stretchr/testify/assert"
)

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
					IsSourceOfFundItr: implement_bank_account.VibAccount{
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
					IsSourceOfFundItr: implement_bank_account.VibAccount{
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
					IsSourceOfFundItr: implement_bank_account.VibAccount{
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

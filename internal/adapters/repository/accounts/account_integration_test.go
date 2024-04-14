//go:build integration

package accounts

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/bqdanh/money_transfer/internal/entities/account"
	"github.com/bqdanh/money_transfer/internal/entities/sof/bank_account"
	"github.com/bqdanh/money_transfer/internal/entities/sof/bank_account/implement_bank_account"
	"github.com/bqdanh/money_transfer/pkg/database"
	"github.com/bqdanh/money_transfer/pkg/osenv"
	"github.com/stretchr/testify/assert"
)

func getSqlDatabaseTest(t *testing.T) (*sql.DB, error) {
	db, err := database.NewMysqlDatabaseConn(&database.Config{
		Address:              osenv.GetStringEnvWithDefault("MYSQL_ADDRESS", "0.0.0.0:3306"),
		User:                 osenv.GetStringEnvWithDefault("MYSQL_USER", "app_user"),
		Passwd:               osenv.GetStringEnvWithDefault("MYSQL_PASSWD", "pwd123"),
		AllowNativePasswords: true,
		DatabaseName:         osenv.GetStringEnvWithDefault("MYSQL_DATABASE_NAME", "money_transfer"),
		MaxOpenConn:          osenv.GetIntEnvWithDefault("MYSQL_MAX_OPEN_CONN", 10),
		MaxIdleConn:          osenv.GetIntEnvWithDefault("MYSQL_MAX_IDLE_CONN", 10),
		ConnMaxLifeTime:      osenv.GetDurationEnvWithDefault("MYSQL_CONN_MAX_LIFE_TIME", 1*time.Minute),
	})
	assert.NoError(t, err)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func TestCreateAccountAndGetAccountByUserID(t *testing.T) {
	//to manage test data, we use accountTestID to create account and delete account after test
	accountTestID := int64(11)

	testcases := []struct {
		name     string
		accounts []account.Account
	}{
		{
			name: "create account success",
			accounts: []account.Account{
				{
					ID:     0,
					UserID: accountTestID,
					Status: account.StatusNormal,
					SourceOfFundData: account.SourceOfFundData{
						IsSourceOfFundItr: implement_bank_account.VibAccount{
							BankAccount: bank_account.BankAccount{
								AccountNumber: "account_number",
								AccountName:   "account_name",
							},
							Status: implement_bank_account.VIBAccountStatusActive,
						},
					},
				},
			},
		},
		{
			name: "create accounts success",
			accounts: []account.Account{
				{
					ID:     0,
					UserID: accountTestID,
					Status: account.StatusNormal,
					SourceOfFundData: account.SourceOfFundData{
						IsSourceOfFundItr: implement_bank_account.VibAccount{
							BankAccount: bank_account.BankAccount{
								AccountNumber: "account_number",
								AccountName:   "account_name",
							},
							Status: implement_bank_account.VIBAccountStatusActive,
						},
					},
				},
				{
					ID:     0,
					UserID: accountTestID,
					Status: account.StatusNormal,
					SourceOfFundData: account.SourceOfFundData{
						IsSourceOfFundItr: implement_bank_account.ACBAccount{
							BankAccount: bank_account.BankAccount{
								AccountNumber: "account_number",
								AccountName:   "account_name",
							},
							Status: implement_bank_account.ACBAccountStatusActive,
						},
					},
				},
				{
					ID:     0,
					UserID: accountTestID,
					Status: account.StatusNormal,
					SourceOfFundData: account.SourceOfFundData{
						IsSourceOfFundItr: implement_bank_account.VcbAccount{
							BankAccount: bank_account.BankAccount{
								AccountNumber: "account_number",
								AccountName:   "account_name",
							},
							Status: implement_bank_account.VCBAccountStatusActive,
						},
					},
				},
			},
		},
	}
	db, err := getSqlDatabaseTest(t)
	if err != nil {
		t.Errorf("create db failed: %v", err)
		return
	}
	repo, err := NewAccountMysqlRepository(db)
	if err != nil {
		t.Errorf("create repo failed: %v", err)
		return
	}
	defer func() {
		err := repo.DeleteAccountByUserID(context.Background(), accountTestID)
		assert.NoError(t, err)
	}()

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			newacs := make([]account.Account, 0)
			defer func() {
				err := repo.DeleteAccountByUserID(context.Background(), tc.accounts[0].UserID)
				assert.NoError(t, err)
			}()

			for _, ac := range tc.accounts {
				newac, err := repo.CreateAccount(context.Background(), ac)
				assert.NoError(t, err)
				if err != nil {
					return
				}
				newacs = append(newacs, newac)
			}
			if len(newacs) != len(tc.accounts) {
				assert.Equal(t, newacs, tc.accounts)
				return
			}
			for idx := range newacs {
				tc.accounts[idx].ID = newacs[idx].ID
			}
			assert.Equal(t, newacs, tc.accounts)

			acs, err := repo.GetAccountsByUserID(context.Background(), tc.accounts[0].UserID)
			assert.NoError(t, err)
			assert.Equal(t, acs, tc.accounts)
		})
	}

}

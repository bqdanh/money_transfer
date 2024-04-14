package accounts

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/bqdanh/money_transfer/internal/adapters/repository/sqlc/moneytransfer"
	"github.com/bqdanh/money_transfer/internal/entities/account"
)

type AccountMysqlRepository struct {
	db *sql.DB
}

func NewAccountMysqlRepository(db *sql.DB) (AccountMysqlRepository, error) {
	if db == nil {
		return AccountMysqlRepository{}, fmt.Errorf("db must not nil")
	}
	return AccountMysqlRepository{
		db: db,
	}, nil
}

func (r AccountMysqlRepository) GetAccountsByUserID(ctx context.Context, userID int64) ([]account.Account, error) {
	q := moneytransfer.New(r.db)
	daAccounts, err := q.GetAccountsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get accounts by user id: %w", err)
	}
	acs, err := fromAccountsDA2AccountsEntity(daAccounts)
	if err != nil {
		return nil, fmt.Errorf("failed to parse accounts: %w", err)
	}
	return acs, nil
}

// CreateAccount create account for user, return account with ID, ID is unique
func (r AccountMysqlRepository) CreateAccount(ctx context.Context, ac account.Account) (account.Account, error) {
	q := moneytransfer.New(r.db)
	daAccount, err := fromAccountEntity2AccountDA(ac)
	if err != nil {
		return account.Account{}, fmt.Errorf("failed to parse account: %w", err)
	}
	result, err := q.InsertAccount(ctx, &moneytransfer.InsertAccountParams{
		UserID:      daAccount.UserID,
		AccountType: daAccount.AccountType,
		AccountData: daAccount.AccountData,
	})
	if err != nil {
		return account.Account{}, fmt.Errorf("insert account: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return account.Account{}, fmt.Errorf("get last insert id: %w", err)
	}
	ac.ID = id
	return ac, nil
}

func (r AccountMysqlRepository) DeleteAccountByUserID(ctx context.Context, userID int64) error {
	q := moneytransfer.New(r.db)
	_, err := q.DeleteAccountByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("delete account by id: %w", err)
	}

	return nil
}

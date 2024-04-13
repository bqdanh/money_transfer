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
	_ = daAccounts
	return nil, fmt.Errorf("not implemented")
}

// CreateAccount create account for user, return account with ID, ID is unique
func (r AccountMysqlRepository) CreateAccount(ctx context.Context, ac account.Account) (account.Account, error) {
	return account.Account{}, fmt.Errorf("not implemented")
}

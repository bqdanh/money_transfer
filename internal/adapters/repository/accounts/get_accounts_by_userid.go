package accounts

import (
	"context"
	"fmt"

	"github.com/bqdanh/money_transfer/internal/adapters/repository/sqlc/moneytransfer"
	"github.com/bqdanh/money_transfer/internal/entities/account"
)

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

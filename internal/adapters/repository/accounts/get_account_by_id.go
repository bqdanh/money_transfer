package accounts

import (
	"context"
	"fmt"

	"github.com/bqdanh/money_transfer/internal/adapters/repository/sqlc/moneytransfer"
	"github.com/bqdanh/money_transfer/internal/entities/account"
)

func (r AccountMysqlRepository) GetAccountsByID(ctx context.Context, accountID int64) (account.Account, error) {
	q := moneytransfer.New(r.db)
	result, err := q.GetAccountByID(ctx, accountID)
	if err != nil {
		return account.Account{}, fmt.Errorf("get account by id: %w", err)
	}
	ac, err := fromAccountDA2AccountEntity(result)
	if err != nil {
		return account.Account{}, fmt.Errorf("failed to parse account: %w", err)
	}
	return ac, nil
}

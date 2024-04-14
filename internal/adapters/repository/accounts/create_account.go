package accounts

import (
	"context"
	"fmt"

	"github.com/bqdanh/money_transfer/internal/adapters/repository/sqlc/moneytransfer"
	"github.com/bqdanh/money_transfer/internal/entities/account"
)

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

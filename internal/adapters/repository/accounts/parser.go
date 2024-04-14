package accounts

import (
	"encoding/json"
	"fmt"

	"github.com/bqdanh/money_transfer/internal/adapters/repository/sqlc/moneytransfer"
	"github.com/bqdanh/money_transfer/internal/entities/account"
)

func fromAccountDA2AccountEntity(daAccount *moneytransfer.Account) (account.Account, error) {
	ac := account.Account{}
	err := json.Unmarshal(daAccount.AccountData, &ac.SourceOfFundData)
	if err != nil {
		return account.Account{}, fmt.Errorf("failed to unmarshal account data: %w", err)
	}
	return ac, nil
}

func fromAccountsDA2AccountsEntity(daAccounts []*moneytransfer.Account) ([]account.Account, error) {
	accounts := make([]account.Account, 0, len(daAccounts))
	for _, daAccount := range daAccounts {
		ac, err := fromAccountDA2AccountEntity(daAccount)
		if err != nil {
			return nil, fmt.Errorf("failed to parse account: %w", err)
		}
		accounts = append(accounts, ac)
	}
	return accounts, nil
}

func fromAccountEntity2AccountDA(ac account.Account) (moneytransfer.Account, error) {
	data, err := json.Marshal(ac)
	if err != nil {
		return moneytransfer.Account{}, fmt.Errorf("failed to marshal account data: %w", err)
	}
	daAccount := moneytransfer.Account{
		ID:          ac.ID,
		UserID:      ac.UserID,
		AccountType: moneytransfer.AccountAccountType(ac.SourceOfFundData.GetSourceOfFundType()),
		AccountData: data,
	}
	return daAccount, nil
}

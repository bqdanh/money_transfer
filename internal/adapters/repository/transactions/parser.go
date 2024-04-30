package transactions

import (
	"encoding/json"
	"fmt"

	"github.com/bqdanh/money_transfer/internal/adapters/repository/sqlc/moneytransfer"
	"github.com/bqdanh/money_transfer/internal/entities/transaction"
)

func fromTransactionDAOToTransaction(da moneytransfer.Transaction) (transaction.Transaction, error) {
	t := transaction.Transaction{}
	err := json.Unmarshal(da.Data, &t)
	if err != nil {
		return transaction.Transaction{}, fmt.Errorf("error unmarshal transaction: %w", err)
	}
	t.ID = da.ID
	return t, nil
}

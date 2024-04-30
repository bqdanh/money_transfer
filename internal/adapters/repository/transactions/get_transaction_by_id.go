package transactions

import (
	"context"
	"fmt"

	"github.com/bqdanh/money_transfer/internal/adapters/repository/sqlc/moneytransfer"
	"github.com/bqdanh/money_transfer/internal/entities/transaction"
)

func (r TransactionMysqlRepository) GetTransactionByID(ctx context.Context, transID int64) (transaction.Transaction, error) {
	q := moneytransfer.New(r.db)
	result, err := q.GetTransactionByID(ctx, transID)
	if err != nil {
		return transaction.Transaction{}, err
	}
	t, err := fromTransactionDAOToTransaction(*result)
	if err != nil {
		return transaction.Transaction{}, fmt.Errorf("error convert transaction dao to transaction: %w", err)
	}
	return t, nil
}

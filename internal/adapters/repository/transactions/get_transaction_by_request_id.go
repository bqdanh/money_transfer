package transactions

import (
	"context"
	"fmt"

	"github.com/bqdanh/money_transfer/internal/adapters/repository/sqlc/moneytransfer"
	"github.com/bqdanh/money_transfer/internal/entities/account"
	"github.com/bqdanh/money_transfer/internal/entities/transaction"
)

func (r TransactionMysqlRepository) GetTransactionByRequestID(ctx context.Context, ac account.Account, requestID string) (transaction.Transaction, error) {
	q := moneytransfer.New(r.db)
	result, err := q.GetTransactionByRequestID(ctx, &moneytransfer.GetTransactionByRequestIDParams{
		AccountID: ac.ID,
		RequestID: requestID,
	})
	if err != nil {
		return transaction.Transaction{}, fmt.Errorf("error get transaction by request id: %w", err)
	}
	t, err := fromTransactionDAOToTransaction(*result)
	if err != nil {
		return transaction.Transaction{}, fmt.Errorf("error convert transaction dao to transaction: %w", err)
	}
	return t, nil
}

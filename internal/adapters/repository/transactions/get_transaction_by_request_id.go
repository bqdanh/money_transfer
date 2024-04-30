package transactions

import (
	"context"
	"encoding/json"
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
	t := transaction.Transaction{}
	err = json.Unmarshal(result.Data, &t)
	if err != nil {
		return transaction.Transaction{}, fmt.Errorf("error unmarshal transaction: %w", err)
	}
	t.ID = result.ID
	return t, nil
}

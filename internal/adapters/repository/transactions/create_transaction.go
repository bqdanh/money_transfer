package transactions

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bqdanh/money_transfer/internal/adapters/repository/sqlc/moneytransfer"
	"github.com/bqdanh/money_transfer/internal/entities/transaction"
)

func (r TransactionMysqlRepository) CreateTransaction(ctx context.Context, t transaction.Transaction) (transaction.Transaction, error) {
	q := moneytransfer.New(r.db)
	bs, err := json.Marshal(t)
	if err != nil {
		return t, fmt.Errorf("error marshal transaction: %w", err)
	}
	result, err := q.CreateTransaction(ctx, &moneytransfer.CreateTransactionParams{
		AccountID:               t.Account.ID,
		Amount:                  t.Amount.String(),
		Version:                 int32(t.Version),
		RequestID:               t.RequestID,
		Description:             t.Description,
		PartnerRefTransactionID: t.GetPartnerRefTransactionID(),
		Status:                  string(t.Status),
		Type:                    string(t.Type),
		Data:                    bs,
	})
	if err != nil {
		return t, fmt.Errorf("error create transaction: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return t, fmt.Errorf("error get last insert id: %w", err)
	}
	t.ID = id
	return t, nil
}

package transactions

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/bqdanh/money_transfer/internal/adapters/repository/sqlc/moneytransfer"
	"github.com/bqdanh/money_transfer/internal/entities/transaction"
	"github.com/bqdanh/money_transfer/pkg/database"
)

func (r TransactionMysqlRepository) UpdateTransaction(ctx context.Context, t transaction.Transaction, evt transaction.Event) error {
	tx, errtx := r.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if errtx != nil {
		return fmt.Errorf("error begin transaction: %w", errtx)
	}
	err := database.ExecuteTransaction(tx, func() error {
		q := moneytransfer.New(tx)
		transDataBs, err := json.Marshal(t)
		if err != nil {
			return fmt.Errorf("error marshal transaction: %w", err)
		}
		err = q.UpdateTransaction(ctx, &moneytransfer.UpdateTransactionParams{
			Amount:                  t.Amount.String(),
			Version:                 int32(t.Version),
			RequestID:               t.RequestID,
			Description:             t.Description,
			PartnerRefTransactionID: t.GetPartnerRefTransactionID(),
			Status:                  string(t.Status),
			Type:                    string(t.Type),
			Data:                    transDataBs,
			ID:                      t.ID,
		})
		if err != nil {
			return fmt.Errorf("error update transaction: %w", err)
		}
		evtbs, err := json.Marshal(evt)
		if err != nil {
			return fmt.Errorf("error marshal event: %w", err)
		}
		err = q.CreateTransactionEvent(ctx, &moneytransfer.CreateTransactionEventParams{
			TransactionID: t.ID,
			Version:       int32(t.Version),
			Data:          evtbs,
			EventName:     string(evt.Name),
		})
		if err != nil {
			return fmt.Errorf("error create transaction event: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error execute transaction: %w", err)
	}
	return nil
}

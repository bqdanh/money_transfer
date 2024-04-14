package accounts

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/bqdanh/money_transfer/internal/adapters/repository/sqlc/moneytransfer"
)

type AccountMysqlRepository struct {
	db *sql.DB
}

func NewAccountMysqlRepository(db *sql.DB) (AccountMysqlRepository, error) {
	if db == nil {
		return AccountMysqlRepository{}, fmt.Errorf("db must not nil")
	}
	return AccountMysqlRepository{
		db: db,
	}, nil
}

func (r AccountMysqlRepository) DeleteAccountByUserID(ctx context.Context, userID int64) error {
	q := moneytransfer.New(r.db)
	_, err := q.DeleteAccountByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("delete account by id: %w", err)
	}

	return nil
}

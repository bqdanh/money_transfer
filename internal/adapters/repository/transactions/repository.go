package transactions

import (
	"database/sql"
	"fmt"
)

type TransactionMysqlRepository struct {
	db *sql.DB
}

func NewTransactionMysqlRepository(db *sql.DB) (TransactionMysqlRepository, error) {
	if db == nil {
		return TransactionMysqlRepository{}, fmt.Errorf("db must not nil")
	}
	return TransactionMysqlRepository{
		db: db,
	}, nil
}

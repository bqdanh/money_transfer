package transactions

import "database/sql"

type TransactionMysqlRepository struct {
	db *sql.DB
}

package database

import (
	"database/sql"
	"fmt"
)

func ExecuteTransaction(tx *sql.Tx, f func() error) error {
	var err error
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		if err2 := tx.Commit(); err2 != nil {
			err = fmt.Errorf("error commit transaction: %w", err2)
		}
	}()
	err = f()
	if err != nil {
		return err
	}
	return nil
}

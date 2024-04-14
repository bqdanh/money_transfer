package users

import (
	"database/sql"
	"fmt"
)

type UserMysqlRepository struct {
	db *sql.DB
}

func NewUserMysqlRepository(db *sql.DB) (UserMysqlRepository, error) {
	if db == nil {
		return UserMysqlRepository{}, fmt.Errorf("db must not nil")
	}
	return UserMysqlRepository{
		db: db,
	}, nil
}

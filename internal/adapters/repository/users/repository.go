package users

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/bqdanh/money_transfer/internal/adapters/repository/sqlc/moneytransfer"
	"github.com/bqdanh/money_transfer/internal/entities/user"
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

// CreateUser persist user into repo with generate user id is unique and make sure UserName must no duplicated
// return User with ID
func (r UserMysqlRepository) CreateUser(ctx context.Context, u user.User) (ru user.User, err error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		return u, fmt.Errorf("database begin tx: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
		if err != nil {
			_ = tx.Rollback()
			err = fmt.Errorf("commit transaction: %w", err)
			return
		}
	}()

	q := moneytransfer.New(tx)
	rs, err := q.InsertUser(ctx, &moneytransfer.InsertUserParams{
		UserName: u.UserName,
		Password: u.Password,
		FullName: u.FullName,
		Phone:    u.Phone,
	})
	if err != nil {
		return u, fmt.Errorf("insert user into msyql got error: %w", err)
	}
	userId, err := rs.LastInsertId()
	if err != nil {
		return u, fmt.Errorf("get last insert id: %w", err)
	}
	u, err = u.WithID(userId)
	if err != nil {
		return u, fmt.Errorf("user with id: %w", err)
	}
	return u, nil
}

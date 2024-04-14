package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/bqdanh/money_transfer/internal/adapters/repository/sqlc/moneytransfer"
	"github.com/bqdanh/money_transfer/internal/entities/exceptions"
	"github.com/bqdanh/money_transfer/internal/entities/user"
)

func (r UserMysqlRepository) GetUserByUsername(ctx context.Context, username string) (user.User, error) {
	q := moneytransfer.New(r.db)
	u, err := q.GetUserByUserName(ctx, username)
	if errors.Is(err, sql.ErrNoRows) {
		return user.User{}, exceptions.NewPreconditionError(
			exceptions.PreconditionReasonUserNotFound,
			exceptions.SubjectUser,
			"user not found",
			map[string]interface{}{
				"username": username,
			},
		)
	}
	if err != nil {
		return user.User{}, fmt.Errorf("get user by username: %w", err)
	}
	if u == nil {
		return user.User{}, exceptions.NewPreconditionError(
			exceptions.PreconditionReasonUserNotFound,
			exceptions.SubjectUser,
			"user not found",
			map[string]interface{}{
				"username": username,
			},
		)
	}
	return fromUserDAOToUser(*u), nil
}

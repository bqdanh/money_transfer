package users

import (
	"github.com/bqdanh/money_transfer/internal/adapters/repository/sqlc/moneytransfer"
	"github.com/bqdanh/money_transfer/internal/entities/user"
)

func fromUserDAOToUser(u moneytransfer.User) user.User {
	return user.User{
		ID:       u.ID,
		UserName: u.UserName,
		Password: u.Password,
		FullName: u.FullName,
		Phone:    u.Phone,
	}
}

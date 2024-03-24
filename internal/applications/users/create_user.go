package users

import (
	"context"
	"fmt"

	"github.com/banhquocdanh/money_transfer/internal/entities/exceptions"
	"github.com/banhquocdanh/money_transfer/internal/entities/user"
)

type CreateUser struct {
	ur userRepository
}

type userRepository interface {
	//CreateUserWithGenerateID persist user into repo with generate user id is unique and make sure UserName must no duplicated
	//return User with ID
	CreateUserWithGenerateID(ctx context.Context, u user.User) (user.User, error)
}

type CreateUserParams struct {
	UserName string
	Password string
	FullName string
	Phone    string
}

const (
	minimumLengthOfUserName = 6
	maximumLengthOfUserName = 15
)

func validateParams(p CreateUserParams) error {
	if len(p.UserName) < minimumLengthOfUserName || len(p.UserName) > maximumLengthOfUserName {
		return exceptions.NewInvalidArgumentError(
			"UserName",
			fmt.Sprintf("length of user name must between %d and %d", minimumLengthOfUserName, maximumLengthOfUserName),
			map[string]interface{}{
				"min_len": minimumLengthOfUserName,
				"max_len": maximumLengthOfUserName,
			},
		)
	}
	if len(p.Password) < minimumLengthOfUserName {

	}
}

func (CreateUser) Handle(ctx context.Context, p CreateUserParams) (user.User, error) {

}

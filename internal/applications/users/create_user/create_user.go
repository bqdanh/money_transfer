package create_user

import (
	"context"
	"fmt"

	"github.com/bqdanh/money_transfer/internal/entities/user"
)

//go:generate mockgen --source=./create_user.go --destination=./mocks.go --package=create_user .

type CreateUser struct {
	ur userRepository
}

type userRepository interface {
	//CreateUser persist user into repo with generate user id is unique and make sure UserName must no duplicated
	//return User with ID
	CreateUser(ctx context.Context, u user.User) (user.User, error)
}

func NewCreateUser(ur userRepository) (CreateUser, error) {
	if ur == nil {
		return CreateUser{}, fmt.Errorf("user repository must not nil")
	}
	return CreateUser{
		ur: ur,
	}, nil
}

type CreateUserParams struct {
	UserName string
	Password string
	FullName string
	Phone    string
}

func (h CreateUser) Handle(ctx context.Context, p CreateUserParams) (user.User, error) {
	u, err := user.CreateUser(p.UserName, p.Password, p.FullName, p.Phone)
	if err != nil {
		return user.User{}, fmt.Errorf("create entity user: %w", err)
	}
	u, err = h.ur.CreateUser(ctx, u)
	if err != nil {
		return user.User{}, fmt.Errorf("repository create user: %w", err)
	}
	return u, nil
}

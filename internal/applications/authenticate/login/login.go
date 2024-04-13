package login

import (
	"context"
	"fmt"

	"github.com/bqdanh/money_transfer/internal/applications/authenticate/generate_user_token"
	"github.com/bqdanh/money_transfer/internal/entities/exceptions"
	"github.com/bqdanh/money_transfer/internal/entities/user"
)

type Login struct {
	userNamePasswordValidator userNamePasswordValidator
	generateUserToken         generate_user_token.GenerateUserToken
}

func NewLogin(u userNamePasswordValidator, g generate_user_token.GenerateUserToken) (Login, error) {
	if u == nil {
		return Login{}, exceptions.NewInvalidArgumentError("UserNamePasswordValidator", "user name password validator must not nil", nil)
	}
	if g == (generate_user_token.GenerateUserToken{}) {
		return Login{}, exceptions.NewInvalidArgumentError("GenerateUserToken", "generate user token must not empty", nil)
	}
	return Login{
		userNamePasswordValidator: u,
		generateUserToken:         g,
	}, nil
}

//go:generate mockgen --source=./login.go --destination=./mocks.go --package=login .

type userNamePasswordValidator interface {
	ValidateUserNamePassword(ctx context.Context, username, password string) (user.User, error)
}

type LoginParams struct {
	Username string
	Password string
}

type LoginResponse struct {
	User  user.User
	Token string
}

func ValidateLoginParams(p LoginParams) error {
	if p.Username == "" {
		return exceptions.NewInvalidArgumentError("Username", "username must not empty", nil)
	}
	if p.Password == "" {
		return exceptions.NewInvalidArgumentError("Password", "password must not empty", nil)
	}
	return nil
}

func (l Login) Handle(ctx context.Context, p LoginParams) (LoginResponse, error) {
	if err := ValidateLoginParams(p); err != nil {
		return LoginResponse{}, err
	}
	u, err := l.userNamePasswordValidator.ValidateUserNamePassword(ctx, p.Username, p.Password)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("validate username password: %w", err)
	}
	token, err := l.generateUserToken.Handle(ctx, u)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("generate user token: %w", err)
	}
	return LoginResponse{
		Token: token,
		User:  u,
	}, nil
}

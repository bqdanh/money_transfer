package login

import (
	"context"
	"fmt"

	"github.com/bqdanh/money_transfer/internal/entities/exceptions"
	"github.com/bqdanh/money_transfer/internal/entities/user"
)

type LoginWithUsernamePassword struct {
	userRepo userRepository
}

func NewLoginWithUsernamePassword(ur userRepository) (LoginWithUsernamePassword, error) {
	if ur == nil {
		return LoginWithUsernamePassword{}, exceptions.NewInvalidArgumentError("UserRepository", "user repository must not nil", nil)
	}
	return LoginWithUsernamePassword{
		userRepo: ur,
	}, nil
}

type LoginWithUsernamePasswordParams struct {
	Username string
	Password string
}

func ValidateLoginWithUsernamePasswordParams(p LoginWithUsernamePasswordParams) error {
	if p.Username == "" {
		return exceptions.NewInvalidArgumentError("Username", "username must not empty", nil)
	}
	if p.Password == "" {
		return exceptions.NewInvalidArgumentError("Password", "password must not empty", nil)
	}
	return nil
}

//go:generate mockgen --source=./login_username_pw.go --destination=./mocks.go --package=login .

type userRepository interface {
	GetUserByUsername(ctx context.Context, username string) (user.User, error)
}

func (h LoginWithUsernamePassword) Handle(ctx context.Context, p LoginWithUsernamePasswordParams) (user.User, error) {
	if err := ValidateLoginWithUsernamePasswordParams(p); err != nil {
		return user.User{}, fmt.Errorf("validate login with username: %w", err)
	}
	u, err := h.userRepo.GetUserByUsername(ctx, p.Username)
	if err != nil {
		return user.User{}, fmt.Errorf("get user by username: %w", err)
	}
	if err := user.ComparePassword(u.Password, p.Password); err != nil {
		return user.User{}, exceptions.NewPreconditionError(exceptions.PreconditionTypePasswordNotMatch, exceptions.SubjectUser, "password not match", nil)
	}
	return u, nil
}

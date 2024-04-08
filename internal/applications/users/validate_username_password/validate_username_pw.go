package validate_username_password

import (
	"context"
	"fmt"

	"github.com/bqdanh/money_transfer/internal/entities/exceptions"
	"github.com/bqdanh/money_transfer/internal/entities/user"
)

type ValidateUsernamePassword struct {
	userRepo userRepository
}

func NewValidateUsernamePassword(ur userRepository) (ValidateUsernamePassword, error) {
	if ur == nil {
		return ValidateUsernamePassword{}, exceptions.NewInvalidArgumentError("UserRepository", "user repository must not nil", nil)
	}
	return ValidateUsernamePassword{
		userRepo: ur,
	}, nil
}

type ValidateUsernamePasswordParams struct {
	Username string
	Password string
}

func ValidateParams(p ValidateUsernamePasswordParams) error {
	if p.Username == "" {
		return exceptions.NewInvalidArgumentError("Username", "username must not empty", nil)
	}
	if p.Password == "" {
		return exceptions.NewInvalidArgumentError("Password", "password must not empty", nil)
	}
	return nil
}

//go:generate mockgen --source=./validate_username_pw.go --destination=./mocks.go --package=validate_username_password .

type userRepository interface {
	GetUserByUsername(ctx context.Context, username string) (user.User, error)
}

func (h ValidateUsernamePassword) Handle(ctx context.Context, p ValidateUsernamePasswordParams) (user.User, error) {
	if err := ValidateParams(p); err != nil {
		return user.User{}, fmt.Errorf("validate params: %w", err)
	}
	u, err := h.userRepo.GetUserByUsername(ctx, p.Username)
	if err != nil {
		return user.User{}, fmt.Errorf("get user by username: %w", err)
	}
	if err := user.ComparePassword(u.Password, p.Password); err != nil {
		return user.User{}, exceptions.NewPreconditionError(exceptions.PreconditionReasonPasswordNotMatch, exceptions.SubjectUser, "password not match", map[string]interface{}{})
	}
	return u, nil
}

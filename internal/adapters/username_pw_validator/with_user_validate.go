package username_pw_validator

import (
	"context"
	"fmt"

	"github.com/bqdanh/money_transfer/internal/applications/users/validate_username_password"
	"github.com/bqdanh/money_transfer/internal/entities/user"
)

type ValidateUserNamePasswordWithUserUseCase struct {
	validator validate_username_password.ValidateUsernamePassword
}

func NewValidateUserNamePasswordWithUserUseCase(v validate_username_password.ValidateUsernamePassword) (ValidateUserNamePasswordWithUserUseCase, error) {
	if v == (validate_username_password.ValidateUsernamePassword{}) {
		return ValidateUserNamePasswordWithUserUseCase{}, fmt.Errorf("validate username password must not empty")
	}
	return ValidateUserNamePasswordWithUserUseCase{
		validator: v,
	}, nil
}

func (v ValidateUserNamePasswordWithUserUseCase) ValidateUserNamePassword(ctx context.Context, username, password string) (user.User, error) {
	u, err := v.validator.Handle(ctx, validate_username_password.ValidateUsernamePasswordParams{
		Username: username,
		Password: password,
	})
	if err != nil {
		return user.User{}, fmt.Errorf("validator handle: %w", err)
	}
	return u, nil
}

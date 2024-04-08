package validate_user_token

import (
	"context"
	"time"

	"github.com/bqdanh/money_transfer/internal/entities/authenticate"
	"github.com/bqdanh/money_transfer/internal/entities/exceptions"
)

type ValidateUserToken struct {
	validator tokenValidator
}

func NewValidateUserToken(v tokenValidator) (ValidateUserToken, error) {
	if v == nil {
		return ValidateUserToken{}, exceptions.NewInvalidArgumentError("validator", "validator must not nil", nil)
	}
	return ValidateUserToken{
		validator: v,
	}, nil
}

//go:generate mockgen --source=./validate_user_token.go --destination=./mocks.go --package=validate_user_token .

type tokenValidator interface {
	ValidateToken(token string) (authenticate.UserAuthenticateData, error)
}

func (v ValidateUserToken) Handle(_ context.Context, token string) (authenticate.UserAuthenticateData, error) {
	if token == "" {
		return authenticate.UserAuthenticateData{}, exceptions.NewPreconditionError(exceptions.PreconditionReasonInvalidToken, exceptions.SubjectAuthentication, "token must not empty", nil)
	}
	u, err := v.validator.ValidateToken(token)
	if err != nil {
		return authenticate.UserAuthenticateData{}, exceptions.NewPreconditionError(
			exceptions.PreconditionReasonInvalidToken,
			exceptions.SubjectAuthentication,
			"token invalid",
			map[string]interface{}{
				"validate": err,
			},
		)
	}

	if u.ExpireAt < time.Now().UnixMilli() {
		return authenticate.UserAuthenticateData{}, exceptions.NewPreconditionError(
			exceptions.PreconditionReasonTokenExpired,
			exceptions.SubjectAuthentication,
			"token expired",
			map[string]interface{}{
				"expired_at": u.ExpireAt,
			},
		)
	}

	return u, nil
}

package generate_user_token

import (
	"context"
	"fmt"
	"time"

	"github.com/bqdanh/money_transfer/internal/entities/authenticate"
	"github.com/bqdanh/money_transfer/internal/entities/exceptions"
	"github.com/bqdanh/money_transfer/internal/entities/user"
)

type GenerateUserToken struct {
	generator     tokenGenerator
	tokenDuration time.Duration
}

func NewGenerateUserToken(g tokenGenerator, d time.Duration) (GenerateUserToken, error) {
	if g == nil {
		return GenerateUserToken{}, exceptions.NewInvalidArgumentError("generator", "generator must not nil", nil)
	}
	if d <= 0 {
		return GenerateUserToken{}, exceptions.NewInvalidArgumentError("duration", "duration must greater than 0", nil)
	}
	return GenerateUserToken{
		generator:     g,
		tokenDuration: d,
	}, nil
}

//go:generate mockgen --source=./generate_user_token.go --destination=./mocks.go --package=generate_user_token .
type tokenGenerator interface {
	GenerateToken(m authenticate.UserAuthenticateData) (string, error)
}

func (g GenerateUserToken) Handle(_ context.Context, u user.User) (string, error) {
	m := authenticate.UserAuthenticateData{
		UserID:       u.ID,
		Username:     u.UserName,
		UserFullName: u.FullName,
		CreatedAt:    time.Now().UnixMilli(),
		ExpireAt:     time.Now().Add(g.tokenDuration).UnixMilli(),
	}
	token, err := g.generator.GenerateToken(m)
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}
	return token, nil
}

package generate_user_token

import (
	"context"
	"fmt"
	"time"

	"github.com/bqdanh/money_transfer/internal/entities/authenticate"
	"github.com/bqdanh/money_transfer/internal/entities/user"
)

type GenerateUserToken struct {
	generator     tokenGenerator
	tokenDuration time.Duration
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

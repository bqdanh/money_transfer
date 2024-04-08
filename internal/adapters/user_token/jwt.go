package user_token

import (
	"encoding/json"
	"fmt"

	"github.com/bqdanh/money_transfer/internal/entities/authenticate"
	"github.com/bqdanh/money_transfer/pkg/jwt"
)

type JWTToken struct {
	generator *jwt.JwtRSAGenerator
	validator *jwt.JwtRSAValiator
}

type Config struct {
	PrivateKeyPath string `json:"private_key_path" mapstructure:"private_key_path"`
	PublicKeyPath  string `json:"public_key_path" mapstructure:"public_key_path"`
}

func NewJWTToken(cfg Config) (JWTToken, error) {
	validator, err := jwt.NewJwtRSAValidatorFromFile(cfg.PublicKeyPath)
	if err != nil {
		return JWTToken{}, fmt.Errorf("new jwt rsa validator from file: %w", err)
	}
	generator, err := jwt.NewJwtRSAGeneratorFromFile(cfg.PrivateKeyPath)
	if err != nil {
		return JWTToken{}, fmt.Errorf("new jwt rsa generator from file: %w", err)
	}

	return JWTToken{
		generator: generator,
		validator: validator,
	}, nil
}

func (j JWTToken) GenerateToken(m authenticate.UserAuthenticateData) (string, error) {
	bs, err := json.Marshal(m)
	if err != nil {
		return "", fmt.Errorf("marshal user authenticate data: %w", err)
	}
	tokenData := map[string]interface{}{}
	if err := json.Unmarshal(bs, &tokenData); err != nil {
		return "", fmt.Errorf("unmarshal user authenticate data: %w", err)
	}
	token, err := j.generator.GenerateToken(tokenData)
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}
	return token, nil
}

func (j JWTToken) ValidateToken(token string) (authenticate.UserAuthenticateData, error) {
	tokenData, err := j.validator.ValidateToken(token)
	if err != nil {
		return authenticate.UserAuthenticateData{}, fmt.Errorf("validate token: %w", err)
	}
	bs, err := json.Marshal(tokenData)
	if err != nil {
		return authenticate.UserAuthenticateData{}, fmt.Errorf("marshal token data: %w", err)
	}
	var u authenticate.UserAuthenticateData
	if err := json.Unmarshal(bs, &u); err != nil {
		return authenticate.UserAuthenticateData{}, fmt.Errorf("unmarshal token data: %w", err)
	}
	return u, nil
}

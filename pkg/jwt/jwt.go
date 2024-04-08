package jwt

import (
	"crypto/rsa"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

type JwtRSAGenerator struct {
	privateKey *rsa.PrivateKey
}

func NewJwtRSAGeneratorFromFile(privateKeyPath string) (*JwtRSAGenerator, error) {
	if len(privateKeyPath) == 0 {
		return nil, fmt.Errorf("private key path is empty")
	}
	privateKey, err := GetPrivateKey(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("get private key: %w", err)
	}

	return &JwtRSAGenerator{
		privateKey: privateKey,
	}, nil
}
func (j *JwtRSAGenerator) GenerateToken(m map[string]interface{}) (string, error) {
	return generateToken(m, j.privateKey)
}

type JwtRSAValiator struct {
	publicKey *rsa.PublicKey
}

func NewJwtRSAValidatorFromFile(publicKeyPath string) (*JwtRSAValiator, error) {
	if len(publicKeyPath) == 0 {
		return nil, fmt.Errorf("public key path is empty")
	}
	publicKey, err := GetPublicKey(publicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("get public key: %w", err)
	}

	return &JwtRSAValiator{
		publicKey: publicKey,
	}, nil
}

func (j *JwtRSAValiator) ValidateToken(t string) (map[string]interface{}, error) {
	return validateToken(t, j.publicKey)
}

func generateToken(m map[string]interface{}, privateKey *rsa.PrivateKey) (string, error) {
	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims(m)).SignedString(privateKey)
	if err != nil {
		return "", err
	}
	return token, nil
}

func validateToken(t string, publicKey *rsa.PublicKey) (map[string]interface{}, error) {
	claims := &jwt.MapClaims{}

	tkn, err := jwt.ParseWithClaims(t, claims, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})
	if err != nil || !tkn.Valid {
		return nil, err
	}

	return *claims, nil
}

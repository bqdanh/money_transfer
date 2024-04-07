package jwt

import (
	"crypto/rsa"

	"github.com/golang-jwt/jwt/v4"
)

type JwtRSAGenerator struct {
	privateKey *rsa.PrivateKey
}

func NewJwtRSAGeneratorFromFile(privateKeyPath string) *JwtRSAGenerator {
	var privateKey *rsa.PrivateKey

	if len(privateKeyPath) > 0 {
		privateKey = GetPrivateKey(privateKeyPath)
	}
	return &JwtRSAGenerator{
		privateKey: privateKey,
	}
}
func (j *JwtRSAGenerator) GenerateToken(m map[string]interface{}) (string, error) {
	return generateToken(m, j.privateKey)
}

type JwtRSAValiator struct {
	publicKey *rsa.PublicKey
}

func NewJwtRSAValiatorFromFile(publicKeyPath string) *JwtRSAValiator {
	var publicKey *rsa.PublicKey

	if len(publicKeyPath) > 0 {
		publicKey = GetPublicKey(publicKeyPath)
	}

	return &JwtRSAValiator{
		publicKey: publicKey,
	}
}

func (j *JwtRSAValiator) ValidateToken(t string) (map[string]interface{}, bool, error) {
	return validateToken(t, j.publicKey)
}

func generateToken(m map[string]interface{}, privateKey *rsa.PrivateKey) (string, error) {
	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims(m)).SignedString(privateKey)
	if err != nil {
		return "", err
	}
	return token, nil
}

func validateToken(t string, publicKey *rsa.PublicKey) (map[string]interface{}, bool, error) {
	claims := &jwt.MapClaims{}

	tkn, err := jwt.ParseWithClaims(t, claims, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})
	if err != nil || !tkn.Valid {
		return nil, false, err
	}

	return *claims, true, nil
}

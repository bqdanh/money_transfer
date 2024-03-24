package user

import (
	"fmt"

	"github.com/bqdanh/money_transfer/internal/entities/exceptions"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID int64
	//TODO: move UserName and Password to Auth entity
	UserName string
	Password string

	FullName string
	Phone    string
}

const (
	minimumLengthOfUserName = 6
	maximumLengthOfUserName = 16

	minimumLengthOfPassword = 1
	maximumLengthOfPassword = 30
)

func CreateUser(userName, password, fullName, phone string) (User, error) {
	if len(userName) < minimumLengthOfUserName || len(userName) > maximumLengthOfUserName {
		return User{}, exceptions.NewInvalidArgumentError(
			"UserName",
			fmt.Sprintf("length of user name must between %d and %d", minimumLengthOfUserName, maximumLengthOfUserName),
			map[string]interface{}{
				"min_len": minimumLengthOfUserName,
				"max_len": maximumLengthOfUserName,
			},
		)
	}
	if len(password) < minimumLengthOfPassword || len(password) > maximumLengthOfPassword {
		return User{}, exceptions.NewInvalidArgumentError(
			"Password",
			fmt.Sprintf("length of password must between %d and %d", minimumLengthOfUserName, maximumLengthOfUserName),
			map[string]interface{}{
				"min_len": minimumLengthOfUserName,
				"max_len": maximumLengthOfUserName,
			},
		)
	}
	hpw, err := hashPassword(password)
	if err != nil {
		return User{}, fmt.Errorf("hash password: %w", err)
	}

	return User{
		ID:       0,
		UserName: userName,
		Password: hpw,
		FullName: fullName,
		Phone:    phone,
	}, nil

}

func hashPassword(pw string) (string, error) {
	bs, err := bcrypt.GenerateFromPassword([]byte(pw), 10)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}
	return string(bs), nil
}

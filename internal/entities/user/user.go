package user

import (
	"fmt"

	"github.com/banhquocdanh/money_transfer/internal/entities/exceptions"
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

}

func hashPassword(pw string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
}

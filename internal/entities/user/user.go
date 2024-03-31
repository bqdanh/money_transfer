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
	if len(fullName) == 0 {
		return User{}, exceptions.NewInvalidArgumentError(
			"FullName",
			"full name must not empty",
			nil,
		)
	}
	if len(phone) > 0 {
		for _, c := range phone {
			if c == '+' {
				continue
			}
			if c < '0' || c > '9' {
				return User{}, exceptions.NewInvalidArgumentError(
					"Phone",
					"phone must be number",
					map[string]interface{}{
						"got": phone,
					},
				)
			}
		}
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

const (
	SubjectUser                        = "user"
	PreconditionTypeCannotChangeUserID = exceptions.PreconditionType("cannot-change-user-id")
)

func (u User) WithID(id int64) (User, error) {
	if u.ID > 0 {
		return u, exceptions.NewPreconditionError(
			PreconditionTypeCannotChangeUserID,
			SubjectUser,
			fmt.Sprintf("cannot change user id"), map[string]interface{}{
				"current_id": u.ID,
			},
		)
	}
	u.ID = id
	return u, nil
}

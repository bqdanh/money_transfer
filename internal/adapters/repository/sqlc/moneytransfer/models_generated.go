// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0

package moneytransfer

import (
	"time"
)

type User struct {
	// is identify user with primary key auto increment
	ID int64 `json:"id"`
	// user name for login, unique
	UserName string `json:"user_name"`
	// password hashed
	Password string `json:"password"`
	// user full name
	FullName string `json:"full_name"`
	// user phone number
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0

package moneytransfer

import (
	"context"
	"database/sql"
)

type Querier interface {
	InsertUser(ctx context.Context, arg *InsertUserParams) (sql.Result, error)
}

var _ Querier = (*Queries)(nil)

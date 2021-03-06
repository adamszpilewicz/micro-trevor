// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0

package db

import (
	"context"
)

type Querier interface {
	Delete(ctx context.Context, id int32) error
	DeleteById(ctx context.Context, id int32) error
	GetAll(ctx context.Context) ([]User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	GetOne(ctx context.Context, id int32) (User, error)
	Insert(ctx context.Context, arg InsertParams) error
	ResetPassword(ctx context.Context, arg ResetPasswordParams) error
	Update(ctx context.Context, arg UpdateParams) error
}

var _ Querier = (*Queries)(nil)

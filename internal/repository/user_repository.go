package repository

import (
	"apiserver/internal/domain"
	"context"
	"errors"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByID(ctx context.Context, id int) (*domain.User, error)
	GetAll(ctx context.Context) ([]*domain.User, error)
	Delete(ctx context.Context, id int) error
}

var ErrNotFound = errors.New("user not found at database")
var ErrNoRowsAffected = errors.New("no rows affected by query")

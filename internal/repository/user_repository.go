package repository

import (
	"context"
	"errors"

	"apiserver/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
	FindByID(ctx context.Context, id int) (*domain.User, error)
	GetAll(ctx context.Context) ([]*domain.User, error)
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, user *domain.User) error
}

var ErrNotFound = errors.New("user not found at database")
var ErrNoRowsAffected = errors.New("no rows affected by query")

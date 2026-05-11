package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"apiserver/internal/domain"
	"apiserver/internal/repository"
)

type UserService interface {
	CreateUserFromBusinessRule(ctx context.Context, userCreateInput CreateUserInput) error
	GetAllUserFromAnyBusinessRule(ctx context.Context) ([]*domain.User, error)
	FindByID(ctx context.Context, id int) (*domain.User, error)
	DeleteUser(ctx context.Context, id int) error
}

type userService struct {
	repository repository.UserRepository
}

type CreateUserInput struct {
	Name string
}

func NewUserService(r repository.UserRepository) UserService {
	return &userService{
		repository: r,
	}
}

func (u *userService) GetAllUserFromAnyBusinessRule(ctx context.Context) ([]*domain.User, error) {
	return u.repository.GetAll(ctx)
}

func (u *userService) CreateUserFromBusinessRule(ctx context.Context, input CreateUserInput) error {

	user := &domain.User{
		Name: input.Name,
		UUID: uuid.New(),
	}

	if err := u.repository.Create(ctx, user); err != nil {
		return err
	}

	return nil
}

func (u *userService) DeleteUser(ctx context.Context, id int) error {
	err := u.repository.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNoRowsAffected) {
			return fmt.Errorf("delete by user: %w", domain.ErrUserNotFound)
		}
		return fmt.Errorf("delete by user: %w", err)
	}
	return nil
}

func (u *userService) FindByID(ctx context.Context, id int) (*domain.User, error) {
	user, err := u.repository.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, fmt.Errorf("find user by id: %w", domain.ErrUserNotFound)
		}
		return nil, fmt.Errorf("find user by id: %w", err)
	}
	return user, nil
}

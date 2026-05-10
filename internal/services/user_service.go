package services

import (
	"apiserver/internal/domain"
	"apiserver/internal/repository"
	"context"

	"github.com/google/uuid"
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
	return u.repository.Delete(ctx, id)
}

func (u *userService) FindByID(ctx context.Context, id int) (*domain.User, error) {
	return u.repository.FindByID(ctx, id)
}

package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"

	"apiserver/internal/domain"
	"apiserver/internal/repository"
)

type UserService interface {
	CreateUserFromBusinessRule(ctx context.Context, userCreateInput UserCreateInput) (*domain.User, error)
	GetAllUserFromAnyBusinessRule(ctx context.Context) ([]*domain.User, error)
	FindByID(ctx context.Context, id int) (*domain.User, error)
	DeleteUser(ctx context.Context, id int) error
	UpdateUser(ctx context.Context, userUpdateInput UserUpdateInput) error
}

type userService struct {
	repository repository.UserRepository
}

type UserUpdateInput struct {
	ID        int
	Name      string
	UpdatedAt time.Time
}

type UserCreateInput struct {
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

func (u *userService) CreateUserFromBusinessRule(ctx context.Context, input UserCreateInput) (*domain.User, error) {

	ctx, span := otel.Tracer("app/user-service").Start(ctx, "user.service.create_user")
	defer span.End()

	span.SetAttributes(
		attribute.String("user.operation", "create"),
		attribute.String("user.name", input.Name),
	)

	user := &domain.User{
		Name: input.Name,
		UUID: uuid.New(),
	}

	user, err := u.repository.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("service error on create user > %w", err)
	}
	span.SetStatus(codes.Ok, "user created")
	return user, nil
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

func (u *userService) UpdateUser(ctx context.Context, input UserUpdateInput) error {
	user, err := u.FindByID(ctx, input.ID)
	if err != nil {
		return fmt.Errorf("user not found for update > %w", err)
	}
	now := time.Now().UTC()
	user.UpdatedAt = &now
	user.Name = input.Name
	if err := u.repository.Update(ctx, user); err != nil {
		return err
	}

	return nil

}

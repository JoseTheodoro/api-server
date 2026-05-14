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

	"apiserver/internal/config/observability/logging"
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
	FirstName *string
	LastName  *string
	Email     *string
	Password  *string
	DateBirth *string
	Genre     *domain.Genre
	UpdatedAt *time.Time
}

type UserCreateInput struct {
	FirstName string
	LastName  string
	Email     string
	DateBirth string
	Genre     domain.Genre
	Password  string
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
	log := logging.NewLogger().LoggerFromContext(ctx)
	defer span.End()

	span.SetAttributes(
		attribute.String("user.operation", "create"),
		attribute.String("user.name", input.FirstName),
	)

	user := &domain.User{
		UUID:      uuid.New(),
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Password:  input.Password,
		DateBirth: input.DateBirth,
		Genre:     input.Genre,
	}

	user, err := u.repository.Create(ctx, user)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "create user failed")
		return nil, fmt.Errorf("service error on create user > %w", err)
	}
	span.SetStatus(codes.Ok, "user created")
	log.Info("user created successful", "user", user)
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

	user = u.fillUserToUpdate(user, &input)

	if err := u.repository.Update(ctx, user); err != nil {
		return err
	}

	return nil

}

func (u *userService) fillUserToUpdate(user *domain.User, i *UserUpdateInput) *domain.User {

	if i.FirstName != nil {
		user.FirstName = *i.FirstName
	}

	if i.LastName != nil {
		user.LastName = *i.LastName
	}

	if i.Email != nil {
		user.Email = *i.Email
	}

	if i.Password != nil {
		user.Password = *i.Password
	}

	if i.Genre != nil {
		user.Genre = *i.Genre
	}

	if i.DateBirth != nil {
		user.DateBirth = *i.DateBirth
	}
	now := time.Now().UTC()
	user.UpdatedAt = &now
	return user
}

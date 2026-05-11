package services

import (
	"apiserver/internal/domain"
	"apiserver/internal/repository"
	"context"
	"errors"
	"testing"
)

type fakeUserRepo struct {
	DeleteFakeFn   func(ctx context.Context, id int) error
	FindByIDFakeFn func(ctx context.Context, id int) (*domain.User, error)
}

func (fr *fakeUserRepo) Delete(ctx context.Context, id int) error {
	return fr.DeleteFakeFn(ctx, id)
}

func (fr *fakeUserRepo) FindByID(ctx context.Context, id int) (*domain.User, error) {
	return fr.FindByIDFakeFn(ctx, id)
}

func (fr *fakeUserRepo) Create(ctx context.Context, user *domain.User) error {
	return nil
}

func (fr *fakeUserRepo) GetAll(ctx context.Context) ([]*domain.User, error) {
	return nil, nil
}

func TestServiceDeleteUser(t *testing.T) {
	var ctx context.Context = context.Background()

	t.Run("success", func(t *testing.T) {
		var errExpexted error
		userFakeRepo := &fakeUserRepo{
			DeleteFakeFn: func(ctx context.Context, id int) error {
				return nil
			},
			FindByIDFakeFn: func(ctx context.Context, id int) (*domain.User, error) {
				return nil, nil
			},
		}

		userService := NewUserService(userFakeRepo)
		errGot := userService.DeleteUser(ctx, 99)
		if errExpexted != errGot {
			t.Fatalf("expected nil error, got %v", errGot)
		}

	})

	t.Run("user not found", func(t *testing.T) {

		userFakeRepo := &fakeUserRepo{
			DeleteFakeFn: func(ctx context.Context, id int) error {
				return repository.ErrNoRowsAffected
			},
		}

		userService := NewUserService(userFakeRepo)
		errGot := userService.DeleteUser(ctx, 10)

		if !errors.Is(errGot, domain.ErrUserNotFound) {
			t.Fatalf("expected domain.ErrUserNotFoud, got: %v", errGot)
		}

	})

	t.Run("no specified error", func(t *testing.T) {
		userFakeRepo := &fakeUserRepo{
			DeleteFakeFn: func(ctx context.Context, id int) error {
				return errors.New("no specified error")
			},
		}

		userService := NewUserService(userFakeRepo)
		errGot := userService.DeleteUser(ctx, 1)
		if errGot == nil {
			t.Fatalf("expected error instance, got: error nil")
		}
	})

}

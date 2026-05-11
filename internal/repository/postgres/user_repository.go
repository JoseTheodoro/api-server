package postgres

import (
	"apiserver/internal/domain"
	"apiserver/internal/repository"
	"context"
	"database/sql"
	"errors"
	"log/slog"
)

type UserRepositoryPostgres struct {
	db *sql.DB
}

func NewUserRepositoryPostgres(db *sql.DB) *UserRepositoryPostgres {
	return &UserRepositoryPostgres{db: db}
}

func (r *UserRepositoryPostgres) Create(ctx context.Context, user *domain.User) error {

	_, err := r.db.ExecContext(ctx, "INSERT INTO users (name, uuid) VALUES ($1, $2)", user.Name, user.UUID.String())
	if err != nil {
		return err
	}

	slog.Info("user created successful", "user", user)
	return nil

}

func (r *UserRepositoryPostgres) GetAll(ctx context.Context) ([]*domain.User, error) {

	var users []*domain.User
	rows, err := r.db.QueryContext(ctx, "SELECT id, uuid, name, created_at FROM users")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		u := &domain.User{}
		if err := rows.Scan(&u.ID, &u.UUID, &u.Name, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepositoryPostgres) FindByID(ctx context.Context, id int) (*domain.User, error) {
	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, "SELECT id, uuid, name, created_at FROM users WHERE id = $1", id).Scan(&user.ID, &user.UUID, &user.Name, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	return user, nil
}

func (r *UserRepositoryPostgres) Delete(ctx context.Context, id int) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM users where id = $1", id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return repository.ErrNoRowsAffected
	}

	return nil

}

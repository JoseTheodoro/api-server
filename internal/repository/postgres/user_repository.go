package postgres

import (
	"context"
	"database/sql"
	"errors"

	"apiserver/internal/domain"
	"apiserver/internal/repository"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type UserRepositoryPostgres struct {
	db *sql.DB
}

func NewUserRepositoryPostgres(db *sql.DB) *UserRepositoryPostgres {
	return &UserRepositoryPostgres{db: db}
}

func (r *UserRepositoryPostgres) Create(ctx context.Context, user *domain.User) (*domain.User, error) {

	ctx, span := otel.Tracer("app/user-repository-postgres").Start(ctx, "user.repo.insert_postgres")
	defer span.End()

	span.SetAttributes(
		attribute.String("user.uuid", user.UUID.String()),
		attribute.String("db.operation", "INSERT"),
		attribute.String("db.sql.table", "users"),
		attribute.String("db.system", "postgres"),
	)

	err := r.db.QueryRowContext(ctx, "INSERT INTO users (uuid, first_name, last_name, email, genre, date_birth, password) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, uuid, first_name, last_name, email, password, date_birth, created_at, updated_at", user.UUID.String(), user.FirstName, user.LastName, user.Email, user.Genre, user.DateBirth, user.Password).
		Scan(&user.ID, &user.UUID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.DateBirth, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "insert failed")
		return nil, err
	}
	span.SetStatus(codes.Ok, "insert OK")
	return user, nil

}

func (r *UserRepositoryPostgres) GetAll(ctx context.Context) ([]*domain.User, error) {

	var users []*domain.User
	rows, err := r.db.QueryContext(ctx, "SELECT id, uuid, first_name, last_name, email, password, date_birth, genre, created_at, updated_at FROM users")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		u := &domain.User{}
		if err := rows.Scan(&u.ID, &u.UUID, &u.FirstName, &u.Email, &u.Password, &u.DateBirth, &u.Genre, &u.CreatedAt, &u.UpdatedAt); err != nil {
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
	err := r.db.QueryRowContext(ctx, "SELECT id, uuid, first_name, last_name, email, date_birth, password, genre, created_at, updated_at FROM users WHERE id = $1", id).Scan(&user.ID, &user.UUID, &user.FirstName, &user.LastName, &user.Email, &user.DateBirth, &user.Password, &user.Genre, &user.CreatedAt, &user.UpdatedAt)
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

func (r *UserRepositoryPostgres) Update(ctx context.Context, user *domain.User) error {
	_, err := r.db.ExecContext(ctx, "UPDATE users SET first_name = $2, updated_at = $3, last_name = $4, email = $5, password = $6, date_birth = $7, genre = $8 WHERE id = $1", user.ID, user.FirstName, user.UpdatedAt, user.LastName, user.Email, user.Password, user.DateBirth, user.Genre)
	if err != nil {
		return err
	}

	return nil
}

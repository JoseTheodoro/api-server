package postgres

import (
	"apiserver/internal/domain"
	"apiserver/internal/repository/postgres/queries"
	"context"
	"database/sql"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type UserRepositoryPostgres struct {
	db *sql.DB
	qq *queries.Queries
}

func NewUserRepositoryPostgres(db *sql.DB) *UserRepositoryPostgres {
	return &UserRepositoryPostgres{db: db, qq: queries.New(db)}
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

	dob, err := time.Parse("2006-01-01", user.DateBirth)
	if err != nil {
		return nil, err
	}

	params := queries.CreateUserParams{
		Uuid:      user.UUID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		DateBirth: dob,
		Password:  user.Password,
		Genre:     queries.Genres(user.Genre),
	}
	created, err := r.qq.CreateUser(ctx, params)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "insert failed")
		return nil, err
	}
	span.SetStatus(codes.Ok, "insert OK")

	user = &domain.User{
		ID:        int(created.ID),
		UUID:      created.Uuid,
		FirstName: created.FirstName,
		LastName:  created.LastName,
		Email:     created.Email,
		Password:  created.Password,
		Genre:     domain.Genre(created.Genre),
		CreatedAt: created.CreatedAt,
		UpdatedAt: &created.UpdatedAt.Time,
		DateBirth: created.DateBirth.Format("2006-01-02"),
	}

	return user, nil

}

func (r *UserRepositoryPostgres) GetAll(ctx context.Context) ([]*domain.User, error) {

	rows, err := r.qq.ListUsers(ctx)
	if err != nil {
		return nil, err
	}
	users := make([]*domain.User, 0, len(rows))

	for _, u := range rows {

		usr := &domain.User{
			ID:        int(u.ID),
			UUID:      u.Uuid,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Email:     u.Email,
			Password:  u.Password,
			Genre:     domain.Genre(u.Genre),
			CreatedAt: u.CreatedAt,
			UpdatedAt: &u.UpdatedAt.Time,
			DateBirth: u.DateBirth.Format("2006-01-02"),
		}
		users = append(users, usr)
	}

	return users, nil

}

func (r *UserRepositoryPostgres) FindByID(ctx context.Context, id int) (*domain.User, error) {
	row, err := r.qq.GetUser(ctx, int64(id))
	if err != nil {
		return nil, err
	}

	var t *time.Time
	if row.UpdatedAt.Valid {
		t = &row.UpdatedAt.Time
	}

	user := &domain.User{
		ID:        int(row.ID),
		UUID:      row.Uuid,
		FirstName: row.FirstName,
		LastName:  row.LastName,
		Email:     row.Email,
		Password:  row.Password,
		Genre:     domain.Genre(row.Genre),
		DateBirth: row.DateBirth.Format("2006-01-02"),
		CreatedAt: row.CreatedAt,
		UpdatedAt: t,
	}

	return user, err

}

func (r *UserRepositoryPostgres) Delete(ctx context.Context, id int) error {
	err := r.qq.DeleteUser(ctx, int64(id))
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepositoryPostgres) Update(ctx context.Context, user *domain.User) error {

	dob, err := time.Parse("2006-01-02", user.DateBirth)
	if err != nil {
		return err
	}

	userParams := queries.UpdateUserParams{
		ID:        int64(user.ID),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		DateBirth: dob,
		Password:  user.Password,
		Genre:     queries.Genres(user.Genre),
		UpdatedAt: sql.NullTime{Time: time.Now().UTC(), Valid: true},
	}
	_, err = r.qq.UpdateUser(ctx, userParams)
	if err != nil {
		return err
	}
	return nil
}

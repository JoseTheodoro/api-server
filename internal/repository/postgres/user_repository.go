package postgres

import (
	"apiserver/internal/domain"
	"apiserver/internal/repository/postgres/queries"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type UserRepositoryPostgres struct {
	db *pgxpool.Pool
	qq *queries.Queries
}

func NewUserRepositoryPostgres(db *pgxpool.Pool) *UserRepositoryPostgres {
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

	params := queries.CreateUserParams{
		Uuid:      user.UUID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		DateBirth: user.DateBirth,
		Password:  user.Password,
		Genre:     user.Genre,
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
		Genre:     created.Genre,
		CreatedAt: created.CreatedAt,
		UpdatedAt: created.UpdatedAt,
		DeletedAt: created.DeletedAt,
		DateBirth: created.DateBirth,
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
			Genre:     u.Genre,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
			DateBirth: u.DateBirth,
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

	user := &domain.User{
		ID:        int(row.ID),
		UUID:      row.Uuid,
		FirstName: row.FirstName,
		LastName:  row.LastName,
		Email:     row.Email,
		Password:  row.Password,
		Genre:     row.Genre,
		DateBirth: row.DateBirth,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
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

	userParams := queries.UpdateUserParams{
		ID:        int64(user.ID),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		DateBirth: user.DateBirth,
		Password:  user.Password,
		Genre:     user.Genre,
		UpdatedAt: user.UpdatedAt,
	}
	_, err := r.qq.UpdateUser(ctx, userParams)
	if err != nil {
		return err
	}
	return nil
}

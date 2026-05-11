package handlers_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5"
	"github.com/testcontainers/testcontainers-go"
	pg "github.com/testcontainers/testcontainers-go/modules/postgres"

	"apiserver/internal/domain"
	"apiserver/internal/http/handlers"
	"apiserver/internal/repository/postgres"
	"apiserver/internal/services"
)

var testCtx context.Context
var testDSN string
var testContainer testcontainers.Container
var pgContainer *pg.PostgresContainer

func TestMain(M *testing.M) {
	setupGlobalState()
	M.Run()
	teardownGlobalState()
}

func setupGlobalState() {
	createPostgresContainer()
	runMigrations(testDSN)
}

func teardownGlobalState() {
	defer func() {
		testcontainers.TerminateContainer(pgContainer)
	}()
}

func TestCreateUser_ValidPayload_ReturnsCreated(T *testing.T) {
	c := postgres.NewConnection(testDSN)

	db, err := c.Connect(testCtx)
	if err != nil {
		T.Fatal("error on connection db", err)
	}
	defer db.Close()

	r := postgres.NewUserRepositoryPostgres(db)
	s := services.NewUserService(r)
	h := handlers.NewUserHandler(s)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /users", h.CreateUser)

	// monta request
	payload := bytes.NewReader([]byte(`{"name": "Testing Name"}`))
	request := httptest.NewRequest("POST", "/users", payload)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, request)

	if rr.Code != http.StatusCreated {
		T.Error("failed create user on endpoint", "statuscode:", rr.Code)
	}

	var count int
	if err := db.QueryRowContext(testCtx, "select count(*) from users where name = $1", "Testing Name").Scan(&count); err != nil {
		T.Fatal("error on executing query", err)
	}

	if count == 0 {
		T.Error("User doesnt exists on table users")
	}

}

func TestCreateUser_InvalidPayload_ReturnsBadRequest(T *testing.T) {
	conn := postgres.NewConnection(testDSN)
	db, err := conn.Connect(testCtx)
	if err != nil {
		log.Fatal("error on connect on postgres container", err)
	}
	r := postgres.NewUserRepositoryPostgres(db)
	s := services.NewUserService(r)
	h := handlers.NewUserHandler(s)

	payload := bytes.NewReader([]byte(`{}`))
	mux := http.NewServeMux()
	mux.HandleFunc("POST /users", h.CreateUser)

	request := httptest.NewRequest("POST", "/users", payload)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, request)

	if w.Code != http.StatusBadRequest {
		T.Error("status code is", w.Code, "expected: ", http.StatusBadRequest)
	}

}

func TestUpdateUser_UserExists_ReturnsNoContent(t *testing.T) {
	method := "PUT"
	pattern := fmt.Sprintf("%s /users/{id}", method)
	conn := postgres.NewConnection(testDSN)
	db, err := conn.Connect(testCtx)
	if err != nil {
		t.Fatal("error to connect postgress", "err", err)
	}

	bigRichard := domain.User{
		UUID: uuid.New(),
		Name: "Big Richard",
	}
	var userIDExists int
	err = db.QueryRowContext(testCtx, "INSERT INTO users (uuid, name) VALUES ($1, $2) RETURNING id", bigRichard.UUID.String(), bigRichard.Name).Scan(&userIDExists)
	if err != nil {
		t.Fatal("error on insert user to update user")
	}

	ur := postgres.NewUserRepositoryPostgres(db)
	us := services.NewUserService(ur)
	uh := handlers.NewUserHandler(us)

	mux := http.NewServeMux()
	mux.HandleFunc(pattern, uh.UpdateUser)

	payload := bytes.NewReader([]byte(`{"name": "User Updated"}`))
	req := httptest.NewRequestWithContext(testCtx, method, fmt.Sprintf("/users/%d", userIDExists), payload)
	writer := httptest.NewRecorder()

	mux.ServeHTTP(writer, req)

	if writer.Code != http.StatusNoContent {
		t.Error("expected status 204, got:", writer.Code)
	}
	userUpdated := domain.User{}
	err = db.QueryRowContext(testCtx, "SELECT id, name, updated_at FROM users where id = $1", userIDExists).Scan(&userUpdated.ID, &userUpdated.Name, &userUpdated.UpdatedAt)
	if err != nil {
		t.Fatal("error executing query for find user updated", err)
	}

	if userUpdated.Name != "User Updated" {
		t.Error("exepected User Update, got:", userUpdated.Name)
	}

}

func createPostgresContainer() {
	testCtx = context.Background()

	pgContainer, err := pg.Run(testCtx,
		"postgres:16-alpine",
		pg.WithDatabase("test"),
		pg.WithUsername("test"),
		pg.WithPassword("test"),
		pg.BasicWaitStrategies(),
	)

	if err != nil {
		log.Fatal("failed to start container:", err)
	}

	testDSN, err = pgContainer.ConnectionString(testCtx, "sslmode=disable")
	if err != nil {
		log.Fatal("error on ConnectionString > err:", err)
	}
}

func runMigrations(dsn string) {
	rel := "../../../internal/repository/postgres/migrations"
	abs, err := filepath.Abs(rel)
	if err != nil {
		log.Fatalf("erro ao resolver path de migrations: %v", err)
	}
	pathSource := fmt.Sprintf("file://%s", abs)

	m, err := migrate.New(pathSource, dsn)
	if err != nil {
		log.Fatal("error on setting migration", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal("error on executing migration", err)
	}

}

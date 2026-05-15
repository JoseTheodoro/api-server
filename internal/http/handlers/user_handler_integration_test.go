package handlers_test

import (
	"apiserver/internal/repository/postgres/queries"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"apiserver/internal/domain"
	"apiserver/internal/http/handlers"
	"apiserver/internal/repository/postgres"
	"apiserver/internal/services"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	pg "github.com/testcontainers/testcontainers-go/modules/postgres"
)

var (
	conn                   *postgres.Connection
	db                     *pgxpool.Pool
	q                      *queries.Queries
	testCtx                context.Context
	err                    error
	testDSN                string
	pgContainer            *pg.PostgresContainer
	repositoryUserPostgres *postgres.UserRepositoryPostgres
	serviceUser            services.UserService
	handleUser             *handlers.UserHandler
	mux                    *http.ServeMux
)

func TestMain(M *testing.M) {
	setupGlobalState()
	M.Run()
	teardownGlobalState()
}

func setupGlobalState() {
	createPostgresContainer()
	db, err = postgres.NewConnection(testDSN).Connect(testCtx)
	if err != nil {
		log.Fatal("error to connect postgres container")
	}
	q = queries.New(db)
	runMigrations(testDSN)
	repositoryUserPostgres = postgres.NewUserRepositoryPostgres(db)
	serviceUser = services.NewUserService(repositoryUserPostgres)
	handleUser = handlers.NewUserHandler(serviceUser)
	mux = http.NewServeMux()
}

func teardownGlobalState() {
	defer db.Close()
	defer func() {
		err := pgContainer.Terminate(testCtx)
		if err != nil {
			log.Fatal(err)
		}
	}()
}

func TestCreateUser_ValidPayload_ReturnsCreated(T *testing.T) {
	mux.HandleFunc("POST /users", handleUser.CreateUser)

	// monta request
	newUser := `{
    "first_name":"Sebastião",
    "last_name":"b",
    "email":"sebastiao@api.com",
    "password":"a",
    "genre": "male",
    "date_birth":"1999-09-11"
	}`
	payload := bytes.NewReader([]byte(newUser))
	request := httptest.NewRequest("POST", "/users", payload)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, request)

	if w.Code != http.StatusOK {
		T.Error("failed create user on endpoint", "statuscode:", w.Code)
	}
	var u domain.User
	if err := json.NewDecoder(w.Body).Decode(&u); err != nil {
		T.Error("failed decoder response user create", err)
	}

	userFound, err := q.GetUser(testCtx, int64(u.ID))
	if err != nil {
		T.Error("error get user on DB", err)
	}

	if userFound.Email != "sebastiao@api.com" {
		T.Errorf("expected user email: sebastiao@api.com, got: %v", userFound.Email)
	}

}

func TestCreateUser_InvalidPayload_ReturnsBadRequest(T *testing.T) {

	payload := bytes.NewReader([]byte(`{}`))
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

	bigRichard := queries.CreateUserParams{
		Uuid:      uuid.New(),
		FirstName: "Big",
		LastName:  "Richard",
		Email:     "big@big.com",
		DateBirth: time.Date(1999, time.April, 23, 10, 10, 10, 10, time.UTC),
		Password:  "123456",
		Genre:     domain.MALE,
	}

	row, pa := q.CreateUser(testCtx, bigRichard)
	if pa != nil {
		t.Fatalf("error to create user on db > %v ", err)
	}

	userCreated := &domain.User{
		ID:    int(row.ID),
		UUID:  row.Uuid,
		Email: row.Email,
	}

	fmt.Printf("userCreated=%v", userCreated)

	mux.HandleFunc(pattern, handleUser.UpdateUser)

	payload := bytes.NewReader([]byte(`{"first_name": "Small"}`))
	req := httptest.NewRequestWithContext(testCtx, method, fmt.Sprintf("/users/%d", userCreated.ID), payload)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Error("expected status 204, got:", w.Code)
	}

	userUpdated := domain.User{}
	userFoundRow, err := q.GetUser(testCtx, int64(userCreated.ID))
	if err != nil {
		t.Errorf("error on get user on db > %v", err)
	}

	userUpdated.FirstName = userFoundRow.FirstName

	if userFoundRow.FirstName != "Small" {
		t.Error("exepected User Small, got:", userUpdated.FirstName)
	}

	if userFoundRow.Email != userCreated.Email {
		t.Errorf("expected email: %s, got: %s", userCreated.Email, userFoundRow.Email)
	}

}

func createPostgresContainer() {
	testCtx = context.Background()
	var err error

	pgContainer, err = pg.Run(testCtx,
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

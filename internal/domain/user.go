package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Genre string

const (
	MALE   Genre = "male"
	FEMALE Genre = "female"
)

type User struct {
	ID        int        `json:"id"`
	UUID      uuid.UUID  `json:"uuid"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	Email     string     `json:"email"`
	DateBirth time.Time  `json:"date_birth"`
	Genre     Genre      `json:"genre"`
	Password  string     `json:"password"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

var ErrUserNotFound = errors.New("user not found")

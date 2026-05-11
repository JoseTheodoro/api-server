package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        int        `json:"id"`
	UUID      uuid.UUID  `json:"uuid"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

var ErrUserNotFound = errors.New("user not found")

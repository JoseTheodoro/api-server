package domain

import (
	"errors"

	"github.com/google/uuid"
)

type User struct {
	ID        int       `json:"id"`
	UUID      uuid.UUID `json:"uuid"`
	Name      string    `json:"name"`
	CreatedAt string    `json:"created_at"`
}

var ErrUserNotFound = errors.New("user not found")

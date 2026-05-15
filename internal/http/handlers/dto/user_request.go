package dto

import (
	"apiserver/internal/domain"
)

type UserCreateRequest struct {
	FirstName string       `json:"first_name"`
	LastName  string       `json:"last_name"`
	Email     string       `json:"email"`
	DateBirth string       `json:"date_birth"`
	Genre     domain.Genre `json:"genre"`
	Password  string       `json:"password"`
}

func (u *UserCreateRequest) Validate() bool {
	if u.FirstName == "" {
		return false
	}
	return true
}

type UserUpdateRequest struct {
	ID        int           `json:"id"`
	FirstName *string       `json:"first_name"`
	LastName  *string       `json:"last_name"`
	Email     *string       `json:"email"`
	Password  *string       `json:"password"`
	DateBirth *string       `json:"date_birth"`
	Genre     *domain.Genre `json:"genre"`
}

func (up *UserUpdateRequest) Validate() bool {
	return true
}

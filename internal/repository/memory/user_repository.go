package memory

import "apiserver/internal/domain"

type UserRepositoryMemory struct{}

func NewUserMemoryRepository() *UserRepositoryMemory {
	return &UserRepositoryMemory{}
}

func (m *UserRepositoryMemory) Get() *domain.User {

	return &domain.User{
		ID:   1,
		Name: "Theodoro",
	}

}

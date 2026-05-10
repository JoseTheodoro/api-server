package file

import "apiserver/internal/domain"

type UserRespositoryFile struct{}

func NewUserRepositoryFile() *UserRespositoryFile {
	return &UserRespositoryFile{}
}

func (f *UserRespositoryFile) Get() *domain.User {
	return &domain.User{
		ID:   1,
		Name: "Theoodoro File",
	}
}

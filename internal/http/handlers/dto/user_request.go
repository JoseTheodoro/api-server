package dto

type UserCreateRequest struct {
	Name string `json:"name"`
}

func (u *UserCreateRequest) Validate() bool {
	if u.Name == "" {
		return false
	}
	return true
}

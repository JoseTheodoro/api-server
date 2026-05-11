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

type UserUpdateRequest struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (up *UserUpdateRequest) Validate() bool {
	if up.ID <= 0 || up.Name == "" {
		return false
	}

	return true
}

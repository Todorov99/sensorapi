package client

type UserDto struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

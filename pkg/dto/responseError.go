package dto

type ResponseError struct {
	ErrMessage string      `json:"error,omitempty"`
	Entity     interface{} `json:"entity,omitempty"`
}

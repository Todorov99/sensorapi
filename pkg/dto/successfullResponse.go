package dto

type SuccessfulResponse struct {
	Message string      `json:"message,omitempty"`
	Model   interface{} `json:"model,omitempty"`
}

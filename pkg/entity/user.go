package entity

type User struct {
	ID        int    `json:"id" mapstructure:"id"`
	UserName  string `json:"user_name" mapstructure:"user_name"`
	FirstName string `json:"first_name" mapstructure:"first_name"`
	LastName  string `json:"last_name" mapstructure:"last_name"`
	Email     string `json:"email" mapstructure:"email"`
	Password  string `json:"pass" mapstructure:"pass"`
}

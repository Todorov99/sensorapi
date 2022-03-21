package dto

type Register struct {
	UserName  string `json:"username" mapstructure:"user_name"`
	FirstName string `json:"firstname" mapstructure:"first_name"`
	LastName  string `json:"lastname" mapstructure:"last_name"`
	Email     string `json:"email" mapstructure:"email"`
	Password  string `json:"password" mapstructure:"pass"`
}

type Login struct {
	UserName string `json:"username" mapstructure:"first_name"`
	Password string `json:"password" mapstructure:"pass"`
}

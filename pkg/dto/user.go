package dto

type Register struct {
	UserName  string `json:"username" bson:"firstname"`
	FirstName string `json:"firstname" bson:"firstname"`
	LastName  string `json:"lastname" bson:"lastname"`
	Email     string `json:"email" bson:"email"`
	Password  string `json:"password" bson:"password"`
}

type Login struct {
	UserName string `json:"username" bson:"firstname"`
	Password string `json:"password" bson:"password"`
}

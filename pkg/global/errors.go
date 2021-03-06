package global

import (
	"errors"
)

var (
	ErrorObjectNotFound             = errors.New("object not found")
	ErrorUserWithUsernameNotExist   = errors.New("user does not exist")
	ErrorDeviceWithNameAlreadyExist = errors.New("device already exists")
)

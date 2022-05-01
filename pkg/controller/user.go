package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Todorov99/sensorapi/pkg/dto"
	"github.com/Todorov99/sensorapi/pkg/service"
)

type userController struct {
	userService service.UserService
}

func NewUserController() *userController {
	return &userController{
		userService: service.NewUserService(),
	}
}

func (l *userController) Login(w http.ResponseWriter, r *http.Request) {
	defer func() {
		r.Body.Close()
	}()

	username, password, ok := r.BasicAuth()
	if !ok {
		response(w, "Failed log in ", fmt.Errorf("invalid user: %s password", username), nil, http.StatusForbidden)
		return
	}

	login := dto.Login{
		UserName: username,
		Password: password,
	}

	token, err := l.userService.Login(r.Context(), login)
	if err != nil {
		login.Password = "****"
		response(w, "Failed to log in", err, login, http.StatusConflict)
		return
	}
	w.Header().Set("Token", token)
	response(w, "Sensor POST query execution.", err, dto.SuccessfulResponse{
		Message: "Successfully logged in",
	}, http.StatusConflict)
}

func (l *userController) Register(w http.ResponseWriter, r *http.Request) {
	defer func() {
		r.Body.Close()
	}()

	user := dto.Register{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		response(w, "Failed decoding ", err, user, http.StatusInternalServerError)
		return
	}

	err = l.userService.Register(r.Context(), user)
	user.Password = "****"
	response(w, "Sensor POST query execution.", err, user, http.StatusConflict)
}

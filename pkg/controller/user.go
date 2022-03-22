package controller

import (
	"encoding/json"
	"net/http"

	"github.com/Todorov99/server/pkg/dto"
	"github.com/Todorov99/server/pkg/service"
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

	login := dto.Login{}
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		response(w, "Failed decoding ", err, login, http.StatusInternalServerError)
		return
	}

	token, err := l.userService.Login(r.Context(), login)
	if err != nil {
		response(w, "Failed to log in", err, login, http.StatusConflict)
		return
	}
	w.Header().Set("Token", token)
	response(w, "Sensor POST query execution.", err, login, http.StatusConflict)
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
	response(w, "Sensor POST query execution.", err, user, http.StatusConflict)
}

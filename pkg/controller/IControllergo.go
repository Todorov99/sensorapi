package controller

import (
	"net/http"
	"os"

	"github.com/Todorov99/sensorcli/pkg/logger"
	"github.com/dgrijalva/jwt-go"
)

var controllerLogger = logger.NewLogrus("controller", os.Stdout)

// IController represent the commont REST verbs for the sensors and measurements
type IController interface {
	GetAll(w http.ResponseWriter, r *http.Request, token *jwt.Token)
	GetByID(w http.ResponseWriter, r *http.Request, token *jwt.Token)
	Post(w http.ResponseWriter, r *http.Request, token *jwt.Token)
	Put(w http.ResponseWriter, r *http.Request, token *jwt.Token)
	Delete(w http.ResponseWriter, r *http.Request, token *jwt.Token)
}

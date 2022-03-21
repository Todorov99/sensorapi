package controller

import (
	"net/http"
	"os"

	"github.com/Todorov99/sensorcli/pkg/logger"
)

var controllerLogger = logger.NewLogrus("controller", os.Stdout)

// IController represent the commont REST verbs for the sensors and measurements
type IController interface {
	GetAll(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	Post(w http.ResponseWriter, r *http.Request)
	Put(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

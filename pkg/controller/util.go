package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/Todorov99/sensorcli/pkg/logger"
	"github.com/Todorov99/server/pkg/models"
	"github.com/gorilla/mux"
)

var controllerLogger = logger.NewLogrus("controller", os.Stdout)

func response(w http.ResponseWriter, respondMessage string, loggMessagge string, err error, model interface{}, statusCode int) {
	controllerLogger.Info(loggMessagge)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		controllerLogger.Error(err)
		responseError := models.ResponseError{
			ErrMessage: err.Error(),
			Entity:     model,
		}
		json.NewEncoder(w).Encode(responseError)
		return
	}

	if respondMessage != "" {
		w.Header().Set("Content-Type", "text/html; charset=UTF-8")
		fmt.Fprintln(w, respondMessage)
	} else {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(model)
	}

}

func getIDFromPathVariable(r *http.Request) string {
	return mux.Vars(r)["id"]
}

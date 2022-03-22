package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Todorov99/server/pkg/dto"
	"github.com/gorilla/mux"
)

func response(w http.ResponseWriter, loggMessagge string, err error, model interface{}, statusCode int) {
	controllerLogger.Info(loggMessagge)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		controllerLogger.Error(err)
		responseError := dto.ResponseError{
			ErrMessage: err.Error(),
			Entity:     model,
		}
		json.NewEncoder(w).Encode(responseError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(model)
}

func getIDFromPathVariable(r *http.Request) int {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		fmt.Println(err)
	}
	return id
}

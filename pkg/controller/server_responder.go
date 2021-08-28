package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func respond(w http.ResponseWriter, respondMessage string, loggMessagge string, err error, model interface{}, statusCode int) {

	controllerLogger.Info(loggMessagge)

	if err != nil {
		controllerLogger.Error(err)
		http.Error(w, err.Error(), statusCode)
		return
	}

	if respondMessage != "" {
		fmt.Fprintln(w, respondMessage)
	}

	json.NewEncoder(w).Encode(model)
}

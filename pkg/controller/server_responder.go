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

	if respondMessage == "Skip" {
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

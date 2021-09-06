package controller

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func getIDFromPathVariable(r *http.Request) string {
	return mux.Vars(r)["id"]
}

func getURLQueryParams(r *http.Request, params ...string) []string {
	keys := r.URL.Query()
	fmt.Println(keys)
	if len(params) == 6 {
		return []string{keys.Get(params[0]), keys.Get(params[1]), keys.Get(params[2]), keys.Get(params[3]),
			keys.Get(params[4]), keys.Get(params[5])}
	}

	return []string{keys.Get(params[0]), keys.Get(params[1]), keys.Get(params[2]), keys.Get(params[3])}
}
package main

import (
	"fmt"
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	json := `{"status": "available", "environment": %q, "version": %q}`
	json = fmt.Sprintf(json, app.config.env, version)

	//Go defaults to Content-Type: text/plain; charset=utf-8" when this is not set
	w.Header().Set("Content-Type", "application/json")

	w.Write([]byte(json))
}

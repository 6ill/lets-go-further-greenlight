package main

import (
	"net/http"
)

func (app *Application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	// Create a map which holds the information that we want to send in the response.
	data := Envelope{
		"status":      "available",
		"system_info": map[string]string {
			"version":     version,
			"environment": app.config.env,
		},
	}

	err := app.writeJSON(w, http.StatusCreated, data, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
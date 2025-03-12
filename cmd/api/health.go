package main

import (
	"net/http"
)

// healthcheckHandler godoc
//
// @Summary      Healthcheck endpoint
// @Description  Returns the status of the application, including environment and version.
// @Tags         health
// @Accept       json
// @Produce      json
// @Success      200     {object}  map[string]string  "Returns the status of the application"
// @Failure      500     {object}  error  "Internal server error"
// @Router       /healthcheck [get]
func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "ok",
		"env":     app.config.env,
		"version": version,
	}

	if err := app.jsonResponse(w, http.StatusOK, data); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

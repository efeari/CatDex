package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *application) healthCheckHandler(c *gin.Context) {
	data := map[string]string{
		"status":  "ok",
		"env":     app.config.env,
		"version": version,
	}
	if err := writeJSON(c.Writer, http.StatusOK, data); err != nil {
		writeJSONError(c.Writer, http.StatusBadRequest, err.Error())
	}
}

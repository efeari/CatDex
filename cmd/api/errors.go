package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *application) internalServerError(c *gin.Context, err error) {
	app.logger.Errorw("internal error", "method", c.Request.Method, "path", c.Request.URL.Path, "error", err.Error())

	writeJSONError(c.Writer, http.StatusInternalServerError, "the server encountered a problem")
}

func (app *application) badRequestResponse(c *gin.Context, err error) {
	app.logger.Warnf("bad request", "method", c.Request.Method, "path", c.Request.URL.Path, "error", err.Error())

	writeJSONError(c.Writer, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundResponse(c *gin.Context, err error) {
	app.logger.Warnf("not found error", "method", c.Request.Method, "path", c.Request.URL.Path, "error", err.Error())

	writeJSONError(c.Writer, http.StatusNotFound, "not found")
}

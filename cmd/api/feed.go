package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *application) getUserFeedHandler(c *gin.Context) {

	//pagination, filters

	ctx := c.Request.Context()

	feed, err := app.store.Cats.GetGlobalFeed(ctx)
	if err != nil {
		writeJSONError(c.Writer, http.StatusInternalServerError, err.Error())
		return
	}

	if err := writeJSON(c.Writer, http.StatusOK, feed); err != nil {
		writeJSONError(c.Writer, http.StatusInternalServerError, err.Error())
	}
}

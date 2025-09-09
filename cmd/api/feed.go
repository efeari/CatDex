package main

import (
	"net/http"

	"github.com/efeari/catdex/internal/store.go"
	"github.com/gin-gonic/gin"
)

func (app *application) getUserFeedHandler(c *gin.Context) {

	fq := store.PaginatedFeedQuery{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
	}

	fq, err := fq.Parse(c)
	if err != nil {
		writeJSONError(c.Writer, http.StatusBadRequest, err.Error())
		return
	}

	if err := Validate.Struct(fq); err != nil {
		writeJSONError(c.Writer, http.StatusBadRequest, err.Error())
		return
	}

	ctx := c.Request.Context()

	feed, err := app.store.Cats.GetGlobalFeed(ctx, fq)
	if err != nil {
		writeJSONError(c.Writer, http.StatusInternalServerError, err.Error())
		return
	}

	if err := writeJSON(c.Writer, http.StatusOK, feed); err != nil {
		writeJSONError(c.Writer, http.StatusInternalServerError, err.Error())
	}
}

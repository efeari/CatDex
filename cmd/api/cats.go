package main

import (
	"errors"
	"net/http"

	"github.com/efeari/catdex/internal/store.go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type catPayload struct {
	Name        string `json:"name" validate:"required,min=3,max=25"`
	Description string `json:"description" validate:"required,min=10,max=255"`
	Location    string `json:"location" validate:"required,min=3,max=50"`
	//PhotoPath   string    `json:"photo_path"`
	UserID uuid.UUID `json:"user_id" validate:"required,uuid"`
}

func (app *application) getCatHandler(c *gin.Context) {
	catIDString := c.Param("catID")
	catID, err := uuid.Parse(catIDString)
	if err != nil {
		writeJSONError(c.Writer, http.StatusInternalServerError, err.Error())
		return
	}

	ctx := c.Request.Context()
	cat, err := app.store.Cats.GetByID(ctx, catID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			writeJSONError(c.Writer, http.StatusNotFound, err.Error())
		default:
			writeJSONError(c.Writer, http.StatusInternalServerError, err.Error())
		}
		return
	}

	if err := writeJSON(c.Writer, http.StatusOK, cat); err != nil {
		writeJSONError(c.Writer, http.StatusInternalServerError, err.Error())
		return
	}
}

func (app *application) createCatHandler(c *gin.Context) {
	var payload catPayload
	if err := readJSON(c.Writer, c.Request, &payload); err != nil {
		writeJSONError(c.Writer, 400, err.Error())
		return
	}

	if err := Validate.Struct(payload); err != nil {
		writeJSONError(c.Writer, http.StatusBadRequest, err.Error())
		return
	}

	cat := &store.Cat{
		Name:        payload.Name,
		Description: payload.Description,
		Location:    payload.Location,
		//PhotoPath:   payload.PhotoPath,
		UserID: payload.UserID,
	}

	ctx := c.Request.Context()

	if err := app.store.Cats.Create(ctx, cat); err != nil {
		writeJSONError(c.Writer, http.StatusInternalServerError, err.Error())
		return
	}

	if err := writeJSON(c.Writer, http.StatusCreated, cat); err != nil {
		writeJSONError(c.Writer, http.StatusInternalServerError, err.Error())
		return
	}
}

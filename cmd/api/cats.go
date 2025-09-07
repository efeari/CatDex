package main

import (
	"errors"
	"net/http"

	"github.com/efeari/catdex/internal/store.go"
	"github.com/efeari/catdex/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type catPayload struct {
	Name        string    `schema:"name" json:"name"`
	Description string    `schema:"description" json:"description"`
	Location    string    `schema:"location" json:"location"`
	UserID      uuid.UUID `schema:"user_id" json:"user_id"`
}

func (app *application) deleteCatHandler(c *gin.Context) {
	ctx := c.Request.Context()

	catIDString := c.Param("catID")
	catID, err := uuid.Parse(catIDString)
	if err != nil {
		writeJSONError(c.Writer, http.StatusInternalServerError, err.Error())
		return
	}

	err = app.store.Cats.DeleteByID(ctx, catID)
	if err != nil {
		writeJSONError(c.Writer, http.StatusInternalServerError, err.Error())
	}
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
	if err := readForm(c.Writer, c.Request, &payload); err != nil {
		writeJSONError(c.Writer, http.StatusBadRequest, err.Error())
		return
	}

	if err := Validate.Struct(payload); err != nil {
		writeJSONError(c.Writer, http.StatusBadRequest, err.Error())
		return
	}

	newUuid := uuid.New().String()

	photoPath, err := utils.HandleCatPhoto(c, newUuid)
	if err != nil {
		writeJSONError(c.Writer, http.StatusInternalServerError, err.Error())
	}

	cat := &store.Cat{
		Name:        payload.Name,
		Description: payload.Description,
		Location:    payload.Location,
		//PhotoPath:   payload.PhotoPath,
		UserID: payload.UserID,
	}

	cat.ID, err = uuid.Parse(newUuid)
	if err != nil {
		writeJSONError(c.Writer, http.StatusInternalServerError, err.Error())
		return
	}

	ctx := c.Request.Context()

	cat.PhotoPath = photoPath

	if err := app.store.Cats.Create(ctx, cat); err != nil {
		writeJSONError(c.Writer, http.StatusInternalServerError, err.Error())
		return
	}

	if err := writeJSON(c.Writer, http.StatusCreated, cat); err != nil {
		writeJSONError(c.Writer, http.StatusInternalServerError, err.Error())
		return
	}
}

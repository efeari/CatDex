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

type updateCatPayload struct {
	Name        *string `schema:"name" json:"name" validate:"omitempty"`
	Description *string `schema:"description" json:"description" validate:"omitempty,max=255"`
	Location    *string `schema:"location" json:"location" validate:"omitempty"`
	LastSeen    *string `json:"last_seen" db:"last_seen" validate:"omitempty"`
}

func (app *application) updateCatHandler(c *gin.Context) {
	cat, ok := getCatsFromCtx(c)
	if !ok {
		writeJSONError(c.Writer, http.StatusInternalServerError, "cat not found in context")
		return
	}

	var payload updateCatPayload
	if err := readJSON(c.Writer, c.Request, &payload); err != nil {
		writeJSONError(c.Writer, http.StatusBadRequest, err.Error())
		return
	}

	if err := Validate.Struct(payload); err != nil {
		writeJSONError(c.Writer, http.StatusBadRequest, err.Error())
		return
	}

	if payload.Name != nil {
		cat.Name = *payload.Name
	}
	if payload.Description != nil {
		cat.Description = *payload.Description
	}
	if payload.Location != nil {
		cat.Location = *payload.Location
	}
	if payload.LastSeen != nil {
		cat.LastSeen = *payload.LastSeen
	}

	if err := app.store.Cats.UpdateByID(c.Request.Context(), cat); err != nil {
		writeJSONError(c.Writer, http.StatusInternalServerError, err.Error())
		return
	}

	if err := writeJSON(c.Writer, http.StatusOK, cat); err != nil {
		writeJSONError(c.Writer, http.StatusInternalServerError, err.Error())
		return
	}
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
	cat, ok := getCatsFromCtx(c)
	if !ok {
		writeJSONError(c.Writer, http.StatusInternalServerError, "cat not found in context")
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
		utils.DeleteCatPhoto(newUuid)
		return
	}

	ctx := c.Request.Context()

	cat.PhotoPath = photoPath

	if err := app.store.Cats.Create(ctx, cat); err != nil {
		writeJSONError(c.Writer, http.StatusInternalServerError, err.Error())
		utils.DeleteCatPhoto(newUuid)
		return
	}

	if err := writeJSON(c.Writer, http.StatusCreated, cat); err != nil {
		writeJSONError(c.Writer, http.StatusInternalServerError, err.Error())
		return
	}
}

func (app *application) catsContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		catIDString := c.Param("catID")
		catID, err := uuid.Parse(catIDString)
		if err != nil {
			writeJSONError(c.Writer, http.StatusBadRequest, "invalid cat ID")
			c.Abort()
			return
		}

		fetchedCat, err := app.store.Cats.GetByID(c.Request.Context(), catID)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				writeJSONError(c.Writer, http.StatusNotFound, err.Error())
			default:
				writeJSONError(c.Writer, http.StatusInternalServerError, err.Error())
			}
			c.Abort()
			return
		}

		// Store the cat in Gin context
		c.Set("cat", fetchedCat)

		c.Next() // continue to handler
	}
}

func getCatsFromCtx(c *gin.Context) (*store.Cat, bool) {
	cat, exists := c.Get("cat")
	if !exists {
		return nil, false
	}
	fetchedCat, ok := cat.(*store.Cat)
	return fetchedCat, ok
}

package main

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"github.com/efeari/catdex/internal/store.go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,max=30"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

func (app *application) registerUserHandler(c *gin.Context) {

	var payload RegisterUserPayload
	if err := readJSON(c.Writer, c.Request, &payload); err != nil {
		app.badRequestResponse(c, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(c, err)
		return
	}

	user := &store.User{
		Username: payload.Username,
		Email:    payload.Email,
	}

	if err := user.Password.Set(payload.Password); err != nil {
		app.internalServerError(c, err)
		return
	}

	ctx := c.Request.Context()

	plainToken := uuid.New().String()
	hash := sha256.Sum256([]byte(plainToken))
	hashedToken := hex.EncodeToString(hash[:])

	err := app.store.Users.CreateAndInvite(ctx, user, hashedToken, app.config.mail.exp)
	if err != nil {
		switch err {
		case store.ErrDuplicateEmail:
			app.badRequestResponse(c, err)
		case store.ErrDuplicateUsername:
			app.badRequestResponse(c, err)
		default:
			app.internalServerError(c, err)
		}
		return
	}

	// For debugging, return the plainToken
	// if err := writeJSON(c.Writer, http.StatusCreated, plainToken); err != nil {
	// 	app.internalServerError(c, err)
	// }

	if err := writeJSON(c.Writer, http.StatusCreated, nil); err != nil {
		app.internalServerError(c, err)
	}
}

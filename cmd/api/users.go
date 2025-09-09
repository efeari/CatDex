package main

import (
	"errors"
	"net/http"

	"github.com/efeari/catdex/internal/store.go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (app *application) getUserHandler(c *gin.Context) {
	user, ok := getUserFromCtx(c)
	if !ok {
		app.internalServerError(c, errors.New("user not found in context"))
		return
	}

	if err := writeJSON(c.Writer, http.StatusOK, user); err != nil {
		app.internalServerError(c, err)
		return
	}
}

func (app *application) usersContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDString := c.Param("userID")
		userID, err := uuid.Parse(userIDString)
		if err != nil {
			app.badRequestResponse(c, err)
			c.Abort()
			return
		}

		fetchedUser, err := app.store.Users.GetByID(c.Request.Context(), userID)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundResponse(c, err)
			default:
				app.internalServerError(c, err)
			}
			c.Abort()
			return
		}

		// Store the user in Gin context
		c.Set("user", fetchedUser)

		c.Next() // continue to handler
	}
}

func getUserFromCtx(c *gin.Context) (*store.User, bool) {
	user, exists := c.Get("user")
	if !exists {
		return nil, false
	}
	fetchedUser, ok := user.(*store.User)
	return fetchedUser, ok
}

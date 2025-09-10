package main

import (
	"net/http"
	"time"

	"github.com/efeari/catdex/internal/store.go"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const version = "0.0.1"

type application struct {
	config config
	store  store.Storage
	logger *zap.SugaredLogger
}

// General configuration
type config struct {
	addr string
	db   dbConfig
	env  string
	mail mailConfig
}

type mailConfig struct {
	exp time.Duration
}

// Database configuration struct
type dbConfig struct {
	addr               string
	maxOpenConnections int
	maxIdleConnections int
	maxIdleTime        string
}

func (app *application) mount() http.Handler {
	//gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	v1 := r.Group("/v1")
	v1.Static("/photos", "../../photos")

	v1.GET("/health", app.healthCheckHandler)

	v1.GET("/cat/:catID", app.catsContextMiddleware(), app.getCatHandler)
	v1.DELETE("/cat/:catID", app.catsContextMiddleware(), app.deleteCatHandler)
	v1.PATCH("/cat/:catID", app.catsContextMiddleware(), app.updateCatHandler)
	v1.POST("/cat", app.createCatHandler)

	v1.GET("/user/:userID", app.usersContextMiddleware(), app.getUserHandler)

	v1.GET("/feed", app.getUserFeedHandler)

	v1.POST("authentication/user", app.registerUserHandler)
	v1.PUT("/users/activate/:token", app.activateUserHandler)
	return r
}

func (app *application) run(mux http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	app.logger.Infow("server has started", "addr", app.config.addr, "env", app.config.env)

	return srv.ListenAndServe()
}

package main

import (
	"log"
	"net/http"
	"time"

	"github.com/efeari/catdex/internal/store.go"
	"github.com/gin-gonic/gin"
)

const version = "0.0.1"

type application struct {
	config config
	store  store.Storage
}

// General configuration
type config struct {
	addr string
	db   dbConfig
	env  string
}

// Database configuration struct
type dbConfig struct {
	addr               string
	maxOpenConnections int
	maxIdleConnections int
	maxIdleTime        string
}

func (app *application) mount() http.Handler {
	r := gin.Default()

	v1 := r.Group("/v1")
	v1.GET("/health", app.healthCheckHandler)
	v1.GET("/cat/:catID", app.getCatHandler)

	v1.POST("/cat", app.createCatHandler)
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

	log.Printf("server has started at %s", app.config.addr)

	return srv.ListenAndServe()
}

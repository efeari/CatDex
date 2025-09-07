package main

import (
	"net/http"

	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

func readForm(w http.ResponseWriter, r *http.Request, data any) error {
	// Limit body size (same as readJSON)
	const maxBytes = 9999999999999 // 1MB
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	// Parse form values
	if err := r.ParseMultipartForm(maxBytes); err != nil {
		return err
	}

	// Disallow unknown fields (default is to ignore extras)
	decoder.IgnoreUnknownKeys(false)

	// Map form values to struct
	return decoder.Decode(data, r.Form)
}

package models

import (
	"time"
)

type Cat struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	DateOfPhoto time.Time `json:"date_of_photo"`
	Location    string    `json:"location"`
	PhotoPath   string    `json:"photo_path"`
	CreatedAt   time.Time `json:"created_at"`
}

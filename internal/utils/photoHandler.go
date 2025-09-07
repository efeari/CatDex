package utils

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// Copies the photo to local, and returns the path
func HandleCatPhoto(c *gin.Context, catID string) (string, error) {

	file, err := c.FormFile("photo")
	if err != nil {
		return "", err
	}

	ext := filepath.Ext(file.Filename)
	if ext == "" {
		ext = ".jpg"
	}
	filename := fmt.Sprintf("%s%s", catID, ext)

	// Ensure folder exists
	err = os.MkdirAll("../../photos", os.ModePerm)
	if err != nil {
		c.Error(errors.New("failed to create photos folder"))
		return "", err
	}

	// Open uploaded file
	src, err := file.Open()
	if err != nil {
		c.Error(errors.New("failed to open uploaded file"))
		return "", err
	}
	defer src.Close()

	// Create destination file
	photoPath := filepath.Join("../../photos", filename)
	out, err := os.Create(photoPath)
	if err != nil {
		c.Error(errors.New("failed to create file"))
		return "", err
	}
	defer out.Close()

	// Copy content
	_, err = io.Copy(out, src)
	if err != nil {
		c.Error(errors.New("failed to save photo"))
		return "", err
	}

	return photoPath, nil
}

package api

import (
	"errors"
	"log"

	"github.com/gin-gonic/gin"
)

var (
	ErrInternalServerError = errors.New("internal server error")

	ErrInvalidRequest     = errors.New("invalid request")
	ErrMissingAuthPayload = errors.New("missing auth payload")
	ErrNotFound           = errors.New("not found")
	ErrPermissionDenied   = errors.New("permission denied")
)

func errorResponse(err error) gin.H {
	if errors.Is(err, ErrInternalServerError) {
		log.Println("Internal server error:", err)
	}

	return gin.H{"error": err.Error()}
}

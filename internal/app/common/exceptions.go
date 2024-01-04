package common

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type APIError struct {
	StatusCode int    `json:"statusCode"`
	Err        error  `json:"error,omitempty"`
	Message    string `json:"message"`
}

func (e *APIError) Error() string {
	return e.Err.Error()
}

func GlobalErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() //Only inside middleware. It executes the pending handlers in the chain inside the calling handler.

		// Check if there are any errors
		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				// Check if it's an APIError
				var apiErr *APIError
				if errors.As(e.Err, &apiErr) {
					c.JSON(apiErr.StatusCode, APIError{StatusCode: apiErr.StatusCode, Message: apiErr.Message})
					return
				}
			}

			// If it's not an APIError, return a generic server error
			c.JSON(http.StatusInternalServerError, APIError{StatusCode: http.StatusInternalServerError, Message: "Internal server error"})
		}
	}
}

// NewAPIError creates an APIError and adds it to the Gin context
func NewAPIError(c *gin.Context, statusCode int, err error, message string) {
	_ = c.Error(&APIError{
		StatusCode: statusCode,
		Err:        err,
		Message:    message,
	})
}

func HandleNoRoute() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotFound, APIError{StatusCode: http.StatusNotFound, Message: "Not found"})
	}
}

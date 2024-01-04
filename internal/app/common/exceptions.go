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

func UnauthorizedError(message string) *APIError {
	return &APIError{
		StatusCode: http.StatusUnauthorized,
		Message:    message,
	}
}

func EmptyParamsError(param string) *APIError {
	return &APIError{
		StatusCode: http.StatusBadRequest,
		Message:    param + " cannot be empty",
	}
}

func FileOwnerMismatch() *APIError {
	return &APIError{
		StatusCode: http.StatusUnauthorized,
		Message:    "token and file owner do not match",
	}
}

func FileCreationMismatch() *APIError {
	return &APIError{
		StatusCode: http.StatusUnauthorized,
		Message:    "can't create a file for another user, file owner and token do not match",
	}
}

func BadRequestError(message string) *APIError {
	return &APIError{
		StatusCode: http.StatusBadRequest,
		Message:    message,
	}
}

func ForwardError(c *gin.Context, apiError *APIError) {
	_ = c.Error(apiError)
}

func HandleNoRoute() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotFound, APIError{StatusCode: http.StatusNotFound, Message: "Not found"})
	}
}

// HandleError abstracts the error handling logic.
func HandleError(c *gin.Context, err error) {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		ForwardError(c, apiErr)
		return
	}

	_ = c.Error(&APIError{
		StatusCode: http.StatusInternalServerError,
		Err:        err,
		Message:    err.Error(),
	})
}

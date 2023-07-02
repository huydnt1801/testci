package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-logr/logr"
	"github.com/labstack/echo/v4"
)

type ErrorResponse struct {
	Code       int    `json:"code"`
	StatusCode int    `json:"-"`
	Message    string `json:"message"`
	Err        error  `json:"-"`
}

func CustomHTTPErrorHandler(logger logr.Logger) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}
		e := &ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    http.StatusText(http.StatusInternalServerError),
			Err:        err,
		}
		if httpErr, ok := err.(interface {
			HTTPStatusCode() int
		}); ok {
			e.StatusCode = httpErr.HTTPStatusCode()
			e.Code = e.StatusCode
			e.Message = err.Error()
		}

		if he, ok := err.(*echo.HTTPError); ok {
			e.StatusCode = he.Code
			e.Code = e.StatusCode
			e.Message = fmt.Sprintf("%v", he.Message)
		}

		logger.Error(e.Err, e.Message)
		if err = c.JSON(e.StatusCode, e); err != nil {
			logger.Error(err, "Can not handle error")
		}
	}
}

// Only for testing purposes
func GetErrorCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	var httpErr interface {
		error
		HTTPStatusCode() int
	}
	if errors.As(err, &httpErr) {
		return httpErr.HTTPStatusCode()
	}

	var echoErr *echo.HTTPError
	if errors.As(err, &echoErr) {
		return echoErr.Code
	}

	return http.StatusInternalServerError
}

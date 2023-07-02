package http

import "net/http"

type NotFoundError struct{}

func (NotFoundError) Error() string {
	return "Not Found"
}

func (NotFoundError) HTTPStatusCode() int {
	return http.StatusNotFound
}

type BadRequestError struct{}

func (BadRequestError) Error() string {
	return "Bad Request"
}

func (BadRequestError) HTTPStatusCode() int {
	return http.StatusBadRequest
}

type UnauthorizedError struct{}

func (UnauthorizedError) Error() string {
	return "Unauthorized"
}

func (UnauthorizedError) HTTPStatusCode() int {
	return http.StatusUnauthorized
}

type ForbiddenError struct{}

func (ForbiddenError) Error() string {
	return "Forbidden"
}

func (ForbiddenError) HTTPStatusCode() int {
	return http.StatusForbidden
}

type ConflictError struct{}

func (ConflictError) Error() string {
	return "Conflict"
}

func (ConflictError) HTTPStatusCode() int {
	return http.StatusConflict
}

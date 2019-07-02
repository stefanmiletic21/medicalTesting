package serverErr

import "errors"

var (
	// Returned as 400 HTTP Status
	ErrBadRequest = errors.New("Bad request")

	// Returned as 401 HTTP Status
	ErrNotAuthenticated = errors.New("Not authenticated")

	// Returned as 403 HTTP Status
	ErrForbidden = errors.New("Forbidden")

	// Returned as 404 HTTP Status
	ErrInvalidAPICall = errors.New("Invalid API call")

	// Returned as 404 HTTP Status
	ErrResourceNotFound = errors.New("Resource not found")

	// Returned as 405 HTTP Status
	ErrMethodNotAllowed = errors.New("Method not allowed")

	// Returned as 500 HTTP Status
	ErrInternal = errors.New("Internal error")
)

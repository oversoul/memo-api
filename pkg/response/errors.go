package response

import "net/http"

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// Error is required by the error interface.
func (e ErrorResponse) Error() string {
	return e.Message
}

// StatusCode is required by routing.HTTPError interface.
func (e ErrorResponse) StatusCode() int {
	return e.Status
}

// InternalServerError creates a new error response representing an internal server error (HTTP 500)
func InternalServerError() ErrorResponse {
	var msg = "We encountered an error while processing your request."
	return ErrorResponse{
		Status:  http.StatusInternalServerError,
		Message: msg,
	}
}

// NotFound creates a new error response representing a resource-not-found error (HTTP 404)
func NotFound() ErrorResponse {
	var msg = "The requested resource was not found."
	return ErrorResponse{
		Status:  http.StatusNotFound,
		Message: msg,
	}
}

// Unauthorized creates a new error response representing an authentication/authorization failure (HTTP 401)
func Unauthorized() ErrorResponse {
	var msg = "You are not authenticated to perform the requested action."
	return ErrorResponse{
		Status:  http.StatusUnauthorized,
		Message: msg,
	}
}

// Forbidden creates a new error response representing an authorization failure (HTTP 403)
func Forbidden() ErrorResponse {
	var msg = "You are not authorized to perform the requested action."
	return ErrorResponse{
		Status:  http.StatusForbidden,
		Message: msg,
	}
}

// BadRequest creates a new error response representing a bad request (HTTP 400)
func BadRequest() ErrorResponse {
	var msg = "Your request is in a bad format."
	return ErrorResponse{
		Status:  http.StatusBadRequest,
		Message: msg,
	}
}

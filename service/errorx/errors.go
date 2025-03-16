// Package errorx exposes Error object that can be returned to clients as structured error.
package errorx

type Code string

type Error struct {
	Code    Code   `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	// It is possible to expose a key to offer rich info about the error for client to work on.
}

// Implementing the error interface of the Golang's error package, so that the Error object can be returned as `error` type
func (e Error) Error() string {
	return string(e.Code)
}

const (
	ServerError    Code = "SERVER_ERROR"
	BadRequestData Code = "BAD_REQUEST_DATA"
	AccessDenied   Code = "ACCESS_DENIED"
	InvalidToken   Code = "INVALID_TOKEN"
	NoContent      Code = "NO_CONTENT"
)

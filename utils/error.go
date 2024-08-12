package utils

import "net/http"

type ErrorCode int

const (
	ErrInternalServer ErrorCode = iota + 1
	ErrBadRequest
	ErrNotFound
	DataNotExist
	ClientCodeEmpty
	BodyError
	EmptyFields
	DatabaseStore
	Duplicate
)

var ErrorMessage = map[ErrorCode]string{
	ErrInternalServer: "Internal server error",
	ErrBadRequest:     "Bad request",
	ErrNotFound:       "Not found",
	DataNotExist:      "Data does not exist in database",
	ClientCodeEmpty:   "Client code cannot be empty",
	BodyError:         "Unable to parse body",
	EmptyFields:       "All fields are required and cannot be empty",
	DatabaseStore:     "Failed to store data in the database",
	Duplicate:         "Duplicate data",
}

var HTTPStatusCodeMap = map[ErrorCode]int{
	ErrInternalServer: http.StatusInternalServerError,
	ErrBadRequest:     http.StatusBadRequest,
	ErrNotFound:       http.StatusNotFound,
	DataNotExist:      http.StatusSeeOther,
	ClientCodeEmpty:   http.StatusInternalServerError,
	BodyError:         http.StatusInternalServerError,
	EmptyFields:       http.StatusBadRequest,
	DatabaseStore:     http.StatusInternalServerError,
	Duplicate:         http.StatusBadRequest,
}

type Error struct {
	Code    ErrorCode
	Message string
}

func NewError(code ErrorCode) *Error {
	return &Error{
		Code:    code,
		Message: ErrorMessage[code],
	}
}

func (e *Error) Error() string {
	return e.Message
}

package httpx

import "net/http"

type AppError struct {
	Status  int    `json:"-"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string { return e.Message }

func BadRequest(code, msg string, err error) *AppError {
	return &AppError{Status: http.StatusBadRequest, Code: code, Message: msg, Err: err}
}

func Unauthorized(msg string, err error) *AppError {
	return &AppError{Status: http.StatusUnauthorized, Code: "Unauthorized", Message: msg, Err: err}
}

func Forbidden(msg string, err error) *AppError {
	return &AppError{Status: http.StatusForbidden, Code: "Forbidden", Message: msg, Err: err}
}

func NotFound(code, msg string, err error) *AppError {
	return &AppError{Status: http.StatusNotFound, Code: code, Message: msg, Err: err}
}

func Conflict(code, msg string, err error) *AppError {
	return &AppError{Status: http.StatusConflict, Code: code, Message: msg, Err: err}
}

func Internal(msg string, err error) *AppError {
	return &AppError{Status: http.StatusInternalServerError, Code: "InternalError", Message: msg, Err: err}
}

package common

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type HTTPException struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors,omitempty"`
}

func (e *HTTPException) Error() string {
	if e == nil {
		return ""
	}
	return fmt.Sprintf("%d: %s", e.Status, e.Message)
}

func (e *HTTPException) Write(w http.ResponseWriter) {
	if e == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.Status)
	_ = json.NewEncoder(w).Encode(e)
}

func NewHTTPException(status int, message string, errs interface{}) *HTTPException {
	if message == "" {
		message = http.StatusText(status)
	}
	return &HTTPException{
		Status:  status,
		Message: message,
		Errors:  errs,
	}
}

func BadRequest(message string, errs ...interface{}) *HTTPException {
	var e interface{}
	if len(errs) > 0 {
		e = errs[0]
	}
	return NewHTTPException(http.StatusBadRequest, message, e)
}

func Unauthorized(message string, errs ...interface{}) *HTTPException {
	var e interface{}
	if len(errs) > 0 {
		e = errs[0]
	}
	return NewHTTPException(http.StatusUnauthorized, message, e)
}

func Forbidden(message string, errs ...interface{}) *HTTPException {
	var e interface{}
	if len(errs) > 0 {
		e = errs[0]
	}
	return NewHTTPException(http.StatusForbidden, message, e)
}

func NotFound(message string, errs ...interface{}) *HTTPException {
	var e interface{}
	if len(errs) > 0 {
		e = errs[0]
	}
	return NewHTTPException(http.StatusNotFound, message, e)
}

func MethodNotAllowed(message string, errs ...interface{}) *HTTPException {
	var e interface{}
	if len(errs) > 0 {
		e = errs[0]
	}
	return NewHTTPException(http.StatusMethodNotAllowed, message, e)
}

func Conflict(message string, errs ...interface{}) *HTTPException {
	var e interface{}
	if len(errs) > 0 {
		e = errs[0]
	}
	return NewHTTPException(http.StatusConflict, message, e)
}

func UnprocessableEntity(message string, errs ...interface{}) *HTTPException {
	var e interface{}
	if len(errs) > 0 {
		e = errs[0]
	}
	return NewHTTPException(http.StatusUnprocessableEntity, message, e)
}

func TooManyRequests(message string, errs ...interface{}) *HTTPException {
	var e interface{}
	if len(errs) > 0 {
		e = errs[0]
	}
	return NewHTTPException(http.StatusTooManyRequests, message, e)
}

func InternalServerError(message string, errs ...interface{}) *HTTPException {
	var e interface{}
	if len(errs) > 0 {
		e = errs[0]
	}
	return NewHTTPException(http.StatusInternalServerError, message, e)
}

func ServiceUnavailable(message string, errs ...interface{}) *HTTPException {
	var e interface{}
	if len(errs) > 0 {
		e = errs[0]
	}
	return NewHTTPException(http.StatusServiceUnavailable, message, e)
}

func GatewayTimeout(message string, errs ...interface{}) *HTTPException {
	var e interface{}
	if len(errs) > 0 {
		e = errs[0]
	}
	return NewHTTPException(http.StatusGatewayTimeout, message, e)
}

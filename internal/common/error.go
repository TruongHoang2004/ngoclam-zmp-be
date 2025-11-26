package common

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/constant"
)

type CodeResponse string

type ErrorResponse struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	TraceID    string `json:"trace_id,omitempty"`
	Detail     string `json:"detail,omitempty"`
	Source     string `json:"source"`
	HTTPStatus int    `json:"http_status"`
}

const (
	//internal
	ErrorCodeBadRequest   CodeResponse = "BAD_REQUEST"
	ErrorCodeUnauthorized CodeResponse = "UNAUTHORIZED"
	ErrorCodeForbidden    CodeResponse = "FORBIDDEN"
	ErrorCodeNotFound     CodeResponse = "NOT_FOUND"
	ErrorCodeConflict     CodeResponse = "CONFLICT"
	ErrorCodeSystemError  CodeResponse = "INTERNAL_SERVER_ERROR"
)

type Source string

const (
	CurrentService Source = constant.ServiceName
)

type Error struct {
	Code       CodeResponse `json:"code"`
	Message    string       `json:"message"`
	TraceID    string       `json:"trace_id,omitempty"`
	Detail     string       `json:"detail"`
	Source     Source       `json:"source"`
	HTTPStatus int          `json:"http_status"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("code:[%s], message:[%s], detail:[%s], source:[%s]", e.Code, e.Message, e.Detail, e.Source)
}

func (e *Error) GetHttpStatus() int {
	return e.HTTPStatus
}

func (e *Error) GetCode() CodeResponse {
	return e.Code
}

func (e *Error) GetMessage() string {
	return e.Message
}

func (e *Error) SetTraceId(traceId string) *Error {
	e.TraceID = fmt.Sprintf("%s:%d", traceId, time.Now().Unix())
	return e
}

func (e *Error) SetHTTPStatus(status int) *Error {
	e.HTTPStatus = status
	return e
}

func (e *Error) SetMessage(msg string) *Error {
	e.Message = msg
	return e
}

func (e *Error) SetDetail(detail string) *Error {
	e.Detail = detail
	return e
}

func (e *Error) GetDetail() string {
	return e.Detail
}

func (e *Error) SetSource(source Source) *Error {
	e.Source = source
	return e
}

func (e *Error) ToJSon() string {
	data, err := json.Marshal(e)
	if err != nil {
		//Todo fix this
		return "marshal error failed"
	}
	return string(data)
}

var (
	// Status 4xx ********

	ErrUnauthorized = func(ctx context.Context) *Error {
		traceId := GetTraceId(ctx)
		return &Error{
			Code:       ErrorCodeUnauthorized,
			Message:    DefaultUnauthorizedMessage,
			TraceID:    traceId,
			Source:     CurrentService,
			HTTPStatus: http.StatusUnauthorized,
		}
	}

	ErrNotFound = func(ctx context.Context, object, status string) *Error {
		traceId := GetTraceId(ctx)
		return &Error{
			Code:       ErrorCodeNotFound,
			Message:    getMsg(object, status),
			TraceID:    traceId,
			HTTPStatus: http.StatusNotFound,
		}
	}

	ErrBadRequest = func(ctx context.Context) *Error {
		traceId := GetTraceId(ctx)
		return &Error{
			Code:       ErrorCodeBadRequest,
			Message:    DefaultBadRequestMessage,
			TraceID:    traceId,
			HTTPStatus: http.StatusBadRequest,
			Source:     CurrentService,
		}
	}

	ErrConflict = func(ctx context.Context, object, status string) *Error {
		traceId := GetTraceId(ctx)
		return &Error{
			Code:       ErrorCodeConflict,
			Message:    getMsg(object, status),
			TraceID:    traceId,
			HTTPStatus: http.StatusConflict,
			Source:     CurrentService,
		}
	}

	// Status 5xx *******

	ErrSystemError = func(ctx context.Context, detail string) *Error {
		traceId := GetTraceId(ctx)
		return &Error{
			Code:       ErrorCodeSystemError,
			Message:    DefaultServerErrorMessage,
			TraceID:    traceId,
			HTTPStatus: http.StatusInternalServerError,
			Source:     CurrentService,
			Detail:     detail,
		}
	}

	ErrForbidden = func(ctx context.Context) *Error {
		traceId := GetTraceId(ctx)
		return &Error{
			Code:       ErrorCodeForbidden,
			Message:    DefaultForbiddenMessage,
			TraceID:    traceId,
			HTTPStatus: http.StatusForbidden,
			Source:     CurrentService,
		}
	}
)

func getMsg(object, status string) string {
	return fmt.Sprintf("%s %s", object, status)
}

const (
	DefaultServerErrorMessage  = "Something has gone wrong, please contact admin"
	DefaultBadRequestMessage   = "Invalid request"
	DefaultUnauthorizedMessage = "Token invalid"
	DefaultForbiddenMessage    = "Forbidden"
	DefauultConflict           = "Conflict"
)

func ConvertErrorToResponse(err *Error) *ErrorResponse {

	return &ErrorResponse{
		Code:       string(err.Code),
		Message:    err.Message,
		TraceID:    err.TraceID,
		Detail:     err.Detail,
		Source:     string(err.Source),
		HTTPStatus: err.HTTPStatus,
	}
}

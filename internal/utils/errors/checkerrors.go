package errors

import (
	"errors"
	"net/http"
)

// ErrorCode представляет код ошибки, который будет возвращаться клиенту
type ErrorCode string

// Константы с кодами ошибок
const (
	// Общие коды ошибок
	CodeUnknownError   ErrorCode = "UNKNOWN_ERROR"
	CodeInvalidRequest ErrorCode = "INVALID_REQUEST"
	CodeInternalError  ErrorCode = "INTERNAL_ERROR"
	CodeNotFound       ErrorCode = "NOT_FOUND"
	CodeUnauthorized   ErrorCode = "UNAUTHORIZED"
	CodeForbidden      ErrorCode = "FORBIDDEN"

	// Пользовательские коды ошибок
	CodeUserNotFound      ErrorCode = "USER_NOT_FOUND"
	CodeUserAlreadyExists ErrorCode = "USER_ALREADY_EXISTS"
	CodeInvalidPassword   ErrorCode = "INVALID_PASSWORD"
	CodeInvalidEmail      ErrorCode = "INVALID_EMAIL"
	CodeInvalidUsername   ErrorCode = "INVALID_USERNAME"

	// Коды ошибок для операций с данными
	CodeDataNotFound ErrorCode = "DATA_NOT_FOUND"
	CodeDataInvalid  ErrorCode = "DATA_INVALID"
	CodeDataConflict ErrorCode = "DATA_CONFLICT"
)

// HTTPStatusMapping сопоставляет коды ошибок с HTTP-статусами
var HTTPStatusMapping = map[ErrorCode]int{
	// Общие коды
	CodeUnknownError:   http.StatusInternalServerError,
	CodeInvalidRequest: http.StatusBadRequest,
	CodeInternalError:  http.StatusInternalServerError,
	CodeNotFound:       http.StatusNotFound,
	CodeUnauthorized:   http.StatusUnauthorized,
	CodeForbidden:      http.StatusForbidden,

	// Пользовательские коды
	CodeUserNotFound:      http.StatusNotFound,
	CodeUserAlreadyExists: http.StatusConflict,
	CodeInvalidPassword:   http.StatusBadRequest,
	CodeInvalidEmail:      http.StatusBadRequest,
	CodeInvalidUsername:   http.StatusBadRequest,

	// Коды для операций с данными
	CodeDataNotFound: http.StatusNotFound,
	CodeDataInvalid:  http.StatusBadRequest,
	CodeDataConflict: http.StatusConflict,
}

// APIError представляет структуру ошибки для API ответов
type APIError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

// Error представляет ошибку с дополнительным контекстом для API
type Error struct {
	Err     error
	Code    ErrorCode
	Message string
}

// New создает новую ошибку с указанным кодом и сообщением
func New(code ErrorCode, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

// NewWithError создает новую ошибку на основе существующей ошибки
func NewWithError(err error, code ErrorCode, message string) *Error {
	return &Error{
		Err:     err,
		Code:    code,
		Message: message,
	}
}

// Error реализует интерфейс error
func (e *Error) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

// ToAPIError конвертирует Error в APIError для ответа клиенту
func (e *Error) ToAPIError() APIError {
	return APIError{
		Code:    e.Code,
		Message: e.Message,
	}
}

// GetHTTPStatus возвращает HTTP статус для ошибки
func (e *Error) GetHTTPStatus() int {
	if status, ok := HTTPStatusMapping[e.Code]; ok {
		return status
	}
	return http.StatusInternalServerError
}

// IsErrorCode проверяет, соответствует ли ошибка указанному коду
func IsErrorCode(err error, code ErrorCode) bool {
	var e *Error
	if errors.As(err, &e) {
		return e.Code == code
	}
	return false
}

package errors

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response представляет стандартный формат ответа API
type Response struct {
	Data    any       `json:"data,omitempty"`
	Error   *APIError `json:"error,omitempty"`
	Success bool      `json:"success"`
}

// SuccessResponse создает успешный ответ с данными
func SuccessResponse(data any) Response {
	return Response{
		Success: true,
		Data:    data,
	}
}

// ErrorResponse создает ответ с ошибкой
func ErrorResponse(err error) Response {
	var customErr *Error
	var apiErr APIError

	if errors.As(err, &customErr) {
		// Если это наша кастомная ошибка, используем информацию из нее
		apiErr = customErr.ToAPIError()
	} else {
		// Если это стандартная ошибка, используем общий код ошибки
		apiErr = APIError{
			Code:    CodeUnknownError,
			Message: err.Error(),
		}
	}

	return Response{
		Success: false,
		Error:   &apiErr,
	}
}

// RespondWithError отправляет ответ с ошибкой через gin.Context
func RespondWithError(c *gin.Context, err error) {
	var customErr *Error
	var statusCode int
	var response Response

	if errors.As(err, &customErr) {
		// Для кастомной ошибки используем соответствующий HTTP-статус
		statusCode = customErr.GetHTTPStatus()
		response = ErrorResponse(err)
	} else {
		// Для стандартной ошибки используем Internal Server Error
		statusCode = http.StatusInternalServerError
		response = ErrorResponse(err)
	}

	c.JSON(statusCode, response)
}

// RespondWithSuccess отправляет успешный ответ через gin.Context
func RespondWithSuccess(c *gin.Context, data any) {
	c.JSON(http.StatusOK, SuccessResponse(data))
}

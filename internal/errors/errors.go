package errors

import "fmt"

// ForumError представляет базовую структуру для всех ошибок форума
type ForumError struct {
	Code    int    // HTTP статус код
	Message string // Сообщение об ошибке
	Err     error  // Вложенная ошибка
}

func (e *ForumError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Константы для типов ошибок
const (
	ErrNotFound          = "not_found error"
	ErrUnauthorized      = "unauthorized error"
	ErrForbidden         = "forbidden error"
	ErrBadRequest        = "bad_request error"
	ErrInternalServer    = "internal_server_error"
	ErrValidation        = "validation_error"
	ErrDuplicate         = "duplicate_error"
	ErrPermissionDenied  = "permission_denied error"
)

// Функции-конструкторы для создания ошибок
func NewNotFoundError(message string, err error) *ForumError {
	return &ForumError{
		Code:    404,
		Message: message,
		Err:     err,
	}
}

func NewUnauthorizedError(message string, err error) *ForumError {
	return &ForumError{
		Code:    401,
		Message: message,
		Err:     err,
	}
}

func NewForbiddenError(message string, err error) *ForumError {
	return &ForumError{
		Code:    403,
		Message: message,
		Err:     err,
	}
}

func NewBadRequestError(message string, err error) *ForumError {
	return &ForumError{
		Code:    400,
		Message: message,
		Err:     err,
	}
}

func NewInternalServerError(message string, err error) *ForumError {
	return &ForumError{
		Code:    500,
		Message: message,
		Err:     err,
	}
}

func NewValidationError(message string, err error) *ForumError {
	return &ForumError{
		Code:    400,
		Message: message,
		Err:     err,
	}
}

func NewDuplicateError(message string, err error) *ForumError {
	return &ForumError{
		Code:    409,
		Message: message,
		Err:     err,
	}
}

func NewPermissionDeniedError(message string, err error) *ForumError {
	return &ForumError{
		Code:    403,
		Message: message,
		Err:     err,
	}
} 
package apperror

import "fmt"

type ErrorType int

const (
	ErrTypeUnknown ErrorType = iota
	ErrorTypeNotFound
	ErrorTypeInternal
	ErrorTypeConflict
	ErrorInvalidArguement
	ErrorDeadlineExceed
	ErrorUnAuthenticated
	ErrorBadRequest
)

type DomainError struct {
	Type    ErrorType
	Message string
	Err     error
}

func (e *DomainError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s:%v", e.Message, e.Err)
	}
	return e.Message
}

func NewNotFoundError(msg string) *DomainError {
	return &DomainError{Type: ErrorTypeNotFound, Message: msg}
}

func NewInternalError(msg string, err error) *DomainError {
	return &DomainError{Type: ErrorTypeConflict, Message: msg, Err: err}
}

func NewConflictError(msg string) *DomainError {
	return &DomainError{
		Type:    ErrorTypeConflict,
		Message: msg,
	}
}

func NewValidationError(msg string) *DomainError {
	return &DomainError{
		Type:    ErrorInvalidArguement,
		Message: msg,
	}
}

func NewDeadlineExceedeError(msg string) *DomainError {
	return &DomainError{
		Type:    ErrorDeadlineExceed,
		Message: msg,
	}
}

func NewUnAuthenticatedError(msg string) *DomainError {
	return &DomainError{
		Type:    ErrorUnAuthenticated,
		Message: msg,
	}

}

func NewBadRequestError(msg string) *DomainError {
	return &DomainError{
		Type:    ErrorBadRequest,
		Message: msg,
	}
}

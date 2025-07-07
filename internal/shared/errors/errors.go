package errors

import "errors"

type DomainError struct {
	ErrType ErrorType
	err     error
}

type ErrorType int

const (
	// common 1 -> 1000
	System ErrorType = iota + 1
	Validation
	NotAuthorized
	Forbidden
	NotFound
	AlreadyExist
	// module specific 1001 -> 2000
)

func (e *DomainError) Error() string {
	if e == nil {
		return ""
	}
	return e.err.Error()
}

func (e *DomainError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.err
}

func (e *DomainError) GetType() ErrorType {
	if e == nil {
		return System
	}
	return e.ErrType
}
func NewDomainError(errType ErrorType, message string) *DomainError {
	return &DomainError{
		ErrType: errType,
		err:     errors.New(message),
	}
}

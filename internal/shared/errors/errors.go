package errors

import "errors"

type DomainError struct {
	ErrCode ErrorCode
	err     error
}

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

func (e *DomainError) ErrorCode() ErrorCode {
	if e == nil {
		return System
	}
	return e.ErrCode
}

func NewDomainError(errCode ErrorCode, message string) *DomainError {
	return &DomainError{
		ErrCode: errCode,
		err:     errors.New(message),
	}
}

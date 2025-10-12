package richerror

import "transaction/pkg/genericcode"

type RichError struct {
	Code      genericcode.Code
	Message   string
	WrapError error
	Data      any
}

func (r RichError) Error() string {
	if r.Message == "" && r.WrapError != nil {
		return r.WrapError.Error()
	}

	return r.Message
}

func (r RichError) GetCode() genericcode.Code {
	if r.WrapError != nil {
		wrapError, ok := r.WrapError.(RichError)
		if ok {
			return wrapError.GetCode()
		}
	}

	if r.Code != 0 {
		return r.Code
	}

	return genericcode.InternalServerError
}

func (r RichError) GetMessage() string {
	if r.WrapError != nil {
		wrapError, ok := r.WrapError.(RichError)
		if ok {
			return wrapError.GetMessage()
		}
	}

	if r.Message != "" {
		return r.Message
	}

	if r.WrapError != nil {
		return r.WrapError.Error()
	}

	return "internal server error"
}

func New(message string) RichError {
	return RichError{
		Code:    genericcode.InternalServerError,
		Message: message,
	}
}

func NewWithCode(code genericcode.Code, message string) RichError {
	return RichError{
		Code:    code,
		Message: message,
	}
}

func Wrap(err error, message string) RichError {
	return RichError{
		Message:   message,
		WrapError: err,
	}
}

func WrapWithCode(err error, code genericcode.Code, message string) RichError {
	return RichError{
		Code:      code,
		Message:   message,
		WrapError: err,
	}
}

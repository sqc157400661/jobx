package errors

import (
	"errors"
	"fmt"
)

type Error struct {
	Code    int32  `json:"code,omitempty"`
	Reason  string `json:"reason,omitempty"`
	Message string `json:"message,omitempty"`
}

func (e *Error) Error() string {
	if e.Reason == "" {
		return fmt.Sprintf("error: code = %d message = %s", e.Code, e.Message)
	}
	return fmt.Sprintf("error: code = %d reason = %s message = %s", e.Code, e.Reason, e.Message)
}

// New returns an error object for the code, message.
func New(code int, message string) *Error {
	e := &Error{
		Code:    int32(code),
		Message: message,
	}
	if v, has := ErrorReasonMap[code]; has {
		e.Reason = v
	}
	return e
}

// Newf New(code fmt.Sprintf(format, a...))
func Newf(code int, format string, a ...interface{}) *Error {
	return New(code, fmt.Sprintf(format, a...))
}

// Code returns the http code for an error.
// It supports wrapped errors.
func Code(err error) int {
	if err == nil {
		return 200
	}
	return int(FromError(err).Code)
}

// Reason returns the reason for a particular error.
// It supports wrapped errors.
func Reason(err error) string {
	if err == nil {
		return UnknownReason
	}
	return FromError(err).Reason
}

// FromError try to convert an error to *Error.
// It supports wrapped errors.
func FromError(err error) *Error {
	if err == nil {
		return nil
	}
	if se := new(Error); errors.As(err, &se) {
		return se
	}
	return New(UnknownCode, err.Error())
}

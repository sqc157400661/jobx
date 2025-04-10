package errors

import (
	"errors"
	"fmt"
)

type Error struct {
	Code    int32  `json:"code,omitempty"`
	Reason  string `json:"reason,omitempty"`
	Message string `json:"message,omitempty"`
	cause   error
}

func (e *Error) Error() string {
	errMsg := fmt.Sprintf("error: code = %d message = %s", e.Code, e.Message)
	if e.Reason != "" {
		errMsg += fmt.Sprintf(" reason = %s", e.Reason)
	}
	if e.cause != nil {
		errMsg += fmt.Sprintf(" cause = %v", e.cause)
	}
	return errMsg
}

// Is reports whether any error in err's chain matches target.
func (e *Error) Is(err error) bool {
	if err == nil {
		return false
	}
	targetErr := new(Error)
	if !errors.As(err, &targetErr) {
		return false
	}
	return targetErr.Code == e.Code
}

// Wrap warp a error msg
func (e *Error) Wrap(err error) *Error {
	if err == nil {
		return e
	}
	e.cause = err
	return e
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

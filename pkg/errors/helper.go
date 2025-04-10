package errors

import "fmt"

func IgnoreBIZExist(err error) error {
	if Code(err) == BIZExist {
		return nil
	}
	return err
}

func IsBIZExist(err error) bool {
	if Code(err) == BIZExist {
		return true
	}
	return false
}

func IgnoreTokenExist(err error) error {
	if Code(err) == TokenExist {
		return nil
	}
	return err
}

func IsTokenExist(err error) bool {
	if Code(err) == TokenExist {
		return true
	}
	return false
}

func IsWaitTimeout(err error) bool {
	if Code(err) == WaitJobTimeout {
		return true
	}
	return false
}

func IgnoreWaitTimeout(err error) error {
	if IsWaitTimeout(err) {
		return nil
	}
	return err
}

func IsQueueFull(err error) bool {
	if Code(err) == LocalQueueFullCode {
		return true
	}
	return false
}
func IsQueueEmpty(err error) bool {
	if Code(err) == ErrEmptyQueueCode {
		return true
	}
	return false
}

func BIZConflict(msg string) *Error {
	return New(BIZExist, msg)
}

func TokenConflict(msg string) *Error {
	return New(TokenExist, msg)
}

func NewNoTaskFoundErrorWithJobID(rootID, jobID int) *Error {
	return New(NoTaskFoundCode, fmt.Sprintf("rootID:%d job:%d", rootID, jobID))
}

func NewNoRunnablePipelineTaskErrorWithJobID(rootID, jobID int) *Error {
	return New(NoRunnablePipelineTaskCode, fmt.Sprintf("rootID:%d job:%d", rootID, jobID))
}

func NewParamError(msg string) *Error {
	return New(InvalidParameterCode, msg)
}

func NewDBError(msg string) *Error {
	return New(InternalDBError, msg)
}

func NewInternalError(msg string) *Error {
	return New(InternalError, msg)
}

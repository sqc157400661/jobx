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

func BIZConflict(msg string) error {
	return New(BIZExist, msg)
}

func TokenConflict(msg string) error {
	return New(TokenExist, msg)
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

func NewNoTaskFoundErrorWithJobID(rootID, jobID int) error {
	return New(NoTaskFoundCode, fmt.Sprintf("rootID:%d job:%d", rootID, jobID))
}

func NewNoRunnablePipelineTaskErrorWithJobID(rootID, jobID int) error {
	return New(NoRunnablePipelineTaskCode, fmt.Sprintf("rootID:%d job:%d", rootID, jobID))
}

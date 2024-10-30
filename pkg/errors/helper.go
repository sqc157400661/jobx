package errors

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

func NoProvider() error {
	return New(NoJobProvider, "no provider")
}

func BIZConflict(msg string) error {
	return New(BIZExist, msg)
}

func NotFoundDefJob() error {
	return New(NotFoundJobDefs, NotFoundJobDefsReason)
}

func TokenConflict(msg string) error {
	return New(TokenExist, msg)
}

func WaitTimeout() error {
	return New(WaitJobTimeout, WaitJobTimeoutReason)
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

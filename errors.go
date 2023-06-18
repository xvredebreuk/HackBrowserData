package hackbrowserdata

import (
	"errors"
)

var (
	ErrBrowserNotExists       = errors.New("browser not exists")
	ErrWrongSecurityCommand   = errors.New("wrong security command")
	ErrNoPasswordInOutput     = errors.New("no password in output")
	ErrCouldNotFindInKeychain = errors.New("could not be find in keychain")
)

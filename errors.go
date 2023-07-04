package hackbrowserdata

import (
	"errors"
)

var (
	ErrBrowserNotExists       = errors.New("browser not exists")
	ErrBrowserNotSupport      = errors.New("browser not support")
	ErrWrongSecurityCommand   = errors.New("wrong security command")
	ErrNoPasswordInOutput     = errors.New("no password in output")
	ErrCouldNotFindInKeychain = errors.New("could not be find in keychain")
	ErrBrowsingDataNotSupport = errors.New("browsing data not support")
	ErrBrowsingDataNotExists  = errors.New("browsing data not exists")
)

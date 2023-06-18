package hackbrowserdata

import (
	"path/filepath"
)

type BrowserOption func(browserOptionsSetter)

type browserOptionsSetter interface {
	setProfilePath(string)

	setDisableAllUsers(bool)

	setStorageName(string)
}

func WithProfilePath(p string) BrowserOption {
	return func(b browserOptionsSetter) {
		b.setProfilePath(filepath.Clean(p))
	}
}

func WithDisableAllUsers(e bool) BrowserOption {
	return func(b browserOptionsSetter) {
		b.setDisableAllUsers(e)
	}
}

func WithStorageName(s string) BrowserOption {
	return func(b browserOptionsSetter) {
		b.setStorageName(s)
	}
}

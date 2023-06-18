package hackbrowserdata

import (
	"path/filepath"
)

type BrowserOption func(browserOptionsSetter)

type browserOptionsSetter interface {
	setProfilePath(string)

	setEnableAllUsers(bool)

	setStorageName(string)
}

func WithProfilePath(p string) BrowserOption {
	return func(b browserOptionsSetter) {
		b.setProfilePath(filepath.Clean(p))
	}
}

func WithEnableAllUsers(e bool) BrowserOption {
	return func(b browserOptionsSetter) {
		b.setEnableAllUsers(e)
	}
}

func WithStorageName(s string) BrowserOption {
	return func(b browserOptionsSetter) {
		b.setStorageName(s)
	}
}

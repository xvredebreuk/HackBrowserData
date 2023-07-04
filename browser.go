package hackbrowserdata

import (
	"github.com/moond4rk/hackbrowserdata/browserdata"
)

type Browser interface {
	Passwords() ([]browserdata.Password, error)

	Cookies() ([]browserdata.Cookie, error)

	ExtractBrowserData(dataTypes []DataType) (map[DataType]interface{}, error)
}

func NewBrowser(b browser, options ...BrowserOption) (Browser, error) {
	opt, ok := defaultBrowserOptions[b]
	if !ok {
		return nil, ErrBrowserNotSupport
	}

	for _, options := range options {
		options(opt)
	}

	if opt.NewBrowserFunc == nil {
		return nil, ErrBrowserNotSupport
	}
	return opt.NewBrowserFunc(opt)
}

type browser string

const (
	Chrome   browser = "chrome"
	Firefox  browser = "firefox"
	Yandex   browser = "yandex"
	Edge     browser = "edge"
	Chromium browser = "chromium"
)

type browserType int

const (
	browserTypeChromium browserType = iota + 1
	browserTypeFirefox
	browserTypeYandex
)

func (b browser) Type() browserType {
	switch b {
	case Firefox:
		return browserTypeFirefox
	case Yandex:
		return browserTypeYandex
	default:
		return browserTypeChromium
	}
}

package hackbrowserdata

type Browser interface {
	BrowserData

	Init() error
}

func NewBrowser(b browser, options ...BrowserOption) (Browser, error) {
	browser := browsers[b]
	if setter, ok := browser.(browserOptionsSetter); ok {
		for _, option := range options {
			option(setter)
		}
	}
	if err := browser.Init(); err != nil {
		return nil, err
	}

	return browser, nil
}

type browser string

type BrowserData interface {
	Passwords() ([]Password, error)

	Cookies() ([]Cookie, error)
}

func (c *chromium) BrowsingData(items []browserDataType) ([]BrowserData, error) {
	for _, item := range items {
		_ = item
	}
	return nil, nil
}

func (c *chromium) AllBrowsingData() ([]BrowserData, error) {
	return nil, nil
}

func (f *firefox) BrowsingData(_ []browserDataType) (BrowserData, error) {
	return nil, nil
}

const (
	Chrome  browser = "chrome"
	Firefox browser = "firefox"
	Yandex  browser = "yandex"
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

var browsers = map[browser]Browser{
	Chrome: &chromium{
		name:          Chrome,
		storage:       chromeStorageName,
		profilePath:   chromeProfilePath,
		supportedData: []browserDataType{TypePassword},
	},
	Firefox: &firefox{
		name:          Firefox,
		profilePath:   firefoxProfilePath,
		supportedData: []browserDataType{TypePassword},
	},
	Yandex: &chromium{},
}

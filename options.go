package hackbrowserdata

type Options struct {
	Name            browser
	Storage         string
	ProfilePath     string
	IsEnableAllUser bool
	DataTypes       []DataType
	NewBrowserFunc  func(*Options) (Browser, error)
}

type BrowserOption func(*Options)

func WithBrowserName(p string) BrowserOption {
	return func(o *Options) {
		o.Name = browser(p)
	}
}

func WithProfilePath(p string) BrowserOption {
	return func(o *Options) {
		o.ProfilePath = p
	}
}

func WithEnableAllUsers(e bool) BrowserOption {
	return func(o *Options) {
		o.IsEnableAllUser = e
	}
}

func WithStorageName(s string) BrowserOption {
	return func(o *Options) {
		o.Storage = s
	}
}

package hackbrowserdata

type Cookie struct {
	Domain     string
	Expiration float64
	Value      string
}

func (c *chromium) Cookies() ([]Cookie, error) {
	return nil, nil
}

func (f *firefox) Cookies() ([]Cookie, error) {
	return nil, nil
}

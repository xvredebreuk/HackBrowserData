package hackbrowserdata

type firefox struct {
	name          string
	storage       string
	profilePath   string
	enableAllUser bool
	masterKey     []byte
}

func (f *firefox) Init() error {
	return nil
}

func (f *firefox) setEnableAllUsers(e bool) {
	f.enableAllUser = e
}

func (f *firefox) setProfilePath(p string) {
	f.profilePath = p
}

func (f *firefox) setStorageName(s string) {
	f.storage = s
}

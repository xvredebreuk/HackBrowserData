package hackbrowserdata

type browserDataType int

const (
	TypePassword browserDataType = iota + 1
	TypeCookie
	TypeHistory
	TypeBookmark
	TypeCreditCard
	TypeDownload
	TypeExtensions
	TypeSessionStorage
	TypeLocalStorage
)

func (i browserDataType) Filename(b browser) string {
	switch b.Type() {
	case browserTypeChromium:
		return i.chromiumFilename()
	case browserTypeFirefox:
		return i.firefoxFilename()
	case browserTypeYandex:
		return i.yandexFilename()
	}
	return ""
}

func (i browserDataType) chromiumFilename() string {
	switch i {
	case TypePassword:
		return "Login Data"
	case TypeCookie:
		return "Cookies"
	case TypeHistory:
	}
	return ""
}

func (i browserDataType) yandexFilename() string {
	switch i {
	case TypePassword:
		return "Login State"
	case TypeCookie:
		return "Cookies"
	case TypeHistory:
	}
	return ""
}

func (i browserDataType) firefoxFilename() string {
	switch i {
	case TypePassword:
		return "logins.json"
	case TypeCookie:
		return "cookies.sqlite"
	case TypeHistory:
	}
	return ""
}

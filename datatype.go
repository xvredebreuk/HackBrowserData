package hackbrowserdata

import (
	"github.com/moond4rk/hackbrowserdata/browserdata"
)

type DataType int

const (
	TypeMasterKey DataType = iota
	TypePassword
	TypeCookie
	TypeHistory
	TypeBookmark
	TypeCreditCard
	TypeDownload
	TypeExtensions
	TypeSessionStorage
	TypeLocalStorage
)

func (i DataType) NewExtractor(browserType browserType, masterKey []byte, datafiles []string) browserdata.Extractor {
	switch i {
	case TypePassword:
		switch browserType {
		case browserTypeChromium:
			return browserdata.NewPassExtractor(masterKey, datafiles, browserdata.DefaultDBHandler, browserdata.ChromiumPassRowsHandler)
		case browserTypeFirefox:
			return browserdata.NewPassExtractor(masterKey, datafiles, browserdata.DefaultJSONHandler, browserdata.FirefoxPassRowsHandler)
		}
	case TypeCookie:
	}
	return nil
}

var (
	defaultDataTypes = []DataType{TypePassword, TypeCookie}
)

const unsupportedType = ""

func (i DataType) Filename(b browser) string {
	switch b.Type() {
	case browserTypeChromium:
		return i.chromiumFilename()
	case browserTypeFirefox:
		return i.firefoxFilename()
	case browserTypeYandex:
		return i.yandexFilename()
	default:
		return unsupportedType
	}
}

func (i DataType) chromiumFilename() string {
	switch i {
	case TypeMasterKey:
		return fileChromiumKey
	case TypePassword:
		return fileChromiumPassword
	case TypeCookie:
		return fileChromiumCookie
	case TypeHistory:
		return fileChromiumHistory
	case TypeBookmark:
		return fileChromiumBookmark
	case TypeCreditCard:
		return fileChromiumCredit
	case TypeDownload:
		return fileChromiumDownload
	case TypeExtensions:
		return fileChromiumExtension
	case TypeSessionStorage:
		return fileChromiumSessionStorage
	case TypeLocalStorage:
		return fileChromiumLocalStorage
	default:
		return unsupportedFile
	}
}

func (i DataType) yandexFilename() string {
	switch i {
	case TypePassword:
		return fileYandexPassword
	case TypeCreditCard:
		return fileYandexCredit
	default:
		return i.chromiumFilename()
	}
}

func (i DataType) firefoxFilename() string {
	switch i {
	case TypeMasterKey:
		return fileFirefoxMasterKey
	case TypePassword:
		return fileFirefoxPassword
	case TypeCookie:
		return fileFirefoxCookie
	case TypeHistory:
		return fileFirefoxData
	case TypeBookmark:
		return fileFirefoxData
	case TypeCreditCard:
		// Firefox does not store credit cards
		return unsupportedFile
	case TypeDownload:
		return fileFirefoxData
	case TypeExtensions:
		return fileFirefoxExtension
	case TypeSessionStorage:
		return fileFirefoxData
	case TypeLocalStorage:
		return fileFirefoxLocalStorage
	default:
		return unsupportedFile
	}
}

const unsupportedFile = "unsupported file"

const (
	fileChromiumKey            = "Local State"
	fileChromiumCredit         = "Web Data"
	fileChromiumPassword       = "Login Data"
	fileChromiumHistory        = "History"
	fileChromiumDownload       = "History"
	fileChromiumCookie         = "Cookies"
	fileChromiumBookmark       = "Bookmarks"
	fileChromiumLocalStorage   = "Local Storage/leveldb"
	fileChromiumSessionStorage = "Session Storage"
	fileChromiumExtension      = "Extensions"

	fileYandexPassword = "Ya Passman Data"
	fileYandexCredit   = "Ya Credit Cards"

	fileFirefoxMasterKey    = "key4.db"
	fileFirefoxCookie       = "cookies.sqlite"
	fileFirefoxPassword     = "logins.json"
	fileFirefoxData         = "places.sqlite"
	fileFirefoxLocalStorage = "webappsstore.sqlite"
	fileFirefoxExtension    = "extensions.json"
)

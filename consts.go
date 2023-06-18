package hackbrowserdata

import (
	"os"
)

const (
	chromeStorageName     = "Chrome"
	chromeBetaStorageName = "Chrome"
	chromiumStorageName   = "Chromium"
	edgeStorageName       = "Microsoft Edge"
	braveStorageName      = "Brave"
	operaStorageName      = "Opera"
	vivaldiStorageName    = "Vivaldi"
	coccocStorageName     = "CocCoc"
	yandexStorageName     = "Yandex"
	arcStorageName        = "Arc"
)

var (
	homeDir, _ = os.UserHomeDir()
)

var (
	chromeProfilePath     = homeDir + "/Library/Application Support/Google/Chrome/Default/"
	chromeBetaProfilePath = homeDir + "/Library/Application Support/Google/Chrome Beta/Default/"
	chromiumProfilePath   = homeDir + "/Library/Application Support/Chromium/Default/"
	edgeProfilePath       = homeDir + "/Library/Application Support/Microsoft Edge/Default/"
	braveProfilePath      = homeDir + "/Library/Application Support/BraveSoftware/Brave-Browser/Default/"
	operaProfilePath      = homeDir + "/Library/Application Support/com.operasoftware.Opera/Default/"
	operaGXProfilePath    = homeDir + "/Library/Application Support/com.operasoftware.OperaGX/Default/"
	vivaldiProfilePath    = homeDir + "/Library/Application Support/Vivaldi/Default/"
	coccocProfilePath     = homeDir + "/Library/Application Support/Coccoc/Default/"
	yandexProfilePath     = homeDir + "/Library/Application Support/Yandex/YandexBrowser/Default/"
	arcProfilePath        = homeDir + "/Library/Application Support/Arc/User Data/Default"

	firefoxProfilePath = homeDir + "/Library/Application Support/Firefox/Profiles/"
)

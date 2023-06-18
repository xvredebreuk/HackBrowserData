package hackbrowserdata

import (
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"time"

	// import sqlite3 driver
	_ "github.com/mattn/go-sqlite3"

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
		return nil, errors.New("password for c.name is not supported")
	}
	var fullPass []Password
	for _, profile := range c.profilePaths {
		passFile := filepath.Join(profile, TypePassword.Filename(c.name))
		if !fileutil.IsFileExists(passFile) {
			fmt.Println(passFile)
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
	if f.masterKey != nil {
		return nil, nil
	}
	return nil, nil
}

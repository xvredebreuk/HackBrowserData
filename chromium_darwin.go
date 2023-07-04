package hackbrowserdata

import (
	"bytes"
	"crypto/sha1"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

var (
	salt = []byte("saltysalt")
)

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

	key := pbkdf2.Key(secret, salt, 1003, 16, sha1.New)
	if len(key) == 0 {
		return ErrWrongSecurityCommand
	}
	c.masterKey = key
	return nil
}

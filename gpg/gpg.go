// Package gpg wraps some functionality of the gpg CLI program.
package gpg

import (
	"fmt"
	"os/exec"
	"strings"
)

// Installed returns true if gpg is installed.
func Installed() bool {
	_, err := exec.LookPath("gpg")
	return err == nil
}

// Version returns the version of gpg currently installed.
func Version() (string, error) {
	out, err := exec.Command("gpg", "--version").Output()
	if err != nil {
		return "", err
	} else if len(out) == 0 {
		return "", fmt.Errorf("gpg: output empty")
	}

	if parts := strings.SplitN(string(out), "\n", 2); len(parts) > 1 {
		return parts[0], nil
	}

	return "", fmt.Errorf("gpg: failed to get version")
}

// Fingerprint returns the fingerprint for the given key ID. It shells out to gpg to do this.
func Fingerprint(keyID string) (string, error) {
	out, err := exec.Command("gpg", "--with-colons", "--fingerprint", keyID).Output()
	if err != nil {
		return "", err
	} else if len(out) == 0 {
		return "", fmt.Errorf("gpg: output empty")
	}

	// The command should return 6 lines of output, and the third should contain the fingerprint. The
	// third line should start with "fpr", and the fingerprint should be the tenth colon-separated field.
	lines := strings.SplitN(string(out), "\n", 4)
	if len(lines) >= 3 && strings.HasPrefix(lines[2], "fpr") {
		fields := strings.Split(lines[2], ":")
		if len(fields) >= 10 && len(fields[9]) > 0 {
			return fields[9], nil
		}
	}

	return "", fmt.Errorf("gpg: failed to get fingerprint")
}

// ExportPrivateKey returns the private key for the given key ID. It shells out to gpg to do this.
func ExportPrivateKey(keyID string) ([]byte, error) {
	out, err := exec.Command("gpg", "--export-secret-keys", keyID).Output()
	if err != nil {
		return []byte{}, err
	} else if len(out) == 0 {
		return []byte{}, fmt.Errorf("gpg: output empty")
	}

	return out, nil
}

// EncryptAES256 encrypts the given data using AES256. To do this it shells out to gpg, which will
// prompt for a passphrase.
func EncryptAES256(plaintext []byte) ([]byte, error) {
	cmd := exec.Command("gpg", "--cipher-algo", "AES256", "--symmetric")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return []byte{}, err
	}

	go func() {
		defer stdin.Close()
		stdin.Write(plaintext)
	}()

	out, err := cmd.Output()
	if err != nil {
		return []byte{}, err
	} else if len(out) == 0 {
		return []byte{}, fmt.Errorf("gpg: output empty")
	}

	return out, nil
}

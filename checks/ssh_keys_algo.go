package checks

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type KeyType string

const (
	ED25519    KeyType = "ED25519"
	ED25519_SK KeyType = "ED25519-SK"
	ECDSA      KeyType = "ECDSA"
	ECDSA_SK   KeyType = "ECDSA-SK"
	DSA        KeyType = "DSA"
	RSA        KeyType = "RSA"
)

type KeyInfo struct {
	strength  int
	signature string
	keyType   KeyType
}

type SSHKeysAlgo struct {
	passed  bool
	sshKey  string
	sshPath string
}

// Name returns the name of the check
func (f *SSHKeysAlgo) Name() string {
	return "SSH keys have sufficient algorithm strength"
}

func parseKeyInfo(output string) KeyInfo {
	parts := strings.Fields(strings.TrimSpace(output))
	if len(parts) < 4 {
		return KeyInfo{}
	}

	strength, _ := strconv.Atoi(parts[0])
	return KeyInfo{
		strength:  strength,
		signature: parts[1],
		keyType:   KeyType(strings.ToUpper(parts[len(parts)-1])),
	}
}

func (f *SSHKeysAlgo) isKeyStrong(path string) bool {
	output, err := exec.Command("ssh-keygen", "-l", "-f", path).Output()
	if err != nil {
		return false
	}

	info := parseKeyInfo(string(output))
	switch info.keyType {
	case RSA:
		return info.strength >= 2048
	case DSA:
		return info.strength >= 8192
	case ECDSA, ECDSA_SK:
		return info.strength >= 521
	case ED25519, ED25519_SK:
		return info.strength >= 256
	default:
		return false
	}
}

// Run executes the check
func (f *SSHKeysAlgo) Run() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	f.sshPath = filepath.Join(home, ".ssh")
	entries, err := os.ReadDir(f.sshPath)
	if err != nil {
		return err
	}

	f.passed = true
	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), ".pub") {
			continue
		}

		pubPath := filepath.Join(f.sshPath, entry.Name())
		privPath := strings.TrimSuffix(pubPath, ".pub")

		if _, err := os.Stat(privPath); os.IsNotExist(err) {
			continue
		}

		if !f.isKeyStrong(pubPath) {
			f.passed = false
			f.sshKey = strings.TrimSuffix(entry.Name(), ".pub")
			break
		}
	}

	return nil
}

// Passed returns the status of the check
func (f *SSHKeysAlgo) Passed() bool {
	return f.passed
}

// CanRun returns whether the check can run
func (f *SSHKeysAlgo) IsRunnable() bool {
	home, err := os.UserHomeDir()
	if err != nil {
		return false
	}

	sshPath := filepath.Join(home, ".ssh")
	if _, err := os.Stat(sshPath); os.IsNotExist(err) {
		return false
	}

	//check if there are any private keys in the .ssh directory
	files, err := os.ReadDir(sshPath)
	if err != nil {
		return false
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".pub") {
			privateKeyPath := filepath.Join(sshPath, strings.TrimSuffix(file.Name(), ".pub"))
			if _, err := os.Stat(privateKeyPath); err == nil {
				return true
			}
		}
	}
	return false
}

// UUID returns the UUID of the check
func (f *SSHKeysAlgo) UUID() string {
	return "ef69f752-0e89-46e2-a644-310429ae5f45"
}

// PassedMessage returns the message to return if the check passed
func (f *SSHKeysAlgo) PassedMessage() string {
	return "SSH keys use strong encryption"
}

// FailedMessage returns the message to return if the check failed
func (f *SSHKeysAlgo) FailedMessage() string {
	return "SSH keys are using weak encryption"
}

// RequiresRoot returns whether the check requires root access
func (f *SSHKeysAlgo) RequiresRoot() bool {
	return false
}

// Status returns the status of the check
func (f *SSHKeysAlgo) Status() string {
	if f.Passed() {
		return f.PassedMessage()
	}

	return "SSH key " + f.sshKey + " is using weak encryption"
}

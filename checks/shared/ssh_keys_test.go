package shared

import (
	"path/filepath"
	"testing"

	sharedG "github.com/ParetoSecurity/pareto-core/shared"
)

const (
	// A valid unencrypted RSA private key for testing purposes.
	unencryptedPrivateKey = `-----BEGIN RSA PRIVATE KEY-----
MIICXgIBAAKBgQDTD2FKNdbA8KCMK0AevYSCzrM9hCjGj6aYAQH2SI2ZCbfYKbj3
eXEM4m2XWxUvnBETHCVyg4c99phXgheaewWv1zJjuAnophD2WzcNUdo6Db+bfO75
h9/AkOHSWXbJSt/m1PUQvxwje2Vwf6YOkjd82tM1rt2EbOCCK/knjJLH3wIDAQAB
AoGAMOmbjmszvbsGOfW8AmPBVd85QsRh/sJDxW5WWhEuX40VAg+JQjDutiGzbCQ7
oLD2dAtN0mAQ85c2bvFDLxXblxz0JJn6Gq2D7EoIOiFLTgXv5JlZGEWNRu8nnf0Q
vDwkyFt6TraaFAIKqv7y/lmmK3CFgb5NlWARsLq+Rg7bByECQQDrvpw+iqAr+hO+
/lmg1sSg72HmllCXRgApX0k5RZXb1YLgxanTNeZ1yYLj8QvsQUOHWkXOcnvqshGE
mcv0gRUjAkEA5THQinEdEFQS2edSPdVJT3PhyIwCBtLZelhPd/8m8iiWPvpwaaCd
gBGyP/rnmfO2AsCz8SZnRodsVN19fPKEFQJBAIBxQvnEV85+G2IVfMnoGgvBQWr7
/P7esdnYw7GDm0nCQ+OpboTYOi900m7U93UKfftyENSRhbhyup6vmPMnnVcCQQDN
RThJRdWJ8kKP9qWpy4TFLDxjqUGHawBsmvtRtavj5oXqEdLsR3XIZhEHTGhxcdzp
yj1fFc4ZVOCpgVYKugmhAkEA2o+Je6TdKlo7P4jNIFPbQmUd9+Y55BBX7Hn6oXeL
V2VqSeaNZGgMuquMF6G0FtIvpkxQ4K5Wrq07mRWuBwLbuw==
-----END RSA PRIVATE KEY-----`

	// An invalid key content to simulate an encrypted or malformed key.
	invalidKey = "this-is-not-a-valid-key"
)

func TestHasPassword(t *testing.T) {
	// Create a temporary directory for test files.
	tmpDir := t.TempDir()

	s := &SSHKeys{}

	t.Run("NonExistentFile", func(t *testing.T) {
		// Provide a file path that does not exist.
		nonExistent := filepath.Join(tmpDir, "nonexistent")
		// Expect true since ReadFile will fail.
		if got := s.hasPassword(nonExistent); got != true {
			t.Errorf("hasPassword() = %v; want true", got)
		}
	})

	t.Run("ValidUnencryptedKey", func(t *testing.T) {
		sharedG.ReadFileMocks = map[string]string{
			"unencrypted": unencryptedPrivateKey,
		}
		// Expect false because the key is unencrypted (no password).
		if got := s.hasPassword("unencrypted"); got != false {
			t.Errorf("hasPassword() = %v; want false", got)
		}
	})

	t.Run("InvalidKeyContent", func(t *testing.T) {
		sharedG.ReadFileMocks = map[string]string{
			"invalid": invalidKey,
		}
		// Expect true because parsing will fail.
		if got := s.hasPassword("invalid"); got != true {
			t.Errorf("hasPassword() = %v; want true", got)
		}
	})
}

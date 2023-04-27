package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidatePasswordErr(t *testing.T) {
	assert.Error(t, ValidatePassword("hunter", "hunter2"))
}

func TestValidatePasswordOk(t *testing.T) {
	assert.Error(t, ValidatePassword("hunter2", "hunter2"))
}

func TestHashPassWordOk(t *testing.T) {
	//hunter2 with 10 rounds of Bcrypt
	password := "hunter2"
	hashedPass, _ := HashPassword(password)

	err := ValidatePassword(password, hashedPass)
	assert.NoError(t, err)
}

func TestHashPasswordTooLongError(t *testing.T) {
	password := "huuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuunter2"
	_, err := HashPassword(password)
	assert.Error(t, err)
}

func TestNewJwtTokerOK(t *testing.T) {
	_, err := NewJwtToker(getTestingPrivateRSAKey(), getTestingPublicRSAKey())
	assert.NoError(t, err)
}

func TestNewJwtTokerPrivKeyError(t *testing.T) {
	_, err := NewJwtToker([]byte("wrong token"), getTestingPublicRSAKey())
	assert.Error(t, err)
}

func TestNewJwtTokerPubKeyError(t *testing.T) {
	_, err := NewJwtToker(getTestingPublicRSAKey(), []byte("wrong token"))
	assert.Error(t, err)
}

func TestCreateTokenOk(t *testing.T) {
	toker, _ := NewJwtToker(getTestingPrivateRSAKey(), getTestingPublicRSAKey())
	token, err := toker.CreateToken("h014", true)
	assert.NoError(t, err)
	assert.NotEqual(t, len(token), 0)
}

func getTestingPrivateRSAKey() []byte {
	return []byte(`-----BEGIN RSA PRIVATE KEY-----
MIICXgIBAAKBgQCLVVr5ESP3ThFsQsoHL4yhiRA3oeLAu+F+SmRmaQUXboNB+sdY
9sZOP9bQro7KCCp1k4kI8EwtecigWzqv01x7Wr77Ul9RyNI1/ZFUbZf+gDAIR6Iq
N8ElQBbppCQicE75iLSIh8no6GHXCJj1v875iozHa2hqYKF6/fFGac29cQIDAQAB
AoGAF1Zfm3IchRKlZm21awiy1GehuL+7vC578WxCbsjOWoNfJtD7TNJgmsCkmWVz
czF08yaYAFBHYiKQ0RMWvFZ5mcxgmtSBfmOPyfwrb1OtdJwDx/SACGMcv3noxj4n
uxH9204Lm9XsrKgRpJ8NcgWypwhY+EB0Pk4nekeLlllfpXECQQDFODQJNfqmOo7Q
672jUR7Q3Urxp651Wa7DQLSBAV7bwl+Trn7RBa4WvHMhRU+rbiG9gWB2JFhyc0ia
kybtVAJ1AkEAtNx0712/B/SBpq2rhTQj0bySvodq9uvGdwwToDuek1aU9/bP0Utj
N72U9LA7uvKttKZz6jkpRBQdAHbem1R3jQJBALBYGO9DfOO16I2WvPKTTmKj/Kcn
sC7uCf48lSnk99S4cI20sWBlG8zopGlTeHFpAHJahM4eoZd0za6pdV0wiSECQQCr
+ZjjZwPX36JMyIT48zxAGgx7SS7nvggIeR5MVYSS21hpdHHltMaSYR27kbwqJsoP
pdtA07uudWWiZGWF08qdAkEAjV6m2VA+esla4x/SrnWMN0EHIqGXhayulGyQ65M+
SEnUZW5g7xBlh3AGSx+SUyAo7dCSL7vClt96H+dbkeNOxg==
-----END RSA PRIVATE KEY-----`)
}

func getTestingPublicRSAKey() []byte {
	return []byte(`-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCLVVr5ESP3ThFsQsoHL4yhiRA3
oeLAu+F+SmRmaQUXboNB+sdY9sZOP9bQro7KCCp1k4kI8EwtecigWzqv01x7Wr77
Ul9RyNI1/ZFUbZf+gDAIR6IqN8ElQBbppCQicE75iLSIh8no6GHXCJj1v875iozH
a2hqYKF6/fFGac29cQIDAQAB
-----END PUBLIC KEY-----`)
}

package keychain

import (
	"crypto/sha256"
	"fmt"
	"github.com/rahul0tripathi/clipd/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKeychainService_GetPassword(t *testing.T) {
	l, err := util.NewLogger()
	if err != nil {
		t.Error(err)
		return
	}
	service, err := NewKeychainService(l)
	if !assert.NoError(t, err, "failed to get new keychain service") {
		return
	}
	password := sha256.Sum256([]byte("hello world"))
	passHex := fmt.Sprintf("%x", password[:])
	err = service.CreateNewSecret([]byte(passHex), nil)
	if err != nil && err != ErrSecretAlreadyExists {
		t.Error(err, "failed to call CreateNewSecret")
		return
	}
	getPassword, err := service.GetPassword(nil)
	if !assert.NoError(t, err, "failed to call GetPassword") {
		return
	}
	fmt.Println(getPassword)
	t.Logf(passHex, string(getPassword))
	assert.Equal(t, passHex, string(getPassword), "passwords do not match")
}

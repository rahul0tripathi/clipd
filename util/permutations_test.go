package util

import (
	"crypto/sha256"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestGetRandomNPermutations(t *testing.T) {
	perms, err := GetRandomNPermutations(500)
	if err != nil {
		t.Error(err, "failed to call GetRandomNPermutations")
		return
	}
	for k, v := range perms.RandomPermutations {
		t.Log("key", k, "\n", "value", fmt.Sprintf("%v", v))
	}
}
func PersistNew(t *testing.T, l *zap.Logger, readOnly bool) ([]string, StorageEngine, error) {
	db, err := NewLevelDB("./misc/leveldb", readOnly, l)
	if err != nil {
		t.Error(err)
		return nil, nil, err
	}
	perms, err := GetRandomNPermutations(2)
	if err != nil {
		t.Error(err, "failed to call GetRandomNPermutations")
		return nil, nil, err
	}
	err = PersistPermutationToStorage(perms, db)
	return perms.Metadata, db, nil
}
func TestPersistPermutationToStorage(t *testing.T) {
	l, err := NewLogger()
	if err != nil {
		t.Error(err)
		return
	}
	_, _, err = PersistNew(t, l, false)
	assert.NoError(t, err, "failed to call PersistPermutationToStorage")

}

func TestNewDefaultPermutationService(t *testing.T) {
	l, err := NewLogger()
	if err != nil {
		t.Error(err)
		return
	}
	db, err := NewLevelDB("./misc/leveldb", true, l)
	if err != nil {
		t.Error(err)
		return
	}
	service, err := NewDefaultPermutationService(db, l)
	if !assert.NoError(t, err, "failed to call NewDefaultPermutationService") {
		return
	}
	key := sha256.Sum256([]byte("helloworld"))
	t.Logf("Existing Key: %x\n", key)
	withKey, err := service.PermuteRandomly(key)
	if !assert.NoError(t, err, "failed to call PermuteWithKey") {
		return
	}
	t.Logf("PremutedKey: %x\n", withKey.Data)

}

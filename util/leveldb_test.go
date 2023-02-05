package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewLevelDBRW(t *testing.T) {
	l, err := NewLogger()
	if err != nil {
		t.Error(err)
		return
	}
	db, err := NewLevelDB("./misc/leveldb", false, l)
	if err != nil {
		t.Error(err)
		return
	}
	err = db.Put("test", []byte("hello world"))
	if err != nil {
		return
	}
	if err != nil {
		t.Error(err)
		return
	}
	val, err := db.Get("test")

	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, "hello world", string(val), "Saved value does not match")
}

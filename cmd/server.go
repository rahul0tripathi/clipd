package main

import (
	"crypto/sha256"
	"fmt"
	"github.com/rahul0tripathi/clipd/pkg/server"
	"github.com/rahul0tripathi/clipd/util"
)

func start(salt string, keyPath string) {
	l, err := util.NewLogger()
	if err != nil {
		panic(err)
	}
	hash := sha256.Sum256([]byte(salt))
	localServer, err := server.NewLocalServer(fmt.Sprintf("%x", hash[:]), keyPath, l)
	if err != nil {
		panic(err)
	}
	err = localServer.Listen()
	if err != nil {
		panic(err)
	}
}

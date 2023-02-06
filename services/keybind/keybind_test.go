package keybind

import (
	"github.com/rahul0tripathi/clipd/util"
	"github.com/stretchr/testify/assert"
	"golang.design/x/hotkey/mainthread"
	"testing"
)

func DefaultKeybindingService_GetDecryptListener() {
	l, err := util.NewLogger()
	if err != nil {
		panic(err)
		return
	}
	service, err := NewDefaultKeybindService(l)
	if !assert.NoError(nil, err, "failed to call NewDefaultKeybindService") {
		return
	}
	listener, err := service.GetDecryptListener()
	if !assert.NoError(nil, err, "failed to call GetDecryptListener") {
		return
	}
	for {
		select {
		case <-listener:
			l.Info("key called")
			return
		}
	}
}
func TestDefaultKeybindingService_GetDecryptListener(t *testing.T) {
	mainthread.Init(DefaultKeybindingService_GetDecryptListener)
}

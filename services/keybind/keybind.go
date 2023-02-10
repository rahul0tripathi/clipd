package keybind

import (
	"errors"
	"go.uber.org/zap"
	"golang.design/x/hotkey"
)

var (
	ErrHandlerFailed = errors.New("failed to reqister handler")
)

type Service interface {
	GetEncryptListener() (writeTo <-chan bool, err error)
	GetDecryptListener() (writeTo <-chan bool, err error)
}

type defaultKeybindingService struct {
	logger *zap.Logger
}

func (k *defaultKeybindingService) attachListener(key hotkey.Key) (<-chan bool, error) {
	hk := hotkey.New([]hotkey.Modifier{hotkey.ModCmd, hotkey.ModShift}, key)
	err := hk.Register()
	if err != nil {
		k.logger.Fatal("hotkey: failed to register hotkey", zap.Error(err))
		return nil, ErrHandlerFailed
	}
	k.logger.Info("hotkey: registered")
	// create unbuffered chan
	write := make(chan bool)
	go func(logger *zap.Logger) {
		logger.Info("starting watcher")
		for {
			select {
			case <-hk.Keydown():
				select {
				case write <- true:
					continue
				default:
					logger.Error("failed to write to subscriber chan")
				}
			}
		}
	}(k.logger)
	return write, nil
}
func (k *defaultKeybindingService) GetEncryptListener() (writeTo <-chan bool, err error) {
	return k.attachListener(hotkey.KeyE)
}
func (k *defaultKeybindingService) GetDecryptListener() (writeTo <-chan bool, err error) {
	return k.attachListener(hotkey.KeyD)
}
func NewDefaultKeybindService(logger *zap.Logger) (Service, error) {
	service := &defaultKeybindingService{logger: logger}
	return service, nil
}

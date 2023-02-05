package keybind

type Service interface {
	Init() error
	AddHandler(keySequence string, writeTo <-chan string)
}

type defaultKeybindingService struct {
}

func NewDefaultKeybindService() (Service, error) {
	return nil, nil
}

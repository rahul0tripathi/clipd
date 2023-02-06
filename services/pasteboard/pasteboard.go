package pasteboard

import (
	"github.com/syndtr/goleveldb/leveldb/errors"
	"golang.design/x/clipboard"
)

var (
	ErrPasteboardServiceNotInitialised = errors.New("ErrPasterboardServiceNotInitalized")
)

type PasteboardService interface {
	ReadFromPasteboard() (data []byte, err error)
	WriteToPasteboard(data []byte) error
}

type defaultPasteboardService struct {
	isInit bool
}

func (d *defaultPasteboardService) ReadFromPasteboard() (data []byte, err error) {
	if !d.isInit {
		return nil, ErrPasteboardServiceNotInitialised
	}
	return clipboard.Read(clipboard.FmtText), nil
}

func (d defaultPasteboardService) WriteToPasteboard(data []byte) error {
	if !d.isInit {
		return ErrPasteboardServiceNotInitialised
	}
	// ignore write chan
	clipboard.Write(clipboard.FmtText, data)
	return nil
}

func NewDefaultPasteboardService() (PasteboardService, error) {
	service := &defaultPasteboardService{}
	err := clipboard.Init()
	if err != nil {
		return nil, err
	}
	service.isInit = true
	return service, nil
}

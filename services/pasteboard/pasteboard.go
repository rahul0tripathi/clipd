package pasteboard

type PasteboardService interface {
	Init() error
	ReadFromPasteboard() (data []byte, err error)
	WriteToPasteboard(data []byte) error
}

type defaultPasteboardService struct {
}

func NewDefaultPasteboardService() (PasteboardService, error) { return nil, nil }

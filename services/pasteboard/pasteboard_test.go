package pasteboard

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefaultPasteboardService_WriteToPasteboard(t *testing.T) {
	service, err := NewDefaultPasteboardService()
	if !assert.NoError(t, err, "failed to create NewDefaultPasteboardService") {
		return
	}
	err = service.WriteToPasteboard([]byte("hello"))
	if !assert.NoError(t, err, "failed to call WriteToPasteboard") {
		return
	}
	data, err := service.ReadFromPasteboard()
	assert.NoError(t, err, "failed to call ReadFromPasteboard")
	assert.Equal(t, data, []byte("hello"), "clipboard contents do not match")
}

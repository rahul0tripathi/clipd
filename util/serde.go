package util

import (
	"encoding/base64"
	"encoding/json"
	"github.com/syndtr/goleveldb/leveldb/errors"
)

func MarshalAndEncode(data interface{}) ([]byte, error) {
	metadataPlain, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return []byte(base64.StdEncoding.EncodeToString(metadataPlain)), nil
}

func DecodeAndUnmarshal(raw []byte, to interface{}) error {
	base64DecodedMetadata, err := base64.StdEncoding.DecodeString(string(raw))
	if err != nil {
		return err
	}
	if len(base64DecodedMetadata) == 0 {
		return errors.New("empty decode")
	}
	return json.Unmarshal(base64DecodedMetadata, to)
}

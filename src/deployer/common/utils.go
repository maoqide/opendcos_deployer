package common

import (
	"encoding/base64"
)

const (
	base64Table = "ABCDEFGHIJKLMNOPQRSTpqrstuvwxyz0123456789+/UVWXYZabcdefghijklmno"
)

var coder = base64.NewEncoding(base64Table)

func Base64Encode(src []byte) []byte {
	return []byte(coder.EncodeToString(src))
}

func Base64Decode(src []byte) ([]byte, error) {
	return coder.DecodeString(string(src))
}

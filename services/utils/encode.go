package utils

import (
	"encoding/base64"
	"encoding/hex"
)

func FromBase64ToStr(b64 string) (string, error) {
	res, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func FromBase64ToBytes(b64 string) ([]byte, error) {
	res, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func ToHex(s string) string {
	hex := hex.EncodeToString([]byte(s))
	result := ""
	for i := 2; i <= len(hex); i += 2 {
		result += "\\x" + hex[i-2:i]
	}
	return result
}

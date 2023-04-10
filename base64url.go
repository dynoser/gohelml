package gohelml

import (
	"encoding/base64"
	"strings"
)

func (h *HELML) base64Uencode(str string) string {
	base64Str := base64.URLEncoding.EncodeToString([]byte(str))
	return strings.TrimRight(base64Str, "=")
}

func (h *HELML) base64Udecode(str string) (string, error) {
	for len(str)%4 != 0 {
		str += "="
	}

	data, err := base64.URLEncoding.DecodeString(str)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

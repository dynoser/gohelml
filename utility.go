package gohelml

import (
	"strings"
)

func (h *HELML) stripcslashes(str string) string {
	controlCharsMap := map[string]string{
		"\\n":  "\n",
		"\\t":  "\t",
		"\\r":  "\r",
		"\\b":  "\b",
		"\\f":  "\f",
		"\\v":  "\v",
		"\\0":  "\x00",
		"\\\\": "\\",
	}
	for k, v := range controlCharsMap {
		str = strings.Replace(str, k, v, -1)
	}
	return str
}

func keys(m map[string]struct{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

package gohelml

import (
	"encoding/base64"
	"fmt"
	"math"
	"regexp"
	"strconv"
)

func (h *HELML) ValueEncoder(value interface{}, spcCh string) (string, error) {
	switch v := value.(type) {
	case string:
		var goodChars bool
		if spcCh == "_" {
			goodChars = regexp.MustCompile("^[ -}]+$").MatchString(v)
		} else {
			goodChars = regexp.MustCompile("^[^\x00-\x1F\x7E-\xFF]+$").MatchString(v)
		}
		if !goodChars || len(v) == 0 {
			return "-" + base64.URLEncoding.EncodeToString([]byte(v)), nil
		} else if spcCh == string(v[0]) || spcCh == string(v[len(v)-1]) || " " == string(v[len(v)-1]) {
			return "'" + v + "'", nil
		} else {
			return spcCh + v, nil
		}
	case bool:
		if v {
			return spcCh + spcCh + "T", nil
		} else {
			return spcCh + spcCh + "F", nil
		}
	case nil:
		return spcCh + spcCh + "U", nil
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return spcCh + spcCh + fmt.Sprintf("%d", v), nil
	case float32, float64:
		f := float64(v.(float64))
		switch {
		case math.IsInf(f, 1):
			return spcCh + spcCh + "INF", nil
		case math.IsInf(f, -1):
			return spcCh + spcCh + "NIF", nil
		case math.IsNaN(f):
			return spcCh + spcCh + "NAN", nil
		default:
			if spcCh == "_" && f != math.Trunc(f) {
				return "-" + base64.URLEncoding.EncodeToString([]byte(strconv.FormatFloat(f, 'f', -1, 64))), nil
			}
			return spcCh + spcCh + strconv.FormatFloat(f, 'f', -1, 64), nil
		}
	default:
		return "", fmt.Errorf("cannot encode value of type %T", value)
	}
}

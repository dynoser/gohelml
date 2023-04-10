package gohelml

import (
	"strconv"
)

func (h *HELML) ValueDecoder(encodedValue string, spc_ch string) interface{} {
	fc := encodedValue[0]
	if spc_ch[0] == fc {
		if encodedValue[:2] != spc_ch+spc_ch {
			return encodedValue[1:]
		}
		slicedValue := encodedValue[2:]
		if val, ok := SPEC_TYPE_VALUES[slicedValue]; ok {
			return val
		}
		if floatValue, err := strconv.ParseFloat(slicedValue, 64); err == nil {
			intValue := int(floatValue)
			if floatValue == float64(intValue) {
				return intValue
			}
			return floatValue
		}
		if h.CUSTOM_FORMAT_DECODER != nil {
			return h.CUSTOM_FORMAT_DECODER(encodedValue, spc_ch)
		}
		return slicedValue
	} else if fc == '"' || fc == '\'' {
		encodedValue = encodedValue[1 : len(encodedValue)-1]
		if fc == '"' {
			return h.stripcslashes(encodedValue)
		}
		return encodedValue
	} else if fc == '-' {
		encodedValue = encodedValue[1:]
	} else if h.CUSTOM_FORMAT_DECODER != nil {
		return h.CUSTOM_FORMAT_DECODER(encodedValue, spc_ch)
	}
	decoded, _ := h.base64Udecode(encodedValue)
	return decoded
}

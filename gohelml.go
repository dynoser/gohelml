package gohelml

import (
	"math"
)

type HELML struct {
	CUSTOM_FORMAT_DECODER func(value string, spc_ch string) interface{}
	CUSTOM_VALUE_DECODER  func(value string, spc_ch string) interface{}
	ADD_LAYERS_KEY        bool `default:"true"`
}

var SPEC_TYPE_VALUES = map[string]interface{}{
	"N":   nil,
	"U":   (*interface{})(nil),
	"T":   true,
	"F":   false,
	"NAN": math.NaN(),
	"INF": math.Inf(1),
	"NIF": math.Inf(-1),
}

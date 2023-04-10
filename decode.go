package gohelml

import (
	"strconv"
	"strings"
)

func (h *HELML) Decode(src_rows string, get_layers *[]string) interface{} {
	valueDecoFun := h.CUSTOM_VALUE_DECODER
	if valueDecoFun == nil {
		valueDecoFun = h.ValueDecoder
	}

	layer_init := "0"
	layer_curr := layer_init
	all_layers := map[string]struct{}{"0": {}}

	if get_layers == nil || len(*get_layers) == 0 {
		get_layers = &[]string{"0"}
	}

	layers := *get_layers
	layers_list := map[string]struct{}{}
	for _, layer := range layers {
		layers_list[layer] = struct{}{}
	}

	lvl_ch := ":"
	spc_ch := " "
	exploder_ch := "\n"

	for _, ch := range []string{"\n", "~", "\r"} {
		if strings.Contains(src_rows, ch) {
			exploder_ch = ch
			if ch == "~" && strings.HasSuffix(src_rows, "~") {
				lvl_ch = "."
				spc_ch = "_"
			}
			break
		}
	}

	str_arr := strings.Split(src_rows, exploder_ch)

	result := make(map[string]interface{})
	stack := []string{}

	min_level := -1

	for _, line := range str_arr {
		line = strings.TrimSpace(line)

		if len(line) == 0 || line[0] == '#' {
			continue
		}

		level := 0
		for level < len(line) && line[level] == lvl_ch[0] {
			level++
		}

		if level > 0 {
			line = line[level:]
		}

		haveDot := strings.Index(line, lvl_ch)
		key := ""
		value := ""
		if haveDot == -1 {
			key = line
			haveDot = 0
		} else {
			key = line[:haveDot]
			value = line[haveDot+1:]
			haveDot = 1
		}

		if min_level < 0 || min_level > level {
			min_level = level
		}

		extra_keys_cnt := len(stack) - (level - min_level)
		if extra_keys_cnt > 0 {
			stack = stack[:len(stack)-extra_keys_cnt]
			layer_curr = layer_init
		}

		var parent interface{} = result
		for _, parentKey := range stack {
			switch parentVal := parent.(type) {
			case map[string]interface{}:
				if childVal, ok := parentVal[parentKey]; ok {
					parent = childVal
				} else {
					break
				}
			case map[int]interface{}:
				if parentKeyInt, err := strconv.Atoi(parentKey); err == nil {
					if childVal, ok := parentVal[parentKeyInt]; ok {
						parent = childVal
					} else {
						break
					}
				} else {
					break
				}
			default:
				break
			}
		}

		if key[0] == '-' {
			if key == "--" || key == "---" {
				switch p := parent.(type) {
				case map[string]interface{}:
					key = strconv.Itoa(len(p))
				case map[int]interface{}:
					key = strconv.Itoa(len(p))
				default:
					key = "0"
				}
			} else if key == "-+" || key == "-++" {
				value = strings.TrimSpace(value)
				if key == "-++" {
					if value != "" {
						layer_init = value
					}
					layer_curr = layer_init
				} else if key == "-+" {

					if value == "" {
						num, err := strconv.Atoi(layer_curr)
						if err == nil {
							layer_curr = strconv.Itoa(num + 1)
						} else {
							layer_curr = layer_init
						}
					} else {
						layer_curr = value
					}
				}
				all_layers[layer_curr] = struct{}{}
				continue
			} else {
				decoded_key, err := h.base64Udecode(key[1:])
				if err == nil {
					key = decoded_key
				}
			}
		}

		var setValue interface{}
		needSet := true
		if value == "" {
			if haveDot == 1 {
				setValue = map[int]interface{}{}
			} else {
				setValue = map[string]interface{}{}
			}
			stack = append(stack, key)
			layer_curr = layer_init
		} else {
			if _, ok := layers_list[layer_curr]; ok {
				setValue = valueDecoFun(value, spc_ch)
			} else {
				needSet = false
			}
		}
		if needSet {
			if parentMapString, ok := parent.(map[string]interface{}); ok {
				parentMapString[key] = setValue
			} else if parentMapInt, ok := parent.(map[int]interface{}); ok {
				if keyInt, err := strconv.Atoi(key); err == nil {
					parentMapInt[keyInt] = setValue
				}
			}
		}
	}

	if h.ADD_LAYERS_KEY == true && len(all_layers) > 1 {
		result["_layers"] = keys(all_layers)
	}

	return result
}

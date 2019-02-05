package tahwil

import (
	"encoding/json"
	"fmt"
)

type Value struct {
	Refid uint64      `json:"refid"`
	Kind  string      `json:"kind"`
	Value interface{} `json:"value"`
}

// An InvalidValueError describes invalid Value state.
type InvalidValueError struct {
	Value interface{}
	Kind string
}

func (e *InvalidValueError) Error() string {
	return fmt.Sprintf("tahwil.Value: invalid value \"%#v\" for kind \"%s\"", e.Value, e.Kind)
}

type InvalidValueKindError struct {
	Kind string
}

func (e *InvalidValueKindError) Error() string {
	return "tahwil.Value: invalid value kind \"" + e.Kind + "\""
}

// fixTypes recursively fixes field types after json.Unmarshal
func fixTypes(kind string, v interface{}) (res interface{}, err error) {
	switch kind {
	default:
		if v == nil {
			return v, nil
		}
		return nil, &InvalidValueKindError{Kind: kind}
	case "string", "bool":
		return v, nil
	case "ref", "int":
		return int(v.(float64)), nil
	case "int8":
		return int8(v.(float64)), nil
	case "int16":
		return int16(v.(float64)), nil
	case "int32":
		return int32(v.(float64)), nil
	case "int64":
		return int64(v.(float64)), nil
	case "uint":
		return uint(v.(float64)), nil
	case "uint8":
		return uint8(v.(float64)), nil
	case "uint16":
		return uint16(v.(float64)), nil
	case "uint32":
		return uint32(v.(float64)), nil
	case "uint64":
		return uint64(v.(float64)), nil
	case "float32":
		return float32(v.(float64)), nil
	case "float64":
		return v.(float64), nil
	case "ptr":
		m, ok := v.(map[string]interface{})
		if !ok {
			return nil, &InvalidValueError{Kind: kind, Value: v}
		}
		iv := &Value{
			Refid: uint64(m["refid"].(float64)),
			Kind:  m["kind"].(string),
		}
		if m["value"] == nil {
			return iv, nil
		}
		iv.Value, err = fixTypes(iv.Kind, m["value"])
		if err != nil {
			return nil, err
		}
		return iv, nil
	case "struct", "map":
		m, ok := v.(map[string]interface{})
		if !ok {
			return nil, &InvalidValueError{Kind: kind, Value: v}
		}
		for k, mv := range m {
			m[k], err = fixTypes("ptr", mv)
			if err != nil {
				return nil, err
			}
		}
		return m, nil
	case /*"array", */"slice": // TODO: support array?
		m, ok := v.([]interface{})
		if !ok {
			return nil, &InvalidValueError{Kind: kind, Value: v}
		}
		for k, mv := range m {
			m[k], err = fixTypes("ptr", mv)
			if err != nil {
				return nil, err
			}
		}
		return m, nil
	}
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (v *Value) UnmarshalJSON(b []byte) error {
	// Ignore null, like in the main JSON package.
	if string(b) == "null" {
		return nil
	}

	type valueT struct {
		Refid uint64      `json:"refid"`
		Kind  string      `json:"kind"`
		Value interface{} `json:"value"`
	}
	innerV := &valueT{}
	err := json.Unmarshal(b, innerV)
	if err != nil {
		return err
	}

	v.Refid = innerV.Refid
	v.Kind = innerV.Kind
	v.Value, err = fixTypes(innerV.Kind, innerV.Value)
	if err != nil {
		return err
	}

	return nil
}

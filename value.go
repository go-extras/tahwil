package tahwil

import (
	"encoding/json"
	"fmt"
)

type Value struct {
	Refid uint64 `json:"refid"`
	Kind  Kind   `json:"kind"`
	Value any    `json:"value"`
}

// An InvalidValueError describes invalid Value state.
type InvalidValueError struct {
	Value any
	Kind  Kind
}

func (e *InvalidValueError) Error() string {
	return fmt.Sprintf("tahwil.Value: invalid value %T(%#v) for kind %#v", e.Value, e.Value, e.Kind)
}

type InvalidValueKindError struct {
	Kind Kind
}

func (e *InvalidValueKindError) Error() string {
	return "tahwil.Value: invalid value kind \"" + string(e.Kind) + "\""
}

func fixPtr(kind Kind, v any) (any, error) {
	m, ok := v.(map[string]any)
	if !ok {
		return nil, &InvalidValueError{Kind: kind, Value: v}
	}
	iv := &Value{
		Refid: uint64(m["refid"].(float64)),
		Kind:  Kind(m["kind"].(string)),
	}
	if m["value"] == nil {
		return iv, nil
	}
	var err error
	iv.Value, err = fixTypes(iv.Kind, m["value"])
	if err != nil {
		return nil, err
	}
	return iv, nil
}

func fixStructOrMap(kind Kind, v any) (any, error) {
	m, ok := v.(map[string]any)
	if !ok {
		return nil, &InvalidValueError{Kind: kind, Value: v}
	}
	var err error
	for k, mv := range m {
		m[k], err = fixTypes(Ptr, mv)
		if err != nil {
			return nil, err
		}
	}
	return m, nil
}

func fixSlice(kind Kind, v any) (any, error) {
	m, ok := v.([]any)
	if !ok {
		return nil, &InvalidValueError{Kind: kind, Value: v}
	}
	var err error
	for k, mv := range m {
		m[k], err = fixTypes(Ptr, mv)
		if err != nil {
			return nil, err
		}
	}
	return m, nil
}

// fixTypes recursively fixes field types after json.Unmarshal
//
//nolint:gocyclo // go lacks generics and as such there is no further way to optimize it
func fixTypes(kind Kind, v any) (res any, err error) {
	switch kind {
	case String, Bool:
		return v, nil
	case Ref, Int:
		return int(v.(float64)), nil
	case Int8:
		return int8(v.(float64)), nil
	case Int16:
		return int16(v.(float64)), nil
	case Int32:
		return int32(v.(float64)), nil
	case Int64:
		return int64(v.(float64)), nil
	case Uint:
		return uint(v.(float64)), nil
	case Uint8:
		return uint8(v.(float64)), nil
	case Uint16:
		return uint16(v.(float64)), nil
	case Uint32:
		return uint32(v.(float64)), nil
	case Uint64:
		return uint64(v.(float64)), nil
	case Float32:
		return float32(v.(float64)), nil
	case Float64:
		return v.(float64), nil
	case Ptr:
		return fixPtr(kind, v)
	case Struct, Map:
		return fixStructOrMap(kind, v)
	case /*Array, */ Slice: // TODO: support array?
		return fixSlice(kind, v)
	}

	if v == nil {
		return nil, nil
	}

	return nil, &InvalidValueKindError{Kind: kind}
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (v *Value) UnmarshalJSON(b []byte) error {
	// Ignore null, like in the main JSON package.
	if string(b) == "null" {
		return nil
	}

	type valueT struct {
		Refid uint64 `json:"refid"`
		Kind  string `json:"kind"`
		Value any    `json:"value"`
	}
	innerV := &valueT{}
	err := json.Unmarshal(b, innerV)
	if err != nil {
		return err
	}

	v.Refid = innerV.Refid
	v.Kind = Kind(innerV.Kind)
	v.Value, err = fixTypes(Kind(innerV.Kind), innerV.Value)
	if err != nil {
		return err
	}

	return nil
}

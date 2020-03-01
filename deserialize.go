package tahwil

import (
	"fmt"
	"reflect"
)

type UnmapperError struct {
	text string
}

func (e *UnmapperError) Error() string {
	return "tahwil.FromValue: " + e.text
}

// An InvalidUnmapperKindError describes an invalid argument passed to FromValue.
type InvalidUnmapperKindError struct {
	Expected string
	Kind     string
}

func (e *InvalidUnmapperKindError) Error() string {
	if e.Expected == "" {
		return "tahwil.FromValue: unsupported kind (" + e.Kind + ")"
	}
	return "tahwil.FromValue: unexpected kind (expected: " + e.Expected + ", got: " + e.Kind + ")"
}

type valueUnmapper struct {
	// refs contains pointers to reference values during deserialization
	// can be used both forward and backward lookups
	refs map[uint64]reflect.Value
	// fieldTagCache holds type => json:<tag> => field
	// e.g. if a struct Struct has a field that is called FieldName
	// and it has a struct tag `json:"field_name", filedTagCache will hold
	// [<Struct>]["field_name"]["FieldName"]
	fieldTagCache map[reflect.Type]map[string]string
}

func newValueUnmapper() *valueUnmapper {
	return &valueUnmapper{
		refs:          make(map[uint64]reflect.Value),
		fieldTagCache: make(map[reflect.Type]map[string]string),
	}
}

// fieldByTag returns field name for a given type and a tag name.
// If no tag is found, it will return the name of the field
func (vu *valueUnmapper) fieldByTag(t reflect.Type, key string) string {
	if vu.fieldTagCache[t] != nil {
		return vu.fieldTagCache[t][key]
	}

	vu.fieldTagCache[t] = make(map[string]string)
	for i := 0; i < t.NumField(); i++ {
		ft := t.Field(i)
		k := ft.Tag.Get("json")
		if k == "" {
			k = ft.Name
		}
		if k == "-" || k == "_" {
			continue
		}
		vu.fieldTagCache[t][k] = ft.Name
	}

	return vu.fieldTagCache[t][key]
}

// fills v with the values from data
func (vu *valueUnmapper) fromValue(data *Value, v reflect.Value) (err error) {
	vu.refs[data.Refid] = v

	switch data.Kind {
	default:
		return &InvalidUnmapperKindError{Kind: data.Kind}
	case "bool":
		if v.Kind() != reflect.Bool {
			return &InvalidUnmapperKindError{Expected: "bool", Kind: v.Kind().String()}
		}

		if fval, ok := data.Value.(bool); ok {
			v.SetBool(fval)
			return
		}

		return &InvalidValueError{
			Value: data.Value,
			Kind:  data.Kind,
		}
	case "int", "int8", "int16", "int32", "int64":
		switch v.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			switch vv := data.Value.(type) {
			case int:
				if !v.OverflowInt(int64(vv)) {
					v.SetInt(int64(vv))
					return
				}
			case int8:
				if !v.OverflowInt(int64(vv)) {
					v.SetInt(int64(vv))
					return
				}
			case int16:
				if !v.OverflowInt(int64(vv)) {
					v.SetInt(int64(vv))
					return
				}
			case int32:
				if !v.OverflowInt(int64(vv)) {
					v.SetInt(int64(vv))
					return
				}
			case int64:
				if !v.OverflowInt(vv) {
					v.SetInt(vv)
					return
				}
			}
			return &InvalidValueError{
				Value: data.Value,
				Kind:  data.Kind,
			}
		}
		return &InvalidUnmapperKindError{Expected: "int|int8|int16|int32|int64", Kind: v.Kind().String()}
	case "uint", "uint8", "uint16", "uint32", "uint64":
		switch v.Kind() {
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			switch vv := data.Value.(type) {
			case uint:
				if !v.OverflowUint(uint64(vv)) {
					v.SetUint(uint64(vv))
					return
				}
			case uint8:
				if !v.OverflowUint(uint64(vv)) {
					v.SetUint(uint64(vv))
					return
				}
			case uint16:
				if !v.OverflowUint(uint64(vv)) {
					v.SetUint(uint64(vv))
					return
				}
			case uint32:
				if !v.OverflowUint(uint64(vv)) {
					v.SetUint(uint64(vv))
					return
				}
			case uint64:
				if !v.OverflowUint(vv) {
					v.SetUint(vv)
					return
				}
			}
			return &InvalidValueError{
				Value: data.Value,
				Kind:  data.Kind,
			}
		}
		return &InvalidUnmapperKindError{Expected: "uint|uint8|uint16|uint32|uint64", Kind: v.Kind().String()}
	case "float32", "float64":
		switch v.Kind() {
		case reflect.Float32, reflect.Float64:
			switch vv := data.Value.(type) {
			case float32:
				if !v.OverflowFloat(float64(vv)) {
					v.SetFloat(float64(vv))
					return
				}
			case float64:
				if !v.OverflowFloat(vv) {
					v.SetFloat(vv)
					return
				}
			}
			return &InvalidValueError{
				Value: data.Value,
				Kind:  data.Kind,
			}
		}
		return &InvalidUnmapperKindError{Expected: "float32|float64", Kind: v.Kind().String()}
	case /*"array", */ "slice": // TODO: how to deal with array?
		if data.Value == nil {
			return
		}
		var sl reflect.Value
		mi, ok := data.Value.([]interface{})
		if ok {
			sl = reflect.MakeSlice(v.Type(), len(mi), len(mi))
		}
		mv, ok := data.Value.([]*Value)
		if ok {
			sl = reflect.MakeSlice(v.Type(), len(mv), len(mv))
		}
		if mi == nil && mv == nil {
			return &InvalidValueError{Value: data.Value, Kind: data.Kind}
		}

		v.Set(sl)
		for i := 0; i < v.Len(); i++ {
			var x *Value
			el := v.Index(i)
			if mv != nil {
				x = mv[i]
			} else {
				x = mi[i].(*Value)
			}
			err = vu.fromValue(x, el)
			if err != nil {
				return
			}
		}
	case "map":
		if data.Value == nil {
			return
		}
		if v.Kind() != reflect.Map {
			return &InvalidUnmapperKindError{Expected: "map", Kind: v.Kind().String()}
		}
		var keys []reflect.Value
		mi, ok := data.Value.(map[string]interface{})
		if ok {
			rm := reflect.ValueOf(mi)
			keys = rm.MapKeys()
		}
		mv, ok := data.Value.(map[string]*Value)
		if ok {
			rm := reflect.ValueOf(mv)
			keys = rm.MapKeys()
		}
		if mi == nil && mv == nil {
			return &InvalidValueError{Value: data.Value, Kind: data.Kind}
		}
		v.Set(reflect.MakeMap(v.Type()))
		for _, key := range keys {
			var x *Value
			if mv != nil {
				x = mv[key.String()]
			} else {
				x = mi[key.String()].(*Value)
			}
			f := reflect.New(v.Type().Elem()).Elem()
			err = vu.fromValue(x, f)
			if err != nil {
				return
			}
			fmt.Println(key, f)
			v.SetMapIndex(key, f)
			fmt.Println(v)
		}
	case "ptr":
		if data.Value == nil {
			return
		}

		el := v.Elem()
		if !el.IsValid() {
			t := v.Type()
			telm := t.Elem()
			elm := reflect.New(telm)
			v.Set(elm)
			el = v.Elem()
		}
		x := data.Value.(*Value)
		err = vu.fromValue(x, el)
		if err != nil {
			return
		}
	case "string":
		if v.Kind() != reflect.String {
			return &InvalidUnmapperKindError{Expected: "string", Kind: v.Kind().String()}
		}

		if fval, ok := data.Value.(string); ok {
			v.SetString(fval)
			return
		}

		return &InvalidValueError{
			Value: data.Value,
			Kind:  data.Kind,
		}
	case "struct":
		if v.Kind() == reflect.Interface {
			v = v.Elem()
		}
		if v.Kind() != reflect.Struct {
			return &InvalidUnmapperKindError{Expected: "struct", Kind: v.Kind().String()}
		}
		var keys []reflect.Value
		mi, ok := data.Value.(map[string]interface{})
		if ok {
			rm := reflect.ValueOf(mi)
			keys = rm.MapKeys()
		}
		mv, ok := data.Value.(map[string]*Value)
		if ok {
			rm := reflect.ValueOf(mv)
			keys = rm.MapKeys()
		}
		if mi == nil && mv == nil {
			return &InvalidValueError{Value: data.Value, Kind: data.Kind}
		}

		for _, key := range keys {
			tagName := key.String()
			keyName := vu.fieldByTag(v.Type(), tagName)
			if keyName == "" {
				keyName = tagName
			}
			f := v.FieldByName(keyName)
			var x *Value
			if mv != nil {
				x = mv[tagName]
			} else {
				x = mi[tagName].(*Value)
			}
			if f.IsValid() {
				err = vu.fromValue(x, f)
				if err != nil {
					return err
				}
			}
		}
	case "ref":
		ref := data.Value.(*Reference)
		if refv, ok := vu.refs[ref.Refid]; ok {
			v.Set(refv)
			return
		}
		err = vu.fromValue(ref.Value, v)
		if err != nil {
			return err
		}
	}

	return
}

func FromValue(data *Value, v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &UnmapperError{text: "value must be non-nil Pointer"}
	}

	resolver := NewResolver(data)
	vu := newValueUnmapper()
	if err := resolver.Resolve(); err != nil {
		return &UnmapperError{text: err.Error()}
	}
	if resolver.HasUnresolved() {
		return &UnmapperError{text: "can't resolve all refs, invalid input"}
	}

	return vu.fromValue(data, rv)
}

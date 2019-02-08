package tahwil

import (
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
		if fval, ok := data.Value.(bool); ok {
			v.SetBool(fval)
		}
	case "int", "int8", "int16", "int32", "int64":
		if fval, ok := data.Value.(int64); ok && !v.OverflowInt(fval) {
			v.SetInt(fval)
		}
	case "uint", "uint8", "uint16", "uint32", "uint64":
		if fval, ok := data.Value.(uint64); ok && !v.OverflowUint(fval) {
			v.SetUint(fval)
		}
	case "float32", "float64":
		if fval, ok := data.Value.(float64); ok && !v.OverflowFloat(fval) {
			v.SetFloat(fval)
		}
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
		for _, key := range keys {
			f := v.MapIndex(key)
			var x *Value
			if mv != nil {
				x = mv[key.String()]
			} else {
				x = mi[key.String()].(*Value)
			}
			if f.IsValid() {
				// A Value can be changed only if it is
				// addressable and was not obtained by
				// the use of unexported struct fields.
				if f.CanSet() {
					err = vu.fromValue(x, f)
					if err != nil {
						return
					}
				}
			}
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
		if fval, ok := data.Value.(string); ok {
			v.SetString(fval)
		}
	case "struct":
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
				// TODO: maybe error?
				continue
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
	resolver.Resolve()
	if resolver.HasUnresolved() {
		return &UnmapperError{text: "can't resolve all refs, invalid input"}
	}

	return vu.fromValue(data, rv)
}

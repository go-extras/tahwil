package tahwil

import (
	"reflect"
	"strings"
)

type UnmapperError struct {
	text  string
	cause error
}

func (e *UnmapperError) Error() string {
	if e.cause != nil {
		return "tahwil.FromValue: " + e.cause.Error()
	}
	return "tahwil.FromValue: " + e.text
}

func (e *UnmapperError) Unwrap() error {
	return e.cause
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
	for _, ft := range reflect.VisibleFields(t) {
		if !ft.IsExported() || ft.Anonymous {
			continue
		}
		k := ft.Tag.Get("json")
		if k != "" {
			k, _, _ = strings.Cut(k, ",")
		}
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

func (vu *valueUnmapper) fromBoolValue(data *Value, v reflect.Value) error {
	if v.Kind() != reflect.Bool {
		return &InvalidUnmapperKindError{Expected: string(Bool), Kind: v.Kind().String()}
	}

	if fval, ok := data.Value.(bool); ok {
		v.SetBool(fval)
		return nil
	}

	return &InvalidValueError{
		Value: data.Value,
		Kind:  data.Kind,
	}
}

//nolint:dupl // false positive!
func (vu *valueUnmapper) fromIntValue(data *Value, v reflect.Value) error {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch vv := data.Value.(type) {
		case int:
			if !v.OverflowInt(int64(vv)) {
				v.SetInt(int64(vv))
				return nil
			}
		case int8:
			if !v.OverflowInt(int64(vv)) {
				v.SetInt(int64(vv))
				return nil
			}
		case int16:
			if !v.OverflowInt(int64(vv)) {
				v.SetInt(int64(vv))
				return nil
			}
		case int32:
			if !v.OverflowInt(int64(vv)) {
				v.SetInt(int64(vv))
				return nil
			}
		case int64:
			if !v.OverflowInt(vv) {
				v.SetInt(vv)
				return nil
			}
		}
		return &InvalidValueError{
			Value: data.Value,
			Kind:  data.Kind,
		}
	}
	return &InvalidUnmapperKindError{Expected: "int|int8|int16|int32|int64", Kind: v.Kind().String()}
}

//nolint:dupl // false positive!
func (vu *valueUnmapper) fromUintValue(data *Value, v reflect.Value) error {
	switch v.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		switch vv := data.Value.(type) {
		case uint:
			if !v.OverflowUint(uint64(vv)) {
				v.SetUint(uint64(vv))
				return nil
			}
		case uint8:
			if !v.OverflowUint(uint64(vv)) {
				v.SetUint(uint64(vv))
				return nil
			}
		case uint16:
			if !v.OverflowUint(uint64(vv)) {
				v.SetUint(uint64(vv))
				return nil
			}
		case uint32:
			if !v.OverflowUint(uint64(vv)) {
				v.SetUint(uint64(vv))
				return nil
			}
		case uint64:
			if !v.OverflowUint(vv) {
				v.SetUint(vv)
				return nil
			}
		}
		return &InvalidValueError{
			Value: data.Value,
			Kind:  data.Kind,
		}
	}
	return &InvalidUnmapperKindError{Expected: "uint|uint8|uint16|uint32|uint64", Kind: v.Kind().String()}
}

func (vu *valueUnmapper) fromFloatValue(data *Value, v reflect.Value) error {
	switch v.Kind() {
	case reflect.Float32, reflect.Float64:
		switch vv := data.Value.(type) {
		case float32:
			if !v.OverflowFloat(float64(vv)) {
				v.SetFloat(float64(vv))
				return nil
			}
		case float64:
			if !v.OverflowFloat(vv) {
				v.SetFloat(vv)
				return nil
			}
		}
		return &InvalidValueError{
			Value: data.Value,
			Kind:  data.Kind,
		}
	}
	return &InvalidUnmapperKindError{Expected: "float32|float64", Kind: v.Kind().String()}
}

func (vu *valueUnmapper) fromSliceValue(data *Value, v reflect.Value) error {
	if data.Value == nil {
		return nil
	}
	var sl reflect.Value
	mi, ok := data.Value.([]any)
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
		err := vu.fromValue(x, el)
		if err != nil {
			return err
		}
	}

	return nil
}

func (vu *valueUnmapper) fromMapValue(data *Value, v reflect.Value) error {
	if data.Value == nil {
		return nil
	}
	if v.Kind() != reflect.Map {
		return &InvalidUnmapperKindError{Expected: string(Map), Kind: v.Kind().String()}
	}
	var keys []reflect.Value
	mi, ok := data.Value.(map[string]any)
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
		err := vu.fromValue(x, f)
		if err != nil {
			return err
		}
		v.SetMapIndex(key, f)
	}

	return nil
}

func (vu *valueUnmapper) fromPtrValue(data *Value, v reflect.Value) error {
	if data.Value == nil {
		return nil
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
	err := vu.fromValue(x, el)
	if err != nil {
		return err
	}

	return nil
}

func (vu *valueUnmapper) fromStringValue(data *Value, v reflect.Value) error {
	if v.Kind() != reflect.String {
		return &InvalidUnmapperKindError{Expected: string(String), Kind: v.Kind().String()}
	}

	if fval, ok := data.Value.(string); ok {
		v.SetString(fval)
		return nil
	}

	return &InvalidValueError{
		Value: data.Value,
		Kind:  data.Kind,
	}
}

func (vu *valueUnmapper) fromStructValue(data *Value, v reflect.Value) error {
	if v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return &InvalidUnmapperKindError{Expected: string(Struct), Kind: v.Kind().String()}
	}
	var keys []reflect.Value
	mi, ok := data.Value.(map[string]any)
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
			err := vu.fromValue(x, f)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (vu *valueUnmapper) fromRefValue(data *Value, v reflect.Value) error {
	ref := data.Value.(*Reference)
	if refv, ok := vu.refs[ref.Refid]; ok {
		v.Set(refv)
		return nil
	}
	err := vu.fromValue(ref.Value, v)
	if err != nil {
		return err
	}

	return nil
}

// fills v with the values from data
func (vu *valueUnmapper) fromValue(data *Value, v reflect.Value) error {
	if data.Refid != 0 {
		vu.refs[data.Refid] = v
	}

	switch data.Kind {
	case Bool:
		return vu.fromBoolValue(data, v)
	case Int, Int8, Int16, Int32, Int64:
		return vu.fromIntValue(data, v)
	case Uint, Uint8, Uint16, Uint32, Uint64:
		return vu.fromUintValue(data, v)
	case Float32, Float64:
		return vu.fromFloatValue(data, v)
	case /*Array, */ Slice: // TODO: how to deal with array?
		return vu.fromSliceValue(data, v)
	case Map:
		return vu.fromMapValue(data, v)
	case Ptr:
		return vu.fromPtrValue(data, v)
	case String:
		return vu.fromStringValue(data, v)
	case Struct:
		return vu.fromStructValue(data, v)
	case Ref:
		return vu.fromRefValue(data, v)
	}

	return &InvalidUnmapperKindError{Kind: string(data.Kind)}
}

func FromValue(data *Value, v any) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &UnmapperError{text: "value must be non-nil Pointer"}
	}

	resolver := NewResolver(data)
	vu := newValueUnmapper()
	if err := resolver.Resolve(); err != nil {
		return &UnmapperError{cause: err}
	}
	if resolver.HasUnresolved() {
		return &UnmapperError{text: "can't resolve all refs, invalid input"}
	}

	return vu.fromValue(data, rv)
}

package tahwil

import (
	"fmt"
	"reflect"
)

// An InvalidMapperKindError describes an invalid argument passed to ToValue.
type InvalidMapperKindError struct {
	Kind string
}

func (e *InvalidMapperKindError) Error() string {
	if e.Kind == "" {
		return "tahwil.ToValue: empty kind"
	}
	return "tahwil.ToValue: unsupported kind (" + e.Kind + ")"
}

type valueMapper struct {
	// references during serialization
	refs map[uintptr]uint64
	// refid that was last generated
	lastRefid uint64
}

func newValueMapper() *valueMapper {
	return &valueMapper{
		refs:      make(map[uintptr]uint64),
		lastRefid: 0,
	}
}

func (vm *valueMapper) nextRefid() uint64 {
	vm.lastRefid += 1
	return vm.lastRefid
}

func (vm *valueMapper) toValueSlice(v reflect.Value) (result []*Value, err error) {
	result = make([]*Value, v.Len(), v.Len())
	for i := 0; i < v.Len(); i++ {
		result[i], err = vm.toValue(v.Index(i).Interface())
		if err != nil {
			return nil, err
		}
	}
	return
}

func (vm *valueMapper) toValueMap(v reflect.Value) (result map[string]*Value, err error) {
	result = make(map[string]*Value)
	kind := v.Kind()
	if kind == reflect.Map {
		keys := v.MapKeys()
		if len(keys) == 0 {
			return
		}

		for _, idx := range keys {
			i := idx.Interface()
			val := v.MapIndex(idx)
			resIdx := fmt.Sprintf("%v", i)
			result[resIdx], err = vm.toValue(val.Interface())
			if err != nil {
				return nil, err
			}
		}
		return
	}

	if kind == reflect.Struct {
		for i := 0; i < v.NumField(); i++ {
			ft := v.Type().Field(i)
			k := ft.Tag.Get("json")
			if k == "" {
				k = ft.Name
			}
			if k == "-" || k == "_" {
				continue
			}

			f := v.Field(i)
			if !f.CanInterface() {
				continue
			}
			result[k], err = vm.toValue(f.Interface())
			if err != nil {
				return nil, err
			}
		}
		return
	}

	return nil, &InvalidMapperKindError{Kind: kind.String()}
}

func (vm *valueMapper) toValue(i interface{}) (result *Value, err error) {
	result = &Value{}

	v := reflect.ValueOf(i)
	kind := v.Kind()

	if kind == reflect.Chan || kind == reflect.Func || kind == reflect.Interface {
		return nil, &InvalidMapperKindError{Kind: kind.String()}
	}

	if kind == reflect.Ptr {
		if refid, ok := vm.refs[v.Pointer()]; ok {
			result.Refid = vm.nextRefid()
			result.Kind = "ref"
			result.Value = refid
			return
		}
		if v.IsNil() {
			result.Refid = vm.nextRefid()
			result.Kind = "ptr"
			result.Value = nil
			p := v.Pointer()
			vm.refs[p] = result.Refid
			return
		}

		result.Refid = vm.nextRefid()
		result.Kind = "ptr"
		p := v.Pointer()
		vm.refs[p] = result.Refid
		val, err := vm.toValue(v.Elem().Interface())
		if err != nil {
			return nil, err
		}
		result.Value = val
		return result, nil
	}

	if kind == reflect.Slice || kind == reflect.Array {
		result.Refid = vm.nextRefid()
		result.Kind = kind.String()
		result.Value, err = vm.toValueSlice(v)
		if err != nil {
			return nil, err
		}
		return
	}

	if kind == reflect.Map || kind == reflect.Struct {
		result.Refid = vm.nextRefid()
		result.Kind = kind.String()
		result.Value, err = vm.toValueMap(v)
		if err != nil {
			return nil, err
		}
		return
	}

	result.Refid = vm.nextRefid()
	result.Kind = v.Type().Name()
	if result.Kind == "" {
		return nil, &InvalidMapperKindError{Kind: ""}
	}
	result.Value = v.Interface()

	return
}

func ToValue(i interface{}) (*Value, error) {
	v := reflect.ValueOf(i)
	if v.Kind() != reflect.Ptr {
		i = &i
	}

	vm := newValueMapper()

	return vm.toValue(i)
}

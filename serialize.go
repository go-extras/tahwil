package tahwil

import (
	"fmt"
	"reflect"
	"strings"
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

func (vm *valueMapper) saveRef(v reflect.Value) uint64 {
	refid := vm.nextRefid()
	vm.refs[v.Pointer()] = refid
	return refid
}

func (vm *valueMapper) nextRefid() uint64 {
	vm.lastRefid++
	return vm.lastRefid
}

func (vm *valueMapper) toValueSlice(v reflect.Value) (result []*Value, err error) {
	result = make([]*Value, v.Len())
	for i := 0; i < v.Len(); i++ {
		result[i], err = vm.toValue(v.Index(i))
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
			return result, nil
		}

		for _, idx := range keys {
			i := idx.Interface()
			val := v.MapIndex(idx)
			resIdx := fmt.Sprintf("%v", i)
			result[resIdx], err = vm.toValue(val)
			if err != nil {
				return nil, err
			}
		}
		return result, nil
	}

	if kind == reflect.Struct {
		for i := 0; i < v.NumField(); i++ {
			ft := v.Type().Field(i)
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

			f := v.Field(i)
			if !f.CanInterface() {
				continue
			}
			result[k], err = vm.toValue(f)
			if err != nil {
				return nil, err
			}
		}
		return result, nil
	}

	return nil, &InvalidMapperKindError{Kind: kind.String()}
}

func (vm *valueMapper) ptrToValue(v reflect.Value) (result *Value, err error) {
	result = &Value{}

	if refid, ok := vm.refs[v.Pointer()]; ok {
		result.Refid = vm.nextRefid()
		result.Kind = Ref
		result.Value = refid
		return result, nil
	}

	result.Refid = vm.saveRef(v)
	result.Kind = Ptr

	if v.IsNil() || v.Elem().Interface() == nil {
		// nil values a final, no further elements
		result.Value = nil
		return result, nil
	}

	// recursively proceed to the pointer value
	var val *Value
	val, err = vm.toValue(v.Elem())
	if err != nil {
		return nil, err
	}
	result.Value = val
	return result, err
}

func (vm *valueMapper) sliceToValue(v reflect.Value, kind reflect.Kind) (result *Value, err error) {
	result = &Value{}

	result.Refid = vm.nextRefid()
	result.Kind = Kind(kind.String())
	result.Value, err = vm.toValueSlice(v)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (vm *valueMapper) mapOrStructToValue(v reflect.Value, kind reflect.Kind) (result *Value, err error) {
	result = &Value{}

	result.Refid = vm.nextRefid()
	result.Kind = Kind(kind.String())
	// not only maps can be set here, but also slices as they
	// can be represented as a map[fieldName]value
	result.Value, err = vm.toValueMap(v)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (vm *valueMapper) scalarToValue(v reflect.Value) (result *Value, err error) {
	result = &Value{}

	// here we process the remaining kinds ("simple" ones)
	result.Refid = vm.nextRefid()
	result.Kind = Kind(v.Type().Name())
	if result.Kind == "" {
		return nil, &InvalidMapperKindError{Kind: ""}
	}
	result.Value = v.Interface()

	return result, nil
}

func (vm *valueMapper) toValue(v reflect.Value) (result *Value, err error) {
	kind := v.Kind()

	// special case for interface (internally interfaces act similarly to pointers,
	// but we don't want to store them like pointers)
	if kind == reflect.Interface {
		v = v.Elem()
		kind = v.Kind()
	}

	switch kind {
	case reflect.Chan, reflect.Func, reflect.Uintptr, reflect.Array, reflect.UnsafePointer:
		return nil, &InvalidMapperKindError{Kind: kind.String()}
	case reflect.Ptr:
		return vm.ptrToValue(v)
	// case reflect.Array: // TODO: implement array?
	//	fallthrough
	case reflect.Slice:
		return vm.sliceToValue(v, kind)
	case reflect.Struct, reflect.Map:
		return vm.mapOrStructToValue(v, kind)
	default:
		return vm.scalarToValue(v)
	}
}

// ToValue transforms i to *Value.
// NOTES:
//   - (*Value).Kind will be set to the reflected value kind (see reflect.Kind).
//   - each value will get its own (*Value).Refid that is a counter incremented
//     with every next (underlying value).
//   - each "simple" type (int, string, bool) will be stored in (*Value).Value
//   - each type that holds other values inside (ptr, map, slice, struct) will
//     produce a further *Value that will be stored in (*Value).Value.
//   - the transformation process will continue until all the non-simple types are processed.
//   - non-serializable types (func, chan) will lead to a mapping error.
//   - there are unsupported serializable types: array, complex[64,128], unsafe pointer.
//   - ptr type will produce *Value with an underlying value.
//   - nil ptr will result in (*Value).Value set to nil.
//   - each non-nil pointer Refid is stored in a Refid map. This map is used
//     to break circular references (when transforming a pointer the Refid map is being checked,
//     and if the pointer is already on the list, (*Value).Kind is set to a special "ref" type
//     and (*Value).Value is set to the Refid of the previously transformed value).
//   - struct will produce *Value whose Value property will be set to the map
//     of exported struct fields (keys will correspond to the field name or to the json tag value,
//     values will be *Value, with the underlying values of the fields), if a field is
//     exported but its json tag value is set to "_" or "-", it will be ignored.
//   - map will produce a map of *Value with the key names that correspond to the original
//     map keys, and the values will be *Value, with corresponding map values transformed.
//   - slice will produce a slice of *Value in the same order like the original slice has
//   - i is expected to be a pointer, but if it's not, a pointer from it will be created,
//     it means that even for "simple" types the resulting (*Value).Value will hold *Value
//     that will represent the original value (mapped to Value{})
//
// The result is *Value and an error, if there was a mapping error.
func ToValue(i any) (*Value, error) {
	v := reflect.ValueOf(i)
	if v.Kind() != reflect.Ptr {
		v = reflect.ValueOf(&i)
	}
	vm := newValueMapper()
	return vm.toValue(v)
}

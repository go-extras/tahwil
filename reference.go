package tahwil

import (
	"fmt"
)

type Reference struct {
	Refid uint64
	Value *Value
}

type ResolverError struct {
	Value *Value
	Kind  string
	Type  string
}

func (e *ResolverError) Error() string {
	if e.Value == nil {
		return "tahwil.Resolver: nil *Value"
	}
	if e.Kind == Ref && e.Value == e.Value.Value {
		return "tahwil.Resolver: *Value == (*Value).Value"
	}

	return fmt.Sprintf("tahwil.Resolver: invalid *Value.Value type: Kind=%q, Type=T%q", e.Kind, e.Type)
}

type Resolver struct {
	data *Value
	// unresolved references during deserialization
	unresolvedRefs map[uint64]*Reference
	// resolved references during deserialization
	resolvedRefs map[uint64]*Value
}

func NewResolver(data *Value) *Resolver {
	return &Resolver{
		data:           data,
		unresolvedRefs: make(map[uint64]*Reference),
		resolvedRefs:   make(map[uint64]*Value),
	}
}

func (r *Resolver) Resolve() error {
	return r.resolve(r.data)
}

func (r *Resolver) HasUnresolved() bool {
	return len(r.unresolvedRefs) > 0
}

func (r *Resolver) Unresolved() []uint64 {
	result := make([]uint64, 0, len(r.unresolvedRefs))
	for k := range r.unresolvedRefs {
		result = append(result, k)
	}
	return result
}

func (r *Resolver) resolvePtr(v *Value) error {
	if v.Value == nil {
		return nil
	}
	if v == v.Value {
		return &ResolverError{
			Value: v,
			Kind:  Ptr,
		}
	}
	iv := v.Value.(*Value)
	return r.resolve(iv)
}

func (r *Resolver) resolveWIthSubvalues(v *Value) error {
	switch val := v.Value.(type) {
	case map[string]any:
		for _, mv := range val {
			iv := mv.(*Value)
			if err := r.resolve(iv); err != nil {
				return err
			}
		}
	case map[string]*Value:
		for _, mv := range val {
			iv := mv
			if err := r.resolve(iv); err != nil {
				return err
			}
		}
	case []any:
		for _, mv := range val {
			iv := mv.(*Value)
			if err := r.resolve(iv); err != nil {
				return err
			}
		}
	case []*Value:
		for _, mv := range val {
			iv := mv
			if err := r.resolve(iv); err != nil {
				return err
			}
		}
	default:
		if v.Value == nil {
			return nil
		}

		return &ResolverError{
			Value: v,
			Type:  fmt.Sprintf("%T", val),
			Kind:  v.Kind,
		}
	}

	return nil
}

func (r *Resolver) refFromValue(v *Value) (uint64, error) {
	var signed int64
	var isSigned bool

	switch vv := v.Value.(type) {
	case float32:
		signed, isSigned = int64(vv), true
	case float64:
		signed, isSigned = int64(vv), true
	case int:
		signed, isSigned = int64(vv), true
	case int8:
		signed, isSigned = int64(vv), true
	case int16:
		signed, isSigned = int64(vv), true
	case int32:
		signed, isSigned = int64(vv), true
	case int64:
		signed, isSigned = vv, true
	case uint:
		return uint64(vv), nil
	case uint8:
		return uint64(vv), nil
	case uint16:
		return uint64(vv), nil
	case uint32:
		return uint64(vv), nil
	case uint64:
		return vv, nil
	default:
		return 0, &ResolverError{Value: v, Kind: Ref, Type: Uint64}
	}

	if isSigned && signed < 0 {
		return 0, &ResolverError{Value: v, Kind: Ref, Type: Uint64}
	}
	return uint64(signed), nil //nolint:gosec // bounds checked above
}

func (r *Resolver) resolveRef(v *Value) error {
	refid, err := r.refFromValue(v)
	if err != nil {
		return err
	}

	iv := r.resolvedRefs[refid]
	if iv != nil {
		v.Value = &Reference{
			Refid: refid,
			Value: iv,
		}
		return nil
	}
	ref := r.unresolvedRefs[refid]
	if ref == nil {
		ref = &Reference{
			Refid: refid,
			Value: nil,
		}
		r.unresolvedRefs[refid] = ref
	}
	v.Value = ref
	if v == v.Value {
		return &ResolverError{
			Value: v,
			Kind:  Ref,
		}
	}

	return nil
}

func (r *Resolver) resolve(v *Value) error {
	if v == nil {
		return &ResolverError{Value: v}
	}

	r.resolvedRefs[v.Refid] = v
	if ref, ok := r.unresolvedRefs[v.Refid]; ok {
		ref.Value = v
		// ok, we resolved it, remove it from the unresolved map
		delete(r.unresolvedRefs, v.Refid)
	}

	switch v.Kind {
	case Ptr:
		return r.resolvePtr(v)
	case Struct, Map, Array, Slice:
		return r.resolveWIthSubvalues(v)
	case Ref:
		return r.resolveRef(v)
	}

	return nil
}

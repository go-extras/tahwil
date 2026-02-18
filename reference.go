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
	Kind  Kind
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

type signedOrFloat interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64
}

func toUint64Signed[T signedOrFloat](v T) (uint64, bool) {
	i := int64(v)
	if i < 0 {
		return 0, false
	}
	return uint64(i), true //nolint:gosec // bounds checked above
}

func (r *Resolver) refFromValue(v *Value) (uint64, error) {
	switch vv := v.Value.(type) {
	case float32:
		if u, ok := toUint64Signed(vv); ok {
			return u, nil
		}
	case float64:
		if u, ok := toUint64Signed(vv); ok {
			return u, nil
		}
	case int:
		if u, ok := toUint64Signed(vv); ok {
			return u, nil
		}
	case int8:
		if u, ok := toUint64Signed(vv); ok {
			return u, nil
		}
	case int16:
		if u, ok := toUint64Signed(vv); ok {
			return u, nil
		}
	case int32:
		if u, ok := toUint64Signed(vv); ok {
			return u, nil
		}
	case int64:
		if u, ok := toUint64Signed(vv); ok {
			return u, nil
		}
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
		return 0, &ResolverError{Value: v, Kind: Ref, Type: string(Uint64)}
	}
	return 0, &ResolverError{Value: v, Kind: Ref, Type: string(Uint64)}
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

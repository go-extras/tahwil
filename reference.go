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
	Kind string
	Type string
}

func (e *ResolverError) Error() string {
	if e.Value == nil {
		return fmt.Sprintf("tahwil.Resolver: nil *Value")
	}
	if e.Kind == "ref" && e.Value == e.Value.Value {
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
	default:
		return nil
	case "ptr":
		if v.Value == nil {
			return nil
		}
		if v == v.Value {
			return &ResolverError{
				Value: v,
				Kind: "ptr",
			}
		}
		iv := v.Value.(*Value)
		return r.resolve(iv)
	case "struct", "map", "array", "slice":
		switch val := v.Value.(type) {
		case map[string]interface{}:
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
		case []interface{}:
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
				Type: fmt.Sprintf("%T", val),
				Kind: v.Kind,
			}
		}
	case "ref":
		refid, ok := v.Value.(uint64)
		if !ok {
			return &ResolverError{
				Value: v,
				Kind: "ref",
				Type: "float64",
			}
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
				Kind: "ref",
			}
		}
	}

	return nil
}

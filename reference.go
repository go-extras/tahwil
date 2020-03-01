package tahwil

import (
	"fmt"
	"reflect"
)

type Reference struct {
	Refid uint64
	Value *Value
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

func (r *Resolver) Resolve() {
	r.resolve(r.data)
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

func (r *Resolver) resolve(v *Value) {
	r.resolvedRefs[v.Refid] = v
	if ref, ok := r.unresolvedRefs[v.Refid]; ok {
		ref.Value = v
		delete(r.unresolvedRefs, v.Refid)
	}
	switch v.Kind {
	default:
	case "ptr":
		if v.Value == nil {
			return
		}
		iv := v.Value.(*Value)
		r.resolve(iv)
	case "struct", "map":
		if v.Value == nil {
			return
		}
		if m, ok := v.Value.(map[string]interface{}); ok {
			for _, mv := range m {
				iv := mv.(*Value)
				r.resolve(iv)
			}
		} else if m, ok := v.Value.(map[string]*Value); ok {
			for _, mv := range m {
				iv := mv
				r.resolve(iv)
			}
		} else {
			panic(fmt.Errorf("unexpected type: %#+v", reflect.TypeOf(v.Value)))
		}
	case /*"array", */ "slice":
		if v.Value == nil {
			return
		}
		if m, ok := v.Value.([]interface{}); ok {
			for _, mv := range m {
				iv := mv.(*Value)
				r.resolve(iv)
			}
		} else if m, ok := v.Value.([]*Value); ok {
			for _, mv := range m {
				iv := mv
				r.resolve(iv)
			}
		} else {
			panic(fmt.Errorf("unexpected type: %#+v", reflect.TypeOf(v.Value)))
		}
	case "ref":
		refid := uint64(v.Value.(uint64))
		iv := r.resolvedRefs[refid]
		if iv != nil {
			v.Value = &Reference{
				Refid: refid,
				Value: iv,
			}
			return
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
	}
}

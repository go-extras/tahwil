package tahwil

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
	for k, _ := range r.unresolvedRefs {
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
		return
	case "ptr":
		if v.Value == nil {
			return
		}
		iv := v.Value.(*Value)
		r.resolve(iv)
		return
	case "struct", "map":
		for _, mv := range v.Value.(map[string]interface{}) {
			iv := mv.(*Value)
			r.resolve(iv)
		}
	case "array", "slice":
		for _, mv := range v.Value.([]interface{}) {
			iv := mv.(*Value)
			r.resolve(iv)
		}
	case "ref":
		refid := uint64(v.Value.(float64))
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

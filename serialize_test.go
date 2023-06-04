package tahwil_test

import (
	"encoding/json"
	"reflect"
	"testing"
	"unsafe"

	"github.com/go-extras/tahwil"
)

type parentSerT struct {
	Name     string
	Children []*childSerT
}

type childSerT struct {
	Name   string
	Parent *parentSerT
}

type selfRefT struct {
	Name string
	Self *selfRefT
}

type interfaceST struct {
	Value        any
	NoSerialize1 bool `json:"-"`
	NoSerialize2 bool `json:"_"`
}

type valueTest struct {
	in  any
	out *tahwil.Value
	err any
}

func valueTests() []valueTest {
	result := make([]valueTest, 0)

	result = append(result, valueTest{
		in: true,
		out: &tahwil.Value{
			Refid: 1,
			Kind:  tahwil.Ptr,
			Value: &tahwil.Value{
				Refid: 2,
				Kind:  tahwil.Bool,
				Value: true,
			},
		},
	})

	result = append(result, valueTest{
		in: nil,
		out: &tahwil.Value{
			Refid: 1,
			Kind:  tahwil.Ptr,
			Value: nil,
		},
	})

	result = append(result, valueTest{
		in: "test",
		out: &tahwil.Value{
			Refid: 1,
			Kind:  tahwil.Ptr,
			Value: &tahwil.Value{
				Refid: 2,
				Kind:  tahwil.String,
				Value: "test",
			},
		},
	})

	s := "test"
	result = append(result, valueTest{
		in: &s,
		out: &tahwil.Value{
			Refid: 1,
			Kind:  tahwil.Ptr,
			Value: &tahwil.Value{
				Refid: 2,
				Kind:  tahwil.String,
				Value: "test",
			},
		},
	})

	result = append(result, valueTest{
		in: &interfaceST{Value: "test"},
		out: &tahwil.Value{
			Refid: 1,
			Kind:  tahwil.Ptr,
			Value: &tahwil.Value{
				Refid: 2,
				Kind:  tahwil.Struct,
				Value: map[string]*tahwil.Value{
					"Value": {
						Refid: 3,
						Kind:  tahwil.String,
						Value: "test",
					},
				},
			},
		},
	})

	result = append(result, valueTest{
		in: any("test"),
		out: &tahwil.Value{
			Refid: 1,
			Kind:  tahwil.Ptr,
			Value: &tahwil.Value{
				Refid: 2,
				Kind:  tahwil.String,
				Value: "test",
			},
		},
	})

	result = append(result, valueTest{
		in: int(47),
		out: &tahwil.Value{
			Refid: 1,
			Kind:  tahwil.Ptr,
			Value: &tahwil.Value{
				Refid: 2,
				Kind:  tahwil.Int,
				Value: int(47),
			},
		},
	})
	result = append(result, valueTest{
		in: int8(47),
		out: &tahwil.Value{
			Refid: 1,
			Kind:  tahwil.Ptr,
			Value: &tahwil.Value{
				Refid: 2,
				Kind:  tahwil.Int8,
				Value: int8(47),
			},
		},
	})
	result = append(result, valueTest{
		in: int16(47),
		out: &tahwil.Value{
			Refid: 1,
			Kind:  tahwil.Ptr,
			Value: &tahwil.Value{
				Refid: 2,
				Kind:  tahwil.Int16,
				Value: int16(47),
			},
		},
	})
	result = append(result, valueTest{
		in: int32(47),
		out: &tahwil.Value{
			Refid: 1,
			Kind:  tahwil.Ptr,
			Value: &tahwil.Value{
				Refid: 2,
				Kind:  tahwil.Int32,
				Value: int32(47),
			},
		},
	})
	result = append(result, valueTest{
		in: int64(47),
		out: &tahwil.Value{
			Refid: 1,
			Kind:  tahwil.Ptr,
			Value: &tahwil.Value{
				Refid: 2,
				Kind:  tahwil.Int64,
				Value: int64(47),
			},
		},
	})

	result = append(result, valueTest{
		in: uint(47),
		out: &tahwil.Value{
			Refid: 1,
			Kind:  tahwil.Ptr,
			Value: &tahwil.Value{
				Refid: 2,
				Kind:  tahwil.Uint,
				Value: uint(47),
			},
		},
	})
	result = append(result, valueTest{
		in: uint8(47),
		out: &tahwil.Value{
			Refid: 1,
			Kind:  tahwil.Ptr,
			Value: &tahwil.Value{
				Refid: 2,
				Kind:  tahwil.Uint8,
				Value: uint8(47),
			},
		},
	})
	result = append(result, valueTest{
		in: byte(47),
		out: &tahwil.Value{
			Refid: 1,
			Kind:  tahwil.Ptr,
			Value: &tahwil.Value{
				Refid: 2,
				Kind:  tahwil.Uint8,
				Value: uint8(47),
			},
		},
	})
	result = append(result, valueTest{
		in: uint16(47),
		out: &tahwil.Value{
			Refid: 1,
			Kind:  tahwil.Ptr,
			Value: &tahwil.Value{
				Refid: 2,
				Kind:  tahwil.Uint16,
				Value: uint16(47),
			},
		},
	})
	result = append(result, valueTest{
		in: uint32(47),
		out: &tahwil.Value{
			Refid: 1,
			Kind:  tahwil.Ptr,
			Value: &tahwil.Value{
				Refid: 2,
				Kind:  tahwil.Uint32,
				Value: uint32(47),
			},
		},
	})
	result = append(result, valueTest{
		in: float32(47.47),
		out: &tahwil.Value{
			Refid: 1,
			Kind:  tahwil.Ptr,
			Value: &tahwil.Value{
				Refid: 2,
				Kind:  tahwil.Float32,
				Value: float32(47.47),
			},
		},
	})
	result = append(result, valueTest{
		in: float64(47.47),
		out: &tahwil.Value{
			Refid: 1,
			Kind:  tahwil.Ptr,
			Value: &tahwil.Value{
				Refid: 2,
				Kind:  tahwil.Float64,
				Value: float64(47.47),
			},
		},
	})

	result = append(result, valueTest{
		in: &parentSerT{
			Name: "Patrik",
		},
		out: &tahwil.Value{
			Refid: 1,
			Kind:  tahwil.Ptr,
			Value: &tahwil.Value{
				Refid: 2,
				Kind:  tahwil.Struct,
				Value: map[string]*tahwil.Value{
					"Name": {
						Refid: 3,
						Kind:  tahwil.String,
						Value: "Patrik",
					},
					"Children": {
						Refid: 4,
						Kind:  tahwil.Slice,
						Value: []*tahwil.Value{},
					},
				},
			},
		},
	})

	p1 := &parentSerT{
		Name:     "Patrik",
		Children: nil,
	}
	c1 := &childSerT{
		Name:   "Valentine",
		Parent: p1,
	}
	p1.Children = append(p1.Children, c1)
	result = append(result, valueTest{
		in: p1,
		out: &tahwil.Value{
			Refid: 1,
			Kind:  tahwil.Ptr,
			Value: &tahwil.Value{
				Refid: 2,
				Kind:  tahwil.Struct,
				Value: map[string]*tahwil.Value{
					"Name": {
						Refid: 3,
						Kind:  tahwil.String,
						Value: "Patrik",
					},
					"Children": {
						Refid: 4,
						Kind:  tahwil.Slice,
						Value: []*tahwil.Value{
							{
								Refid: 5,
								Kind:  tahwil.Ptr,
								Value: &tahwil.Value{
									Refid: 6,
									Kind:  tahwil.Struct,
									Value: map[string]*tahwil.Value{
										"Name": {
											Refid: 7,
											Kind:  tahwil.String,
											Value: "Valentine",
										},
										"Parent": {
											Refid: 8,
											Kind:  tahwil.Ref,
											Value: uint64(1),
										},
									},
								},
							},
						},
					},
				},
			},
		},
	})

	selfRef := &selfRefT{
		Name: "Klark",
	}
	selfRef.Self = selfRef
	result = append(result, valueTest{
		in: selfRef,
		out: &tahwil.Value{
			Refid: 1,
			Kind:  tahwil.Ptr,
			Value: &tahwil.Value{
				Refid: 2,
				Kind:  tahwil.Struct,
				Value: map[string]*tahwil.Value{
					"Name": {
						Refid: 3,
						Kind:  tahwil.String,
						Value: "Klark",
					},
					"Self": {
						Refid: 4,
						Kind:  tahwil.Ref,
						Value: uint64(1),
					},
				},
			},
		},
	})

	result = append(result, valueTest{
		in: map[string]any{
			"Id": uint64(1),
		},
		out: &tahwil.Value{
			Refid: 1,
			Kind:  tahwil.Ptr,
			Value: &tahwil.Value{
				Refid: 2,
				Kind:  tahwil.Map,
				Value: map[string]*tahwil.Value{
					"Id": {
						Refid: 3,
						Kind:  tahwil.Uint64,
						Value: uint64(1),
					},
				},
			},
		},
	})

	result = append(result, valueTest{
		in:  uintptr(1),
		err: &tahwil.InvalidMapperKindError{Kind: "uintptr"},
	})

	result = append(result, valueTest{
		in:  [4]int{1, 2, 3, 4},
		err: &tahwil.InvalidMapperKindError{Kind: tahwil.Array},
	})

	result = append(result, valueTest{
		in:  make(chan any),
		err: &tahwil.InvalidMapperKindError{Kind: "chan"},
	})

	dummy1 := true
	result = append(result, valueTest{
		in:  unsafe.Pointer(&dummy1),
		err: &tahwil.InvalidMapperKindError{Kind: "unsafe.Pointer"},
	})

	dummy2 := &interfaceST{
		Value:        uintptr(1),
		NoSerialize1: true,
		NoSerialize2: true,
	}
	result = append(result, valueTest{
		in:  &dummy2,
		err: &tahwil.InvalidMapperKindError{Kind: "uintptr"},
	})

	return result
}

func TestToValue(t *testing.T) {
	for i, arg := range valueTests() {
		v, err := tahwil.ToValue(arg.in)
		if err != nil {
			if !reflect.DeepEqual(err, arg.err) {
				t.Errorf("#%d: %#+v", i, err)
			}
			// otherwise the error is expected
		} else if !reflect.DeepEqual(v, arg.out) {
			x, err := json.Marshal(v)
			if err != nil {
				t.Fatalf("#%d: %#+v", i, err)
			}
			y, err := json.Marshal(arg.out)
			if err != nil {
				t.Fatalf("#%d: %#+v", i, err)
			}
			t.Errorf("#%d: mismatch\nhave: %v\nwant: %v", i, string(x), string(y))
		}
	}
}

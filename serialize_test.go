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

type omitemptyT struct {
	Name  string `json:"name,omitempty"`
	Value int    `json:"value,omitempty"`
}

type embeddedBaseT struct {
	Name string `json:"name"`
}

type embeddedOuterT struct {
	embeddedBaseT
	Value int `json:"value"`
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
				Refid: 0,
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
				Refid: 0,
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
				Refid: 0,
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
				Refid: 0,
				Kind:  tahwil.Struct,
				Value: map[string]*tahwil.Value{
					"Value": {
						Refid: 0,
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
				Refid: 0,
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
				Refid: 0,
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
				Refid: 0,
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
				Refid: 0,
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
				Refid: 0,
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
				Refid: 0,
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
				Refid: 0,
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
				Refid: 0,
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
				Refid: 0,
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
				Refid: 0,
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
				Refid: 0,
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
				Refid: 0,
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
				Refid: 0,
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
				Refid: 0,
				Kind:  tahwil.Struct,
				Value: map[string]*tahwil.Value{
					"Name": {
						Refid: 0,
						Kind:  tahwil.String,
						Value: "Patrik",
					},
					"Children": {
						Refid: 0,
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
				Refid: 0,
				Kind:  tahwil.Struct,
				Value: map[string]*tahwil.Value{
					"Name": {
						Refid: 0,
						Kind:  tahwil.String,
						Value: "Patrik",
					},
					"Children": {
						Refid: 0,
						Kind:  tahwil.Slice,
						Value: []*tahwil.Value{
							{
								Refid: 2,
								Kind:  tahwil.Ptr,
								Value: &tahwil.Value{
									Refid: 0,
									Kind:  tahwil.Struct,
									Value: map[string]*tahwil.Value{
										"Name": {
											Refid: 0,
											Kind:  tahwil.String,
											Value: "Valentine",
										},
										"Parent": {
											Refid: 3,
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
				Refid: 0,
				Kind:  tahwil.Struct,
				Value: map[string]*tahwil.Value{
					"Name": {
						Refid: 0,
						Kind:  tahwil.String,
						Value: "Klark",
					},
					"Self": {
						Refid: 2,
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
				Refid: 0,
				Kind:  tahwil.Map,
				Value: map[string]*tahwil.Value{
					"Id": {
						Refid: 0,
						Kind:  tahwil.Uint64,
						Value: uint64(1),
					},
				},
			},
		},
	})

	// json tags with options like omitempty should use only the name part
	result = append(result, valueTest{
		in: &omitemptyT{Name: "test", Value: 42},
		out: &tahwil.Value{
			Refid: 1,
			Kind:  tahwil.Ptr,
			Value: &tahwil.Value{
				Refid: 0,
				Kind:  tahwil.Struct,
				Value: map[string]*tahwil.Value{
					"name": {
						Refid: 0,
						Kind:  tahwil.String,
						Value: "test",
					},
					"value": {
						Refid: 0,
						Kind:  tahwil.Int,
						Value: 42,
					},
				},
			},
		},
	})

	// embedded struct fields should be promoted (flat, not nested)
	result = append(result, valueTest{
		in: &embeddedOuterT{
			embeddedBaseT: embeddedBaseT{Name: "embedded"},
			Value:         99,
		},
		out: &tahwil.Value{
			Refid: 1,
			Kind:  tahwil.Ptr,
			Value: &tahwil.Value{
				Refid: 0,
				Kind:  tahwil.Struct,
				Value: map[string]*tahwil.Value{
					"name": {
						Refid: 0,
						Kind:  tahwil.String,
						Value: "embedded",
					},
					"value": {
						Refid: 0,
						Kind:  tahwil.Int,
						Value: 99,
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
		in: &[3]int{10, 20, 30},
		out: &tahwil.Value{
			Refid: 1,
			Kind:  tahwil.Ptr,
			Value: &tahwil.Value{
				Refid: 0,
				Kind:  tahwil.Array,
				Value: []*tahwil.Value{
					{Refid: 0, Kind: tahwil.Int, Value: 10},
					{Refid: 0, Kind: tahwil.Int, Value: 20},
					{Refid: 0, Kind: tahwil.Int, Value: 30},
				},
			},
		},
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

func TestToValueCompat(t *testing.T) {
	// ToValueCompat should assign refids to all values, including scalars
	v, err := tahwil.ToValueCompat(&parentSerT{Name: "test"})
	if err != nil {
		t.Fatal(err)
	}
	// outer Ptr gets refid 1
	if v.Refid != 1 || v.Kind != tahwil.Ptr {
		t.Fatalf("expected Ptr with refid 1, got %v with refid %d", v.Kind, v.Refid)
	}
	inner := v.Value.(*tahwil.Value)
	// inner Struct gets refid 2 (all values get refids in compat mode)
	if inner.Refid != 2 || inner.Kind != tahwil.Struct {
		t.Fatalf("expected Struct with refid 2, got %v with refid %d", inner.Kind, inner.Refid)
	}
	fields := inner.Value.(map[string]*tahwil.Value)
	// Name field gets refid 3
	if fields["Name"].Refid != 3 {
		t.Errorf("expected Name refid 3, got %d", fields["Name"].Refid)
	}
	// Children field gets refid 4
	if fields["Children"].Refid != 4 {
		t.Errorf("expected Children refid 4, got %d", fields["Children"].Refid)
	}
}

package tahwil_test

import (
	"encoding/json"
	"reflect"
	"testing"
	"unsafe"

	. "github.com/go-extras/tahwil"
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
	Value        interface{}
	NoSerialize1 bool `json:"-"`
	NoSerialize2 bool `json:"_"`
}

type valueTest struct {
	in  interface{}
	out *Value
	err interface{}
}

func valueTests() []valueTest {
	result := make([]valueTest, 0)

	result = append(result, valueTest{
		in: true,
		out: &Value{
			Refid: 1,
			Kind:  Ptr,
			Value: &Value{
				Refid: 2,
				Kind:  Bool,
				Value: true,
			},
		},
	})

	result = append(result, valueTest{
		in: nil,
		out: &Value{
			Refid: 1,
			Kind:  Ptr,
			Value: nil,
		},
	})

	result = append(result, valueTest{
		in: "test",
		out: &Value{
			Refid: 1,
			Kind:  Ptr,
			Value: &Value{
				Refid: 2,
				Kind:  String,
				Value: "test",
			},
		},
	})

	s := "test"
	result = append(result, valueTest{
		in: &s,
		out: &Value{
			Refid: 1,
			Kind:  Ptr,
			Value: &Value{
				Refid: 2,
				Kind:  String,
				Value: "test",
			},
		},
	})

	result = append(result, valueTest{
		in: &interfaceST{Value: "test"},
		out: &Value{
			Refid: 1,
			Kind:  Ptr,
			Value: &Value{
				Refid: 2,
				Kind:  Struct,
				Value: map[string]*Value{
					"Value": {
						Refid: 3,
						Kind:  String,
						Value: "test",
					},
				},
			},
		},
	})

	result = append(result, valueTest{
		in: interface{}("test"),
		out: &Value{
			Refid: 1,
			Kind:  Ptr,
			Value: &Value{
				Refid: 2,
				Kind:  String,
				Value: "test",
			},
		},
	})

	result = append(result, valueTest{
		in: int(47),
		out: &Value{
			Refid: 1,
			Kind:  Ptr,
			Value: &Value{
				Refid: 2,
				Kind:  Int,
				Value: int(47),
			},
		},
	})
	result = append(result, valueTest{
		in: int8(47),
		out: &Value{
			Refid: 1,
			Kind:  Ptr,
			Value: &Value{
				Refid: 2,
				Kind:  Int8,
				Value: int8(47),
			},
		},
	})
	result = append(result, valueTest{
		in: int16(47),
		out: &Value{
			Refid: 1,
			Kind:  Ptr,
			Value: &Value{
				Refid: 2,
				Kind:  Int16,
				Value: int16(47),
			},
		},
	})
	result = append(result, valueTest{
		in: int32(47),
		out: &Value{
			Refid: 1,
			Kind:  Ptr,
			Value: &Value{
				Refid: 2,
				Kind:  Int32,
				Value: int32(47),
			},
		},
	})
	result = append(result, valueTest{
		in: int64(47),
		out: &Value{
			Refid: 1,
			Kind:  Ptr,
			Value: &Value{
				Refid: 2,
				Kind:  Int64,
				Value: int64(47),
			},
		},
	})

	result = append(result, valueTest{
		in: uint(47),
		out: &Value{
			Refid: 1,
			Kind:  Ptr,
			Value: &Value{
				Refid: 2,
				Kind:  Uint,
				Value: uint(47),
			},
		},
	})
	result = append(result, valueTest{
		in: uint8(47),
		out: &Value{
			Refid: 1,
			Kind:  Ptr,
			Value: &Value{
				Refid: 2,
				Kind:  Uint8,
				Value: uint8(47),
			},
		},
	})
	result = append(result, valueTest{
		in: byte(47),
		out: &Value{
			Refid: 1,
			Kind:  Ptr,
			Value: &Value{
				Refid: 2,
				Kind:  Uint8,
				Value: uint8(47),
			},
		},
	})
	result = append(result, valueTest{
		in: uint16(47),
		out: &Value{
			Refid: 1,
			Kind:  Ptr,
			Value: &Value{
				Refid: 2,
				Kind:  Uint16,
				Value: uint16(47),
			},
		},
	})
	result = append(result, valueTest{
		in: uint32(47),
		out: &Value{
			Refid: 1,
			Kind:  Ptr,
			Value: &Value{
				Refid: 2,
				Kind:  Uint32,
				Value: uint32(47),
			},
		},
	})
	result = append(result, valueTest{
		in: float32(47.47),
		out: &Value{
			Refid: 1,
			Kind:  Ptr,
			Value: &Value{
				Refid: 2,
				Kind:  Float32,
				Value: float32(47.47),
			},
		},
	})
	result = append(result, valueTest{
		in: float64(47.47),
		out: &Value{
			Refid: 1,
			Kind:  Ptr,
			Value: &Value{
				Refid: 2,
				Kind:  Float64,
				Value: float64(47.47),
			},
		},
	})

	result = append(result, valueTest{
		in: &parentSerT{
			Name: "Patrik",
		},
		out: &Value{
			Refid: 1,
			Kind:  Ptr,
			Value: &Value{
				Refid: 2,
				Kind:  Struct,
				Value: map[string]*Value{
					"Name": {
						Refid: 3,
						Kind:  String,
						Value: "Patrik",
					},
					"Children": {
						Refid: 4,
						Kind:  Slice,
						Value: []*Value{},
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
		out: &Value{
			Refid: 1,
			Kind:  Ptr,
			Value: &Value{
				Refid: 2,
				Kind:  Struct,
				Value: map[string]*Value{
					"Name": {
						Refid: 3,
						Kind:  String,
						Value: "Patrik",
					},
					"Children": {
						Refid: 4,
						Kind:  Slice,
						Value: []*Value{
							{
								Refid: 5,
								Kind:  Ptr,
								Value: &Value{
									Refid: 6,
									Kind:  Struct,
									Value: map[string]*Value{
										"Name": {
											Refid: 7,
											Kind:  String,
											Value: "Valentine",
										},
										"Parent": {
											Refid: 8,
											Kind:  Ref,
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
		out: &Value{
			Refid: 1,
			Kind:  Ptr,
			Value: &Value{
				Refid: 2,
				Kind:  Struct,
				Value: map[string]*Value{
					"Name": {
						Refid: 3,
						Kind:  String,
						Value: "Klark",
					},
					"Self": {
						Refid: 4,
						Kind:  Ref,
						Value: uint64(1),
					},
				},
			},
		},
	})

	result = append(result, valueTest{
		in: map[string]interface{}{
			"Id": uint64(1),
		},
		out: &Value{
			Refid: 1,
			Kind:  Ptr,
			Value: &Value{
				Refid: 2,
				Kind:  Map,
				Value: map[string]*Value{
					"Id": {
						Refid: 3,
						Kind:  Uint64,
						Value: uint64(1),
					},
				},
			},
		},
	})

	result = append(result, valueTest{
		in:  uintptr(1),
		err: &InvalidMapperKindError{Kind: "uintptr"},
	})

	result = append(result, valueTest{
		in:  [4]int{1, 2, 3, 4},
		err: &InvalidMapperKindError{Kind: Array},
	})

	result = append(result, valueTest{
		in:  make(chan interface{}),
		err: &InvalidMapperKindError{Kind: "chan"},
	})

	dummy1 := true
	result = append(result, valueTest{
		in:  unsafe.Pointer(&dummy1),
		err: &InvalidMapperKindError{Kind: "unsafe.Pointer"},
	})

	dummy2 := &interfaceST{
		Value:        uintptr(1),
		NoSerialize1: true,
		NoSerialize2: true,
	}
	result = append(result, valueTest{
		in:  &dummy2,
		err: &InvalidMapperKindError{Kind: "uintptr"},
	})

	return result
}

func TestToValue(t *testing.T) {
	for i, arg := range valueTests() {
		v, err := ToValue(arg.in)
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

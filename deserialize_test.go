package tahwil_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/go-extras/tahwil"
)

type personT struct {
	Name     string     `json:"name"`
	Parent   *personT   `json:"parent"`
	Children []*personT `json:"children"`
}

type alltypesT struct {
	Bool    bool
	Int     int
	Uint    uint
	Float   float64
	Slice   []int
	Map     map[string]int
	String  string
	Pointer *string
}

type fromValueTest struct {
	in  *tahwil.Value
	out any
	err string
}

func fromValueTests() []fromValueTest {
	result := make([]fromValueTest, 0)

	str := "xxx"
	alltypes := &alltypesT{
		String:  "string",
		Slice:   []int{1, 2},
		Map:     map[string]int{"1": 1, "2": 2},
		Pointer: &str,
		Bool:    true,
		Float:   42.42,
		Int:     42,
		Uint:    42,
	}

	result = append(result, fromValueTest{
		in: &tahwil.Value{
			Refid: 1,
			Kind:  tahwil.Ptr,
			Value: &tahwil.Value{
				Refid: 2,
				Kind:  tahwil.Struct,
				Value: map[string]*tahwil.Value{
					"String": {
						Refid: 3,
						Kind:  tahwil.String,
						Value: "string",
					},
					"Slice": {
						Refid: 4,
						Kind:  tahwil.Slice,
						Value: []*tahwil.Value{
							{
								Refid: 5,
								Kind:  tahwil.Int,
								Value: int64(1),
							},
							{
								Refid: 6,
								Kind:  tahwil.Int,
								Value: int64(2),
							},
						},
					},
					"Map": {
						Refid: 7,
						Kind:  tahwil.Map,
						Value: map[string]*tahwil.Value{
							"1": {
								Refid: 8,
								Kind:  tahwil.Int,
								Value: int64(1),
							},
							"2": {
								Refid: 9,
								Kind:  tahwil.Int,
								Value: int64(2),
							},
						},
					},
					"Pointer": {
						Refid: 10,
						Kind:  tahwil.Ptr,
						Value: &tahwil.Value{
							Refid: 11,
							Kind:  tahwil.String,
							Value: "xxx",
						},
					},
					"Bool": {
						Refid: 12,
						Kind:  tahwil.Bool,
						Value: true,
					},
					"Float": {
						Refid: 13,
						Kind:  tahwil.Float64,
						Value: 42.42,
					},
					"Int": {
						Refid: 14,
						Kind:  tahwil.Int,
						Value: int64(42),
					},
					"Uint": {
						Refid: 15,
						Kind:  tahwil.Uint,
						Value: uint64(42),
					},
				},
			},
		},
		out: alltypes,
	})

	// empty case
	alltypes = &alltypesT{
		String:  "",
		Slice:   nil,
		Map:     nil,
		Pointer: nil,
		Bool:    false,
		Float:   0,
		Int:     0,
		Uint:    0,
	}

	result = append(result, fromValueTest{
		in: &tahwil.Value{
			Refid: 1,
			Kind:  tahwil.Ptr,
			Value: &tahwil.Value{
				Refid: 2,
				Kind:  tahwil.Struct,
				Value: map[string]*tahwil.Value{
					"String": {
						Refid: 3,
						Kind:  tahwil.String,
						Value: "",
					},
					"Slice": {
						Refid: 4,
						Kind:  tahwil.Slice,
						Value: nil,
					},
					"Map": {
						Refid: 7,
						Kind:  tahwil.Map,
						Value: nil,
					},
					"Pointer": {
						Refid: 10,
						Kind:  tahwil.Ptr,
						Value: nil,
					},
					"Bool": {
						Refid: 12,
						Kind:  tahwil.Bool,
						Value: false,
					},
					"Float": {
						Refid: 13,
						Kind:  tahwil.Float64,
						Value: 0.0,
					},
					"Int": {
						Refid: 14,
						Kind:  tahwil.Int,
						Value: int64(0),
					},
					"Uint": {
						Refid: 15,
						Kind:  tahwil.Uint,
						Value: uint64(0),
					},
				},
			},
		},
		out: alltypes,
	})

	result = append(result, fromValueTest{
		in: &tahwil.Value{
			Refid: 1,
			Kind:  tahwil.Ptr,
			Value: &tahwil.Value{
				Refid: 2,
				Kind:  tahwil.Struct,
				Value: map[string]*tahwil.Value{
					"name": {
						Refid: 3,
						Kind:  tahwil.String,
						Value: "Martin",
					},
					"parent": {
						Refid: 4,
						Kind:  tahwil.Ptr,
						Value: nil,
					},
					"children": {
						Refid: 5,
						Kind:  tahwil.Slice,
						Value: nil,
					},
				},
			},
		},
		out: &personT{
			Name: "Martin",
		},
	})

	// json tags with options like omitempty should use only the name part
	result = append(result, fromValueTest{
		in: &tahwil.Value{
			Refid: 1,
			Kind:  tahwil.Ptr,
			Value: &tahwil.Value{
				Refid: 2,
				Kind:  tahwil.Struct,
				Value: map[string]*tahwil.Value{
					"name": {
						Refid: 3,
						Kind:  tahwil.String,
						Value: "test",
					},
					"value": {
						Refid: 4,
						Kind:  tahwil.Int,
						Value: int64(42),
					},
				},
			},
		},
		out: &omitemptyT{Name: "test", Value: 42},
	})

	// embedded struct fields should deserialize from flat keys
	result = append(result, fromValueTest{
		in: &tahwil.Value{
			Refid: 1,
			Kind:  tahwil.Ptr,
			Value: &tahwil.Value{
				Refid: 2,
				Kind:  tahwil.Struct,
				Value: map[string]*tahwil.Value{
					"name": {
						Refid: 3,
						Kind:  tahwil.String,
						Value: "embedded",
					},
					"value": {
						Refid: 4,
						Kind:  tahwil.Int,
						Value: int64(99),
					},
				},
			},
		},
		out: &embeddedOuterT{
			embeddedBaseT: embeddedBaseT{Name: "embedded"},
			Value:         99,
		},
	})

	// array deserialization
	arr := &[3]int{10, 20, 30}
	result = append(result, fromValueTest{
		in: &tahwil.Value{
			Refid: 1,
			Kind:  tahwil.Ptr,
			Value: &tahwil.Value{
				Refid: 2,
				Kind:  tahwil.Array,
				Value: []*tahwil.Value{
					{Refid: 3, Kind: tahwil.Int, Value: int64(10)},
					{Refid: 4, Kind: tahwil.Int, Value: int64(20)},
					{Refid: 5, Kind: tahwil.Int, Value: int64(30)},
				},
			},
		},
		out: arr,
	})

	p1 := &personT{
		Name: "Martin",
		Parent: &personT{
			Name: "Kevin",
		},
	}
	p1.Parent.Children = append(p1.Parent.Children, p1)
	result = append(result, fromValueTest{
		in: &tahwil.Value{
			Refid: 1,
			Kind:  tahwil.Ptr,
			Value: &tahwil.Value{
				Refid: 2,
				Kind:  tahwil.Struct,
				Value: map[string]*tahwil.Value{
					"name": {
						Refid: 3,
						Kind:  tahwil.String,
						Value: "Martin",
					},
					"parent": {
						Refid: 4,
						Kind:  tahwil.Ptr,
						Value: &tahwil.Value{
							Refid: 5,
							Kind:  tahwil.Struct,
							Value: map[string]*tahwil.Value{
								"name": {
									Refid: 6,
									Kind:  tahwil.String,
									Value: "Kevin",
								},
								"children": {
									Refid: 7,
									Kind:  tahwil.Slice,
									Value: []*tahwil.Value{
										{
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
		out: p1,
	})

	str1 := "test"
	result = append(result, fromValueTest{
		in: &tahwil.Value{
			Refid: 1,
			Kind:  tahwil.Ptr,
			Value: &tahwil.Value{
				Refid: 2,
				Kind:  tahwil.String,
				Value: 0,
			},
		},
		out: &str1,
		err: `tahwil.Value: invalid value int(0) for kind "string"`,
	})

	v1 := 0.0
	result = append(result, fromValueTest{
		in: &tahwil.Value{
			Refid: 1,
			Kind:  tahwil.Ptr,
			Value: &tahwil.Value{
				Refid: 2,
				Kind:  tahwil.Float64,
				Value: "xxx",
			},
		},
		out: &v1,
		err: `tahwil.Value: invalid value string("xxx") for kind "float64"`,
	})

	v2 := 0
	result = append(result, fromValueTest{
		in: &tahwil.Value{
			Refid: 1,
			Kind:  tahwil.Ptr,
			Value: &tahwil.Value{
				Refid: 2,
				Kind:  tahwil.String,
				Value: 0,
			},
		},
		out: &v2,
		err: `tahwil.FromValue: unexpected kind (expected: string, got: int)`,
	})

	v3 := ""
	result = append(result, fromValueTest{
		in: &tahwil.Value{
			Refid: 1,
			Kind:  tahwil.Ptr,
			Value: &tahwil.Value{
				Refid: 2,
				Kind:  tahwil.String,
				Value: struct{}{},
			},
		},
		out: &v3,
		err: `tahwil.Value: invalid value struct {}(struct {}{}) for kind "string"`,
	})

	// error inside a slice element should propagate
	sliceTarget := &alltypesT{}
	result = append(result, fromValueTest{
		in: &tahwil.Value{
			Refid: 1,
			Kind:  tahwil.Ptr,
			Value: &tahwil.Value{
				Refid: 2,
				Kind:  tahwil.Struct,
				Value: map[string]*tahwil.Value{
					"String": {Refid: 3, Kind: tahwil.String, Value: ""},
					"Bool":   {Refid: 4, Kind: tahwil.Bool, Value: false},
					"Int":    {Refid: 5, Kind: tahwil.Int, Value: int64(0)},
					"Uint":   {Refid: 6, Kind: tahwil.Uint, Value: uint64(0)},
					"Float":  {Refid: 7, Kind: tahwil.Float64, Value: 0.0},
					"Pointer": {Refid: 8, Kind: tahwil.Ptr, Value: nil},
					"Map":    {Refid: 9, Kind: tahwil.Map, Value: nil},
					"Slice": {
						Refid: 10,
						Kind:  tahwil.Slice,
						Value: []*tahwil.Value{
							{
								Refid: 11,
								Kind:  tahwil.Int,
								Value: "not-an-int", // wrong type: string instead of int
							},
						},
					},
				},
			},
		},
		out: sliceTarget,
		err: `tahwil.Value: invalid value string("not-an-int") for kind "int"`,
	})

	// error inside a map value should propagate
	mapTarget := &alltypesT{}
	result = append(result, fromValueTest{
		in: &tahwil.Value{
			Refid: 1,
			Kind:  tahwil.Ptr,
			Value: &tahwil.Value{
				Refid: 2,
				Kind:  tahwil.Struct,
				Value: map[string]*tahwil.Value{
					"String": {Refid: 3, Kind: tahwil.String, Value: ""},
					"Bool":   {Refid: 4, Kind: tahwil.Bool, Value: false},
					"Int":    {Refid: 5, Kind: tahwil.Int, Value: int64(0)},
					"Uint":   {Refid: 6, Kind: tahwil.Uint, Value: uint64(0)},
					"Float":  {Refid: 7, Kind: tahwil.Float64, Value: 0.0},
					"Pointer": {Refid: 8, Kind: tahwil.Ptr, Value: nil},
					"Slice":  {Refid: 9, Kind: tahwil.Slice, Value: nil},
					"Map": {
						Refid: 10,
						Kind:  tahwil.Map,
						Value: map[string]*tahwil.Value{
							"key1": {
								Refid: 11,
								Kind:  tahwil.Int,
								Value: "not-an-int", // wrong type: string instead of int
							},
						},
					},
				},
			},
		},
		out: mapTarget,
		err: `tahwil.Value: invalid value string("not-an-int") for kind "int"`,
	})

	return result
}

func TestUnmapperError_Unwrap(t *testing.T) {
	// Trigger a resolver error via an invalid ref value type (string instead of numeric).
	data := &tahwil.Value{
		Refid: 1,
		Kind:  tahwil.Ref,
		Value: "not-a-number",
	}
	var target string
	err := tahwil.FromValue(data, &target)
	if err == nil {
		t.Fatal("expected an error, got nil")
	}

	// errors.As should reach the underlying ResolverError
	var resolverErr *tahwil.ResolverError
	if !errors.As(err, &resolverErr) {
		t.Fatalf("expected errors.As to find *ResolverError, got %T: %v", err, err)
	}

	// errors.As should also match UnmapperError itself
	var unmapperErr *tahwil.UnmapperError
	if !errors.As(err, &unmapperErr) {
		t.Fatalf("expected errors.As to find *UnmapperError, got %T: %v", err, err)
	}

	// Unwrap should return the ResolverError directly
	if unmapperErr.Unwrap() == nil {
		t.Fatal("expected Unwrap() to return non-nil error")
	}
}

func TestFromValue(t *testing.T) {
	for i, arg := range fromValueTests() {
		p := reflect.New(reflect.TypeOf(arg.out).Elem()).Interface()
		err := tahwil.FromValue(arg.in, p)
		if err != nil {
			if err.Error() != arg.err {
				t.Fatalf("#%d: unexpected error: %#+v (%s)", i, err, err)
			}
		} else if arg.err != "" {
			t.Fatalf("#%d: expected error %q, but got nil", i, arg.err)
		} else {
			if !reflect.DeepEqual(p, arg.out) {
				t.Errorf("#%d: mismatch\nhave: %#+v\nwant: %#+v", i, reflect.ValueOf(p).Elem(), reflect.ValueOf(arg.out).Elem())
			}
		}
	}
}

func TestUnmarshal(t *testing.T) {
	// string
	strData := &tahwil.Value{
		Refid: 1,
		Kind:  tahwil.Ptr,
		Value: &tahwil.Value{
			Refid: 2,
			Kind:  tahwil.String,
			Value: "hello",
		},
	}
	s, err := tahwil.Unmarshal[string](strData)
	if err != nil {
		t.Fatalf("Unmarshal[string]: %v", err)
	}
	if *s != "hello" {
		t.Errorf("Unmarshal[string] = %q, want %q", *s, "hello")
	}

	// struct
	structData := &tahwil.Value{
		Refid: 1,
		Kind:  tahwil.Ptr,
		Value: &tahwil.Value{
			Refid: 2,
			Kind:  tahwil.Struct,
			Value: map[string]*tahwil.Value{
				"name": {Refid: 3, Kind: tahwil.String, Value: "Alice"},
				"parent": {Refid: 4, Kind: tahwil.Ptr, Value: nil},
				"children": {Refid: 5, Kind: tahwil.Slice, Value: nil},
			},
		},
	}
	p, err := tahwil.Unmarshal[personT](structData)
	if err != nil {
		t.Fatalf("Unmarshal[personT]: %v", err)
	}
	if p.Name != "Alice" {
		t.Errorf("Unmarshal[personT].Name = %q, want %q", p.Name, "Alice")
	}

	// error propagation
	badData := &tahwil.Value{
		Refid: 1,
		Kind:  tahwil.Ptr,
		Value: &tahwil.Value{
			Refid: 2,
			Kind:  tahwil.String,
			Value: 42,
		},
	}
	_, err = tahwil.Unmarshal[string](badData)
	if err == nil {
		t.Fatal("expected error from Unmarshal with bad data, got nil")
	}
}

// TestFromValue_ForwardRef verifies that a Ref pointing to a node encountered
// later in the DFS walk (forward reference) is resolved correctly.
func TestFromValue_ForwardRef(t *testing.T) {
	// children[0] is a forward Ref to refid 7 (the Ptr node
	// wrapping childB), which appears at children[1]. Slices guarantee
	// deterministic iteration order, so children[0] is visited before children[1].
	forwardData := &tahwil.Value{
		Refid: 1,
		Kind:  tahwil.Ptr,
		Value: &tahwil.Value{
			Refid: 2,
			Kind:  tahwil.Struct,
			Value: map[string]*tahwil.Value{
				"name": {Refid: 3, Kind: tahwil.String, Value: "root"},
				"parent": {Refid: 4, Kind: tahwil.Ptr, Value: nil},
				"children": {
					Refid: 5,
					Kind:  tahwil.Slice,
					Value: []*tahwil.Value{
						{
							Refid: 6,
							Kind:  tahwil.Ref,
							Value: uint64(7), // forward ref to the Ptr node at children[1]
						},
						{
							Refid: 7,
							Kind:  tahwil.Ptr,
							Value: &tahwil.Value{
								Refid: 8,
								Kind:  tahwil.Struct,
								Value: map[string]*tahwil.Value{
									"name":     {Refid: 9, Kind: tahwil.String, Value: "childB"},
									"parent":   {Refid: 10, Kind: tahwil.Ptr, Value: nil},
									"children": {Refid: 11, Kind: tahwil.Slice, Value: nil},
								},
							},
						},
					},
				},
			},
		},
	}

	var result personT
	if err := tahwil.FromValue(forwardData, &result); err != nil {
		t.Fatalf("FromValue with forward ref: %v", err)
	}

	if result.Name != "root" {
		t.Errorf("Name = %q, want %q", result.Name, "root")
	}
	if len(result.Children) != 2 {
		t.Fatalf("len(Children) = %d, want 2", len(result.Children))
	}
	// children[0] was a forward Ref to refid 7 (the Ptr at children[1])
	if result.Children[0] == nil {
		t.Fatal("Children[0] is nil, expected resolved forward ref")
	}
	if result.Children[0].Name != "childB" {
		t.Errorf("Children[0].Name = %q, want %q", result.Children[0].Name, "childB")
	}
	if result.Children[1] == nil {
		t.Fatal("Children[1] is nil")
	}
	if result.Children[1].Name != "childB" {
		t.Errorf("Children[1].Name = %q, want %q", result.Children[1].Name, "childB")
	}
}

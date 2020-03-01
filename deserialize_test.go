package tahwil_test

import (
	"reflect"
	"testing"

	. "github.com/go-extras/tahwil"
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
	in  *Value
	out interface{}
	err string
}

func fromValueTests() []fromValueTest {
	result := make([]fromValueTest, 0)

	str := "xxx"
	alltypes := &alltypesT{
		String:  String,
		Slice:   []int{1, 2},
		Map:     map[string]int{"1": 1, "2": 2},
		Pointer: &str,
		Bool:    true,
		Float:   42.42,
		Int:     42,
		Uint:    42,
	}

	result = append(result, fromValueTest{
		in: &Value{
			Refid: 1,
			Kind:  Ptr,
			Value: &Value{
				Refid: 2,
				Kind:  Struct,
				Value: map[string]*Value{
					"String": {
						Refid: 3,
						Kind:  String,
						Value: String,
					},
					"Slice": {
						Refid: 4,
						Kind:  Slice,
						Value: []*Value{
							{
								Refid: 5,
								Kind:  Int,
								Value: int64(1),
							},
							{
								Refid: 6,
								Kind:  Int,
								Value: int64(2),
							},
						},
					},
					"Map": {
						Refid: 7,
						Kind:  Map,
						Value: map[string]*Value{
							"1": {
								Refid: 8,
								Kind:  Int,
								Value: int64(1),
							},
							"2": {
								Refid: 9,
								Kind:  Int,
								Value: int64(2),
							},
						},
					},
					"Pointer": {
						Refid: 10,
						Kind:  Ptr,
						Value: &Value{
							Refid: 11,
							Kind:  String,
							Value: "xxx",
						},
					},
					"Bool": {
						Refid: 12,
						Kind:  Bool,
						Value: true,
					},
					"Float": {
						Refid: 13,
						Kind:  Float64,
						Value: 42.42,
					},
					"Int": {
						Refid: 14,
						Kind:  Int,
						Value: int64(42),
					},
					"Uint": {
						Refid: 15,
						Kind:  Uint,
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
		in: &Value{
			Refid: 1,
			Kind:  Ptr,
			Value: &Value{
				Refid: 2,
				Kind:  Struct,
				Value: map[string]*Value{
					"String": {
						Refid: 3,
						Kind:  String,
						Value: "",
					},
					"Slice": {
						Refid: 4,
						Kind:  Slice,
						Value: nil,
					},
					"Map": {
						Refid: 7,
						Kind:  Map,
						Value: nil,
					},
					"Pointer": {
						Refid: 10,
						Kind:  Ptr,
						Value: nil,
					},
					"Bool": {
						Refid: 12,
						Kind:  Bool,
						Value: false,
					},
					"Float": {
						Refid: 13,
						Kind:  Float64,
						Value: 0.0,
					},
					Int: {
						Refid: 14,
						Kind:  Int,
						Value: int64(0),
					},
					Uint: {
						Refid: 15,
						Kind:  Uint,
						Value: uint64(0),
					},
				},
			},
		},
		out: alltypes,
	})

	result = append(result, fromValueTest{
		in: &Value{
			Refid: 1,
			Kind:  Ptr,
			Value: &Value{
				Refid: 2,
				Kind:  Struct,
				Value: map[string]*Value{
					"name": {
						Refid: 3,
						Kind:  String,
						Value: "Martin",
					},
					"parent": {
						Refid: 4,
						Kind:  Ptr,
						Value: nil,
					},
					"children": {
						Refid: 5,
						Kind:  Slice,
						Value: nil,
					},
				},
			},
		},
		out: &personT{
			Name: "Martin",
		},
	})

	p1 := &personT{
		Name: "Martin",
		Parent: &personT{
			Name: "Kevin",
		},
	}
	p1.Parent.Children = append(p1.Parent.Children, p1)
	result = append(result, fromValueTest{
		in: &Value{
			Refid: 1,
			Kind:  Ptr,
			Value: &Value{
				Refid: 2,
				Kind:  Struct,
				Value: map[string]*Value{
					"name": {
						Refid: 3,
						Kind:  String,
						Value: "Martin",
					},
					"parent": {
						Refid: 4,
						Kind:  Ptr,
						Value: &Value{
							Refid: 5,
							Kind:  Struct,
							Value: map[string]*Value{
								"name": {
									Refid: 6,
									Kind:  String,
									Value: "Kevin",
								},
								"children": {
									Refid: 7,
									Kind:  Slice,
									Value: []*Value{
										{
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
		out: p1,
	})

	str1 := "test"
	result = append(result, fromValueTest{
		in: &Value{
			Refid: 1,
			Kind:  Ptr,
			Value: &Value{
				Refid: 2,
				Kind:  String,
				Value: 0,
			},
		},
		out: &str1,
		err: `tahwil.Value: invalid value int(0) for kind "string"`,
	})

	v1 := 0.0
	result = append(result, fromValueTest{
		in: &Value{
			Refid: 1,
			Kind:  Ptr,
			Value: &Value{
				Refid: 2,
				Kind:  Float64,
				Value: "xxx",
			},
		},
		out: &v1,
		err: `tahwil.Value: invalid value string("xxx") for kind "float64"`,
	})

	v2 := 0
	result = append(result, fromValueTest{
		in: &Value{
			Refid: 1,
			Kind:  Ptr,
			Value: &Value{
				Refid: 2,
				Kind:  String,
				Value: 0,
			},
		},
		out: &v2,
		err: `tahwil.FromValue: unexpected kind (expected: string, got: int)`,
	})

	v3 := ""
	result = append(result, fromValueTest{
		in: &Value{
			Refid: 1,
			Kind:  Ptr,
			Value: &Value{
				Refid: 2,
				Kind:  String,
				Value: struct{}{},
			},
		},
		out: &v3,
		err: `tahwil.Value: invalid value struct {}(struct {}{}) for kind "string"`,
	})

	return result
}

func TestFromValue(t *testing.T) {
	for i, arg := range fromValueTests() {
		p := reflect.New(reflect.TypeOf(arg.out).Elem()).Interface()
		err := FromValue(arg.in, p)
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

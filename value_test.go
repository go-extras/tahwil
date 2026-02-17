package tahwil_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/go-extras/tahwil"
)

type unmarshalJSONTest struct {
	in  string
	out *tahwil.Value
	err any
}

func unmarshalJSONTests() []unmarshalJSONTest {
	res := make([]unmarshalJSONTest, 0)

	res = append(res, unmarshalJSONTest{in: `null`, out: &tahwil.Value{}})
	res = append(res, unmarshalJSONTest{in: `{}`, out: &tahwil.Value{}})
	res = append(res, unmarshalJSONTest{in: `{
		"refid": 1,
		"kind": "string",
		"value": "aaa"
}`, out: &tahwil.Value{
		Refid: 1,
		Kind:  tahwil.String,
		Value: "aaa",
	}})
	res = append(res, unmarshalJSONTest{in: `{
		"refid": 1,
		"kind": "bool",
		"value": true
}`, out: &tahwil.Value{
		Refid: 1,
		Kind:  tahwil.Bool,
		Value: true,
	}})
	res = append(res, unmarshalJSONTest{in: `{
		"refid": 1,
		"kind": "int",
		"value": 1
}`, out: &tahwil.Value{
		Refid: 1,
		Kind:  tahwil.Int,
		Value: 1,
	}})
	res = append(res, unmarshalJSONTest{in: `{
		"refid": 1,
		"kind": "int8",
		"value": 1
}`, out: &tahwil.Value{
		Refid: 1,
		Kind:  tahwil.Int8,
		Value: int8(1),
	}})
	res = append(res, unmarshalJSONTest{in: `{
		"refid": 1,
		"kind": "int16",
		"value": 1
}`, out: &tahwil.Value{
		Refid: 1,
		Kind:  tahwil.Int16,
		Value: int16(1),
	}})
	res = append(res, unmarshalJSONTest{in: `{
		"refid": 1,
		"kind": "int32",
		"value": 1
}`, out: &tahwil.Value{
		Refid: 1,
		Kind:  tahwil.Int32,
		Value: int32(1),
	}})
	res = append(res, unmarshalJSONTest{in: `{
		"refid": 1,
		"kind": "int64",
		"value": 1
}`, out: &tahwil.Value{
		Refid: 1,
		Kind:  tahwil.Int64,
		Value: int64(1),
	}})
	res = append(res, unmarshalJSONTest{in: `{
		"refid": 1,
		"kind": "uint",
		"value": 1
}`, out: &tahwil.Value{
		Refid: 1,
		Kind:  tahwil.Uint,
		Value: uint(1),
	}})
	res = append(res, unmarshalJSONTest{in: `{
		"refid": 1,
		"kind": "uint8",
		"value": 1
}`, out: &tahwil.Value{
		Refid: 1,
		Kind:  tahwil.Uint8,
		Value: uint8(1),
	}})
	res = append(res, unmarshalJSONTest{in: `{
		"refid": 1,
		"kind": "uint16",
		"value": 1
}`, out: &tahwil.Value{
		Refid: 1,
		Kind:  tahwil.Uint16,
		Value: uint16(1),
	}})
	res = append(res, unmarshalJSONTest{in: `{
		"refid": 1,
		"kind": "uint32",
		"value": 1
}`, out: &tahwil.Value{
		Refid: 1,
		Kind:  tahwil.Uint32,
		Value: uint32(1),
	}})
	res = append(res, unmarshalJSONTest{in: `{
		"refid": 1,
		"kind": "uint64",
		"value": 1
}`, out: &tahwil.Value{
		Refid: 1,
		Kind:  tahwil.Uint64,
		Value: uint64(1),
	}})
	res = append(res, unmarshalJSONTest{in: `{
		"refid": 1,
		"kind": "float32",
		"value": 1
}`, out: &tahwil.Value{
		Refid: 1,
		Kind:  tahwil.Float32,
		Value: float32(1),
	}})
	res = append(res, unmarshalJSONTest{in: `{
		"refid": 1,
		"kind": "float64",
		"value": 1
}`, out: &tahwil.Value{
		Refid: 1,
		Kind:  tahwil.Float64,
		Value: float64(1),
	}})
	res = append(res, unmarshalJSONTest{in: `{
		"refid": 1,
		"kind": "ptr",
		"value": {
			"refid": 2,
			"kind": "string",
			"value": "test"
		}
}`, out: &tahwil.Value{
		Refid: 1,
		Kind:  tahwil.Ptr,
		Value: &tahwil.Value{
			Refid: 2,
			Kind:  tahwil.String,
			Value: "test",
		},
	}})
	res = append(res, unmarshalJSONTest{in: `{
		"refid": 1,
		"kind": "ptr",
		"value": {
			"refid": 2,
			"kind": "struct",
			"value": {
				"name": {
					"refid": 3,
					"kind": "string",
					"value": "Arthur"
				},
				"children": {
					"refid": 4,
					"kind": "slice",
					"value": []
				} 
			}
		}
}`, out: &tahwil.Value{
		Refid: 1,
		Kind:  tahwil.Ptr,
		Value: &tahwil.Value{
			Refid: 2,
			Kind:  tahwil.Struct,
			Value: map[string]any{
				"name": &tahwil.Value{
					Refid: 3,
					Kind:  tahwil.String,
					Value: "Arthur",
				},
				"children": &tahwil.Value{
					Refid: 4,
					Kind:  tahwil.Slice,
					Value: []any{},
				},
			},
		},
	}})
	res = append(res, unmarshalJSONTest{in: `{
		"refid": 1,
		"kind": "ptr",
		"value": {
			"refid": 2,
			"kind": "struct",
			"value": {
				"name": {
					"refid": 3,
					"kind": "string",
					"value": "Arthur"
				},
				"parent": {
					"refid": 4,
					"kind": "ptr",
					"value": null
				},
				"children": {
					"refid": 5,
					"kind": "slice",
					"value": [
						{
							"refid": 6,
							"kind": "ptr",
							"value": {
								"refid": 7,
								"kind": "struct",
								"value": {
									"name": {
										"refid": 8,
										"kind": "string",
										"value": "Trillian"
									},
									"parent": {
										"refid": 9,
										"kind": "ref",
										"value": 1
									},
									"children": {
										"refid": 10,
										"kind": "slice",
										"value": []
									}
								}
							}
						}
					]
				} 
			}
		}
}`, out: &tahwil.Value{
		Refid: 1,
		Kind:  tahwil.Ptr,
		Value: &tahwil.Value{
			Refid: 2,
			Kind:  tahwil.Struct,
			Value: map[string]any{
				"name": &tahwil.Value{
					Refid: 3,
					Kind:  tahwil.String,
					Value: "Arthur",
				},
				"parent": &tahwil.Value{
					Refid: 4,
					Kind:  tahwil.Ptr,
					Value: nil,
				},
				"children": &tahwil.Value{
					Refid: 5,
					Kind:  tahwil.Slice,
					Value: []any{
						&tahwil.Value{
							Refid: 6,
							Kind:  tahwil.Ptr,
							Value: &tahwil.Value{
								Refid: 7,
								Kind:  tahwil.Struct,
								Value: map[string]any{
									"name": &tahwil.Value{
										Refid: 8,
										Kind:  tahwil.String,
										Value: "Trillian",
									},
									"parent": &tahwil.Value{
										Refid: 9,
										Kind:  tahwil.Ref,
										Value: 1,
									},
									"children": &tahwil.Value{
										Refid: 10,
										Kind:  tahwil.Slice,
										Value: []any{},
									},
								},
							},
						},
					},
				},
			},
		},
	}})
	res = append(res, unmarshalJSONTest{
		in: `{
			"refid": 1,
			"kind": "chan",
			"value": "aaa"
		}`,
		out: &tahwil.Value{
			Refid: 1,
			Kind:  "chan",
			Value: "aaa",
		},
		err: &tahwil.InvalidValueKindError{Kind: "chan"},
	})
	res = append(res, unmarshalJSONTest{
		in: `{
			"refid": 1,
			"kind": "dummy",
			"value": "aaa"
		}`,
		out: &tahwil.Value{
			Refid: 1,
			Kind:  "chan",
			Value: "aaa",
		},
		err: &tahwil.InvalidValueKindError{Kind: "dummy"},
	})
	res = append(res, unmarshalJSONTest{
		in: `{
			"refid": 1,
			"kind": "complex64",
			"value": "aaa"
		}`,
		out: &tahwil.Value{
			Refid: 1,
			Kind:  "complex64",
			Value: "aaa",
		},
		err: &tahwil.InvalidValueKindError{Kind: "complex64"},
	})
	res = append(res, unmarshalJSONTest{
		in: `{
			"refid": 1,
			"kind": "complex128",
			"value": "aaa"
		}`,
		out: &tahwil.Value{
			Refid: 1,
			Kind:  "complex128",
			Value: "aaa",
		},
		err: &tahwil.InvalidValueKindError{Kind: "complex128"},
	})
	res = append(res, unmarshalJSONTest{
		in: `{
			"refid": 1,
			"kind": "uintptr",
			"value": "aaa"
		}`,
		err: &tahwil.InvalidValueKindError{Kind: "uintptr"},
	})
	res = append(res, unmarshalJSONTest{
		in: `{
			"refid": 1,
			"kind": "byte",
			"value": "aaa"
		}`,
		err: &tahwil.InvalidValueKindError{Kind: "byte"},
	})
	res = append(res, unmarshalJSONTest{
		in: `{
			"refid": 1,
			"kind": "rune",
			"value": "aaa"
		}`,
		out: &tahwil.Value{
			Refid: 1,
			Kind:  "rune",
			Value: "aaa",
		},
		err: &tahwil.InvalidValueKindError{Kind: "rune"},
	})
	res = append(res, unmarshalJSONTest{
		in: `{
			"refid": 1,
			"kind": "string", 
			"value": {"error"}
		}`,
		err: "invalid character '}' after object key",
	})
	res = append(res, unmarshalJSONTest{
		in: `{
			"refid": 1,
			"kind": "ptr", 
			"value": "invalid"
		}`,
		err: &tahwil.InvalidValueError{Kind: tahwil.Ptr, Value: "invalid"},
	})
	res = append(res, unmarshalJSONTest{
		in: `{
			"refid": 1,
			"kind": "ptr", 
			"value": {
				"refid": 1,
				"kind": "ptr",
				"value": "invalid"
			}
		}`,
		err: &tahwil.InvalidValueError{Kind: tahwil.Ptr, Value: "invalid"},
	})
	res = append(res, unmarshalJSONTest{
		in: `{
			"refid": 1,
			"kind": "struct", 
			"value": "invalid"
		}`,
		err: &tahwil.InvalidValueError{Kind: tahwil.Struct, Value: "invalid"},
	})
	res = append(res, unmarshalJSONTest{
		in: `{
			"refid": 1,
			"kind": "struct", 
			"value": {
				"arg": {
					"refid": 1,
					"kind": "struct", 
					"value": "invalid"
				}
			}
		}`,
		err: &tahwil.InvalidValueError{Kind: tahwil.Struct, Value: "invalid"},
	})
	res = append(res, unmarshalJSONTest{
		in: `{
			"refid": 1,
			"kind": "map", 
			"value": "invalid"
		}`,
		err: &tahwil.InvalidValueError{Kind: tahwil.Map, Value: "invalid"},
	})
	res = append(res, unmarshalJSONTest{
		in: `{
			"refid": 1,
			"kind": "map", 
			"value": {
				"arg": {
					"refid": 1,
					"kind": "map", 
					"value": "invalid"
				}
			}
		}`,
		err: &tahwil.InvalidValueError{Kind: tahwil.Map, Value: "invalid"},
	})
	res = append(res, unmarshalJSONTest{
		in: `{
			"refid": 1,
			"kind": "slice", 
			"value": "invalid"
		}`,
		err: &tahwil.InvalidValueError{Kind: tahwil.Slice, Value: "invalid"},
	})
	res = append(res, unmarshalJSONTest{
		in: `{
			"refid": 1,
			"kind": "slice", 
			"value": [
				{
					"refid": 1,
					"kind": "slice", 
					"value": "invalid"
				}
			]
		}`,
		err: &tahwil.InvalidValueError{Kind: tahwil.Slice, Value: "invalid"},
	})

	return res[0:len(res):len(res)]
}

func TestValue_UnmarshalJSON(t *testing.T) {
	for i, arg := range unmarshalJSONTests() {
		v := &tahwil.Value{}
		err := json.Unmarshal([]byte(arg.in), v)
		if err != nil {
			if serr, ok := err.(*json.SyntaxError); ok {
				if serr.Error() != arg.err {
					t.Fatalf("UnmarshalJSON: %v", err)
				}
			} else if !reflect.DeepEqual(arg.err, err) {
				t.Fatalf("UnmarshalJSON: %v", err)
			}
		} else if !reflect.DeepEqual(arg.out, v) {
			t.Errorf("#%d: mismatch\nhave: %#+v\nwant: %#+v", i, v, arg.out)
			continue
		}
	}
}

func TestInvalidValueKindError_Error(t *testing.T) {
	err := &tahwil.InvalidValueKindError{Kind: "invalid"}
	expected := "tahwil.Value: invalid value kind \"" + err.Kind + "\""
	if err.Error() != expected {
		t.Errorf("mismatch\nhave: %#+v\nwant: %#+v", err.Error(), expected)
	}
}

func TestInvalidValueError_Error(t *testing.T) {
	err := &tahwil.InvalidValueError{
		Value: "val",
		Kind:  "invalid",
	}
	expected := "tahwil.Value: invalid value string(\"val\") for kind \"invalid\""
	if err.Error() != expected {
		t.Errorf("mismatch\nhave: %#+v\nwant: %#+v", err.Error(), expected)
	}
}

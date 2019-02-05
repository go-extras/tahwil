package tahwil

import (
	"encoding/json"
	"reflect"
	"testing"
)

type unmarshalJSONTest struct {
	in  string
	out *Value
	err interface{}
}

func unmarshalJSONTests() []unmarshalJSONTest {
	res := make([]unmarshalJSONTest, 0)

	res = append(res, unmarshalJSONTest{in: `null`, out: &Value{}})
	res = append(res, unmarshalJSONTest{in: `{}`, out: &Value{}})
	res = append(res, unmarshalJSONTest{in: `{
		"refid": 1,
		"kind": "string",
		"value": "aaa"
}`, out: &Value{
		Refid: 1,
		Kind:  "string",
		Value: "aaa",
	}})
	res = append(res, unmarshalJSONTest{in: `{
		"refid": 1,
		"kind": "bool",
		"value": true
}`, out: &Value{
		Refid: 1,
		Kind:  "bool",
		Value: true,
	}})
	res = append(res, unmarshalJSONTest{in: `{
		"refid": 1,
		"kind": "int",
		"value": 1
}`, out: &Value{
		Refid: 1,
		Kind:  "int",
		Value: 1,
	}})
	res = append(res, unmarshalJSONTest{in: `{
		"refid": 1,
		"kind": "int8",
		"value": 1
}`, out: &Value{
		Refid: 1,
		Kind:  "int8",
		Value: int8(1),
	}})
	res = append(res, unmarshalJSONTest{in: `{
		"refid": 1,
		"kind": "int16",
		"value": 1
}`, out: &Value{
		Refid: 1,
		Kind:  "int16",
		Value: int16(1),
	}})
	res = append(res, unmarshalJSONTest{in: `{
		"refid": 1,
		"kind": "int32",
		"value": 1
}`, out: &Value{
		Refid: 1,
		Kind:  "int32",
		Value: int32(1),
	}})
	res = append(res, unmarshalJSONTest{in: `{
		"refid": 1,
		"kind": "int64",
		"value": 1
}`, out: &Value{
		Refid: 1,
		Kind:  "int64",
		Value: int64(1),
	}})
	res = append(res, unmarshalJSONTest{in: `{
		"refid": 1,
		"kind": "uint",
		"value": 1
}`, out: &Value{
		Refid: 1,
		Kind:  "uint",
		Value: uint(1),
	}})
	res = append(res, unmarshalJSONTest{in: `{
		"refid": 1,
		"kind": "uint8",
		"value": 1
}`, out: &Value{
		Refid: 1,
		Kind:  "uint8",
		Value: uint8(1),
	}})
	res = append(res, unmarshalJSONTest{in: `{
		"refid": 1,
		"kind": "uint16",
		"value": 1
}`, out: &Value{
		Refid: 1,
		Kind:  "uint16",
		Value: uint16(1),
	}})
	res = append(res, unmarshalJSONTest{in: `{
		"refid": 1,
		"kind": "uint32",
		"value": 1
}`, out: &Value{
		Refid: 1,
		Kind:  "uint32",
		Value: uint32(1),
	}})
	res = append(res, unmarshalJSONTest{in: `{
		"refid": 1,
		"kind": "uint64",
		"value": 1
}`, out: &Value{
		Refid: 1,
		Kind:  "uint64",
		Value: uint64(1),
	}})
	res = append(res, unmarshalJSONTest{in: `{
		"refid": 1,
		"kind": "float32",
		"value": 1
}`, out: &Value{
		Refid: 1,
		Kind:  "float32",
		Value: float32(1),
	}})
	res = append(res, unmarshalJSONTest{in: `{
		"refid": 1,
		"kind": "float64",
		"value": 1
}`, out: &Value{
		Refid: 1,
		Kind:  "float64",
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
}`, out: &Value{
		Refid: 1,
		Kind:  "ptr",
		Value: &Value{
			Refid: 2,
			Kind:  "string",
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
}`, out: &Value{
		Refid: 1,
		Kind:  "ptr",
		Value: &Value{
			Refid: 2,
			Kind:  "struct",
			Value: map[string]interface{}{
				"name": &Value{
					Refid: 3,
					Kind:  "string",
					Value: "Arthur",
				},
				"children": &Value{
					Refid: 4,
					Kind:  "slice",
					Value: []interface{}{},
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
}`, out: &Value{
		Refid: 1,
		Kind:  "ptr",
		Value: &Value{
			Refid: 2,
			Kind:  "struct",
			Value: map[string]interface{}{
				"name": &Value{
					Refid: 3,
					Kind:  "string",
					Value: "Arthur",
				},
				"parent": &Value{
					Refid: 4,
					Kind:  "ptr",
					Value: nil,
				},
				"children": &Value{
					Refid: 5,
					Kind:  "slice",
					Value: []interface{}{
						&Value{
							Refid: 6,
							Kind:  "ptr",
							Value: &Value{
								Refid: 7,
								Kind:  "struct",
								Value: map[string]interface{}{
									"name": &Value{
										Refid: 8,
										Kind:  "string",
										Value: "Trillian",
									},
									"parent": &Value{
										Refid: 9,
										Kind:  "ref",
										Value: 1,
									},
									"children": &Value{
										Refid: 10,
										Kind:  "slice",
										Value: []interface{}{},
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
		out: &Value{
			Refid: 1,
			Kind:  "chan",
			Value: "aaa",
		},
		err: &InvalidValueKindError{Kind: "chan"},
	})
	res = append(res, unmarshalJSONTest{
		in: `{
			"refid": 1,
			"kind": "dummy",
			"value": "aaa"
		}`,
		out: &Value{
			Refid: 1,
			Kind:  "chan",
			Value: "aaa",
		},
		err: &InvalidValueKindError{Kind: "dummy"},
	})
	res = append(res, unmarshalJSONTest{
		in: `{
			"refid": 1,
			"kind": "complex64",
			"value": "aaa"
		}`,
		out: &Value{
			Refid: 1,
			Kind:  "complex64",
			Value: "aaa",
		},
		err: &InvalidValueKindError{Kind: "complex64"},
	})
	res = append(res, unmarshalJSONTest{
		in: `{
			"refid": 1,
			"kind": "complex128",
			"value": "aaa"
		}`,
		out: &Value{
			Refid: 1,
			Kind:  "complex128",
			Value: "aaa",
		},
		err: &InvalidValueKindError{Kind: "complex128"},
	})
	res = append(res, unmarshalJSONTest{
		in: `{
			"refid": 1,
			"kind": "uintptr",
			"value": "aaa"
		}`,
		err: &InvalidValueKindError{Kind: "uintptr"},
	})
	res = append(res, unmarshalJSONTest{
		in: `{
			"refid": 1,
			"kind": "byte",
			"value": "aaa"
		}`,
		err: &InvalidValueKindError{Kind: "byte"},
	})
	res = append(res, unmarshalJSONTest{
		in: `{
			"refid": 1,
			"kind": "rune",
			"value": "aaa"
		}`,
		out: &Value{
			Refid: 1,
			Kind:  "rune",
			Value: "aaa",
		},
		err: &InvalidValueKindError{Kind: "rune"},
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
		err: &InvalidValueError{Kind: "ptr", Value: "invalid"},
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
		err: &InvalidValueError{Kind: "ptr", Value: "invalid"},
	})
	res = append(res, unmarshalJSONTest{
		in: `{
			"refid": 1,
			"kind": "struct", 
			"value": "invalid"
		}`,
		err: &InvalidValueError{Kind: "struct", Value: "invalid"},
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
		err: &InvalidValueError{Kind: "struct", Value: "invalid"},
	})
	res = append(res, unmarshalJSONTest{
		in: `{
			"refid": 1,
			"kind": "map", 
			"value": "invalid"
		}`,
		err: &InvalidValueError{Kind: "map", Value: "invalid"},
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
		err: &InvalidValueError{Kind: "map", Value: "invalid"},
	})
	res = append(res, unmarshalJSONTest{
		in: `{
			"refid": 1,
			"kind": "slice", 
			"value": "invalid"
		}`,
		err: &InvalidValueError{Kind: "slice", Value: "invalid"},
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
		err: &InvalidValueError{Kind: "map", Value: "invalid"},
	})

	return res[0:len(res):len(res)]
}

func TestValue_UnmarshalJSON(t *testing.T) {
	for i, arg := range unmarshalJSONTests() {
		v := &Value{}
		err := v.UnmarshalJSON([]byte(arg.in))
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
	err := &InvalidValueKindError{Kind:"invalid"}
	expected := "tahwil.Value: invalid value kind \"" + err.Kind + "\""
	if err.Error() != expected {
		t.Errorf("mismatch\nhave: %#+v\nwant: %#+v", err.Error(), expected)
	}
}
package tahwil

import (
	"reflect"
	"testing"
)

type unmarshalJSONTest struct {
	in  string
	out *Value
}

func unmarshalJSONTests() []unmarshalJSONTest {
	res := make([]unmarshalJSONTest, 0)

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

	return res[0:len(res):len(res)]
}

func TestValue_UnmarshalJSON(t *testing.T) {
	for i, arg := range unmarshalJSONTests() {
		v := &Value{}
		err := v.UnmarshalJSON([]byte(arg.in))
		if err != nil {
			t.Fatalf("UnmarshalJSON: %v", err)
		}
		if !reflect.DeepEqual(arg.out, v) {
			t.Errorf("#%d: mismatch\nhave: %#+v\nwant: %#+v", i, v, arg.out)
			continue
		}
	}
}

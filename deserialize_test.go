package tahwil

import (
	"reflect"
	"testing"
)

type personT struct {
	Name     string     `json:"name"`
	Parent   *personT   `json:"parent"`
	Children []*personT `json:"children"`
}

type fromValueTest struct {
	in  *Value
	out *personT
	err error
}

func fromValueTests() []fromValueTest {
	result := make([]fromValueTest, 0)

	result = append(result, fromValueTest{
		in: &Value{
			Refid: 1,
			Kind:  "ptr",
			Value: &Value{
				Refid: 2,
				Kind:  "struct",
				Value: map[string]*Value{
					"name": {
						Refid: 3,
						Kind:  "string",
						Value: "Martin",
					},
					"parent": {
						Refid: 4,
						Kind:  "ptr",
						Value: nil,
					},
					"children": {
						Refid: 5,
						Kind:  "slice",
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
			Kind:  "ptr",
			Value: &Value{
				Refid: 2,
				Kind:  "struct",
				Value: map[string]*Value{
					"name": {
						Refid: 3,
						Kind:  "string",
						Value: "Martin",
					},
					"parent": {
						Refid: 4,
						Kind:  "ptr",
						Value: &Value{
							Refid: 5,
							Kind:  "struct",
							Value: map[string]*Value{
								"name": {
									Refid: 6,
									Kind:  "string",
									Value: "Kevin",
								},
								"children": {
									Refid: 7,
									Kind:  "slice",
									Value: []*Value{
										{
											Refid: 8,
											Kind:  "ref",
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

	return result
}

func TestFromValue(t *testing.T) {
	for i, arg := range fromValueTests() {
		p := &personT{}
		err := FromValue(arg.in, p)
		if err != nil {
			t.Fatalf("#%d: %#+v", i, err)
		}
		if !reflect.DeepEqual(p, arg.out) {
			t.Errorf("#%d: mismatch\nhave: %#+v\nwant: %#+v", i, p, arg.out)
		}
	}
}

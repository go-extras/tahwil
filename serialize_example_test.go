package tahwil_test

import (
	"encoding/json"
	"fmt"

	"github.com/go-extras/tahwil"
)

type SerializedPerson struct {
	Name     string
	Parent   *SerializedPerson
	Children []*SerializedPerson
}

func ExampleToValue() {
	parent := &SerializedPerson{
		Name: "Arthur",
		Children: []*SerializedPerson{
			{
				Name: "Ford",
			},
			{
				Name: "Trillian",
			},
		},
	}
	parent.Children[0].Parent = parent
	parent.Children[1].Parent = parent
	v, err := tahwil.ToValue(parent)
	if err != nil {
		panic(err)
	}
	res, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(res))
	//nolint:lll // valid JSON
	// Output: {"refid":1,"kind":"ptr","value":{"refid":0,"kind":"struct","value":{"Children":{"refid":0,"kind":"slice","value":[{"refid":3,"kind":"ptr","value":{"refid":0,"kind":"struct","value":{"Children":{"refid":0,"kind":"slice","value":[]},"Name":{"refid":0,"kind":"string","value":"Ford"},"Parent":{"refid":4,"kind":"ref","value":1}}}},{"refid":5,"kind":"ptr","value":{"refid":0,"kind":"struct","value":{"Children":{"refid":0,"kind":"slice","value":[]},"Name":{"refid":0,"kind":"string","value":"Trillian"},"Parent":{"refid":6,"kind":"ref","value":1}}}}]},"Name":{"refid":0,"kind":"string","value":"Arthur"},"Parent":{"refid":2,"kind":"ptr","value":null}}}}
}

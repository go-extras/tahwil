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
	// Output: {"refid":1,"kind":"ptr","value":{"refid":2,"kind":"struct","value":{"Children":{"refid":5,"kind":"slice","value":[{"refid":6,"kind":"ptr","value":{"refid":7,"kind":"struct","value":{"Children":{"refid":10,"kind":"slice","value":[]},"Name":{"refid":8,"kind":"string","value":"Ford"},"Parent":{"refid":9,"kind":"ref","value":1}}}},{"refid":11,"kind":"ptr","value":{"refid":12,"kind":"struct","value":{"Children":{"refid":15,"kind":"slice","value":[]},"Name":{"refid":13,"kind":"string","value":"Trillian"},"Parent":{"refid":14,"kind":"ref","value":1}}}}]},"Name":{"refid":3,"kind":"string","value":"Arthur"},"Parent":{"refid":4,"kind":"ptr","value":null}}}}
}

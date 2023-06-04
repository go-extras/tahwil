package tahwil_test

import (
	"encoding/json"
	"fmt"

	"github.com/go-extras/tahwil"
)

type Person struct {
	Name     string    `json:"name"`
	Parent   *Person   `json:"parent"`
	Children []*Person `json:"children"`
}

func prepareData() []byte {
	parent := &Person{
		Name: "Arthur",
		Children: []*Person{
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
	return res
}

func ExampleFromValue() {
	data := &tahwil.Value{}
	res := prepareData()
	err := json.Unmarshal(res, data)
	if err != nil {
		panic(err)
	}
	person := &Person{}
	err = tahwil.FromValue(data, person)
	if err != nil {
		panic(err)
	}
	fmt.Printf(`Name: %s
Children:
    - %s
    -- parent name: %s
    - %s
    -- parent name: %s
`, person.Name,
		person.Children[0].Name,
		person.Children[0].Parent.Name,
		person.Children[1].Name,
		person.Children[1].Parent.Name)
	// Output: Name: Arthur
	//Children:
	//     - Ford
	//     -- parent name: Arthur
	//     - Trillian
	//     -- parent name: Arthur
}

package main

import (
	"encoding/json"
	"fmt"

	"github.com/go-extras/tahwil"
)

type Person struct {
	Name     string
	Parent   *Person
	Children []*Person
}

func main() {
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
	fmt.Println(string(res))
}

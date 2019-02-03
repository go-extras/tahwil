# tahwil
Cyclic graph structures serialization library written in go.

## How can it be useful?

Sometimes you need to serialize a structure that has circular references.
This library lets you transform your cyclic graph to a tree and then serialize to json.

## How to use it?

### Encoding

```go
package main

import (
	"errors"
	"fmt"
	"encoding/json"

	"github.com/go-extras/tahwil"
)

type Person struct {
	Name string
	Parent *Person
	Children []*Person
}

func main() {
	parent := &Person{
		Name: "Arthur",
		Children: []*Person{
			&Person{
				Name: "Ford",
			},
			&Person{
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
	fmt.Println(res) 
}

```

The output will be one-line equivalent of the following JSON:
```json
    {
      "refid": 1,
      "kind": "ptr",
      "value": {
        "refid": 2,
        "kind": "struct",
        "value": {
          "Children": {
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
                    "Children": {
                      "refid": 10,
                      "kind": "slice",
                      "value": []
                    },
                    "Name": {
                      "refid": 8,
                      "kind": "string",
                      "value": "Ford"
                    },
                    "Parent": {
                      "refid": 9,
                      "kind": "ref",
                      "value": 1
                    }
                  }
                }
              },
              {
                "refid": 11,
                "kind": "ptr",
                "value": {
                  "refid": 12,
                  "kind": "struct",
                  "value": {
                    "Children": {
                      "refid": 15,
                      "kind": "slice",
                      "value": []
                    },
                    "Name": {
                      "refid": 13,
                      "kind": "string",
                      "value": "Trillian"
                    },
                    "Parent": {
                      "refid": 14,
                      "kind": "ref",
                      "value": 1
                    }
                  }
                }
              }
            ]
          },
          "Name": {
            "refid": 3,
            "kind": "string",
            "value": "Arthur"
          },
          "Parent": {
            "refid": 4,
            "kind": "ptr",
            "value": null
          }
        }
      }
    }
```

### Decoding

```go
package main

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
			&Person{
				Name: "Ford",
			},
			&Person{
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

func main() {
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
}
```

This should output:

```
Name: Arthur
Children:
    - Ford
	-- parent name: Arthur
    - Trillian
	-- parent name: Arthur
```

As you can see, Arthur is displayed here 3 times - first as a main person, and then as a parent of the both children.

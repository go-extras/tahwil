# tahwil

A Go library for serializing cyclic graph structures to JSON.

[![Build](https://github.com/go-extras/tahwil/actions/workflows/test.yml/badge.svg)](https://github.com/go-extras/tahwil/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-extras/tahwil)](https://goreportcard.com/report/github.com/go-extras/tahwil)
[![Go Reference](https://pkg.go.dev/badge/github.com/go-extras/tahwil.svg)](https://pkg.go.dev/github.com/go-extras/tahwil)
[![License](https://img.shields.io/github/license/go-extras/tahwil)](LICENSE)

## Overview

Tahwil transforms cyclic graph structures into serializable trees, enabling JSON encoding of complex data structures with circular references. The library preserves structural integrity by converting cycles into explicit references, making it ideal for serializing interconnected domain models, graph databases, or any data structure where objects reference each other.

## Features

- **Cycle Detection**: Automatically identifies and handles circular references
- **Type Safety**: Preserves Go type information during serialization and deserialization
- **Comprehensive Type Support**: Works with structs, slices, maps, pointers, and all primitive types
- **Bidirectional**: Full encoding and decoding support
- **Zero Dependencies**: Built using only the Go standard library

## Installation

```bash
go get github.com/go-extras/tahwil
```

## Quick Start

```go
import "github.com/go-extras/tahwil"

// Serialize
value, err := tahwil.ToValue(myStruct)
jsonData, err := json.Marshal(value)

// Deserialize
var value tahwil.Value
json.Unmarshal(jsonData, &value)
err := tahwil.FromValue(&value, &myStruct)
```

## Usage


### Encoding

Transform a cyclic structure into a serializable format:

```go
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
			{Name: "Ford"},
			{Name: "Trillian"},
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
```

<details>
<summary>Example output (formatted for readability)</summary>

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

</details>

### Decoding

Deserialize back into your original structure:

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

func main() {
	// Assuming you have serialized JSON data
	var value tahwil.Value
	err := json.Unmarshal(jsonData, &value)
	if err != nil {
		panic(err)
	}
	
	var person Person
	err = tahwil.FromValue(&value, &person)
	if err != nil {
		panic(err)
	}
	
	// Circular references are preserved
	fmt.Printf("Name: %s\n", person.Name)
	fmt.Printf("First child: %s\n", person.Children[0].Name)
	fmt.Printf("First child's parent: %s\n", person.Children[0].Parent.Name)
}
```

**Output:**
```
Name: Arthur
First child: Ford
First child's parent: Arthur
```

The circular reference is preserved—Arthur appears as both the root person and as the parent of his children.

## Supported Types

The library handles the following Go types:

**Primitives:**
- `string`, `bool`
- `int`, `int8`, `int16`, `int32`, `int64`
- `uint`, `uint8`, `uint16`, `uint32`, `uint64`
- `float32`, `float64`

**Complex types:**
- `ptr` (pointers)
- `struct`
- `map`
- `slice`

**Special:**
- `ref` (internal type for representing circular references)

All complex types must contain only supported types.

## Use Cases

- **Domain Models**: Serialize interconnected business objects with bidirectional relationships
- **Graph Databases**: Export and import graph structures while preserving topology
- **Caching**: Store complex object graphs in JSON-based caches
- **API Responses**: Send graph structures over HTTP without losing referential integrity
- **Configuration**: Persist and restore complex configurations with shared references

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

package tahwil

// Kind represents the type kind stored in a Value.
type Kind string

const (
	Ref Kind = "ref"

	Bool Kind = "bool"

	Int   Kind = "int"
	Int8  Kind = "int8"
	Int16 Kind = "int16"
	Int32 Kind = "int32"
	Int64 Kind = "int64"

	Uint   Kind = "uint"
	Uint8  Kind = "uint8"
	Uint16 Kind = "uint16"
	Uint32 Kind = "uint32"
	Uint64 Kind = "uint64"

	Float32 Kind = "float32"
	Float64 Kind = "float64"

	String Kind = "string"
	Struct Kind = "struct"
	Slice  Kind = "slice"
	Array  Kind = "array"
	Map    Kind = "map"
	Ptr    Kind = "ptr"
)

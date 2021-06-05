package enki

import (
	"fmt"
	"io"
)

const (
	Int   = "int"
	Int8  = "int8"
	Int16 = "int16"
	Int32 = "int32"
	Int64 = "int64"

	Uint    = "uint"
	Uint8   = "uint8"
	Uint16  = "uint16"
	Uint32  = "uint32"
	Uint64  = "uint64"
	Uintptr = "uintptr"

	Float32 = "float32"
	Float64 = "float64"

	Complex64  = "complex64"
	Complex128 = "complex128"

	Byte    = "byte"
	Rune    = "rune"
	Boolean = "bool"
	Str     = "string"
)

// Materializer builder that can be materialized
type Materializer interface {
	Materialize()
}

// Primitive primitive builder
type Primitive interface {
	NewLine()
	Err() error
}

// Statement builder, that can build any statement
type Statement interface {
	fmt.Stringer
	Primitive
	Line(s string, args ...interface{})
}

// FunctionDef builder for function definition
type FunctionDef interface {
	fmt.Stringer
	Materializer
	Name(name string) FunctionDef
	Params(params ...string) FunctionDef
	Returns(results ...string) FunctionDef
}

// Function builder for free functions
type Function interface {
	FunctionDef
	Body(def func(sb Statement)) Function
}

// Method buildr for method
type Method interface {
	Function
	Receiver(def string) Method
}

type MethodDef func(FunctionDef)
type Methods []MethodDef

// Type builder for types and interfaces
type Type interface {
	fmt.Stringer
	Materializer
	Name(name string) Type
	Is(anotherType string) Type
	Struct(fields func(Statement)) Type
	Interface(methods Methods) Type
}

// File builder for files
type File interface {
	Materializer
	GeneratedBy(tool string) File
	Package(pkg string) File
	AddImport(alias, path string) File
	Type(name string) Type
	Function(name string) FunctionDef
	Write(dest io.Writer) error
	Create(fileName string) error
	Append(fileName string) error
}

// Ptr pointer creation function
func Ptr(name string) string {
	return "*" + name
}

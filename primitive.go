package enki

import (
	"strings"
)

// Primitive primitive builder
type Primitive interface {
	// NewLine add new empty line.
	NewLine()
}

type primitiveBuilder struct {
	inner *strings.Builder
	err   error
}

// NewPrimitive creates new primitive builder
func NewPrimitive() *primitiveBuilder {
	return &primitiveBuilder{inner: &strings.Builder{}}
}

// NewLine add new empty line.
func (pb *primitiveBuilder) NewLine() {
	_, err := pb.inner.WriteString("\n")
	if err != nil {
		pb.err = err
	}
}

func (pb *primitiveBuilder) String() string {
	return pb.inner.String()
}

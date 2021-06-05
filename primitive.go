package enki

import (
	"strings"
)

type primitiveBuilder struct {
	inner *strings.Builder
	err   error
}

func NewPrimitive() *primitiveBuilder {
	return &primitiveBuilder{inner: &strings.Builder{}}
}

func (pb *primitiveBuilder) NewLine() {
	_, err := pb.inner.WriteString("\n")
	if err != nil {
		pb.err = err
	}
}

func (pb *primitiveBuilder) Err() error {
	return pb.err
}

func (pb *primitiveBuilder) String() string {
	return pb.inner.String()
}

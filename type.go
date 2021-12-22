package enki

// Type builder for types and interfaces
type Type interface {
	Block
	// Name type name.
	Name(name string) Type
	// Is type redeclaration.
	Is(anotherType string) Type
	// Struct type as struct.
	Struct(fields ...Statement) Type
	// Interface type as interface.
	Interface(methods ...Function) Type
}

type typeBuilder struct {
	*statement
	name    string
	is      string
	fields  []Statement
	methods []Function
}

// T creates new type builder
func T(name string) *typeBuilder {
	return &typeBuilder{
		name:      name,
		statement: Stmt(),
	}
}

func (tb *typeBuilder) Name(name string) Type {
	tb.name = name
	return tb
}

func (tb *typeBuilder) Is(anotherType string) Type {
	tb.is = anotherType
	return tb
}

func (tb *typeBuilder) Struct(fields ...Statement) Type {
	tb.fields = fields
	return tb
}

func (tb *typeBuilder) Interface(methods ...Function) Type {
	tb.methods = methods
	return tb
}

func (tb *typeBuilder) materialize() string {
	if tb.err != nil {
		return ""
	}
	switch {
	case tb.fields != nil:
		tb.Line("type @1 struct {", tb.name)
		for _, field := range tb.fields {
			tb.Print(field.materialize())
		}
		tb.Line("}")
	case len(tb.methods) > 0:
		tb.Line("type @1 interface {", tb.name)
		for _, method := range tb.methods {
			tb.Print(method.materialize())
		}
		tb.Line("}")
	default:
		tb.Line("type @1 @2", tb.name, tb.is)
	}

	return tb.String()
}

func (tb *typeBuilder) String() string {
	return tb.statement.String()
}

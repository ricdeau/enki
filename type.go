package enki

type typeBuilder struct {
	*functionBuilder
	name    string
	is      string
	fields  func(Statement)
	methods Methods
}

// NewType creates new type builder
func NewType() *typeBuilder {
	return &typeBuilder{functionBuilder: NewFunction()}
}

func (tb *typeBuilder) Name(name string) Type {
	tb.name = name
	return tb
}

func (tb *typeBuilder) Is(anotherType string) Type {
	tb.is = anotherType
	return tb
}

func (tb *typeBuilder) Struct(fields func(Statement)) Type {
	tb.fields = fields
	return tb
}

func (tb *typeBuilder) Interface(methods Methods) Type {
	tb.methods = methods
	return tb
}

func (tb *typeBuilder) Materialize() {
	if tb.Err() != nil {
		return
	}
	switch {
	case tb.fields != nil:
		tb.Line("type @1 struct {", tb.name)
		tb.fields(tb)
		tb.Line("}")
	case len(tb.methods) > 0:
		tb.Line("type @1 interface {", tb.name)
		for _, method := range tb.methods {
			method(tb.functionBuilder)
			tb.Line("@1@2@3", tb.functionBuilder.name, tb.materializeParams(), tb.materializeReturns())
			tb.functionBuilder.reset()
		}
		tb.Line("}")
	default:
		tb.Line("type @1 @2", tb.name, tb.is)
	}
	tb.reset()
}

func (tb *typeBuilder) String() string {
	return tb.functionBuilder.String()
}

func (tb *typeBuilder) reset() {
	tb.name = ""
	tb.is = ""
	tb.fields = nil
	tb.methods = nil
	tb.functionBuilder.reset()
}

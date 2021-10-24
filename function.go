package enki

import (
	"strings"
)

// Function builder for free functions
type Function interface {
	Block
	Name(name string) Function
	Params(params ...string) Function
	Returns(results ...string) Function
	Body(body ...Statement) Function
}

// Method builder for method
type Method interface {
	Function
	Receiver(def string) Method
}

type functionBuilder struct {
	*statement
	isDef    bool
	name     string
	receiver string
	params   []string
	returns  []string
	body     []Statement
}

// F creates new function.
func F(name string) Function {
	return &functionBuilder{
		statement: Stmt(),
		name:      name,
	}
}

// Def creates interface method definition.
func Def(name string) Function {
	return &functionBuilder{
		statement: Stmt(),
		name:      name,
		isDef:     true,
	}
}

// M creates new method.
func M(name string) Method {
	return F(name).(Method)
}

func (fb *functionBuilder) Name(name string) Function {
	fb.name = name
	return fb
}

func (fb *functionBuilder) Params(params ...string) Function {
	fb.params = params
	return fb
}

func (fb *functionBuilder) Returns(results ...string) Function {
	fb.returns = results
	return fb
}

func (fb *functionBuilder) Body(body ...Statement) Function {
	fb.body = body
	return fb
}

func (fb *functionBuilder) Receiver(def string) Method {
	fb.receiver = def
	return fb
}

func (fb *functionBuilder) materialize() string {
	if fb.err != nil {
		return ""
	}
	var openBracket, closeBracket string
	if fb.body != nil {
		openBracket, closeBracket = " {", "}"
	}
	receiver := fb.materializeReceiver()
	params := fb.materializeParams()
	returns := fb.materializeReturns()
	// func<receiver><name><params><returns><openBracket?>
	if fb.isDef {
		fb.Line("@1@2@3", fb.name, params, returns)
	} else {
		fb.Line("func@1 @2@3@4@5", receiver, fb.name, params, returns, openBracket)
	}
	if fb.body != nil {
		for _, stmt := range fb.body {
			fb.Print(stmt.materialize())
		}
		fb.Line(closeBracket)
	}

	return fb.statement.materialize()
}

func (fb *functionBuilder) String() string {
	return fb.statement.String()
}

func (fb *functionBuilder) materializeReturns() (result string) {
	switch {
	case len(fb.returns) == 0:
	case len(fb.returns) > 1 || strings.Contains(fb.returns[0], " "):
		result = " (" + strings.Join(fb.returns, ", ") + ")"
	default:
		result = " " + fb.returns[0]
	}
	return
}

func (fb *functionBuilder) materializeParams() string {
	if len(fb.params) == 0 {
		return "()"
	}
	return "(" + strings.Join(fb.params, ", ") + ")"
}

func (fb *functionBuilder) materializeReceiver() string {
	if fb.receiver == "" {
		return ""
	}
	return " (" + fb.receiver + ")"
}

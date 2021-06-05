package enki

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewFunction(t *testing.T) {
	s := NewFunction()
	require.NotNil(t, s.inner)
	require.NoError(t, s.err)
}

type functionBuilder struct {
	*statementBuilder
	name     string
	receiver string
	params   []string
	returns  []string
	body     func(def Statement)
}

// NewFunction creates new function builder
func NewFunction() *functionBuilder {
	return &functionBuilder{statementBuilder: NewStatement()}
}

func (fb *functionBuilder) Name(name string) FunctionDef {
	fb.name = name
	return fb
}

func (fb *functionBuilder) Params(params ...string) FunctionDef {
	fb.params = params
	return fb
}

func (fb *functionBuilder) Returns(results ...string) FunctionDef {
	fb.returns = results
	return fb
}

func (fb *functionBuilder) Body(def func(sb Statement)) Function {
	fb.body = def
	return fb
}

func (fb *functionBuilder) Receiver(def string) Method {
	fb.receiver = def
	return fb
}

func (fb *functionBuilder) Materialize() {
	if fb.Err() != nil {
		return
	}
	var openBracket, closeBracket string
	if fb.body != nil {
		openBracket, closeBracket = " {", "}"
	}
	receiver := fb.materializeReceiver()
	params := fb.materializeParams()
	returns := fb.materializeReturns()
	// func<receiver><name><params><returns><openBracket?>
	fb.Line("func@1 @2@3@4@5", receiver, fb.name, params, returns, openBracket)
	if fb.body != nil {
		fb.body(fb)
		fb.Line(closeBracket)
	}
}

func (fb *functionBuilder) String() string {
	return fb.statementBuilder.String()
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

func (fb *functionBuilder) reset() {
	fb.name = ""
	fb.receiver = ""
	fb.params = nil
	fb.returns = nil
	fb.body = nil
}

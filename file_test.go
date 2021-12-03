package enki

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	Int     = "int"
	Str     = "string"
	Float32 = "float32"
	Float64 = "float64"
)

func TestNewFile(t *testing.T) {
	s := NewFile()
	require.NotNil(t, s)
	require.NotNil(t, s.statement)
	require.NotNil(t, s.inner)
}

func Test_fileBuilder_WriteNew(t *testing.T) {
	f := NewFile()
	f.Package("enki")
	f.GeneratedBy("enki")
	f.Import(".", "fmt")
	f.Line("// Number type redefined")

	f.Add(T("Number").Is(Int))
	f.NewLine()

	f.Add(T("NumStruct").Struct(
		Field("a, b " + Int),
	))
	f.NewLine()

	f.Add(T("Num").Interface(
		Def("AsNumber").Returns(Str),
		Def("Sum").Params(Int).Returns(Int),
	))
	f.NewLine()

	f.Add(M("AsNumber").Receiver("n Number").Body(
		Stmt().Line("return Sprint(n)"),
	).Returns(Str))
	f.NewLine()

	f.Add(M("Sum").Receiver("n Number").Params("x int").Body(
		Stmt().Line("return int(n) + x"),
	).Returns(Int))
	f.NewLine()

	f.Add(M("Sum").Receiver("ns NumStruct").Body(
		Stmt().Line("return ns.a + ns.b"),
	).Returns(Int))
	f.NewLine()

	f.Add(F("sum").Params("a, b " + Float32).Body(
		Stmt().Line("return @1(a + b)", Float64),
	).Returns("s " + Float64))

	buf := &bytes.Buffer{}

	err := f.Write(buf)

	require.NoError(t, err)

	t.Log(buf.String())
}

package enki

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_fileBuilder_WriteNew(t *testing.T) {
	f := NewFile()
	f.Package("enki")
	f.GeneratedBy("enki")
	f.AddImport("", "fmt")
	f.Line("// Number type redefined")
	f.Type("Number").Is(Int).Materialize()
	f.NewLine()
	f.Type("NumStruct").Struct(func(s Statement) {
		s.Line("a, b " + Int)
	}).Materialize()
	f.NewLine()
	f.Type("Num").Interface(Methods{
		func(f FunctionDef) {
			f.Name("AsNumber").Returns(Str)
		},
		func(f FunctionDef) {
			f.Name("Sum").Params(Int).Returns(Int)
		},
	}).Materialize()
	f.NewLine()
	f.Function("AsNumber").Returns(Str).(Method).Receiver("n Number").Body(func(sb Statement) {
		sb.Line("return fmt.Sprint(n)")
	}).Materialize()
	f.NewLine()
	f.Function("Sum").Params("x int").Returns(Int).(Method).Receiver("n Number").Body(func(sb Statement) {
		sb.Line("return int(n) + x")
	}).Materialize()
	f.NewLine()
	f.Function("Sum").Returns(Int).(Method).Receiver("ns NumStruct").Body(func(sb Statement) {
		sb.Line("return ns.a + ns.b")
	}).Materialize()
	f.NewLine()
	f.Function("sum").Params("a, b " + Float32).Returns("s " + Float64).(Function).Body(func(sb Statement) {
		sb.Line("return @1(a + b)", Float64)
	}).Materialize()

	err := f.Create("file.gen.go")
	require.NoError(t, err)
}

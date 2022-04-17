package enki

import (
	"bytes"
	"io"
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

func TestFileCreate(t *testing.T) {
	f := NewFile()
	f.Package("enki")
	f.GeneratedBy("enki")
	f.Import(".", "fmt")
	f.Import("", "sync")

	f.Consts(Field(`DEF = "default"`))

	f.Vars(Field(`wg sync.WaitGroup`))

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
		Stmt().Line("wg.Add(1)"),
		Stmt().Line("defer wg.Done()"),
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

	err := f.GoFmt(true).GoImports(true).Create("file.gen.go")

	require.NoError(t, err)
}

type errWriter struct{}

func (e errWriter) Write(_ []byte) (n int, err error) {
	return 0, io.EOF
}

func TestFile_Write(t *testing.T) {
	tests := []struct {
		name    string
		file    func() File
		dest    io.Writer
		wantErr bool
	}{
		{
			name: "success",
			file: func() File {
				f := NewFile()
				f.Package("pkg")
				return f
			},
			dest:    bytes.NewBufferString(""),
			wantErr: false,
		},
		{
			name: "error format",
			file: func() File {
				f := NewFile()
				return f
			},
			dest:    bytes.NewBufferString(""),
			wantErr: true,
		},
		{
			name: "error write",
			file: func() File {
				f := NewFile()
				f.Package("pkg")
				return f
			},
			dest:    errWriter{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.file().Write(tt.dest); (err != nil) != tt.wantErr {
				t.Errorf("wantErr = %v but go %v", tt.wantErr, err)
			}
		})
	}
}

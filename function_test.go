package enki

import (
	"go/format"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFunction(t *testing.T) {
	type args struct {
		name     string
		receiver string
		params   []string
		returns  []string
		body     []Statement
	}
	tests := []struct {
		name string
		def  args
		want string
	}{
		{
			name: "method",
			def: args{
				name:     "SomeFunc",
				receiver: "s SomeStruct",
				params:   []string{"a int", "b string"},
				returns:  []string{"res interface{}", "err error"},
				body: []Statement{
					Stmt().Line("return fmt.Sprint(a) + b, nil"),
				},
			},
			want: "func (s SomeStruct) SomeFunc(a int, b string) (res interface{}, err error) {\nreturn fmt.Sprint(a) + b, nil\n}\n",
		},
		{
			name: "function",
			def: args{
				name:    "SomeFunc",
				params:  []string{"a int", "b string"},
				returns: []string{"error"},
				body: []Statement{
					Stmt().Line("fmt.Print(b)"),
					Stmt().Line("return nil"),
				},
			},
			want: "func SomeFunc(a int, b string) error {\nfmt.Print(b)\nreturn nil\n}\n",
		},
		{
			name: "function def",
			def: args{
				name:    "SomeFunc",
				params:  []string{"a int", "b string"},
				returns: []string{"error"},
			},
			want: "func SomeFunc(a int, b string) error\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := M(tt.def.name).Receiver(tt.def.receiver).Params(tt.def.params...).Body(tt.def.body...).Returns(tt.def.returns...)
			got := builder.materialize()
			require.Equal(t, tt.want, got)
			source, err := format.Source([]byte(got))
			require.NoError(t, err)
			t.Log("\n" + string(source))
		})
	}
}

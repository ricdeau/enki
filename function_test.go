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
		body     func(statement Statement)
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
				body: func(s Statement) {
					s.Line("return fmt.Sprint(a) + b, nil")
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
				body: func(s Statement) {
					s.Line("fmt.Print(b)")
					s.Line("return nil")
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
			builder := NewFunction().Name(tt.def.name).Params(tt.def.params...).Returns(tt.def.returns...).(Method).
				Receiver(tt.def.receiver).
				Body(tt.def.body)
			builder.Materialize()
			got := builder.String()
			require.Equal(t, tt.want, got)
			source, err := format.Source([]byte(got))
			require.NoError(t, err)
			t.Log("\n" + string(source))
		})
	}
}

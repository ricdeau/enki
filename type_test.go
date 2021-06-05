package enki

import (
	"go/format"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewType(t *testing.T) {
	s := NewType()
	require.NotNil(t, s.inner)
	require.NoError(t, s.err)
}

func Test_typeBuilder_String(t *testing.T) {
	tests := []struct {
		name string
		got  Type
		want string
	}{
		{
			name: "another type",
			got:  NewType().Name("SomeStruct").Is("string"),
			want: "type SomeStruct string\n",
		},
		{
			name: "struct",
			got: NewType().Name("SomeStruct").Struct(func(s Statement) {
				s.Line("id string")
				s.Line("Time time.Time")
			}),
			want: "type SomeStruct struct {\nid string\nTime time.Time\n}\n",
		},
		{
			name: "interface",
			got: NewType().Name("SomeStruct").Interface(Methods{
				func(f FunctionDef) {
					f.Name("GetId").Returns("string")
				},
				func(f FunctionDef) {
					f.Name("SetId").Params("id string")
				},
			}),
			want: "type SomeStruct interface {\nGetId() string\nSetId(id string)\n}\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.got.Materialize()
			got := tt.got.String()
			require.Equal(t, tt.want, got)
			source, err := format.Source([]byte(got))
			require.NoError(t, err)
			t.Log("\n" + string(source))
		})
	}
}

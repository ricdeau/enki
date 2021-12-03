package enki

import (
	"go/format"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewType(t *testing.T) {
	s := T("")
	require.NotNil(t, s.inner)
	require.NoError(t, s.err)
}

func Test_typeBuilder_String(t *testing.T) {
	withErr := &typeBuilder{statement: &statement{primitiveBuilder: &primitiveBuilder{err: io.EOF}}}
	tests := []struct {
		name    string
		typeDef Type
		want    string
	}{
		{
			name:    "another type",
			typeDef: T("SomeStruct").Is("string"),
			want:    "type SomeStruct string\n",
		},
		{
			name:    "new name",
			typeDef: T("").Name("SomeStruct").Is("string"),
			want:    "type SomeStruct string\n",
		},
		{
			name: "struct",
			typeDef: T("SomeStruct").Struct(
				Field("id string"),
				Field("Time time.Time"),
			),
			want: "type SomeStruct struct {\nid string\nTime time.Time\n}\n",
		},
		{
			name: "interface",
			typeDef: T("SomeStruct").Interface(
				Def("GetId").Returns("string"),
				Def("SetId").Params("id string"),
			),
			want: "type SomeStruct interface {\nGetId() string\nSetId(id string)\n}\n",
		},
		{
			name:    "empty if error",
			typeDef: withErr.Name("A").Is("B"),
			want:    "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.typeDef.materialize()
			require.Equal(t, tt.want, got)
			source, err := format.Source([]byte(got))
			require.NoError(t, err)
			t.Log("\n" + string(source))
		})
	}
}

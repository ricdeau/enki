package enki

import (
	"go/format"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewStatement(t *testing.T) {
	s := Stmt()
	require.NotNil(t, s.inner)
	require.NoError(t, s.err)
}

func Test_statementBuilder_Line(t *testing.T) {
	tests := []struct {
		name    string
		hasErr  bool
		input   string
		want    string
		args    []interface{}
		wantErr bool
	}{
		{
			name:    "success",
			input:   "func @1@2 (@3, @4 string) @5",
			want:    "func SomeFunc1 (a, b string) int64\n",
			args:    []interface{}{"SomeFunc", 1, "a", "b", "int64"},
			wantErr: false,
		},
		{
			name:    "error invalid arg num",
			input:   "func @10",
			args:    []interface{}{"SomeFunc"},
			wantErr: true,
		},
		{
			name:    "error",
			hasErr:  true,
			input:   "func @1",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := Stmt()
			if tt.hasErr {
				sb.err = io.EOF
			}
			sb.Line(tt.input, tt.args...)
			got := sb.String()
			if (sb.err != nil) != tt.wantErr {
				t.Errorf("Line() want wantErr = %v, but typeDef %v", tt.wantErr, sb.err)
			}
			if sb.err != nil {
				t.Log(sb.err)
			} else {
				require.Equal(t, tt.want, got)
				source, err := format.Source([]byte(got))
				require.NoError(t, err)
				t.Log("\n" + string(source))
			}
		})
	}
}

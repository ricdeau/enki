package enki

import (
	"go/format"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_statementBuilder_Line(t *testing.T) {
	tests := []struct {
		name    string
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
			name:    "error",
			input:   "func @10",
			args:    []interface{}{"SomeFunc"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := NewStatement()
			sb.Line(tt.input, tt.args...)
			got := sb.String()
			if (sb.Err() != nil) != tt.wantErr {
				t.Errorf("Line() want wantErr = %v, but got %v", tt.wantErr, sb.Err())
			}
			if sb.Err() != nil {
				t.Log(sb.Err())
			} else {
				require.Equal(t, tt.want, got)
				source, err := format.Source([]byte(got))
				require.NoError(t, err)
				t.Log("\n" + string(source))
			}
		})
	}
}

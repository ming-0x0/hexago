package status

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Parallel()

	type args struct {
		value int64
	}

	tests := []struct {
		name      string
		args      args
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "Valid_Value_1_Replied",
			args:      args{value: 1},
			assertion: assert.NoError,
		},
		{
			name:      "Valid_Value_2_Unreplied",
			args:      args{value: 2},
			assertion: assert.NoError,
		},
		{
			name:      "Invalid_Value_3",
			args:      args{value: 3},
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := New(tc.args.value)
			tc.assertion(t, err)
			if err == nil {
				assert.Equal(t, tc.args.value, got.Value())
			}
		})
	}
}

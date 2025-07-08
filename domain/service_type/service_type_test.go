package service_type

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Parallel()

	type args struct {
		value int
	}

	tests := []struct {
		name      string
		args      args
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "valid value 1 (TuyenDung)",
			args:      args{value: 1},
			assertion: assert.NoError,
		},
		{
			name:      "valid value 2 (LienHe)",
			args:      args{value: 2},
			assertion: assert.NoError,
		},
		{
			name:      "valid value 3 (KhoaHoc)",
			args:      args{value: 3},
			assertion: assert.NoError,
		},
		{
			name:      "invalid value",
			args:      args{value: 4},
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

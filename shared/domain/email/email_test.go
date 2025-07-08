package email

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Parallel()

	type args struct {
		value string
	}

	tests := []struct {
		name      string
		args      args
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "valid email",
			args:      args{value: "test@example.com"},
			assertion: assert.NoError,
		},
		{
			name:      "invalid email",
			args:      args{value: "invalid-email"},
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

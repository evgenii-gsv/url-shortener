package random

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRandomString(t *testing.T) {
	tests := []struct {
		name string
		size int
	}{
		{
			name: "size = 1",
			size: 1,
		},
		{
			name: "size = 5",
			size: 5,
		},
		{
			name: "size = 10",
			size: 10,
		},
		{
			name: "size = 20",
			size: 20,
		},
		{
			name: "size = 30",
			size: 30,
		},
		{
			name: "size = 0",
			size: 0,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			str1 := NewRandomString(tc.size)
			str2 := NewRandomString(tc.size)

			assert.Len(t, str1, tc.size)
			assert.Len(t, str2, tc.size)

			if tc.size > 0 {
				assert.NotEqual(t, str1, str2)
			}
		})
	}
}

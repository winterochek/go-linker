package shorten_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/winterochek/go-linker/internal/shorten"
	"testing"
)

func TestShorten(t *testing.T) {
	t.Run("Returns an alphanumeric short identifier", func(t *testing.T) {
		type testCase struct {
			id       uint32
			expected string
		}

		testCases := []testCase{
			{
				id:       1024,
				expected: "Ey",
			},
			{
				id:       0,
				expected: "",
			},
		}

		for _, tc := range testCases {
			actual := shorten.Shorten(tc.id)
			assert.Equal(t, tc.expected, actual)
		}
	})
	t.Run("Is idempotent", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			assert.Equal(t, "Ey", shorten.Shorten(1024))
		}
	})
}

package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAllInAlphabet(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		Text     string
		Alphabet string
		Expected bool
	}{
		{RandomString(10, KeyboardCharacters), KeyboardCharacters, true},
		{RandomString(13, SpecialCharacters), AlphaNumericCharacters, false},
	}
	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.Text, func(t *testing.T) {
			t.Parallel()
			result := AllInAlphabet(tc.Text, tc.Alphabet)
			require.Equal(t, tc.Expected, result)
		})
	}
}

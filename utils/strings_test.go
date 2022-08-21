package utils

import (
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAllInAlphabet(t *testing.T) {
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
			result := AllInAlphabet(tc.Text, tc.Alphabet)
			require.Equal(t, tc.Expected, result)
		})
	}
}

func TestPrettyStruct(t *testing.T) {
	type Fruit struct {
		Name  string `json:"name"`
		Color string `json:"color"`
	}

	fruit := Fruit{
		Name:  "Strawberry",
		Color: "red",
	}
	res, err := PrettyStruct(fruit)
	if err != nil {
		log.Fatal(err)
	}

	require.Equal(t, `{
    "name": "Strawberry",
    "color": "red"
}`, res)
}

package utils

import "strings"

// AllInAlphabet checks if all characters of s are in alphabet
func AllInAlphabet(s string, alphabet string) bool {
	for _, r := range s {
		if !strings.ContainsRune(alphabet, r) {
			return false
		}
	}
	return true
}

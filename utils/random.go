package utils

import (
	"crypto/rand"
	"fmt"
)

// RandomHeximalString generates a random string consisting of heximal digits
func RandomHeximalString(n int) string {
	if n <= 0 {
		return ""
	}
	k := n >> 1
	if (n & 1) == 1 {
		k++
	}
	b := make([]byte, k)
	_, _ = rand.Read(b)
	s := fmt.Sprintf("%x", b)
	return s[:n]
}

// RandomDecimalString generates a random string consisting of decimal digits
func RandomDecimalString(n int) string {
	if n <= 0 {
		return ""
	}
	b := make([]byte, n)
	_, _ = rand.Read(b)
	for i := range b {
		b[i] = 48 + b[i]%10
	}
	return string(b)
}

// RandomString generates a random string consisting of characters in the provided alphabet
func RandomString(n int, alphabet string) string {
	if n <= 0 {
		return ""
	}

	r := []rune(alphabet)
	k := byte(len(r))

	s := make([]rune, n)
	b := make([]byte, n)

	_, _ = rand.Read(b)
	for i := range b {
		s[i] = r[b[i]%k]
	}

	return string(s)
}

// RandomNaturalNumber generates a random natural number that is less than n
func RandomNaturalNumber(n int) int {
	if n <= 0 {
		return 0
	}
	b := make([]byte, 4)
	_, _ = rand.Read(b)
	v := int(b[0])
	v |= int(b[1]) << 8
	v |= int(b[2]) << 16
	v |= int(b[3]) << 24
	if v < 0 {
		v = -v
	}
	return v % n
}

// IsDecimalString checks if string contains only decimal digits or not
func IsDecimalString(s string) bool {
	return CompiledRegex.DecimalString.MatchString(s)
}

// IsHeximalString checks if string contains only heximal digits or not
func IsHeximalString(s string) bool {
	return CompiledRegex.HeximalString.MatchString(s)
}

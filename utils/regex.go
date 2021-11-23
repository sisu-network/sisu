package utils

import "regexp"

var (
	DecimalRegex = "^[0-9]+$"
	HeximalRegex = "^[0-9a-f]+$"
)

// CompiledRegex contains some compiled common regex
var CompiledRegex = struct {
	DecimalString *regexp.Regexp
	HeximalString *regexp.Regexp
	Spaces        *regexp.Regexp
}{
	DecimalString: regexp.MustCompile(DecimalRegex),
	HeximalString: regexp.MustCompile(HeximalRegex),
	Spaces:        regexp.MustCompile(`\s+`),
}

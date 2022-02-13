package utils

import "math/big"

// Some alphabets
const (
	UpperEnglishLetters    = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	LowerEnglishLetters    = "abcdefghijklmnopqrstuvwxyz"
	UpperUnicodeLetters    = UpperEnglishLetters + "ÀÁẢÃẠĂẰẮẲẴẶÂẦẤẨẪẬÈÉẺẼẸÊỀẾỂỄỆÌÍỈĨỊÒÓỎÕỌÔỒỐỔỖỘƠỜỚỞỠỢÙÚỦŨỤƯỪỨỬỮỰ"
	LowerUnicodeLetters    = LowerEnglishLetters + "àáảãạăằắẳẵặâầấẩẫậèéẻẽẹêềếểễệìíỉĩịòóỏõọôồốổỗộơờớởỡợùúủũụưừứửữự"
	DecimalDigits          = "0123456789"
	SpecialCharacters      = "`~!@#$%^&*()-_=+[{}]\\|:;\"'<,.>/?"
	EnglishLetters         = UpperEnglishLetters + LowerEnglishLetters
	UnicodeLetters         = UpperUnicodeLetters + LowerUnicodeLetters
	AlphaNumericCharacters = EnglishLetters + DecimalDigits
	KeyboardCharacters     = AlphaNumericCharacters + SpecialCharacters
)

const (
	// A value used to replace floating unit. Instead of using float/double, we can use a big
	// integer to represent a floating number with smallest unit to be 10^9
	DecinmalUnit = 1_000_000_000
)

var (
	EthToWei = big.NewInt(1_000_000_000_000_000_000) // 10 ^ 18
	Gwei     = big.NewInt(1_000_000_000)
)

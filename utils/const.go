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
	SisuDecimalBase = big.NewInt(1_000_000_000_000_000_000) // 10 ^ 18
	EthToWei        = big.NewInt(1_000_000_000_000_000_000) // 10 ^ 18
	GweiToWei       = big.NewInt(1_000_000_000)
	Gwei            = big.NewInt(1_000_000_000)
	ZeroBigInt      = big.NewInt(0)

	OnePointSixEthToWei = big.NewInt(1_600_000_000_000_000_000) // 1.6 * 10 ^ 18
)

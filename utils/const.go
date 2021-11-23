package utils

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

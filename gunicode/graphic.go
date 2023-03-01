package gunicode

// IsWordCharacters tests a rune if it is a word character [a-zA-Z0-9] used in regular expression.
func IsWordCharacters(r rune) bool {
	return IsAlphabet(r) || (r >= '0' && r <= '9')
}

// IsAlphabet tests a rune if it is a character in alphabet [a-zA-Z].
func IsAlphabet(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

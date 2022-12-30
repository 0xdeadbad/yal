package lexer

func IsBase16(c byte) bool {
	return IsBase10(c) || c >= 'a' && c <= 'f' || c >= 'A' && c <= 'F'
}

func IsBase10(c byte) bool {
	return c >= '0' && c <= '9'
}

func IsBase8(c byte) bool {
	return c >= '0' && c <= '8'
}

func IsBase2(c byte) bool {
	return c == '0' || c == '1'
}

func IsAlpha(c byte) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
}

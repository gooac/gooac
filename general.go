
package gooa

import (
	"fmt"
	"unicode"
)

// Position Object
type Position struct {
	line 	int
	column 	int
	index 	*int
}

func (p Position) Copy() Position {
	return p
}

func (p Position) Fancy() string {
	return fmt.Sprintf("%v:%v", p.line, p.column)
}

// Tokenization Helpers
func IsNumeric(b byte) bool {
	return unicode.IsDigit(rune(b))
}

func IsHexNum(r byte) bool {
	b := unicode.ToLower(rune(r))
	return IsNumeric(r) || b == 'a' || b == 'b' || b == 'c' || b == 'd' || b == 'e' || b == 'f'
}

func IsWhitespace(b byte) bool {
	return unicode.IsSpace(rune(b))
}

func IsAlpha(b byte) bool {
	return unicode.IsLetter(rune(b))
}

func IsAlnum(b byte) bool {
	return IsAlpha(b) || IsNumeric(b)
}

func GetKeyword(s string) (bool, KeywordType) {
	t, val := keywordTypes[s]

	return val, t
}

func IsValidIdentChar(b byte) bool {
	return IsAlnum(b) || b == '_'
}

// Parser Helpers
func IsLiteralType(t TokenType) bool {
	return IsNumberLiteralType(t) || IsStringLiteralType(t)
}

func IsNumberLiteralType(t TokenType) bool {
	return t == TokenNumber ||
	t == TokenHexNumber ||
	t == TokenSciNot
}

func IsStringLiteralType(t TokenType) bool {
	return t == TokenString ||
	t == TokenMLString
}
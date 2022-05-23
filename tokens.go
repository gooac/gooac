package gooa

import "fmt"

type TokenType int
const (
	TokenEOF TokenType = iota	// <EOF>

	TokenIdent					// Any Identifier
	TokenNumber					// Any Number
	TokenHexNumber				// Any (hex) Number
	TokenString					// Any String
	TokenComment 				// Comment Starter (--)

	TokenKeyword				// if, then, end, ect

	TokenPeriod					// .
	TokenConcat					// ..
	TokenVariadic				// ...

	TokenLength					// #
	TokenComma					// ,

	TokenLParen					// (
	TokenRParen					// )
	TokenLBrac					// [
	TokenRBrac					// ]
	TokenLCurl					// {
	TokenRCurl					// }
	
	TokenColon					// :
	TokenLabel					// ::
	
	TokenAdd					// +
	TokenSub					// -
	TokenMul					// *
	TokenDiv 					// /
	TokenCarot					// ^
	TokenModulo					// %
	
	TokenEq						// =
	TokenLt						// <
	TokenGt						// >
	TokenLtEq					// <=
	TokenGtEq					// >=
	TokenNot					// !, ~
	TokenNotEq					// !=

								// Custom Tokens:
	TokenAttr					// $
)

var symbolLookups = map[string]TokenType {
	"." 		: TokenPeriod,
	"#" 		: TokenLength,
	"," 		: TokenComma,
	"(" 		: TokenLParen,
	")" 		: TokenRParen,
	"[" 		: TokenLBrac,
	"]" 		: TokenRBrac,
	"{" 		: TokenLCurl,
	"}" 		: TokenRCurl,
	":" 		: TokenColon,
	"+" 		: TokenAdd,
	"-" 		: TokenSub,
	"*" 		: TokenMul,
	"/" 		: TokenDiv,
	"^" 		: TokenCarot,
	"%" 		: TokenModulo,
	"=" 		: TokenEq,
	"<" 		: TokenLt,
	">" 		: TokenGt,
	"$" 		: TokenAttr,
	"!"			: TokenNot,
	"~"			: TokenNot,
	"..." 		: TokenVariadic,
	".." 		: TokenConcat,
	"::" 		: TokenLabel,
	"<=" 		: TokenLtEq,
	">=" 		: TokenGtEq,
}

var tokenNames = map[TokenType]string {
	TokenEOF: 			"<EOF>",
	TokenIdent: 		"<Identifier>",
	TokenKeyword: 		"<Keyword>",
	TokenNumber: 		"<Number>",
	TokenHexNumber:		"<HexNum>",
	TokenString: 		"<String>",
	TokenComment: 		"<Comment>",
	TokenPeriod: 		".",
	TokenConcat: 		"..",
	TokenVariadic: 		"...",
	TokenLength: 		"#",
	TokenComma: 		",",
	TokenLParen: 		"(",
	TokenRParen: 		")",
	TokenLBrac: 		"[",
	TokenRBrac: 		"]",
	TokenLCurl: 		"{",
	TokenRCurl: 		"}",
	TokenColon: 		":",
	TokenLabel: 		"::",
	TokenAdd: 			"+",
	TokenSub: 			"-",
	TokenMul: 			"*",
	TokenDiv: 			"/",
	TokenCarot: 		"^",
	TokenModulo: 		"%",
	TokenEq: 			"=",
	TokenLt: 			"<",
	TokenGt: 			">",
	TokenLtEq: 			"<=",
	TokenGtEq: 			">=",
	TokenNotEq:			"!=",
	TokenAttr: 			"$",
	TokenNot: 			"!",
}

type Token struct {
	toktype 		TokenType
	value 			string

	position 		TokenizerPosition

	tokenspecific 	interface{}
}

func (t *Token) Print() {
	if t.tokenspecific != nil {
		fmt.Printf("{Token \"%v\" = \"%v\" (%v)}\n", tokenNames[t.toktype]/* , t.position.line, t.position.column */, t.value, t.tokenspecific)
	} else {
		fmt.Printf("{Token \"%v\" = \"%v\"}\n", tokenNames[t.toktype]/* , t.position.line, t.position.column */, t.value)
	}
}
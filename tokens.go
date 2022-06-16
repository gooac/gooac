
package gooa

import (
	"fmt"
)

// Token Types
type TokenType int
const (
	TokenEOF TokenType = iota	// <EOF>

	TokenIdent					// Any Identifier
	TokenNumber					// Any Integer
	TokenHexNumber				// Any Hexadecimal Number Including "0x" 
	TokenSciNot					// Any Integer with Scientific Notation Exponentiation
	TokenString					// Any Normal String
	TokenMLString				// Any Multiline String
	TokenComment 				// Comment Starter (--)

	TokenKeyword				// if, then, end, ect

	TokenPeriod					// .
	TokenConcat					// ..
	TokenVariadic				// ...

	TokenLength					// #
	TokenComma					// ,
	TokenSemiColon				// ;

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
	TokenIsEq 					// ==
	TokenNot					// !, ~
	TokenNotEq					// !=

	TokenAnd 					// 'and'
	TokenOr 					// 'or'

	TokenTrue					// true
	TokenFalse					// false
	TokenNil					// nil

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
	";"			: TokenSemiColon,

	"true"		: TokenTrue,
	"false"		: TokenFalse,
	"nil"		: TokenNil,

	"..." 		: TokenVariadic,
	".." 		: TokenConcat,
	"::" 		: TokenLabel,
	"<=" 		: TokenLtEq,
	">=" 		: TokenGtEq,
	"~=" 		: TokenNotEq,
	"!=" 		: TokenNotEq,
	"==" 		: TokenIsEq,
}

var tokenNames = map[TokenType]string {
	TokenEOF: 			"<EOF>",
	TokenIdent: 		"<Identifier>",
	TokenKeyword: 		"<Keyword>",
	TokenNumber: 		"<Number>",
	TokenHexNumber: 	"<HexNumber>",
	TokenSciNot: 		"<SciNot>",
	TokenMLString: 		"<MultilineString>",
	TokenString: 		"<String>",
	TokenComment: 		"<Comment>",
	TokenPeriod: 		".",
	TokenConcat: 		"..",
	TokenVariadic: 		"...",
	TokenLength: 		"#",
	TokenComma: 		",",
	TokenSemiColon:		";",
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
	TokenIsEq: 			"==",
	TokenNotEq:			"!=",
	TokenAttr: 			"$",
	TokenNot: 			"!",

	TokenAnd: 			"and",
	TokenOr: 			"or",

	TokenTrue:			"true",
	TokenFalse:			"false",
	TokenNil:			"nil",
}

// Keyword Types
type KeywordType int
const (
	KeywordEmpty KeywordType = iota
	KeywordAnd	
	KeywordBreak
	KeywordDo
	KeywordElse
	KeywordElseif
	KeywordEnd
	KeywordFor
	KeywordFunction
	KeywordIf
	KeywordIn
	KeywordLocal
	KeywordNot
	KeywordOr
	KeywordRepeat
	KeywordReturn
	KeywordThen
	KeywordUntil
	KeywordWhile
	KeywordGoto
	KeywordTrue
	KeywordFalse
	KeywordNil
	KeywordContinue
)

var keywordTypes = map[string]KeywordType {
	"and"		: KeywordAnd,
	"break"		: KeywordBreak,
	"do"		: KeywordDo,
	"else"		: KeywordElse,
	"elseif"	: KeywordElseif,
	"elif"		: KeywordElseif,
	"end"		: KeywordEnd,
	"for"		: KeywordFor,
	"function"	: KeywordFunction,
	"fn"		: KeywordFunction,
	"if"		: KeywordIf,
	"in"		: KeywordIn,
	"local"		: KeywordLocal,
	"not"		: KeywordNot,
	"or"		: KeywordOr,
	"repeat"	: KeywordRepeat,
	"return"	: KeywordReturn,
	"then"		: KeywordThen,
	"until"		: KeywordUntil,
	"while"		: KeywordWhile,
	"goto"		: KeywordGoto,

	"true"		: KeywordTrue,
	"false"		: KeywordFalse,
	"nil"		: KeywordNil,

	"continue"	: KeywordContinue,
}

var keywordTypeValues = map[KeywordType]string {
	KeywordAnd: 		"and",
	KeywordBreak: 		"break",
	KeywordDo: 			"do",
	KeywordElse: 		"else",
	KeywordElseif: 		"elseif",
	KeywordEnd: 		"end",
	KeywordFor: 		"for",
	KeywordFunction: 	"function",
	KeywordIf: 			"if",
	KeywordIn: 			"in",
	KeywordLocal: 		"local",
	KeywordNot: 		"not",
	KeywordOr: 			"or",
	KeywordRepeat: 		"repeat",
	KeywordReturn: 		"return",
	KeywordThen: 		"then",
	KeywordUntil: 		"until",
	KeywordWhile: 		"while",
	KeywordGoto: 		"goto",
	KeywordTrue: 		"true",
	KeywordFalse: 		"false",
	KeywordNil: 		"nil",
	KeywordContinue: 	"continue",
} 

// Token Struct
type Token struct {
	Token 		TokenType				// What type of token it is
	Value 		string					// Actual value of the token
	Invalid 	bool 					// Is the token a valid token at all?
	
	Position	Position				// Copy of tokenizers starting positional data
	EndPos		Position				// Copy of tokenizers ending position
	
	WhiteSpace	int						// Whitespace preceding token
	Newlines 	int 					// Newlines after token

	Keyword 	KeywordType				// Type of keyword if a keyword at all	

	Special 	string					// Special value for strings
}

func (self *Token) ToASTValue() ASTValue {
	return ASTValue{
		Token: *self,
	}
}

func (self *Token) Print() {
	// fmt.Printf("tok(%v: ", tokenNames[self.Token])
	// print("[")
	// for k, v := range []byte(self.Value) {
	// 	print(v)
		
	// 	if k != len(self.Value) - 1 {
	// 		print(", ")
	// 	}
	// }
	// print("] = ")

	fmt.Printf("<Token|%s-%s| is '%s' = \"%s\">", self.Position.Fancy(), self.EndPos.Fancy(), tokenNames[self.Token], self.Value)
}

func (self *Token) Is(t Token) bool {
	return self.Invalid == t.Invalid &&
	self.Token == t.Token &&
	self.Value == t.Value 
}
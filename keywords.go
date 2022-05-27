
package gooa

type KeywordType int
const (
	KeywordAnd KeywordType = iota
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
)

var keywordLookup = map[string]TokenType {
	"and"		: TokenKeyword,
	"break"		: TokenKeyword,
	"do"		: TokenKeyword,
	"else"		: TokenKeyword,
	"elseif"	: TokenKeyword,
	"end"		: TokenKeyword,
	"for"		: TokenKeyword,
	"function"	: TokenKeyword,
	"if"		: TokenKeyword,
	"in"		: TokenKeyword,
	"local"		: TokenKeyword,
	"not"		: TokenKeyword,
	"or"		: TokenKeyword,
	"repeat"	: TokenKeyword,
	"return"	: TokenKeyword,
	"then"		: TokenKeyword,
	"until"		: TokenKeyword,
	"while"		: TokenKeyword,
	"goto"		: TokenKeyword,
	"false"		: TokenKeyword,
	"true"		: TokenKeyword,
	"nil"		: TokenKeyword,
}

var keywordTypes = map[string]KeywordType {
	"and"		: KeywordAnd,
	"break"		: KeywordBreak,
	"do"		: KeywordDo,
	"else"		: KeywordElse,
	"elseif"	: KeywordElseif,
	"end"		: KeywordEnd,
	"for"		: KeywordFor,
	"function"	: KeywordFunction,
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

	"false"		: KeywordTrue,
	"true"		: KeywordFalse,
	"nil"		: KeywordNil,
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
	KeywordTrue: 		"false",
	KeywordFalse: 		"true",
	KeywordNil: 		"nil",
}
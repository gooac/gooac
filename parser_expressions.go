package gooa

type ParserExpression interface {}

type ParserValue struct {
	tokens 		[]Token
}

type ParserExprDefine struct {
	tokens 		[]Token

	name 		string
	value 		ParserExpression
}

type ParserIdentifier struct {
	tokens 		[]*Token

	qualified 	string
	ismethod 	bool
}

package gooa

const (
	ParserErrorExpectedSymbol			= "Expected symbol '%v', got '%v'"
	ParserErrorExpectedExpression		= "Expected expression, got '%v'"
	ParserErrorExpectedIdentifier		= "Expected identifier, got '%v'"
	ParserErrorExpectedKeyword			= "Expected keyword '%v', got '%v'"
	ParserErrorExpectedLiteral			= "Expected literal, got '%v'"
	ParserErrorExpectedToken			= "Expected '%v', got '%v'"

	ParserErrorUnexpectedKeyword 		= "Unexpected keyword '%v'"
	ParserErrorUnexpectedX				= "Unexpected %v"
	ParserErrorUnexpectedEnd			= "Unexpected end of statement"

	ParserErrorExpectedArgumentName 	= "Expected argument name in function argument list, got '%v'"
	ParserErrorMissingEnd				= "Expected 'end' on %v block, got '<EOF>', did you miss an end?"
	ParserErrorAssigningToMethod 		= "Attempting to assign to Method Call (%v)"

	ParserErrorElseMustBeLast			= "'else' block must be the last declared block in an if statement"
	ParserErrorElseAlreadyDeclared		= "'else' block already declared in if statement"

	ParserErrorNumberUnexpectedHexNum	= "Cannot use Hexadecimal numbers on the right hand side of the decimal"
)

package gooa

// Errors for when a tokenization segment encounters an error
type TokenizationError string
var (
	TokErrNone 						TokenizationError = "No Error"
	TokErrMalformedNumber 			TokenizationError = "Malformed Number '%v'"
	TokErrMalformedHexLiteral 		TokenizationError = "Malformed Hex Literal '%v'"
	TokErrMalformedSciNotLiteral 	TokenizationError = "Malformed number in Scientific Notation '%v'"

	TokErrUnfinishedString 			TokenizationError = "Expected '%v', got <EOF>"
	TokErrUnfinishedMLString 		TokenizationError = "Expected \"%v\", got <EOF>"

	TokErrUnknownSymbol 			TokenizationError = "Unexpected Symbol '%v'"

)

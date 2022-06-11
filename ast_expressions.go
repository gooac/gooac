package gooa

// _ = 123			TODO ERROR, dont allow _ assignment

type NodeType string
const (
	NodeInvalid 				NodeType = "<!!  INVALID NODE  !!>"
	NodeProgram 				NodeType = "ProgramRoot"
	NodeComment 				NodeType = "Comment"

	NodeBool					NodeType = "Boolean"
	NodeNil						NodeType = "NilValue"
	NodeLiteral 				NodeType = "Literal"
	NodeLabel 					NodeType = "Label"
	NodeLength					NodeType = "Length"
	NodeNegate 					NodeType = "Negate"
	NodeNot						NodeType = "Not"
	NodeVariadicResolve			NodeType = "Variadic"

	NodeTable					NodeType = "Table"
	NodeTableMapValue			NodeType = "TableMapValue"
	NodeTableArrayValue			NodeType = "TableArrayValue"
	
	NodeVariableAssignment 		NodeType = "VarAssign"
	NodeLocalVariableAssignment NodeType = "LocalAssign"
	NodeLocalVariableStub 		NodeType = "LocalStub"

	NodeVariableNameList		NodeType = "NameList"
	NodeVariableValList 		NodeType = "ValueList"

	NodeIdentifierNormal 		NodeType = "IdentNormal"
	NodeIdentifierColon			NodeType = "IdentColon"
	NodeIdentifier				NodeType = "Identifier"

	NodeIndex					NodeType = "Indexing"

	NodeMemberIdent				NodeType = "DotIndex"
	NodeMemberMeth				NodeType = "ColonIndex"
	NodeMemberExpr				NodeType = "BracketIndex"

	NodeCall					NodeType = "Call"
	NodeMethodCall				NodeType = "MethodCall"
	
	NodeBinaryExpression		NodeType = "BinaryExpr"

	NodeAnonymousFunction		NodeType = "AnonFunc"
	NodeLocalFunction 			NodeType = "LocalFunction"
	NodeFunction 				NodeType = "Function"
	NodeArgumentList 			NodeType = "ArgList"
	NodeArgumentListOmitted 	NodeType = "OmittedArgList"
	NodeArgumentNormal			NodeType = "ArgNorm"
	NodeArgumentVariadic 		NodeType = "ArgVariadic"
	NodeNamedArgumentDef		NodeType = "ArgListMember"
	NodeNamedArgumentVariadic 	NodeType = "ArgListMemberVariadic"

	NodeReturn					NodeType = "Return"
	NodeArbitraryScope			NodeType = "ArbitraryScope"

	NodeIf						NodeType = "If"
	NodeIfScope					NodeType = "IfScope"
	NodeElseScope				NodeType = "ElseScope"
	NodeElseIfScope				NodeType = "ElseIfScope"

	NodeBreak 					NodeType = "Break"
	NodeContinue				NodeType = "Continue"

	NodeForIterator 			NodeType = "ForIterator"
	NodeForIteratorArgs			NodeType = "ForIteratorArgs"

	NodeForI					NodeType = "ForI"
	NodeForIIncr				NodeType = "ForIIncr"

	NodeRepeat					NodeType = "RepeatUntil"

	NodeWhile					NodeType = "WhileLoop"

	NodeGoto					NodeType = "Goto"
)
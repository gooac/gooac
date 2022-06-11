package gooa

type NodeType string
const (
	NodeProgram 				NodeType = "Program"
	NodeComment 				NodeType = "Comment"

	NodeBool					NodeType = "Bool"
	NodeNil						NodeType = "Nil"
	NodeLiteral 				NodeType = "Literal"
	NodeLabel 					NodeType = "Label"
	NodeLength					NodeType = "Length"
	NodeNegate 					NodeType = "Negate"
	NodeNot						NodeType = "Not"
	NodeVariadicResolve			NodeType = "VariadicResolve"

	NodeTable					NodeType = "Table"
	NodeTableMapValue			NodeType = "TableMapValue"
	NodeTableArrayValue			NodeType = "TableArrayValue"
	
	NodeVariableAssignment 		NodeType = "VariableAssignment"
	NodeLocalVariableAssignment NodeType = "LocalVariableAssignment"
	NodeLocalVariableStub 		NodeType = "LocalVariableStub"

	NodeVariableNameList		NodeType = "VariableNameList"
	NodeVariableValList 		NodeType = "VariableValList"

	NodeIdentSegNorm 			NodeType = "IdentifierNormal"
	NodeIdentSegColon			NodeType = "IdentifierColon"
	
	NodeIdentifier				NodeType = "Identifier"
	NodeIdentifierMethod		NodeType = "IdentifierMethod"

	NodeIndex					NodeType = "Index"

	NodeMemberIdent				NodeType = "MemberIdent"
	NodeMemberMeth				NodeType = "MemberMeth"
	NodeMemberExpr				NodeType = "MemberExpr"

	NodeCallArgs				NodeType = "CallArguments"
	NodeCall					NodeType = "Call"
	NodeMethodCall				NodeType = "MethodCall"
	
	NodeBinaryExpression		NodeType = "BinaryExpression"

	NodeAnonymousFunction		NodeType = "AnonymousFunction"
	NodeLocalFunction 			NodeType = "LocalFunction"
	NodeFunction 				NodeType = "Function"
	NodeArgumentList 			NodeType = "ArgumentList"
	NodeArgumentListOmitted 	NodeType = "ArgumentListOmitted"
	NodeArgumentNormal			NodeType = "ArgumentNormal"
	NodeArgumentVariadic 		NodeType = "ArgumentVariadic"
	NodeNamedArgumentDef		NodeType = "NamedArgumentDef"
	NodeNamedArgumentVariadic 	NodeType = "NamedArgumentVariadic"

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

	NodeRepeat					NodeType = "Repeat"

	NodeWhile					NodeType = "While"

	NodeGoto					NodeType = "Goto"
)
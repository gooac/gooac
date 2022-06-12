package gooa

import "fmt"

type Compiler struct {
	ast 		*AST
	output 		string

	err 		ErrorHandler
	lastloops 	[]*ASTNode
	hascontin	bool

	middleware 	*MiddlewareHandler

	luv 		string
}

func (self *Compiler) Reset() {
	self.ast = nil
	self.output = ""
	self.err = nil
	self.lastloops = []*ASTNode{}
	self.luv = "5.1"
}

func (self *Compiler) Compile(ast *AST, err *ErrorHandler) (string, bool) {
	self.Reset()

	self.err = *err
	self.ast = ast

	self.err.SetErrorRealm(ErrorRealmCompiler)

	l := len(ast.root.body) - 1
	for k, v := range ast.root.body {
		if self.err.ShouldImmediatelyStop() {
			break
		}

		switch v.nodetype {
		case NodeComment,
		NodeVariableAssignment,
		NodeLocalVariableStub,
		NodeLocalVariableAssignment,
		NodeCall,
		NodeMethodCall,
		NodeLocalFunction,
		NodeFunction,
		NodeReturn,
		NodeRepeat,
		NodeIf,
		NodeForI,
		NodeForIterator,
		NodeWhile,
		NodeGoto:
			self.output += self.CompileNode(v)
		default:
			self.err.Error(ErrorFatal, CompilerErrUnexpected, v.nodetype)
		}


		if k != l {
			self.output += ";"
		}
	}

	self.err.Dump()
	
	return self.middleware.PostCompile(self.output), self.err.ShouldStop()
}

func (self *Compiler) CompileNode(n *ASTNode) string {
	str := ""

	switch n.nodetype {
	case NodeLiteral:							str += self.CompileLiteral(n)
	case NodeIf:								str += self.CompileIf(n)
	case NodeIdentifierMethod, NodeIdentifier:	str += self.CompileIdentifier(n)
	case NodeMemberExpr: 						str += self.CompileMemberExpr(n)
	case NodeCall, NodeMethodCall:				str += self.CompileCall(n)
	case NodeVariableAssignment: 				str += self.CompileVarAssign(n)
	case NodeLocalVariableStub: 				str += self.CompileLocalStubs(n)
	case NodeLocalVariableAssignment: 			str += self.CompileLocalAssign(n)
	case NodeLabel: 							str += self.CompileLabel(n)
	case NodeGoto: 								str += self.CompileGoto(n)
	case NodeBool: 								str += self.CompileBool(n)
	case NodeNil: 								str += self.CompileNil(n)
	case NodeLength: 							str += self.CompileLength(n)
	case NodeNegate: 							str += self.CompileNegate(n)
	case NodeNot: 								str += self.CompileNot(n)
	case NodeVariadicResolve: 					str += self.CompileVar(n)
	case NodeTable: 							str += self.CompileTable(n)
	case NodeAnonymousFunction: 				str += self.CompileAnonFunc(n)
	case NodeLocalFunction: 					str += self.CompileLocalFunc(n)
	case NodeFunction: 							str += self.CompileFunc(n)
	case NodeBinaryExpression: 					str += self.CompileBinaryExpr(n)
	case NodeForI: 								str += self.CompileForI(n)
	case NodeForIterator: 						str += self.CompileForIter(n)
	case NodeContinue: 							str += self.CompileContinue(n)
	case NodeWhile: 							str += self.CompileWhile(n)
	case NodeRepeat: 							str += self.CompileRepeat(n)
	case NodeBreak: 							str += self.CompileBreak(n)
	}

	// fmt.Printf("[%v] %v\n", n.nodetype, str)
	return str
}

func (self *Compiler) Exists(n *ASTNode, index int) bool {
	if (len(n.body) - 1) <= index {
		return true
	}

	return false
}

func (self *Compiler) PushLoop(n *ASTNode) {
	self.lastloops = append(self.lastloops, n)
}

func (self *Compiler) PopLoop(n *ASTNode) string {
	self.lastloops = self.lastloops[:len(self.lastloops)-1]

	if self.hascontin {
		self.luv = "5.3/LUAJit"
		self.hascontin = false
		return fmt.Sprintf("::cont_%p::", n)
	}

	return ""
}
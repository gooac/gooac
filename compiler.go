package gooa

import "fmt"

type Compiler struct {
	ast 		*AST
	output 		string

	err 		ErrorHandler
}

func (self *Compiler) Reset() {
	self.ast = nil
	self.output = ""
	self.err = nil
}

func (self *Compiler) Compile(ast *AST, err *ErrorHandler) (string, bool) {
	self.err = *err
	self.ast = ast

	self.err.SetErrorRealm(ErrorRealmCompiler)

	for _, v := range ast.root.body {
		if self.err.ShouldImmediatelyStop() {
			break
		}

		self.output += self.CompileNode(v)
	}

	self.err.Dump()

	return "", true
}

func (self *Compiler) CompileNode(n *ASTNode) string {
	str := ""

	switch n.nodetype {
	case NodeLiteral:
		str += self.CompileLiteral(n)
	case NodeIdentifierMethod, NodeIdentifier:
		str += self.CompileIdentifier(n)
	case NodeCall, NodeMethodCall:
		str += self.CompilerCall(n)
	}

	fmt.Printf("[%v] %v\n", n.nodetype, str)
	return str
}

func (self *Compiler) Exists(n *ASTNode, index int) bool {
	if (len(n.body) - 1) <= index {
		return true
	}

	return false
}
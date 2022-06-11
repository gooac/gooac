
package gooa

const (
	CompilerErrUnexpected string		= "Unexpected expression '%v'"
	CompilerErrUnexpectedHex string		= "Unexpected Hex Literal '%v' near '%v'"
)

func (self *Compiler) Expect(n *ASTNode, a... NodeType) bool {
	if !self.IsNode(n, a...) {
		self.err.Error(ErrorFatal, CompilerErrUnexpected, n.nodetype)
		return true
	}

	return false
}


func (self *Compiler) Exclude(n *ASTNode, a... NodeType) bool {
	if self.IsNode(n, a...) {
		self.err.Error(ErrorFatal, CompilerErrUnexpected, n.nodetype)
		return true
	}

	return false
}

func (self *Compiler) IsNode(n *ASTNode, a... NodeType) bool {
	for _, v := range a {
		if n.nodetype == v {
			return true
		}
	}

	return false
}
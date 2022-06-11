
package gooa

// Primitive Converters
func (self *Compiler) CompileLiteral(n *ASTNode) string {
	s := ""

	first := n.values["value"].tokens[0]

	if first.token == TokenString {
		s += first.special + first.value + first.special
	} else if first.token == TokenMLString {
		s += "[" + first.special + "[" + 
		first.value + 
		"]" + first.special + "]"
	} else {
		for k, v := range n.values["value"].tokens {
			if k != 0 && v.token == TokenHexNumber {
				self.err.Error(ErrorGeneral, CompilerErrUnexpectedHex, v.value, s)
			}
			
			s += v.value
		}
	}

	return s
}

func (self *Compiler) CompileIdentifier(n *ASTNode) string {
	s := ""
	for k, v := range n.body {
		if v.nodetype == NodeIdentSegNorm {
			if k != 0 {
				s += "."
			}
			
			s += v.values["value"].token.value
		} else if v.nodetype == NodeIdentSegColon {
			s += ":" + v.values["value"].token.value
		} else {
			s += "[" + self.CompileNode(v) + "]"
		}
	}

	return s
}

// Real Expression Converters
func (self *Compiler) CompilerCall(n *ASTNode) string {
	n.Dump(1)

	s := ""

	s += self.CompileNode(n.body[0])

	if self.Exists(n, 1) {
		next := n.body[1]

		s += self.CompilerCallArguments(next)
	}

	return s
}

func (self *Compiler) CompilerCallArguments(n *ASTNode) string {
	s := "("

	for k, v := range n.body {
		if k != 0 {
			s += ","
		}

		s += self.CompileNode(v)
	}

	s += ")"

	return s
}
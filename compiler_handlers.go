package gooa

import "fmt"

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

func (self *Compiler) CompileMemberExpr(n *ASTNode) string {
	s := self.CompileNode(n.body[0])

	s += "["
	s += self.CompileNode(n.body[1])
	s += "]"

	return s
}

func (self *Compiler) CompileScope(n *ASTNode, skip int) string {
	s := ""

	for k, v := range n.body {
		if k < skip {
			continue
		}

		s += self.CompileNode(v) + ";"
	}

	return s
}

// Real Expression Converters
func (self *Compiler) CompileCall(n *ASTNode) string {
	s := ""

	s += self.CompileNode(n.body[0])

	if self.Exists(n, 1) {
		next := n.body[1]

		s += self.CompileCallArguments(next)
	}

	return s
}

func (self *Compiler) CompileCallArguments(n *ASTNode) string {
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

func (self *Compiler) CompileIf(n *ASTNode) string {
	elsed := false
	var err string

	s := "if("
	s += self.CompileNode(n.body[0])
	s += ")then;"

	for k, v := range n.body {
		if k == 0 {
			continue
		}

		if v.nodetype == NodeIfScope {
			s += self.CompileScope(v, 0)
		} else if v.nodetype == NodeElseScope {
			if len(v.body) == 0 {
				continue
			}

			if elsed {
				err = CompilerErrMultipleElse
			}
			
			elsed = true
			s += "else;" + self.CompileScope(v, 0)
		} else {
			if elsed {
				err = CompilerErrElseifIsLast
			}

			s += "elseif(" + self.CompileNode(v.body[0]) + ")then;"

			for k, vv := range v.body {
				if k == 0 {
					continue
				}

				s += self.CompileNode(vv) + ";" 
			}
		}
	}

	if err != "" {
		self.err.Error(ErrorFatal, err)
	}

	s += "end;"
	return s
}

func (self *Compiler) CompileVarAssign(n *ASTNode) string {
	s := ""

	names := n.body[0]
	values := n.body[1]

	for k, v := range names.body {
		if k != 0 {
			s += ","
		}

		s += self.CompileNode(v)
	}

	
	s += "="

	for k,v := range values.body {
		if k != 0 {
			s += ","
		}

		s += self.CompileNode(v)
	}

	return s
}

func (self *Compiler) CompileLocalStubs(n *ASTNode) string {
	s := "local "

	for k, v := range n.body[0].body {
		if k != 0 {
			s += ","
		}

		s += self.CompileNode(v)
	}

	return s
}

func (self *Compiler) CompileLocalAssign(n *ASTNode) string {
	s := "local "

	for k, v := range n.body[0].body {
		if k != 0 {
			s += ","
		}

		s += self.CompileNode(v)
	}

	s += "="

	for k, v := range n.body[1].body {
		if k != 0 {
			s += ","
		}

		s += self.CompileNode(v)
	}

	return s
}

func (self *Compiler) CompileLabel(n *ASTNode) string {
	return "::" + n.values["label"].token.value + "::"
}

func (self *Compiler) CompileGoto(n *ASTNode) string {
	return "goto " + n.values["label"].token.value
}

func (self *Compiler) CompileBool(n *ASTNode) string {
	return n.values["value"].token.value
}

func (self *Compiler) CompileNil(n *ASTNode) string {
	return "nil"
}

func (self *Compiler) CompileLength(n *ASTNode) string {
	return "#" + self.CompileNode(n.body[0])
}

func (self *Compiler) CompileNegate(n *ASTNode) string {
	return "-" + self.CompileNode(n.body[0])
}

func (self *Compiler) CompileNot(n *ASTNode) string {
	return "not " + self.CompileNode(n.body[0])
}

func (self *Compiler) CompileVar(n *ASTNode) string {
	return "..."
}

func (self *Compiler) CompileTable(n *ASTNode) string {
	s := "({"

	ln := len(n.body) - 1
	append := ","
	for k, v := range n.body {
		if k == ln {
			append = ""
		}

		if v.nodetype == NodeTableArrayValue {
			s += "" + self.CompileNode(v.body[0]) + "" + append
		} else if v.nodetype == NodeTableMapValue {
			pref := self.CompileNode(v.body[0])

			if v.body[0].nodetype != NodeIdentifier {
				pref = "[" + pref + "]"
			}
			
			s += pref + "=" + self.CompileNode(v.body[1]) + "" + append
		}
	}

	return s + "})"
}

func (self *Compiler) CompileAnonFunc(n *ASTNode) string {
	s := "(function"

	s += self.CompileFuncArgs(n.body[0])

	for k, v := range n.body {
		if k == 0 {
			continue
		}

		s += self.CompileNode(v) + ";"
	}

	s += "end)"
	return s
}

func (self *Compiler) CompileFuncArgs(n *ASTNode) string {
	if n.nodetype == NodeArgumentListOmitted {
		return "()"
	}

	s := "("
	post := ")"
	
	length := len(n.body) - 1
	apnd := ","
	for k, v := range n.body {
		if k == length {
			apnd = ""
		}
		
		if v.nodetype == NodeArgumentNormal {
			s += v.values["name"].token.value + apnd
		} else if v.nodetype == NodeArgumentVariadic {
			s += "..." + apnd
		} else if v.nodetype == NodeNamedArgumentDef {
			name := v.values["name"].token.value
			s += name + apnd

			post += name + "=" + name + " or " + self.CompileNode(v.body[0]) + ";"
		} else if v.nodetype == NodeNamedArgumentVariadic {
			s += "..." + apnd
			
			post += "local " + v.values["name"].token.value + "={...};"
		}
	}

	return s + post
}

func (self *Compiler) CompileLocalFunc(n *ASTNode) string {
	return "local " + self.CompileFunc(n)
}

func (self *Compiler) CompileFunc(n *ASTNode) string {
	s := "function "

	s += self.CompileIdentifier(n.body[0])
	s += self.CompileFuncArgs(n.body[1])

	for k, v := range n.body {
		if k <= 1 {
			continue
		}

		s += self.CompileNode(v) + ";"
	}

	s += "end;"

	return s
}

func (self *Compiler) CompileBinaryExpr(n *ASTNode) string {
	return "(" + self.CompileNode(n.body[0]) + " " + n.values["operator"].token.value + " "  + self.CompileNode(n.body[1]) + ")"
}

func (self *Compiler) CompileForI(n *ASTNode) string {
	s := "for "
	skip := 2

	s += n.values["identifier"].token.value + "="

	s += self.CompileNode(n.body[0])

	s += ","

	s += self.CompileNode(n.body[1])

	if (len(n.body) - 1) >= 2 && n.body[2].nodetype == NodeForIIncr {
		s += "," + self.CompileNode(n.body[2].body[0])

		skip = 3
	}

	s += " do;"

	self.PushLoop(n)
	s += self.CompileScope(n, skip)
	s += self.PopLoop(n)

	s += "end;"

	return s
}

func (self *Compiler) CompileForIter(n *ASTNode) string {
	s := "for "

	l := len(n.body[0].body) - 1
	for k, v := range n.body[0].body {
		s += v.values["value"].token.value
		
		if k != l {
			s += ","
		}
	}

	s += " in "

	s += self.CompileNode(n.body[1])

	if n.body[1].nodetype != NodeCall && n.body[1].nodetype != NodeMethodCall {
		s += " "
	}
	s += "do;"

	self.PushLoop(n)
	s += self.CompileScope(n, 2)
	s += self.PopLoop(n)

	
	s += "end;"

	return s
}

func (self *Compiler) CompileContinue(n *ASTNode) string {
	last := self.lastloops[len(self.lastloops) - 1]

	self.hascontin = true
	
	return fmt.Sprintf("goto cont_%p", last)
}
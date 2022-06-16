package gooa

import "fmt"

// Primitive Converters
func (self *Compiler) CompileLiteral(n *ASTNode) string {
	s := ""

	first := n.Values["value"].tokens[0]

	if first.token == TokenString {
		s += first.special + first.value + first.special
	} else if first.token == TokenMLString {
		s += "[" + first.special + "[" + 
		first.value + 
		"]" + first.special + "]"
	} else {
		for k, v := range n.Values["value"].tokens {
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
	for k, v := range n.Body {
		if v.Nodetype == NodeIdentSegNorm {
			if k != 0 {
				s += "."
			}
			
			s += v.Values["value"].token.value
		} else if v.Nodetype == NodeIdentSegColon {
			s += ":" + v.Values["value"].token.value
		} else {
			s += "[" + self.CompileNode(v) + "]"
		}
	}

	return s
}

func (self *Compiler) CompileMemberExpr(n *ASTNode) string {
	s := self.CompileNode(n.Body[0])

	s += "["
	s += self.CompileNode(n.Body[1])
	s += "]"

	return s
}

func (self *Compiler) CompileScope(n *ASTNode, skip int) string {
	s := ""

	for k, v := range n.Body {
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

	s += self.CompileNode(n.Body[0])

	if self.Exists(n, 1) {
		next := n.Body[1]

		s += self.CompileCallArguments(next)
	}

	return s
}

func (self *Compiler) CompileCallArguments(n *ASTNode) string {
	s := "("

	for k, v := range n.Body {
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
	s += self.CompileNode(n.Body[0])
	s += ")then;"

	for k, v := range n.Body {
		if k == 0 {
			continue
		}

		if v.Nodetype == NodeIfScope {
			s += self.CompileScope(v, 0)
		} else if v.Nodetype == NodeElseScope {
			if len(v.Body) == 0 {
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

			s += "elseif(" + self.CompileNode(v.Body[0]) + ")then;"

			for k, vv := range v.Body {
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

	names := n.Body[0]
	Values := n.Body[1]

	for k, v := range names.Body {
		if k != 0 {
			s += ","
		}

		s += self.CompileNode(v)
	}

	
	s += "="

	for k,v := range Values.Body {
		if k != 0 {
			s += ","
		}

		s += self.CompileNode(v)
	}

	return s
}

func (self *Compiler) CompileLocalStubs(n *ASTNode) string {
	s := "local "

	for k, v := range n.Body[0].Body {
		if k != 0 {
			s += ","
		}

		s += self.CompileNode(v)
	}

	return s
}

func (self *Compiler) CompileLocalAssign(n *ASTNode) string {
	s := "local "

	for k, v := range n.Body[0].Body {
		if k != 0 {
			s += ","
		}

		s += self.CompileNode(v)
	}

	s += "="

	for k, v := range n.Body[1].Body {
		if k != 0 {
			s += ","
		}

		s += self.CompileNode(v)
	}

	return s
}

func (self *Compiler) CompileLabel(n *ASTNode) string {
	return "::" + n.Values["label"].token.value + "::"
}

func (self *Compiler) CompileGoto(n *ASTNode) string {
	return "goto " + n.Values["label"].token.value
}

func (self *Compiler) CompileBool(n *ASTNode) string {
	return n.Values["value"].token.value
}

func (self *Compiler) CompileNil(n *ASTNode) string {
	return "nil"
}

func (self *Compiler) CompileLength(n *ASTNode) string {
	return "#" + self.CompileNode(n.Body[0])
}

func (self *Compiler) CompileNegate(n *ASTNode) string {
	return "-" + self.CompileNode(n.Body[0])
}

func (self *Compiler) CompileNot(n *ASTNode) string {
	return "not " + self.CompileNode(n.Body[0])
}

func (self *Compiler) CompileVar(n *ASTNode) string {
	return "..."
}

func (self *Compiler) CompileTable(n *ASTNode) string {
	s := "({"

	ln := len(n.Body) - 1
	append := ","
	for k, v := range n.Body {
		if k == ln {
			append = ""
		}

		if v.Nodetype == NodeTableArrayValue {
			s += "" + self.CompileNode(v.Body[0]) + "" + append
		} else if v.Nodetype == NodeTableMapValue {
			pref := self.CompileNode(v.Body[0])

			if v.Body[0].Nodetype != NodeIdentifier {
				pref = "[" + pref + "]"
			}
			
			s += pref + "=" + self.CompileNode(v.Body[1]) + "" + append
		}
	}

	return s + "})"
}

func (self *Compiler) CompileAnonFunc(n *ASTNode) string {
	s := "(function"

	s += self.CompileFuncArgs(n.Body[0])

	for k, v := range n.Body {
		if k == 0 {
			continue
		}

		s += self.CompileNode(v) + ";"
	}

	s += "end)"
	return s
}

func (self *Compiler) CompileFuncArgs(n *ASTNode) string {
	if n.Nodetype == NodeArgumentListOmitted {
		return "()"
	}

	s := "("
	post := ")"
	
	length := len(n.Body) - 1
	apnd := ","
	for k, v := range n.Body {
		if k == length {
			apnd = ""
		}
		
		if v.Nodetype == NodeArgumentNormal {
			s += v.Values["name"].token.value + apnd
		} else if v.Nodetype == NodeArgumentVariadic {
			s += "..." + apnd
		} else if v.Nodetype == NodeNamedArgumentDef {
			name := v.Values["name"].token.value
			s += name + apnd

			post += name + "=" + name + " or " + self.CompileNode(v.Body[0]) + ";"
		} else if v.Nodetype == NodeNamedArgumentVariadic {
			s += "..." + apnd
			
			post += "local " + v.Values["name"].token.value + "={...};"
		}
	}

	return s + post
}

func (self *Compiler) CompileLocalFunc(n *ASTNode) string {
	return "local " + self.CompileFunc(n)
}

func (self *Compiler) CompileFunc(n *ASTNode) string {
	s := "function "

	s += self.CompileIdentifier(n.Body[0])
	s += self.CompileFuncArgs(n.Body[1])

	for k, v := range n.Body {
		if k <= 1 {
			continue
		}

		s += self.CompileNode(v) + ";"
	}

	s += "end;"

	return s
}

func (self *Compiler) CompileBinaryExpr(n *ASTNode) string {
	return "(" + self.CompileNode(n.Body[0]) + " " + n.Values["operator"].token.value + " "  + self.CompileNode(n.Body[1]) + ")"
}

func (self *Compiler) CompileForI(n *ASTNode) string {
	s := "for "
	skip := 2

	s += n.Values["identifier"].token.value + "="

	s += self.CompileNode(n.Body[0])

	s += ","

	s += self.CompileNode(n.Body[1])

	if (len(n.Body) - 1) >= 2 && n.Body[2].Nodetype == NodeForIIncr {
		s += "," + self.CompileNode(n.Body[2].Body[0])

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

	l := len(n.Body[0].Body) - 1
	for k, v := range n.Body[0].Body {
		s += v.Values["value"].token.value
		
		if k != l {
			s += ","
		}
	}

	s += " in "

	s += self.CompileNode(n.Body[1])

	if n.Body[1].Nodetype != NodeCall && n.Body[1].Nodetype != NodeMethodCall {
		s += " "
	}
	s += "do;"

	self.PushLoop(n)
	s += self.CompileScope(n, 2)
	s += self.PopLoop(n)

	
	s += "end;"

	return s
}

func (self *Compiler) CompileWhile(n *ASTNode) string {
	s := "while "

	s += self.CompileNode(n.Body[0])

	s += " do;"

	self.PushLoop(n)
	s += self.CompileScope(n, 1)
	s += self.PopLoop(n)

	s += "end;"

	return s
}

func (self *Compiler) CompileRepeat(n *ASTNode) string {
	s := "repeat "

	self.PushLoop(n)
	s += self.CompileScope(n.Body[1], 0)
	s += self.PopLoop(n)

	s += "until " + self.CompileNode(n.Body[0])
	
	return s
}


func (self *Compiler) CompileContinue(n *ASTNode) string {
	if len(self.lastloops) == 0 {
		self.err.Error(ErrorGeneral, CompilerErrUsedOutsideLoop, "continue")
		return ""
	}

	last := self.lastloops[len(self.lastloops) - 1]

	self.hascontin = true
	
	return fmt.Sprintf("goto cont_%p", last)
}	

func (self *Compiler) CompileBreak(n *ASTNode) string {
	if len(self.lastloops) == 0 {
		self.err.Error(ErrorGeneral, CompilerErrUsedOutsideLoop, "break")
		return ""
	}

	return "break"
}	
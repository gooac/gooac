package gooa

import (
	"fmt"
	"strings"
)

type AST struct {
	root    *ASTNode
	curnode *ASTNode
	err     ErrorHandler
	last    *ASTNode
}

type ASTNode struct {
	nodetype NodeType
	body     []*ASTNode

	values map[string]ASTValue
	parent *ASTNode

	callee *ASTNode
	trailing *ASTNode
}

type ASTValue struct {
	token 		Token
	tokens 		[]Token
}

func CreateAST(err ErrorHandler) *AST {
	ast := &AST{
		root: &ASTNode{
			nodetype: NodeProgram,
			body:     []*ASTNode{},
			values:   map[string]ASTValue{},
		},
	}

	ast.curnode = ast.root
	ast.err = err

	return ast
}

func (self *AST) OpenScope(n *ASTNode) {
	n.parent = self.curnode
	self.curnode = n
}

func (self *AST) CloseScope() *ASTNode {
	if self.curnode == self.root {
		self.err.Error(ErrorFatal, "Attemping to pop root node!")
		return nil
	}

	old := self.curnode
	self.curnode = self.curnode.parent
	
	return old
}

func (self *AST) Add(n *ASTNode) {
	n.parent = self.curnode

	self.last = n
	self.curnode.body = append(self.curnode.body, n)
}

func (self *ASTNode) Dump(lvl int) {
	fmt.Printf("%v┌%v", strings.Repeat("│", lvl), string(self.nodetype))

	if self.callee != nil {
		fmt.Printf(" (callee: %v)", self.callee)
	}
	println()

	if self.values != nil {
		ind := strings.Repeat("│", lvl)
		cnt := 0
		for k, v := range self.values {
			cnt++

			ch := "├"
			if (cnt == len(self.values)) && len(self.body) == 0 {
				ch = "└"
			}
			
			print(ind, ch, k, ": ")
			v.Print()
		}
	}

	if self.trailing != nil {
		ind := strings.Repeat("│", lvl)

		if self.trailing.nodetype != NodeIndex {
			fmt.Printf("%v]%v\n", ind, self.trailing)
		} else {
			fmt.Printf("%v] ", ind)
			for _, v := range self.trailing.body {
				if v.nodetype == NodeIdentifierNormal {
					print(".", v.values["value"].token.value)
				} else if v.nodetype == NodeIdentifierColon {
					print(":", v.values["value"].token.value)
				} else {
					fmt.Printf("[%v]", v.nodetype)
				}
			}
			println()
		}
	}

	for _, v := range self.body {
		if v == nil {
			println("ATTEMPING TO READ NIL ", v)
			continue
		}
		v.Dump(lvl + 1)
	}
}

// DEBUG: Print everything in the AST structure recusively
func (self *AST) Dump() {
	self.root.Dump(0)
}

// DEBUG: Format and print the ASTValue
func (self *ASTValue) Print() {
	defer fmt.Printf(">\n")
	
	if self.tokens != nil {
		fmt.Printf("<ASTValue(tokens): ")

		for k, v := range self.tokens {
			v.Print()
			
			if k != len(self.tokens) - 1 {
				fmt.Printf(", ")
			}
		}
		return
	}

	fmt.Printf("<ASTValue: ")
	self.token.Print()
}
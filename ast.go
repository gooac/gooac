package gooa

import (
	"fmt"
	"strings"
)

type AST struct {
	RootNode    *ASTNode
	CurNode *ASTNode
	LastNode    *ASTNode
	err     ErrorHandler
}

type ASTNode struct {
	Nodetype NodeType
	Body     []*ASTNode

	Values map[string]ASTValue
	Parent *ASTNode
}

type ASTValue struct {
	token 		Token
	tokens 		[]Token
}

func CreateAST(err ErrorHandler) *AST {
	ast := &AST{
		RootNode: &ASTNode{
			Nodetype: NodeProgram,
			Body:     []*ASTNode{},
			Values:   map[string]ASTValue{},
		},
	}

	ast.CurNode = ast.RootNode
	ast.err = err

	return ast
}

func (self *AST) OpenScope(n *ASTNode) {
	n.Parent = self.CurNode
	self.CurNode = n
}

func (self *AST) CloseScope() *ASTNode {
	if self.CurNode == self.RootNode {
		self.err.Error(ErrorFatal, "Attemping to pop RootNode node!")
		return nil
	}

	old := self.CurNode
	self.CurNode = self.CurNode.Parent
	
	return old
}

func (self *AST) Add(n *ASTNode) {
	n.Parent = self.CurNode

	self.LastNode = n
	self.CurNode.Body = append(self.CurNode.Body, n)
}

func (self *ASTNode) Dump(lvl int) {
	fmt.Printf("%v┌%v", strings.Repeat("│", lvl), string(self.Nodetype))

	println()

	if self.Values != nil {
		ind := strings.Repeat("│", lvl)
		cnt := 0
		for k, v := range self.Values {
			cnt++

			ch := "├"
			if (cnt == len(self.Values)) && len(self.Body) == 0 {
				ch = "└"
			}
			
			print(ind, ch, k, ": ")
			v.Print()
		}
	}

	for _, v := range self.Body {
		if v == nil {
			println("ATTEMPING TO READ NIL ", v)
			continue
		}
		v.Dump(lvl + 1)
	}
}

// DEBUG: Print everything in the AST structure recusively
func (self *AST) Dump() {
	self.RootNode.Dump(0)
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
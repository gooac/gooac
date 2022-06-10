package gooa

type IdentType int

const (
	IdentNormal IdentType = iota
	IdentTable
	IdentMethod
	IdentSubscripted
)

func (self *Parser) HandleIdentStub() *ASTNode {
	sym := self.Consume()
	ident := self.Consume()

	expr := &ASTNode{
		nodetype: NodeIdentifierNormal,
		values: map[string]ASTValue{
			"value": {
				token: ident,
			},
		},
	}

	if ident.token != TokenIdent {
		self.err.Error(ErrorFatal, ParserErrorExpectedIdentifier, ident.value)
		return expr
	}

	if sym.token == TokenColon {
		expr.nodetype = NodeIdentifierColon
	}

	return expr
}

func (self *Parser) HandleAssignmentIdent() *ASTNode {
	id := self.QualifyIdent()

	return self.HandleTrails(&id)
}

// Versatile 'TokenIdent' handler
func (self *Parser) QualifyIdent() ASTNode {
	node := ASTNode{
		nodetype: NodeIdentifier,
	}

	initial := self.Consume()

	if initial.token != TokenIdent {
		self.err.Error(ErrorFatal, ParserErrorExpectedIdentifier, initial.value)
		return node
	}

	node.body = []*ASTNode{
		{nodetype: NodeIdentifierNormal,
			values: map[string]ASTValue{
				"value": {
					token: initial,
				},
			},
		},
	}

	self.HandleIndexing(&node)

	return node
}

// Index Handler
func (self *Parser) HandleIndexing(node *ASTNode) *ASTNode {
	for {
		sym := self.Peek()

		if sym.token == TokenLBrac {
			self.Consume()
			expr := self.ExpectExpression(self.Peek())

			if self.err.ShouldImmediatelyStop() {
				return node
			}

			node.body = append(node.body, expr)

			rbrac := self.Consume()
			if rbrac.token != TokenRBrac {
				self.err.Error(ErrorFatal, ParserErrorExpectedSymbol, "]", rbrac.value)
			}

			continue
		} else if (sym.token != TokenPeriod) && (sym.token != TokenColon) {
			break
		}

		node.body = append(node.body, self.HandleIdentStub())
	}

	return node
}

// Handle Index + Create NodeIndex
func (self *Parser) NewHandleIndexing() *ASTNode {
	p := self.Peek()

	switch p.token {
	case TokenLBrac, TokenPeriod, TokenColon:
		node := &ASTNode{
			nodetype: NodeIndex,
			body: []*ASTNode{},
		}

		self.HandleIndexing(node)
		return node
	default:
		return nil
	}
}

// Handle Index + Create Node + Add onto existing node
func (self *Parser) HandleIndexingNode(n *ASTNode) *ASTNode {
	indexing := self.NewHandleIndexing()
	
	if indexing != nil {
		n.trailing = indexing
	}

	return n
}

var Priorities = map[TokenType]struct {
	left  int
	right int
}{
	TokenAdd: {6, 6},
	TokenSub: {6, 6},

	TokenMul:    {7, 7},
	TokenModulo: {7, 7},
	TokenDiv:    {7, 7},

	TokenCarot:  {10, 9},
	TokenConcat: {5, 4},

	TokenIsEq: {3, 3},
	TokenLt:   {3, 3},
	TokenLtEq: {3, 3},

	TokenNotEq: {3, 3},
	TokenGt:    {3, 3},
	TokenGtEq:  {3, 3},

	TokenAnd: {2, 2},
	TokenOr:  {1, 1},
}

// Binary Expression Handler
func (self *Parser) HandleBinExpr(lhs_ast *ASTNode) *ASTNode {
	p := self.Peek()

	prio, is := Priorities[p.token]

	if !is {
		return lhs_ast
	}

	self.Consume()

	rhs_ast := self.ExpectExpression(self.Peek())

	if self.err.ShouldImmediatelyStop() {
		return nil
	}

	body := []*ASTNode{
		lhs_ast,
		rhs_ast,
	}

	// I have absolutely no clue if priorities are implemented
	// properly, as i have absolutely no idea what they are
	// supposed to indicate :) i know im a dumdum
	// just... verify this later it isnt important now
	if (lhs_ast.nodetype == NodeBinaryExpression) &&
		(Priorities[lhs_ast.values["operator"].token.token].left > prio.left) {

		body[0] = rhs_ast
		body[1] = lhs_ast
	}

	node := &ASTNode{
		nodetype: NodeBinaryExpression,
		body:     body,
		values: map[string]ASTValue{
			"operator": {
				token: p,
			},
		},
	}

	return node
}

// Wraps Binary Expression Handler
// Provide a literal token to be converted
// into a ASTNode
func (self *Parser) HandleBinExprLiteral(lit Token) *ASTNode {
	toks := []Token{
		self.Consume(),
	}

	if IsNumberLiteralType(lit.token) && (self.Peek().token == TokenPeriod) {
		prd := self.Consume()

		next := self.Peek()

		if IsNumberLiteralType(next.token) {
			self.Consume()
			toks = append(toks, next)
		} else if next.token == TokenHexNumber {
			self.Consume()
			self.err.Error(ErrorFatal, ParserErrorNumberUnexpectedHexNum)
			toks = append(toks, next)
		} else {
			toks = append(toks, prd)
		}
	}

	return self.HandleBinExpr(&ASTNode{
		nodetype: NodeLiteral,
		values: map[string]ASTValue{
			"value": {
				tokens: toks,
			},
		},
	})
}

// Handles all forms of trails
func (self *Parser) HandleTrails(n *ASTNode) *ASTNode {
	for {
		p := self.Peek()

		switch p.token {
		case TokenLParen:
			self.Consume()
	
			callargs := self.HandleCall()
		
			body := []*ASTNode{
				n,
			}

			body = append(body, callargs...)

			n = &ASTNode{
				nodetype: NodeCall,
				body: body,
			}
		
			continue
		case TokenString, TokenMLString:
			n = &ASTNode{
				nodetype: NodeCall,
				body: []*ASTNode{
					n,
					{
						nodetype: NodeLiteral,
						values: map[string]ASTValue{
							"value": {
								tokens: []Token{
									self.Consume(),
								},
							},
						},
					},
				},
			}
		case TokenLCurl:
			node := self.HandleTable(self.Peek())
			n = &ASTNode{
				nodetype: NodeCall,
				body: []*ASTNode{
					n,
					&node,
				},
			}

		// Identifier Resolution
		case TokenColon:
		
			self.Consume()

			id := self.Peek()

			if id.token != TokenIdent {
				self.err.Error(ErrorFatal, ParserErrorExpectedIdentifier, id.value)
				break
			}

			n = &ASTNode{
				nodetype: NodeMemberMeth,
				values: map[string]ASTValue{
					"ident":{
						token: self.Consume(),
					},
				},
				callee: n,
			}

			paren := self.Consume()
			
			if paren.token != TokenLParen {
				self.err.Error(ErrorFatal, ParserErrorExpectedSymbol, "(",)
				break
			}

			callargs := self.HandleCall()

			n = &ASTNode{
				nodetype: NodeMethodCall,
				body: callargs,
				callee: n,
			}

			continue
		case TokenPeriod:
			self.Consume()

			id := self.Peek()

			if id.token != TokenIdent {
				self.err.Error(ErrorFatal, ParserErrorExpectedIdentifier, id.value)
				break
			}

			n = &ASTNode{
				nodetype: NodeMemberIdent,
				values: map[string]ASTValue{
					"ident":{
						token: self.Consume(),
					},
				},
				body: []*ASTNode{
					n,
				},
			}
			continue
		case TokenLBrac:
			self.Consume()

			pp := self.Peek()
			expr := self.ExpectExpression(pp)

			if expr == nil {
				self.err.Error(ErrorFatal, ParserErrorExpectedExpression, pp.token)
				break
			}

			pp = self.Consume()

			if pp.token != TokenRBrac {
				self.err.Error(ErrorFatal, ParserErrorExpectedSymbol, "]", pp.value)
				break
			}

			n = &ASTNode{
				nodetype: NodeMemberExpr,
				body: []*ASTNode{
					n,
					expr,
				},
			}


		}

		break
	}

	return self.HandleBinExpr(n)
}

package gooa

func (self *Parser) ExpectExpression(p Token) *ASTNode {
	expr := self.GetExpression()

	if expr == nil {
		self.err.Error(ErrorFatal, ParserErrorExpectedExpression, p.value)
	}

	return expr
}

func (self *Parser) GetExpression() *ASTNode {
	for {
		semi := self.Peek()

		if semi.token == TokenSemiColon {
			self.Consume()
			continue
		} else {
			break
		}
	}

	p := self.Peek()

	if IsLiteralType(p.token) {
		bin := self.HandleBinExprLiteral(p)

		return bin
	}

	var node *ASTNode

	switch p.token {
	case TokenIdent:
		node = self.HandleAssignmentIdent()
	case TokenLCurl:
		fptrs := self.HandleTable(p)
		node = &fptrs
	case TokenLParen:
		self.Consume()
		expr := self.ExpectExpression(self.Peek())

		rpar := self.Peek()
		if rpar.token != TokenRParen {
			self.err.Error(ErrorFatal, ParserErrorExpectedSymbol, ")", rpar.value)
		}

		self.Consume()

		node = expr

	case TokenTrue, TokenFalse:
		self.Consume()
		node = &ASTNode{
			nodetype: NodeBool,
			values: map[string]ASTValue{
				"value": {
					token: p,
				},
			},
		}
	case TokenNil:
		self.Consume()
		node = &ASTNode{
			nodetype: NodeNil,
		}

	case TokenKeyword:
		node = self.HandleExpressionKeyword(p)

	case TokenSub:
		self.Consume()
		node = &ASTNode{
			nodetype: NodeNegate,
			body: []*ASTNode{
				self.ExpectExpression(self.Peek()),
			},
		}

	case TokenLength:
		self.Consume()
		node = &ASTNode{
			nodetype: NodeLength,
			body: []*ASTNode{
				self.ExpectExpression(self.Peek()),
			},
		}

	case TokenNot:
		self.Consume()
		node = &ASTNode{
			nodetype: NodeNot,
			body: []*ASTNode{
				self.ExpectExpression(self.Peek()),
			},
		}

	case TokenVariadic:
		self.Consume()

		node = &ASTNode{
			nodetype: NodeVariadicResolve,
		}

	case TokenPeriod:
		node = &ASTNode{
			nodetype: NodeLiteral,
			values: map[string]ASTValue{
				"value": {
					tokens: []Token{
						self.Consume(),
						self.Consume(),
					},
				},
			},
		}

	default:
		self.err.Error(ErrorFatal, ParserErrorUnexpectedX, self.Consume().value)
		return nil
	}

	if node == nil {
		return &ASTNode{
			nodetype: NodeProgram,
		}
	}

	return self.HandleTrails(node)
}

func (self *Parser) HandleExpressionKeyword(t Token) *ASTNode {
	switch t.kwtype {
	case KeywordFunction:
		return self.HandleLambda(t)
	}

	return nil
}

func (self *Parser) HandleLocal(t Token) *ASTNode {
	pp := self.Peek()
	if (pp.token == TokenKeyword) && (pp.kwtype == KeywordFunction) {
		self.Consume()
		node := self.HandleFunction(self.Peek())

		node.nodetype = NodeLocalFunction

		self.ast.Add(node)

		return node
	}

	idents := ASTNode{
		nodetype: NodeVariableNameList,
		body:     []*ASTNode{},
	}

	for {
		p := self.Peek()

		if p.token != TokenIdent {
			self.err.Error(ErrorFatal, ParserErrorExpectedIdentifier, p.value)
			break
		}

		ident := self.QualifyIdent()
		idents.body = append(idents.body, &ident)

		n := self.Peek()
		if n.token == TokenComma {
			self.Consume()
			continue
		} else {
			break
		}
	}

	if self.err.ShouldImmediatelyStop() {
		return nil
	}

	next := self.Peek()

	if next.token != TokenEq {
		node := &ASTNode{
			nodetype: NodeLocalVariableStub,

			body: []*ASTNode{
				&idents,
			},
		}

		self.ast.Add(node)

		return node
	}

	self.Consume()

	node := &ASTNode{
		nodetype: NodeVariableValList,
		body:     []*ASTNode{},
	}

	for {
		p := self.Peek()

		var expr *ASTNode

		if p.token == TokenIdent {
			expr = self.HandleAssignmentIdent()
		} else {
			expr = self.ExpectExpression(p)
		}

		node.body = append(node.body, expr)

		n := self.Peek()
		if n.token == TokenComma {
			self.Consume()
			continue
		} else {
			break
		}
	}

	self.ast.Add(&ASTNode{
		nodetype: NodeLocalVariableAssignment,
		body: []*ASTNode{
			&idents,
			node,
		},
	})

	if self.err.ShouldImmediatelyStop() {
		return nil
	}

	return node
}

func (self *Parser) HandleTable(t Token) ASTNode {
	self.Consume()
	node := ASTNode{
		nodetype: NodeTable,
		body:     []*ASTNode{},
		values:   map[string]ASTValue{},
	}

	for {
		p := self.Peek()

		if p.token == TokenRCurl {
			break
		}

		switch p.token {
		case TokenLBrac:
			self.Consume()

			expr := self.ExpectExpression(self.Peek())

			self.ExpectToken(TokenRBrac)

			val := ASTNode{
				nodetype: NodeTableMapValue,
				body: []*ASTNode{
					expr,
				},
			}

			self.ExpectToken(TokenEq)

			var valexpr *ASTNode

			if self.Peek().token == TokenIdent {
				valexpr = self.HandleAssignmentIdent()
			} else {
				valexpr = self.ExpectExpression(self.Peek())
			}

			val.body = append(val.body, valexpr)
			node.body = append(node.body, &val)
		case TokenIdent:
			n := self.PeekSome(1)

			if n.token == TokenEq {
				self.Consume()
				self.Consume()

				val := ASTNode{
					nodetype: NodeTableMapValue,
					body: []*ASTNode{
						{
							nodetype: NodeIdentifier,
							body: []*ASTNode{
								{
									nodetype: NodeIdentSegNorm,
									values: map[string]ASTValue{
										"value": {
											token: p,
										},
									},
								},
							},
						},
					},
				}

				var valexpr *ASTNode

				if self.Peek().token == TokenIdent {
					valexpr = self.HandleAssignmentIdent()
				} else {
					valexpr = self.ExpectExpression(self.Peek())
				}

				val.body = append(val.body, valexpr)
				node.body = append(node.body, &val)
				break
			}

			fallthrough
		default:
			var expr *ASTNode

			if self.Peek().token == TokenIdent {
				expr = self.HandleAssignmentIdent()
			} else {
				expr = self.ExpectExpression(p)
			}

			val := ASTNode{
				nodetype: NodeTableArrayValue,
				body: []*ASTNode{
					expr,
				},
			}

			node.body = append(node.body, &val)
		}

		n := self.Peek()
		if n.token == TokenComma || n.token == TokenSemiColon {
			self.Consume()
			continue
		} else {
			break
		}
	}

	p := self.Peek()

	if p.token != TokenRCurl {
		self.err.Error(ErrorFatal, ParserErrorExpectedSymbol, "}", p.value)
	}

	self.Consume()

	return node
}

func (self *Parser) HandleIdentifier(t Token) *ASTNode {
	ident := self.QualifyIdent()
	node := self.HandleTrails(&ident)

	p := self.Peek()

	switch p.token {
	case TokenEq, TokenComma:
		return self.HandleAssignment(*node)
	}

	return node
}

func (self *Parser) HandleComment(t Token) {
	if self.ast.last != nil && self.ast.last.nodetype == NodeComment && t.newline > 1 {
		new := []Token{}

		new = append(new, self.ast.last.values["comments"].tokens...)
		new = append(new, t)

		*self.ast.last = ASTNode{
			nodetype: NodeComment,
			values: map[string]ASTValue{
				"comments": {
					tokens: new,
				},
			},
		}
		return
	}

	n := &ASTNode{
		nodetype: NodeComment,
		values: map[string]ASTValue{
			"comments": {
				tokens: []Token{
					t,
				},
			},
		},
	}

	self.ast.Add(n)
}

func (self *Parser) HandleCall() *ASTNode {
	nodes := []*ASTNode{}

	for {
		p := self.Peek()

		if p.token == TokenRParen {
			break
		}

		var expr *ASTNode

		if p.token == TokenIdent {
			id := self.QualifyIdent()
			expr = self.HandleTrails(&id)
		} else {
			expr = self.ExpectExpression(p)
		}

		nodes = append(nodes, expr)
		n := self.Peek()

		if n.token == TokenComma {
			self.Consume()
			continue
		} else {
			break
		}
	}

	p := self.Peek()
	if p.token != TokenRParen {
		self.err.Error(ErrorFatal, ParserErrorExpectedSymbol, ")", p.value)
	}

	self.Consume()

	return &ASTNode{
		nodetype: NodeCallArgs,
		body: nodes,
	}
}

func (self *Parser) HandleLParen(n ASTNode) {
	self.Consume()

	args := self.HandleCall()

	body := []*ASTNode{
		&n,
	}

	body = append(body, args)

	replacement := ASTNode{
		nodetype: NodeCall,
		body:     body,
	}

	*self.ast.last = replacement
}

func (self *Parser) HandleLBrac(n ASTNode) {
	self.Consume()

	expr := self.ExpectExpression(self.Peek())

	p := self.Peek()

	if p.token != TokenRBrac {
		self.err.Error(ErrorFatal, ParserErrorExpectedKeyword, "]", p.value)
	}

	replacement := ASTNode{
		nodetype: NodeMemberExpr,
		body: []*ASTNode{
			&n,
			expr,
		},
	}

	*self.ast.last = replacement
}

func (self *Parser) HandleAssignment(t ASTNode) *ASTNode {
	names := ASTNode{
		nodetype: NodeVariableNameList,
		body: []*ASTNode{
			&t,
		},
	}
	node := ASTNode{
		nodetype: NodeVariableValList,
		body:     []*ASTNode{},
	}

	for {
		sym := self.Peek()

		if sym.token == TokenComma {
			self.Consume()

			p := self.QualifyIdent()
			nn := self.HandleTrails(&p)

			if nn.nodetype != NodeIdentifier && nn.nodetype != NodeMemberExpr {
				self.err.Error(ErrorFatal, ParserErrorExpectedIdentifier, nn.nodetype)
				break
			}

			names.body = append(names.body, nn)

			continue
		} else {
			break
		}
	}

	p := self.Peek()

	if p.token != TokenEq {
		self.err.Error(ErrorFatal, ParserErrorExpectedSymbol, "=", p.value)
		return &node
	}

	self.Consume()

	for {
		if self.err.ShouldImmediatelyStop() {
			break
		}

		p := self.Peek()

		var expr *ASTNode

		if p.token == TokenIdent {
			expr = self.HandleAssignmentIdent()
		} else {
			expr = self.ExpectExpression(p)
		}

		node.body = append(node.body, expr)

		if self.Peek().token == TokenComma {
			self.Consume()
			continue
		}

		break
	}

	self.HandleBinExpr(&node)

	assign := &ASTNode{
		nodetype: NodeVariableAssignment,
		body: []*ASTNode{
			&names,
			&node,
		},
	}

	return assign
}

func (self *Parser) HandleLabel(t Token) *ASTNode {
	self.Consume()

	ident := self.ExpectToken(TokenIdent)

	if ident == nil {
		return nil
	}

	node := &ASTNode{
		nodetype: NodeLabel,
		values: map[string]ASTValue{
			"label": {
				token: *ident,
			},
		},
	}

	self.ExpectToken(TokenLabel)

	return node
}
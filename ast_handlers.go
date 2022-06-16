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
			Nodetype: NodeBool,
			Values: map[string]ASTValue{
				"value": {
					token: p,
				},
			},
		}
	case TokenNil:
		self.Consume()
		node = &ASTNode{
			Nodetype: NodeNil,
		}

	case TokenKeyword:
		node = self.HandleExpressionKeyword(p)

	case TokenSub:
		self.Consume()
		node = &ASTNode{
			Nodetype: NodeNegate,
			Body: []*ASTNode{
				self.ExpectExpression(self.Peek()),
			},
		}

	case TokenLength:
		self.Consume()
		node = &ASTNode{
			Nodetype: NodeLength,
			Body: []*ASTNode{
				self.ExpectExpression(self.Peek()),
			},
		}

	case TokenNot:
		self.Consume()
		node = &ASTNode{
			Nodetype: NodeNot,
			Body: []*ASTNode{
				self.ExpectExpression(self.Peek()),
			},
		}

	case TokenVariadic:
		self.Consume()

		node = &ASTNode{
			Nodetype: NodeVariadicResolve,
		}

	case TokenPeriod:
		node = &ASTNode{
			Nodetype: NodeLiteral,
			Values: map[string]ASTValue{
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
			Nodetype: NodeProgram,
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

		node.Nodetype = NodeLocalFunction

		self.ast.Add(node)

		return node
	}

	idents := ASTNode{
		Nodetype: NodeVariableNameList,
		Body:     []*ASTNode{},
	}

	for {
		p := self.Peek()

		if p.token != TokenIdent {
			self.err.Error(ErrorFatal, ParserErrorExpectedIdentifier, p.value)
			break
		}

		ident := self.QualifyIdent()
		idents.Body = append(idents.Body, &ident)

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
			Nodetype: NodeLocalVariableStub,

			Body: []*ASTNode{
				&idents,
			},
		}

		self.ast.Add(node)

		return node
	}

	self.Consume()

	node := &ASTNode{
		Nodetype: NodeVariableValList,
		Body:     []*ASTNode{},
	}

	for {
		p := self.Peek()

		var expr *ASTNode

		if p.token == TokenIdent {
			expr = self.HandleAssignmentIdent()
		} else {
			expr = self.ExpectExpression(p)
		}

		node.Body = append(node.Body, expr)

		n := self.Peek()
		if n.token == TokenComma {
			self.Consume()
			continue
		} else {
			break
		}
	}

	self.ast.Add(&ASTNode{
		Nodetype: NodeLocalVariableAssignment,
		Body: []*ASTNode{
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
		Nodetype: NodeTable,
		Body:     []*ASTNode{},
		Values:   map[string]ASTValue{},
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
				Nodetype: NodeTableMapValue,
				Body: []*ASTNode{
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

			val.Body = append(val.Body, valexpr)
			node.Body = append(node.Body, &val)
		case TokenIdent:
			n := self.PeekSome(1)

			if n.token == TokenEq {
				self.Consume()
				self.Consume()

				val := ASTNode{
					Nodetype: NodeTableMapValue,
					Body: []*ASTNode{
						{
							Nodetype: NodeIdentifier,
							Body: []*ASTNode{
								{
									Nodetype: NodeIdentSegNorm,
									Values: map[string]ASTValue{
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

				val.Body = append(val.Body, valexpr)
				node.Body = append(node.Body, &val)
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
				Nodetype: NodeTableArrayValue,
				Body: []*ASTNode{
					expr,
				},
			}

			node.Body = append(node.Body, &val)
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
	if self.ast.LastNode != nil && self.ast.LastNode.Nodetype == NodeComment && t.newline > 1 {
		new := []Token{}

		new = append(new, self.ast.LastNode.Values["comments"].tokens...)
		new = append(new, t)

		*self.ast.LastNode = ASTNode{
			Nodetype: NodeComment,
			Values: map[string]ASTValue{
				"comments": {
					tokens: new,
				},
			},
		}
		return
	}

	n := &ASTNode{
		Nodetype: NodeComment,
		Values: map[string]ASTValue{
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
		Nodetype: NodeCallArgs,
		Body: nodes,
	}
}

func (self *Parser) HandleLParen(n ASTNode) {
	self.Consume()

	args := self.HandleCall()

	Body := []*ASTNode{
		&n,
	}

	Body = append(Body, args)

	replacement := ASTNode{
		Nodetype: NodeCall,
		Body:     Body,
	}

	*self.ast.LastNode = replacement
}

func (self *Parser) HandleLBrac(n ASTNode) {
	self.Consume()

	expr := self.ExpectExpression(self.Peek())

	p := self.Peek()

	if p.token != TokenRBrac {
		self.err.Error(ErrorFatal, ParserErrorExpectedKeyword, "]", p.value)
	}

	replacement := ASTNode{
		Nodetype: NodeMemberExpr,
		Body: []*ASTNode{
			&n,
			expr,
		},
	}

	*self.ast.LastNode = replacement
}

func (self *Parser) HandleAssignment(t ASTNode) *ASTNode {
	names := ASTNode{
		Nodetype: NodeVariableNameList,
		Body: []*ASTNode{
			&t,
		},
	}
	node := ASTNode{
		Nodetype: NodeVariableValList,
		Body:     []*ASTNode{},
	}

	for {
		sym := self.Peek()

		if sym.token == TokenComma {
			self.Consume()

			p := self.QualifyIdent()
			nn := self.HandleTrails(&p)

			if nn.Nodetype != NodeIdentifier && nn.Nodetype != NodeMemberExpr {
				self.err.Error(ErrorFatal, ParserErrorExpectedIdentifier, nn.Nodetype)
				break
			}

			names.Body = append(names.Body, nn)

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

		node.Body = append(node.Body, expr)

		if self.Peek().token == TokenComma {
			self.Consume()
			continue
		}

		break
	}

	self.HandleBinExpr(&node)

	assign := &ASTNode{
		Nodetype: NodeVariableAssignment,
		Body: []*ASTNode{
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
		Nodetype: NodeLabel,
		Values: map[string]ASTValue{
			"label": {
				token: *ident,
			},
		},
	}

	self.ExpectToken(TokenLabel)

	return node
}
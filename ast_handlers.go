package gooa

func (self *Parser) ExpectExpression(p Token) *ASTNode {
	expr := self.GetExpression()

	if expr == nil {
		self.err.Error(ErrorFatal, ParserErrorExpectedExpression, p.Value)
	}

	return expr
}

func (self *Parser) GetExpression() *ASTNode {
	for {
		semi := self.Peek()

		if semi.Token == TokenSemiColon {
			self.Consume()
			continue
		} else {
			break
		}
	}

	p := self.Peek()

	if IsLiteralType(p.Token) {
		bin := self.HandleBinExprLiteral(p)

		return bin
	}

	var node *ASTNode

	switch p.Token {
	case TokenIdent:
		node = self.HandleAssignmentIdent()
	case TokenLCurl:
		fptrs := self.HandleTable(p)
		node = &fptrs
	case TokenLParen:
		self.Consume()
		expr := self.ExpectExpression(self.Peek())

		rpar := self.Peek()
		if rpar.Token != TokenRParen {
			self.err.Error(ErrorFatal, ParserErrorExpectedSymbol, ")", rpar.Value)
		}

		self.Consume()

		node = expr

	case TokenTrue, TokenFalse:
		self.Consume()
		node = &ASTNode{
			Nodetype: NodeBool,
			Values: map[string]ASTValue{
				"value": {
					Token: p,
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
					Tokens: []Token{
						self.Consume(),
						self.Consume(),
					},
				},
			},
		}

	default:
		self.err.Error(ErrorFatal, ParserErrorUnexpectedX, self.Consume().Value)
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
	switch t.Keyword {
	case KeywordFunction:
		return self.HandleLambda(t)
	}

	return nil
}

func (self *Parser) HandleLocal(t Token) *ASTNode {
	pp := self.Peek()
	if (pp.Token == TokenKeyword) && (pp.Keyword == KeywordFunction) {
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

		if p.Token != TokenIdent {
			self.err.Error(ErrorFatal, ParserErrorExpectedIdentifier, p.Value)
			break
		}

		ident := self.QualifyIdent()
		idents.Body = append(idents.Body, &ident)

		n := self.Peek()
		if n.Token == TokenComma {
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

	if next.Token != TokenEq {
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

		if p.Token == TokenIdent {
			expr = self.HandleAssignmentIdent()
		} else {
			expr = self.ExpectExpression(p)
		}

		node.Body = append(node.Body, expr)

		n := self.Peek()
		if n.Token == TokenComma {
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

		if p.Token == TokenRCurl {
			break
		}

		switch p.Token {
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

			if self.Peek().Token == TokenIdent {
				valexpr = self.HandleAssignmentIdent()
			} else {
				valexpr = self.ExpectExpression(self.Peek())
			}

			val.Body = append(val.Body, valexpr)
			node.Body = append(node.Body, &val)
		case TokenIdent:
			n := self.PeekSome(1)

			if n.Token == TokenEq {
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
											Token: p,
										},
									},
								},
							},
						},
					},
				}

				var valexpr *ASTNode

				if self.Peek().Token == TokenIdent {
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

			if self.Peek().Token == TokenIdent {
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
		if n.Token == TokenComma || n.Token == TokenSemiColon {
			self.Consume()
			continue
		} else {
			break
		}
	}

	p := self.Peek()

	if p.Token != TokenRCurl {
		self.err.Error(ErrorFatal, ParserErrorExpectedSymbol, "}", p.Value)
	}

	self.Consume()

	return node
}

func (self *Parser) HandleIdentifier(t Token) *ASTNode {
	ident := self.QualifyIdent()
	node := self.HandleTrails(&ident)

	p := self.Peek()

	switch p.Token {
	case TokenEq, TokenComma:
		return self.HandleAssignment(*node)
	}

	return node
}

func (self *Parser) HandleComment(t Token) {
	if self.ast.LastNode != nil && self.ast.LastNode.Nodetype == NodeComment && t.Newlines > 1 {
		new := []Token{}

		new = append(new, self.ast.LastNode.Values["comments"].Tokens...)
		new = append(new, t)

		*self.ast.LastNode = ASTNode{
			Nodetype: NodeComment,
			Values: map[string]ASTValue{
				"comments": {
					Tokens: new,
				},
			},
		}
		return
	}

	n := &ASTNode{
		Nodetype: NodeComment,
		Values: map[string]ASTValue{
			"comments": {
				Tokens: []Token{
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

		if p.Token == TokenRParen {
			break
		}

		var expr *ASTNode

		if p.Token == TokenIdent {
			id := self.QualifyIdent()
			expr = self.HandleTrails(&id)
		} else {
			expr = self.ExpectExpression(p)
		}

		nodes = append(nodes, expr)
		n := self.Peek()

		if n.Token == TokenComma {
			self.Consume()
			continue
		} else {
			break
		}
	}

	p := self.Peek()
	if p.Token != TokenRParen {
		self.err.Error(ErrorFatal, ParserErrorExpectedSymbol, ")", p.Value)
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

	if p.Token != TokenRBrac {
		self.err.Error(ErrorFatal, ParserErrorExpectedKeyword, "]", p.Value)
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

		if sym.Token == TokenComma {
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

	if p.Token != TokenEq {
		self.err.Error(ErrorFatal, ParserErrorExpectedSymbol, "=", p.Value)
		return &node
	}

	self.Consume()

	for {
		if self.err.ShouldImmediatelyStop() {
			break
		}

		p := self.Peek()

		var expr *ASTNode

		if p.Token == TokenIdent {
			expr = self.HandleAssignmentIdent()
		} else {
			expr = self.ExpectExpression(p)
		}

		node.Body = append(node.Body, expr)

		if self.Peek().Token == TokenComma {
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
				Token: *ident,
			},
		},
	}

	self.ExpectToken(TokenLabel)

	return node
}
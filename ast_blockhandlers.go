
package gooa

func (self *Parser) ExpectKeyword(kw KeywordType) bool {
	next := self.Consume()

	if next.Keyword != kw {
		self.err.Error(ErrorFatal, ParserErrorExpectedKeyword, keywordTypeValues[kw], next.Value)
		return true
	}

	return false
}

func (self *Parser) ExpectToken(tty TokenType) *Token {
	next := self.Consume()

	if next.Token != tty {
		self.err.Error(ErrorFatal, ParserErrorExpectedToken, tokenNames[tty], next.Value)
		return &InvalidToken
	}

	return &next
}

func (self *Parser) HandleUntil(kw KeywordType, msg string) {
	start := self.Peek().Position.Copy()

	for {
		br, last := self.HandleParsing()
		
		if last == kw {
			break
		}

		if br {
			self.err.Error(ErrorFatal, ParserErrorMissingEnd, msg + "->(" + start.Fancy() + ")")
			break
		}
	}
}

func (self *Parser) HandleUntilEnd(msg string) {
	self.HandleUntil(KeywordEnd, msg)
}

func (self *Parser) HandleArgList() *ASTNode {
	lpar := self.Peek()

	if lpar.Token != TokenLParen {
		return &ASTNode{
			Nodetype: NodeArgumentListOmitted,
		}
	}

	self.Consume()

	args := []*ASTNode{}

	for {
		name := self.Peek()

		if name.Token == TokenVariadic {
			self.Consume()
			args = append(args, &ASTNode{
				Nodetype: NodeArgumentVariadic,
			})

			break
		} else if name.Token == TokenRParen {
			break
		} else if name.Token != TokenIdent {
			self.err.Error(ErrorFatal, ParserErrorExpectedArgumentName, name.Value)
			break
		}

		self.Consume()

		p := self.Peek()
		if p.Token == TokenComma || p.Token == TokenRParen {
			args = append(args, &ASTNode{
				Nodetype: NodeArgumentNormal,
				Values: map[string]ASTValue{
					"name": {
						Token: name,
					},
				},
			})

			if p.Token == TokenRParen {
				break
			}

			self.Consume()
			continue
		} else if (p.Token != TokenEq) {
			self.err.Error(ErrorFatal, ParserErrorExpectedSymbol, ",", p.Value)
			break
		}

		self.Consume()
		
		next := self.Peek() 

		if next.Token == TokenVariadic {
			args = append(args, &ASTNode{
				Nodetype: NodeNamedArgumentVariadic,
				Values: map[string]ASTValue{
					"name": {
						Token: name,
					},
					"variadic": {
						Token: self.Consume(),
					},
				},
			})
			break
		}

		expr := self.ExpectExpression(self.Peek())
		if expr == nil || self.err.ShouldImmediatelyStop() {
			break
		}
		
		args = append(args, &ASTNode{
			Nodetype: NodeNamedArgumentDef,
			Body: []*ASTNode{
				expr,
			},
			Values: map[string]ASTValue{
				"name": {
					Token: name,
				},
			},
		})

		p = self.Peek()
	
		if p.Token == TokenComma {
			self.Consume()
			continue
		} else {
			break
		}
	}

	if self.Peek().Token != TokenRParen {
		self.err.Error(ErrorFatal, ParserErrorExpectedSymbol, ")", self.Peek().Value)
	}

	self.Consume()

	return &ASTNode{
		Nodetype: NodeArgumentList,
		Body: args,
	}
}

func (self *Parser) HandleLambda(t Token) *ASTNode {
	self.Consume()
	
	args := self.HandleArgList()

	if self.err.ShouldImmediatelyStop() {
		return nil
	}

	node := &ASTNode{
		Nodetype: NodeAnonymousFunction,
		Body: []*ASTNode{
			args,
		},
	}

	self.ast.OpenScope(node)

	self.HandleUntilEnd("lambda")

	return self.ast.CloseScope()
}

func (self *Parser) HandleFunction(t Token) *ASTNode {
	ident := self.QualifyIdent()
	arglist := self.HandleArgList()

	node := &ASTNode{
		Nodetype: NodeFunction,
		Body: []*ASTNode{
			&ident,
			arglist,
		},
	}

	self.ast.OpenScope(node)

	self.HandleUntilEnd("function")

	return self.ast.CloseScope()
}

func (self *Parser) HandleReturn(t Token) *ASTNode {
	Body := []*ASTNode{}

	for {
		expr := self.GetExpression()

		if expr == nil {
			break
		}

		Body = append(Body, expr)

		if self.Peek().Token == TokenComma {
			self.Consume()
			continue
		}

		break
	}

	return &ASTNode{
		Nodetype: NodeReturn,
		Body: Body,
	}
}

func (self *Parser) HandleDo(t Token) *ASTNode {
	node := &ASTNode{
		Nodetype: NodeArbitraryScope,
		Body: []*ASTNode{},
	}

	self.ast.OpenScope(node)

	self.HandleUntilEnd("do")

	return self.ast.CloseScope()
}

func (self *Parser) HandleIf(t Token) *ASTNode {
	expr := self.ExpectExpression(self.Peek())

	ifnode := &ASTNode{
		Nodetype: NodeIf,
		Body: []*ASTNode{
			expr,
		},
	}

	if self.ExpectKeyword(KeywordThen) {
		return ifnode
	}

	haselse := false
	scopenode := &ASTNode{
		Nodetype: NodeIfScope,
		Body: []*ASTNode{},
	}

	self.ast.OpenScope(scopenode)

	for {
		br, last := self.HandleParsing()
		
		if last == KeywordEnd || last == KeywordElse || last == KeywordElseif {
			ifnode.Body = append(ifnode.Body, self.ast.CloseScope())
		}

		if last == KeywordEnd {
			break
		}

		if br {
			self.err.Error(ErrorFatal, ParserErrorMissingEnd, "if")
			break
		}

		if last == KeywordElse {
			if haselse {
				self.err.Error(ErrorFatal, ParserErrorElseAlreadyDeclared)
			}
			haselse = true

			elsenode := &ASTNode{
				Nodetype: NodeElseScope,
				Body: []*ASTNode{},
			}

			self.ast.OpenScope(elsenode)
		} else if last == KeywordElseif {
			expr := self.ExpectExpression(self.Peek())

			if self.ExpectKeyword(KeywordThen) {
				return ifnode
			}

			elseifnode := &ASTNode{
				Nodetype: NodeElseIfScope,
				Body: []*ASTNode{
					expr,
				},
			}
			
			self.ast.OpenScope(elseifnode)
		}
	}

	return ifnode
}

func (self *Parser) HandleFor(t Token) *ASTNode {
	next := self.Peek()

	if next.Token != TokenIdent {
		self.err.Error(ErrorFatal, ParserErrorExpectedIdentifier, next.Value)
		return nil
	}

	nextsome := self.PeekSome(1)

	if nextsome.Token == TokenEq {
		return self.HandleForI(t)
	}

	node := &ASTNode{
		Nodetype: NodeForIterator,
	}

	args := &ASTNode{
		Nodetype: NodeForIteratorArgs,
		Body: []*ASTNode{},
	}

	for {
		ident := self.Peek()
		
		if ident.Token == TokenIdent {
			args.Body = append(args.Body, &ASTNode{
				Nodetype: NodeIdentSegNorm,
				Values: map[string]ASTValue{
					"value": {
						Token: self.Consume(),
					},
				},
			})

			p := self.Peek()

			if p.Token == TokenComma {
				self.Consume()
				continue
			}

			break
		}
	}

	node.Body = []*ASTNode{
		args,
	}

	if self.ExpectKeyword(KeywordIn) {
		self.err.Error(ErrorFatal, ParserErrorExpectedKeyword, "in", self.Peek().Value)
		return node
	}

	expr := self.ExpectExpression(self.Peek())
	
	node.Body = append(node.Body, expr)

	if self.ExpectKeyword(KeywordDo) {
		self.err.Error(ErrorFatal, ParserErrorExpectedKeyword, "do", self.Peek().Value)
		return node
	}
	
	self.ast.OpenScope(node)

	self.HandleUntilEnd("for")

	return self.ast.CloseScope()
}

func (self *Parser) HandleForI(t Token) *ASTNode {
	ident := self.Consume()

	self.Consume() // already validated eq is there

	identexpr := self.ExpectExpression(self.Peek())

	comma := self.Consume()

	if comma.Token != TokenComma {
		self.err.Error(ErrorFatal, ParserErrorExpectedSymbol, ",", comma.Value)
		return nil
	}
	

	var increxpr *ASTNode
	toexpr := self.ExpectExpression(self.Peek())
	comma = self.Peek()

	if comma.Token == TokenComma {
		self.Consume()
		increxpr = &ASTNode{
			Nodetype: NodeForIIncr,
			Body: []*ASTNode{
				self.ExpectExpression(self.Peek()),
			},
		}
	}

	self.ExpectKeyword(KeywordDo)

	node := &ASTNode{
		Nodetype: NodeForI,
		Body: []*ASTNode{
			identexpr,
			toexpr,
		},

		Values: map[string]ASTValue{
			"identifier":{
				Token: ident,
			},
		},
	}

	if increxpr != nil {
		node.Body = append(node.Body, increxpr)
	}

	self.ast.OpenScope(node)

	self.HandleUntilEnd("fori")

	return self.ast.CloseScope()
}

func (self *Parser) HandleRepeat(t Token) *ASTNode {
	node := &ASTNode{
		Nodetype: NodeRepeat,
	}

	scope := &ASTNode{
		Nodetype: NodeArbitraryScope,
		Body: []*ASTNode{},
	}

	self.ast.OpenScope(scope)

	self.HandleUntil(KeywordUntil, "repeat")

	self.ast.CloseScope()

	expr := self.ExpectExpression(self.Peek())

	node.Body = []*ASTNode{
		expr,
		scope,
	}

	return node
}

func (self *Parser) HandleWhile(t Token) *ASTNode {
	node := &ASTNode{
		Nodetype: NodeWhile,
		Body: []*ASTNode{
			self.ExpectExpression(self.Peek()),
		},
	}

	if self.ExpectKeyword(KeywordDo) {
		return node
	}

	self.ast.OpenScope(node)

	self.HandleUntilEnd("while")

	return self.ast.CloseScope()
}

func (self *Parser) HandleGoto(t Token) *ASTNode {
	ident := self.ExpectToken(TokenIdent)

	if ident == nil {
		return nil
	}
	
	node := &ASTNode{
		Nodetype: NodeGoto,
		Values: map[string]ASTValue{
			"label": {
				Token: *ident,
			},
		},
	}

	return node
}
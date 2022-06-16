
package gooa

func (self *Parser) ExpectKeyword(kw KeywordType) bool {
	next := self.Consume()

	if next.kwtype != kw {
		self.err.Error(ErrorFatal, ParserErrorExpectedKeyword, keywordTypeValues[kw], next.value)
		return true
	}

	return false
}

func (self *Parser) ExpectToken(tty TokenType) *Token {
	next := self.Consume()

	if next.token != tty {
		self.err.Error(ErrorFatal, ParserErrorExpectedToken, tokenNames[tty], next.value)
		return &InvalidToken
	}

	return &next
}

func (self *Parser) HandleUntil(kw KeywordType, msg string) {
	start := self.Peek().position.Copy()

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

	if lpar.token != TokenLParen {
		return &ASTNode{
			Nodetype: NodeArgumentListOmitted,
		}
	}

	self.Consume()

	args := []*ASTNode{}

	for {
		name := self.Peek()

		if name.token == TokenVariadic {
			self.Consume()
			args = append(args, &ASTNode{
				Nodetype: NodeArgumentVariadic,
			})

			break
		} else if name.token == TokenRParen {
			break
		} else if name.token != TokenIdent {
			self.err.Error(ErrorFatal, ParserErrorExpectedArgumentName, name.value)
			break
		}

		self.Consume()

		p := self.Peek()
		if p.token == TokenComma || p.token == TokenRParen {
			args = append(args, &ASTNode{
				Nodetype: NodeArgumentNormal,
				Values: map[string]ASTValue{
					"name": {
						token: name,
					},
				},
			})

			if p.token == TokenRParen {
				break
			}

			self.Consume()
			continue
		} else if (p.token != TokenEq) {
			self.err.Error(ErrorFatal, ParserErrorExpectedSymbol, ",", p.value)
			break
		}

		self.Consume()
		
		next := self.Peek() 

		if next.token == TokenVariadic {
			args = append(args, &ASTNode{
				Nodetype: NodeNamedArgumentVariadic,
				Values: map[string]ASTValue{
					"name": {
						token: name,
					},
					"variadic": {
						token: self.Consume(),
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
					token: name,
				},
			},
		})

		p = self.Peek()
	
		if p.token == TokenComma {
			self.Consume()
			continue
		} else {
			break
		}
	}

	if self.Peek().token != TokenRParen {
		self.err.Error(ErrorFatal, ParserErrorExpectedSymbol, ")", self.Peek().value)
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

		if self.Peek().token == TokenComma {
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

	if next.token != TokenIdent {
		self.err.Error(ErrorFatal, ParserErrorExpectedIdentifier, next.value)
		return nil
	}

	nextsome := self.PeekSome(1)

	if nextsome.token == TokenEq {
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
		
		if ident.token == TokenIdent {
			args.Body = append(args.Body, &ASTNode{
				Nodetype: NodeIdentSegNorm,
				Values: map[string]ASTValue{
					"value": {
						token: self.Consume(),
					},
				},
			})

			p := self.Peek()

			if p.token == TokenComma {
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
		self.err.Error(ErrorFatal, ParserErrorExpectedKeyword, "in", self.Peek().value)
		return node
	}

	expr := self.ExpectExpression(self.Peek())
	
	node.Body = append(node.Body, expr)

	if self.ExpectKeyword(KeywordDo) {
		self.err.Error(ErrorFatal, ParserErrorExpectedKeyword, "do", self.Peek().value)
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

	if comma.token != TokenComma {
		self.err.Error(ErrorFatal, ParserErrorExpectedSymbol, ",", comma.value)
		return nil
	}
	

	var increxpr *ASTNode
	toexpr := self.ExpectExpression(self.Peek())
	comma = self.Peek()

	if comma.token == TokenComma {
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
				token: ident,
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
				token: *ident,
			},
		},
	}

	return node
}
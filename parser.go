package gooa

var InvalidToken = Token{
	Token: TokenEOF,
	Value: "<EOF>",
	Invalid: true,
}

// Parser Construct
type Parser struct {
	tokens 		[]Token
	ast 		*AST

	index 		int
	err 		ErrorHandler
	unhandled 	[]Token

	middleware 	*MiddlewareHandler
}

// Reset Parse State
func (self *Parser) Reset() {
	self.tokens = nil
	self.index 	= 0
	self.unhandled = []Token{}
}

func (self *Parser) Parse(toks []Token, err *ErrorHandler) (AST, bool) {
	self.Reset()
	self.err = *err
	self.err.SetErrorRealm(ErrorRealmParser)
	self.ast = CreateAST(self.err)

	self.tokens = toks

	for {
		valid, _ := self.HandleParsing()
		if valid {
			break
		}
	}

	self.err.Dump()

	return *self.middleware.PostParse(self.ast), self.err.ShouldStop()
}

// Ignore the given token
// push to the unhandled slice
func (self *Parser) Ignore(t Token) {
	self.unhandled = append(self.unhandled, t)
}

// Peek at current token
func (self *Parser) Peek() Token {
	if self.Stop(self.index) {
		return InvalidToken
	}

	tok := self.tokens[self.index]

	if tok.Token == TokenComment {
		self.Consume()
		return self.Peek()
	}

	self.err.SetPosition(&tok.Position)

	return tok
}

// Peek at arbitrary token
func (self *Parser) PeekSome(amt int) Token {
	if self.Stop(self.index + amt) {
		return InvalidToken
	}

	tok := self.tokens[self.index + amt]

	self.err.SetPosition(&tok.Position)

	return tok
}

// Should the parser stop
func (self *Parser) Stop(index int) bool {
	return index >= len(self.tokens)
}

// Consume Token
func (self *Parser) Consume() Token {
	if self.Stop(self.index) {
		return InvalidToken
	}

	b := self.tokens[self.index]
	self.index++

	self.err.SetPosition(&b.Position)

	return b
}

// Pop last block expression
func (self *Parser) EndPop() {
	if self.ast.CurNode == self.ast.RootNode {
		self.err.Error(ErrorFatal, ParserErrorUnexpectedEnd)
		return
	}

	// self.ast.Pop()
	// return self.ast.PopValue()
}

// Handle Keyword Tokens
func (self *Parser) HandleKeyword(t Token) (bool, *ASTNode) {
	switch t.Keyword {
	case KeywordLocal:
		self.HandleLocal(t)

	case KeywordEnd:
		self.EndPop()
		return true, nil
	case KeywordFunction:
		fn := self.HandleFunction(t)
		return false, fn
	
	case KeywordReturn:
		return false, self.HandleReturn(t)

	case KeywordDo:
		return false, self.HandleDo(t)
	
	case KeywordIf:
		return false, self.HandleIf(t)
	
	case KeywordFor:
		return false, self.HandleFor(t)

	case KeywordRepeat:
		return false, self.HandleRepeat(t)

	case KeywordBreak:
		return false, &ASTNode{
			Nodetype: NodeBreak,
		}

	case KeywordContinue:
		return false, &ASTNode{
			Nodetype: NodeContinue,
		}

	case KeywordWhile:
		return false, self.HandleWhile(t)

	case KeywordGoto:
		return false, self.HandleGoto(t)

	// Passthroughs for controlled scopes
	case KeywordElse, KeywordElseif, KeywordUntil: return false, nil

	default:
		self.err.Error(ErrorFatal, ParserErrorUnexpectedKeyword, keywordTypeValues[t.Keyword])
	}

	return false, nil
}

func (self *Parser) HandleParsing() (bool, KeywordType) {
	if self.err.ShouldImmediatelyStop() {
		return true, KeywordEmpty
	}
	
	if self.Stop(self.index) {
		return true, KeywordEmpty
	}
	
	curtok := self.PeekSome(0)	
	
	if self.middleware.ParserHandleToken(self, curtok) {
		return false, KeywordEmpty
	}

	if curtok.Token == TokenEOF {
		return true, KeywordEmpty
	}

	if IsLiteralType(curtok.Token) {
		self.ast.Add(self.HandleBinExprLiteral(self.Consume()))
		return false, KeywordEmpty
	}

	switch curtok.Token {
	case TokenKeyword:
		_, val := self.HandleKeyword(self.Consume())
		
		if val != nil {
			self.ast.Add(val)
		}
		
		return false, curtok.Keyword
	case TokenIdent:
		n := self.HandleIdentifier(curtok)

		if n != nil {
			self.ast.Add(n)
		}

	case TokenLabel:
		n := self.HandleLabel(curtok)

		if n != nil {
			self.ast.Add(n)
		}

	case TokenLParen:
		self.ast.Add(self.ExpectExpression(self.Peek()))

	case TokenSemiColon:
		self.Consume()

	case TokenComment:
		self.HandleComment(self.Consume())

	default:
		self.err.Error(ErrorFatal, ParserErrorUnexpectedX, tokenNames[curtok.Token])
	}

	return false, KeywordEmpty
}
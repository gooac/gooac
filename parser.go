package gooa

import "fmt"

type Parser struct {
	parsing []Token
	parsed  []ParserExpression

	index   int
	islocal bool
}

func (self *Parser) Parse(toks []Token) {
	self.Reset()
	defer self.Reset()

	self.parsing = toks
	self.parsed = []ParserExpression{}

	for {
		tok, err := self.Peek(1)

		if err {
			break
		}

		switch tok.toktype {
		case TokenKeyword:
			self.HandleKeyword(tok)
		case TokenIdent:
			print("Handling ident")
			tok.Print()
			self.HandleIdent(tok)
		default:
			print("Ignored Token: ")
			tok.Print()
		}

		self.Consume(1)
	}
}

func (self *Parser) Reset() {
	self.parsing = nil
	self.parsed = nil
	self.index = 0
}

func (self *Parser) Peek(amt int) (*Token, bool) {
	if self.Stop() {
		return nil, true
	}

	return &self.parsing[self.index], false
}

func (self *Parser) Stop() bool {
	return self.index >= len(self.parsing)
}

func (self *Parser) Consume(amt int) {
	self.index += amt
}

func (self *Parser) Error(message string) {
	panic("[Parser] " + message)
}

func (self *Parser) TestNextType(ty TokenType) (bool, *Token) {
	next, err := self.Peek(1)

	if err {
		return false, nil
	}

	fmt.Printf("%v == %v\n", next.toktype, ty)
	if next.toktype == ty {
		return true, next
	}

	return false, nil
}

func (self *Parser) TestNextValue(val string) (bool, *Token) {
	next, err := self.Peek(1)

	if err {
		return false, nil
	}

	if next.value == val {
		return true, next
	}

	return false, nil
}

func (self *Parser) TestNextKeyword(ty KeywordType) (bool, *Token) {
	next, err := self.Peek(1)

	if err {
		return false, nil
	}

	if next.toktype != TokenKeyword {
		return false, nil
	}

	if next.keywordtype == ty {
		return true, next
	}

	return false, nil
}

func (self *Parser) TestError(is bool, msg string) bool {
	if !is {
		self.Error(msg)
	}

	return !is
}
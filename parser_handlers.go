package gooa

import "fmt"

func (self *Parser) HandleKeyword(token *Token) {
	switch token.keywordtype {
	case KeywordLocal:
		self.islocal = true
	}
}

func (self *Parser) HandleIdent(token *Token) {
	id := ParserIdentifier{}

	id.tokens = append(id.tokens, token)

	self.Consume(1)

	for {
		isprd, tokprd := self.TestNextType(TokenPeriod)
		// iscol, tokcol := self.TestNextType(TokenColon)

		if isprd {
			id.tokens = append(id.tokens, tokprd)

			isid, identtok := self.TestNextType(TokenIdent)

			if !isid {
				fmt.Printf("My balls itch so insanely badly: %v, %v", tokprd.value, identtok.value)

				return 
			}

			id.tokens = append(id.tokens, identtok)
		} else {
			print("Breaking: ")
			if tokprd != nil {
				tokprd.Print()
			} else {
				println("Null Token")
			}
			break
		}
	}

	for k, v := range id.tokens {
		fmt.Printf("Token %v: %v\n", k, v)
	}
}

func (self *Parser) HandleFunctionDef(islocal bool) {}
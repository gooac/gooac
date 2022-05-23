package gooa

import "fmt"

// import "fmt"

type TokenizerPosition struct {
	line		int
	column		int
}

type Tokenizer struct {
	tokenizing 	string
	tokens 		[]Token
	index 		int

	pos 		TokenizerPosition
}

func (t *Tokenizer) Reset(to string) {
	t.tokenizing 	= to
	t.tokens 		= nil
	t.index 		= 0
	t.pos 			= TokenizerPosition{1, 0}
}

func (self *Tokenizer) Tokenize(s string) []Token {
	self.Reset(s)
	defer self.Reset("")

	var err error;
	for {
		if self.Stop() {break}
		if err != nil {break}

		ch := self.tokenizing[self.index]

		switch ch {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9': 		
						err = self.HandleNumber		(ch)
		case '.': 		err = self.HandleVariadics	(ch)
		case '>':		err = self.HandleNextPredef	(ch, '=', TokenGt, TokenGtEq)
		case '<': 		err = self.HandleNextPredef	(ch, '=', TokenLt, TokenLtEq)
		case '!', '~': 	err = self.HandleNextPredef	(ch, '=', TokenNot, TokenNotEq)
		case ':': 		err = self.HandleNextPredef	(ch, ':', TokenColon, TokenLabel)
		case '"', '\'':	err = self.HandleString(ch) 
		case '[':		err = self.HandleMultilineStrings()
		case '-': 		err = self.HandleComment(ch, '-')
		default:
			val, valid := symbolLookups[string(ch)]

			if valid {
				err = self.HandleSymbol(ch, val)
			} else {
				err = self.HandleIdent(ch)
			}
		}

		self.Consume()
	}

	
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	return self.tokens
}

func (self *Tokenizer) Consume() {
	if self.Stop() {return}
	self.index += 1
}

func (self *Tokenizer) Peek(next int) (byte, bool) {
	if (self.index + next) >= len(self.tokenizing) {
		return '_', true
	}

	return self.tokenizing[self.index + next], false
}

func (self Tokenizer) Stop() bool {
	return self.index >= len(self.tokenizing)	
}
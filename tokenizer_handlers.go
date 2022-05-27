package gooa

import (
	"errors"
	"unicode"
)

func (self *Tokenizer) IsWhitespace(ch byte) bool {
	r := rune(ch)

	return unicode.IsSpace(r) || r == ';'
}

func (self *Tokenizer) HandleNumber(s byte) error {
	num := string(s)
	tobreak := false
	isfloat := false
	ishex	:= false
	err 	:= false
	toktyp 	:= TokenNumber

	for {
		if tobreak {
			break
		}

		ch, err := self.Peek(1)

		if err {
			break
		}

		switch ch {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			num = num + string(ch)
			self.Consume()
		case '.':
			if isfloat || ishex {
				err = true
			}
			num = num + string(ch)

			isfloat = true

			self.Consume()
		case 'x':
			if ishex || isfloat || (num != "0") {
				err = true
			}
			
			ishex = true
			toktyp = TokenHexNumber
		
			num = num + string(ch)

			self.Consume()

		case 'a', 'b', 'c', 'd', 'e', 'f', 
		'A', 'B', 'C', 'D', 'E', 'F':
			if !ishex || isfloat {
				err = true
			}

			num = num + string(ch)
			self.Consume()

		default:
			tobreak = true
		}
	}

	nt := Token{
		toktype: toktyp,
		value:   num,

		position: self.pos,
	}

	self.tokens = append(self.tokens, nt)

	if err {
		return errors.New("Malformed number '" + num + "'")
	}
	
	return nil
}

func (self *Tokenizer) HandleVariadics(ch byte) error {
	next, err := self.Peek(1)
	next2, err2 := self.Peek(2)

	nt := Token{
		toktype: TokenPeriod,

		value:    string(ch),
		position: self.pos,
	}

	if next == '.' && !err {
		self.Consume()
		nt.toktype = TokenConcat
		nt.value += "."

		if next2 == '.' && !err2 {
			self.Consume()
			nt.toktype = TokenVariadic
			nt.value += "."
		}
	}

	self.tokens = append(self.tokens, nt)

	return nil
}

func (self *Tokenizer) HandleNextPredef(ch byte, next byte, toktype TokenType, elsetype TokenType) error {
	peek, err := self.Peek(1)

	nt := Token{
		toktype: toktype,
		value:   string(ch),

		position: self.pos,
	}

	if (peek == next) && !err {
		nt.toktype = elsetype
		self.Consume()
		nt.value += string(next)
	}

	self.tokens = append(self.tokens, nt)

	return nil
}

func (self *Tokenizer) HandleSymbol(ch byte, tok TokenType) error {
	nt := Token{
		toktype: tok,
		value:   string(ch),

		position: self.pos,
	}

	self.tokens = append(self.tokens, nt)

	return nil
}

func (self *Tokenizer) HandleIdent(ch byte) error {
	id := string(ch)

	if self.IsWhitespace(ch) {
		return nil
	}

	for {
		char, err := self.Peek(1)

		if err {
			break
		}

		r := rune(char)
		if (unicode.IsNumber(r) || unicode.IsLetter(r)) && !self.IsWhitespace(char) {
			id += string(r)
			self.Consume()
		} else {
			break
		}
	}

	_, iskw := keywordLookup[id]

	nt := Token{
		toktype: TokenIdent,
		value:   id,

		position: self.pos,
	}
	
	if iskw {
		nt.toktype = TokenKeyword
		nt.keywordtype = keywordTypes[id]
	}

	self.tokens = append(self.tokens, nt)

	return nil
}

func (self *Tokenizer) HandleString(ch byte) error {
	str := string("")

	for {
		char, err := self.Peek(1)

		if err {
			return errors.New("Unfinished String")
		}

		if char == ch {
			self.Consume()
			break
		} else {
			str += string(char)
		}

		self.Consume()
	}

	self.tokens = append(self.tokens, Token{
		toktype: TokenString,

		value:    str,
		position: self.pos,

		tokenspecific: string(ch),
	})

	return nil
}

func (self *Tokenizer) HandleMultilineStrings() error {
	str := string("")

	fp, _ := self.Peek(1)

	if fp == '[' {
		self.Consume()
	} else {
		if fp == '=' {
			println("Multiline strings with = in between brackets isnt supported yet!")
		}

		self.tokens = append(self.tokens, Token{
			toktype: TokenLBrac,
	
			value:    "[",
			position: self.pos,
		})

		self.Consume()
		return nil
	}

	for {
		p, e := self.Peek(1)

		if e {
			return errors.New("Unfinished multiline string")
		}

		if p == ']' {
			np, _ := self.Peek(2)

			if np == ']' {
				self.Consume()
				self.Consume()
				break
			}
		}

		str += string(p)
		self.Consume()
	}

	self.tokens = append(self.tokens, Token{
		toktype: TokenString,

		value:    str,
		position: self.pos,
	})

	return nil
}

// TODO:
// func (self *Tokenizer) HandleMultilineStartStop(out byte) (bool, int) {}

func (self *Tokenizer) HandleComment(ch byte, expecting byte) error{
	pe, err := self.Peek(1)

	if pe != expecting || err {
		self.tokens = append(self.tokens, Token{
			toktype:  TokenSub,
	
			value:    string(ch),
			position: self.pos,
		})
	}
	
	self.Consume()

	pe, err = self.Peek(1)
	pe2, err := self.Peek(2)

	comment := string("")
	
	if pe == '[' && pe2 == '[' {
		self.Consume()
		self.Consume()
		comment = self.HandleMultilineComment()
	} else {
		for {
			p, err := self.Peek(1)
	
			if p == '\n' || err {
				break
			} else {
				self.Consume()
	
				comment += string(p)
			}
		}
	}

	self.tokens = append(self.tokens, Token{
		toktype:  TokenComment,

		value:    comment,
		position: self.pos,
	})

	return nil
}

func (self *Tokenizer) HandleMultilineComment() string {
	comment := string("")

	for {
		p, _ := self.Peek(1)
		p2, _ := self.Peek(2)

		if p == ']' && p2 == ']' {
			self.Consume()
			self.Consume()
			break
		} else {
			comment += string(p)
			self.Consume()
		}
	}

	return comment
}
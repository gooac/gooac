package gooa

import "strings"

// Internal: Handle Multiline String Equal Sequences
func (self *Tokenizer) HandleEqSeq(offset int) int {
	amt := 0
	
	for {
		_, p := self.PeekSome(amt + offset)

		if p == '=' {
			amt++
			continue
		}

		break
	}

	return amt
}

// Internal: Handle Token Absorption
// Aka, >= becoming 'geq' instead of 'gt' and 'eq' 
func (self *Tokenizer) HandleAbsorb(ch byte, next byte, tokch TokenType, toknext TokenType) {
	tok := &Token{
		Token: tokch,
		Value: string(ch),
		Position: self.position.Copy(),
		EndPos: self.position.Copy(),
		WhiteSpace: self.WhiteSpace,
	}
	defer self.Append(tok)

	self.WhiteSpace = 0

	err, p := self.Peek()

	if err != nil {
		return
	}

	if p != next {
		return
	}

	self.Consume()
	tok.Token = toknext
	tok.Value += string(next)
	tok.EndPos = self.position.Copy()
}

// Internal: Handle Number Literals
func (self *Tokenizer) HandleNumber(ch byte, ws int) (TokenizationError, Token) {
	num := string(ch)
	spos := self.position.Copy()
	token := TokenNumber
	var realerr TokenizationError = TokErrNone

	for {
		err, s := self.Peek()
	
		if err != nil {
			break
		}

		if IsWhitespace(s) {
			break
		}

		if (s == 'x' || s == 'X') && num == "0" {
			self.Consume()
			num += string(s)
			
			err, hex := self.HandleHexNumber()

			if err != TokErrNone {
				realerr = err
			}

			token = TokenHexNumber
			num += hex

			break
		}

		if s == 'e' || s == 'E' {
			self.Consume()
			num += string(s)
			token = TokenSciNot

			for {
				e, n := self.Peek()

				if e != nil {
					break
				}

				if IsNumeric(n) {
					num += string(n)
					self.Consume()
				} else if IsAlpha(n) {
					num += string(n)
					self.Consume()

					realerr = TokErrMalformedSciNotLiteral
				} else {
					break
				}
			}

			break
		}

		if IsNumeric(s) {
			num += string(s)
			self.Consume()
		} else if IsAlpha(s) {
			realerr = TokErrMalformedNumber

			num += string(s)
			self.Consume()
		} else { 
			break
		}
	}

	tok := Token{
		Token:		token,
		Value: 		num,
		Position: 	spos,
		EndPos: 	self.position.Copy(),
		WhiteSpace: 	ws,
	}

	return realerr, tok
}

// Internal: Handle Hex Literals
func (self *Tokenizer) HandleHexNumber() (TokenizationError, string) {
	num := ""
	err := TokErrNone

	for {
		e, p := self.Peek()
	
		if e != nil {
			break
		}

		if IsWhitespace(p) {
			break
		}

		hn := IsHexNum(p)
		al := IsAlnum(p)
		if al && hn {
			self.Consume()
			num += string(p)
		} else if al {
			self.Consume()
			num += string(p)

			err = TokErrMalformedHexLiteral
		} else {
			break
		}
	}

	if num == "" {
		err = TokErrMalformedHexLiteral
	}

	return err, num
}


// Internal: Handle Identifiers
func (self *Tokenizer) HandleIdentifier(ch byte, ws int) {
	spos := self.position.Copy()
	id := string(ch)

	for {
		err, s := self.Peek()

		if err != nil {
			break
		}

		if IsValidIdentChar(s) {
			id += string(s)
			self.Consume()
		} else {
			break
		}
	}

	tok := Token{
		Token:		TokenIdent,
		Value: 		id,
		Position: 	spos,
		EndPos: 	self.position.Copy(),
		WhiteSpace: 	ws,
	}

	Keyword, valid := keywordTypes[id]

	if valid {
		tok.Token = TokenKeyword
		tok.Keyword = Keyword
	}

	switch Keyword {
	case KeywordAnd: 	tok.Token = TokenAnd
	case KeywordOr: 	tok.Token = TokenOr
	case KeywordTrue:	tok.Token = TokenTrue
	case KeywordFalse:	tok.Token = TokenFalse
	case KeywordNil:	tok.Token = TokenNil
	case KeywordNot:	tok.Token = TokenNot
	}

	self.Append(&tok)
}

// Internal: Handle Regular Strings
func (self *Tokenizer) HandleString(starter byte, ws int) (TokenizationError, byte) {
	str := string("")
	spos := self.position.Copy()
	badstr := false

	for {
		err, p := self.Peek()

		if err != nil {
			return TokErrUnfinishedString, starter
		} else if p == '\\' {
			str += string(p)
			self.Consume()

			_, p := self.Consume()
			
			if p == '\r' {
				self.Consume()
				str += "\\n"
				continue
			}

			str += string(p)
		} else if p == '\n' {
			badstr = true
			str += string(p)
			self.Consume()
		} else if p != starter {
			str += string(p)
			self.Consume()
		} else {
			self.Consume()
			break
		}
	}

	self.Append(&Token{
		Token:		TokenString,
		Value: 		str,
		Position: 	spos,
		EndPos: 	self.position.Copy(),
		WhiteSpace: 	ws,
		Special: 	string(starter),
	})

	// This allows consumation of the entire string before
	// breaking, otherwise would throw 2 errors
	if badstr {
		return TokErrUnfinishedString, starter
	}

	return TokErrNone, '_'
}

// Internal: Handle Multi-line Strings
func (self *Tokenizer) HandleMultilineString(ch byte, ws int) (TokenizationError, int) {
	ierr, ipeek := self.Peek()
	start := self.position.Copy()

	if (ierr != nil) || (ipeek != '[' && ipeek != '=') {
		self.Append(&Token{
			Token: TokenLBrac,
			Value: string(ch),
			Position: start,
			EndPos: start,
		})

		return TokErrNone, 0
	}


	amt := self.HandleEqSeq(0)
	self.ConsumeAmount(amt)
	
	str := string("")

	for {
		err, p := self.Peek()

		if err != nil {
			return TokErrUnfinishedMLString, amt
		}

		if p == ']' {
			newamt := self.HandleEqSeq(1)
			err2, p2 := self.PeekSome(newamt + 1)

			if err2 != nil {
				return TokErrUnfinishedMLString, amt
			}

			if newamt != amt {
				str += string(p)
				self.Consume()
				continue
			}

			if p2 == ']' {
				self.Consume()
				self.ConsumeAmount(newamt)
				break
			}
		}
		
		if p == '\r' {
			self.Consume()
			continue
		}

		str += string(p)
		self.Consume()
	}

	self.Append(&Token{
		Token: TokenMLString,
		Value: string(str),
		Position: start,
		EndPos: self.position.Copy(),
		Special: strings.Repeat("=", amt),
	})

	return TokErrNone, 0
}

// Internal: Handle Variadic Arguments and Concatenation Operators
func (self *Tokenizer) HandleVariadics(pd byte, ws int) {
	pos := self.position.Copy()
	tok := &Token{
		Token: TokenPeriod,
		Value: ".",
		Position: pos,
		EndPos: pos,
	}
	defer self.Append(tok)

	err, p := self.Peek()
	if err != nil {
		return
	}

	if p != '.' {
		return
	}

	self.Consume()
	tok.Token = TokenConcat
	tok.Value = ".."
	tok.EndPos = self.position.Copy()

	err, p = self.Peek()
	if err != nil {
		return
	}

	if p != '.' {
		return
	}
	
	self.Consume()
	tok.Token = TokenVariadic
	tok.Value = "..."
	tok.EndPos = self.position.Copy()
}

// Internal: Handle Regular Comments
func (self *Tokenizer) HandleComment(ch byte, ws int) {
	start := self.position.Copy()
	nexterr, nextp := self.Peek()

	if nextp == '>' {
		self.Consume()
		self.Append(&Token{
			Token: TokenColon,
			Value: string(ch) + ">",
			Position: start,
			EndPos: self.position.Copy(),
		})
		return
	} else if (nexterr != nil) || nextp != '-' {
		self.Append(&Token{
			Token: TokenSub,
			Value: string(ch),
			Position: start,
			EndPos: start,
		})

		return
	}

	self.Consume()
	nexterr, nextp = self.Peek()
	if nextp == '[' {
		amt := self.HandleEqSeq(1)
		_, np := self.PeekSome(amt + 1)

		if np == '[' {
			self.ConsumeAmount(1 + amt)
			self.HandleMultilineComment(&start, amt)
			return
		}
	}

	cmt := string("")
	
	for {
		err, p := self.Peek()

		if (err != nil) || (p == '\n') || (p == '\r') {
			break
		}
	
		cmt += string(p)
		self.Consume()
	}

	self.Append(&Token{
		Token: TokenComment,
		Value: cmt,
		Position: start,
		EndPos: self.position.Copy(),
	})
}

// Internal: Handle CStyle COmments
func (self *Tokenizer) HandleCStyleComment(ch byte, ws int) {
	start := self.position.Copy()
	nexterr, nextp := self.Peek()

	if (nexterr != nil) || (nextp != '/' && nextp != '*') {
		self.Append(&Token{
			Token: TokenDiv,
			Value: string(ch),
			Position: start,
			EndPos: start,
		})

		return
	}

	cmt := ""

	self.Consume()
	if nextp == '/' {
		for {
			err, p := self.Peek()

			if (err != nil) || (p == '\n') || (p == '\r') {
				break
			}
		
			cmt += string(p)
			self.Consume()
		}
	} else {
		for {
			err, p := self.Peek()
			err2, pn := self.PeekSome(1)

			if (err != nil) || (err2 != nil) {
				break
			}

			if p == '\r' {
				self.Consume()
				continue
			}

			if p == '*' && pn == '/' {
				self.Consume()
				self.Consume()
				break
			}
		
			cmt += string(p)
			self.Consume()
		}
	}

	self.Append(&Token{
		Token: TokenComment,
		Value: cmt,
		Position: start,
		EndPos: self.position.Copy(),
	})
}

// Internal: Handle Multiline Comments
func (self *Tokenizer) HandleMultilineComment(start *Position, eqamt int) {
	cmt := string("")

	for {
		err, p := self.Peek()

		if err != nil {
			break
		}

		if p == '\r' {
			self.Consume()
			continue
		}

		if p == ']' {
			amt := self.HandleEqSeq(1)
			_, np := self.PeekSome(amt + 1)

			if np == ']' && amt == eqamt {
				self.ConsumeAmount(1 + amt)
				break
			}
		}

		cmt += string(p)
		self.Consume()
	}

	self.Append(&Token{
		Token: TokenComment,
		Value: cmt,
		Position: *start,
		EndPos: self.position.Copy(),
	})
}

// Internal: Handle Arbitrary Symbols
func (self *Tokenizer) HandleSymbol(sym byte) TokenizationError {
	val, valid := symbolLookups[string(sym)]

	if valid {
		self.Append(&Token{
			Token: val,
			Value: string(sym),
			Position: self.position.Copy(),
			EndPos: self.position.Copy(),
		})
		return TokErrNone
	}

	return TokErrUnknownSymbol
}
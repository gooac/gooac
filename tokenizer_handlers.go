
package gooa

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
		token: tokch,
		value: string(ch),
		position: self.position.Copy(),
		endpos: self.position.Copy(),
		wspace: self.wspace,
	}
	defer self.Append(tok)

	self.wspace = 0

	err, p := self.Peek()

	if err != nil {
		return
	}

	if p != next {
		return
	}

	self.Consume()
	tok.token = toknext
	tok.value += string(next)
	tok.endpos = self.position.Copy()
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
		token:		token,
		value: 		num,
		position: 	spos,
		endpos: 	self.position.Copy(),
		wspace: 	ws,
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
		token:		TokenIdent,
		value: 		id,
		position: 	spos,
		endpos: 	self.position.Copy(),
		wspace: 	ws,
	}

	kwtype, valid := keywordTypes[id]

	if valid {
		tok.token = TokenKeyword
		tok.kwtype = kwtype
	}

	switch kwtype {
	case KeywordAnd: 	tok.token = TokenAnd
	case KeywordOr: 	tok.token = TokenOr
	case KeywordTrue:	tok.token = TokenTrue
	case KeywordFalse:	tok.token = TokenFalse
	case KeywordNil:	tok.token = TokenNil
	case KeywordNot:	tok.token = TokenNot
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
		token:		TokenString,
		value: 		str,
		position: 	spos,
		endpos: 	self.position.Copy(),
		wspace: 	ws,
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
			token: TokenLBrac,
			value: string(ch),
			position: start,
			endpos: start,
		})

		return TokErrNone, 0
	}

	
	amt := self.HandleEqSeq(0)
	self.Consume()
	
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
		token: TokenMLString,
		value: string(str),
		position: start,
		endpos: self.position.Copy(),
	})

	return TokErrNone, 0
}

// Internal: Handle Variadic Arguments and Concatenation Operators
func (self *Tokenizer) HandleVariadics(pd byte, ws int) {
	pos := self.position.Copy()
	tok := &Token{
		token: TokenPeriod,
		value: ".",
		position: pos,
		endpos: pos,
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
	tok.token = TokenConcat
	tok.value = ".."
	tok.endpos = self.position.Copy()

	err, p = self.Peek()
	if err != nil {
		return
	}

	if p != '.' {
		return
	}
	
	self.Consume()
	tok.token = TokenVariadic
	tok.value = "..."
	tok.endpos = self.position.Copy()
}

// Internal: Handle Regular Comments
func (self *Tokenizer) HandleComment(ch byte, ws int) {
	start := self.position.Copy()
	nexterr, nextp := self.Peek()

	// Feat: CallArrow Syntax
	if nextp == '>' {
		self.Consume()
		self.Append(&Token{
			token: TokenColon,
			value: string(ch) + ">",
			position: start,
			endpos: self.position.Copy(),
		})
		return
	} else if (nexterr != nil) || nextp != '-' {
		self.Append(&Token{
			token: TokenSub,
			value: string(ch),
			position: start,
			endpos: start,
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
		token: TokenComment,
		value: cmt,
		position: start,
		endpos: self.position.Copy(),
	})
}

// Internal: Handle CStyle COmments
func (self *Tokenizer) HandleCStyleComment(ch byte, ws int) {
	start := self.position.Copy()
	nexterr, nextp := self.Peek()

	if (nexterr != nil) || (nextp != '/' && nextp != '*') {
		self.Append(&Token{
			token: TokenDiv,
			value: string(ch),
			position: start,
			endpos: start,
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
		token: TokenComment,
		value: cmt,
		position: start,
		endpos: self.position.Copy(),
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
		token: TokenComment,
		value: cmt,
		position: *start,
		endpos: self.position.Copy(),
	})
}

// Internal: Handle Arbitrary Symbols
func (self *Tokenizer) HandleSymbol(sym byte) TokenizationError {
	val, valid := symbolLookups[string(sym)]

	if valid {
		self.Append(&Token{
			token: val,
			value: string(sym),
			position: self.position.Copy(),
			endpos: self.position.Copy(),
		})
		return TokErrNone
	}

	return TokErrUnknownSymbol
}
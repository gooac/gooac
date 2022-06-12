package gooa

import (
	"errors"
	"regexp"
	"strings"
)

// Main construct to tokenize Gooa code
type Tokenizer struct {
	tokenizing 	[]byte 					// Whats being currently tokenized
	tokens 		[]Token					// Output tokens
	index 		int						// Current index being read of tokenization string 
	position 	*Position				// Position the tokenizers at
	err 		ErrorHandler			// ErrorHandler to be used to output errors and warn messages.
	wspace 		int 					// Current whitespace
	regex 		*regexp.Regexp  		// Shebang Stripper
	last 		*Token					// Last token to be appended
	middleware 	*MiddlewareHandler	 	// All instantiated middleware for the tokenizer
}

// Resets the tokenizer state variables
func (self *Tokenizer) Reset() {
	self.tokenizing = nil
	self.tokens 	= nil
	self.index 		= 0
	self.position 	= nil
	self.err 		= nil
	self.wspace 	= 0
}

// Tokenizes given byte slice into Tokens
func (self *Tokenizer) Tokenize(str []byte, e *ErrorHandler) ([]Token, bool) {
	self.Reset()

	if self.regex == nil {
		self.regex = regexp.MustCompile(`\A\#.*`)
	}

	self.tokenizing	= self.middleware.PreTokenize(self.regex.ReplaceAll(str, []byte{}))
	self.tokens 	= []Token{}
	self.index 		= 0
	self.position	= &Position{1, 0, &self.index}
	self.wspace 	= 0

	if e == nil {
		self.err = &BaseErrorHandler{}
	} else {
		self.err = *e
	}

	self.err.SetPosition(self.position)
	self.err.SetErrorRealm(ErrorRealmTokenizer)

	for {
		conserr, b := self.Consume()

		if conserr != nil {
			self.Append(&Token{
				token: TokenEOF,
				value: "<EOF>",
				position: self.position.Copy(),
				endpos: self.position.Copy(),
			})

			break
		}

		if self.middleware.TokenizerHandleByte(self, b) {
			continue
		}

		// Ignore newlines
		if b == '\n' {
			if self.last != nil {
				self.last.newline++
			}
			continue
		} else if b == '\r' {
			continue
		// Whitespace Incrementor
		} else if IsWhitespace(b) {
			self.wspace++
			continue
		
		// Check if byte is numeric and run number handler
		} else if IsNumeric(b) {
			err, tok := self.HandleNumber(b, self.wspace)
			self.wspace = 0

			self.Append(&tok)

			if err != TokErrNone {
				self.err.Error(ErrorGeneral, string(err), tok.value)
			}
			continue
		
		// Check if byte is a valid identifier starter, if so run ident handler
		} else if IsAlpha(b) || b == '_' {
			self.HandleIdentifier(b, self.wspace)
			self.wspace = 0

			continue
		}
		
		switch b {
		// String handlers
		case '\'', '"': 
			err, starter := self.HandleString(b, self.wspace)

			if err != TokErrNone {
				self.err.Error(ErrorFatal, string(err), string(starter))
			}
		case '[': 
			err, cnt := self.HandleMultilineString(b, self.wspace)

			if err != TokErrNone {
				self.err.Error(ErrorFatal, string(err), "]" + strings.Repeat("=", cnt) + "]")
			}
		
		// Comment Handlers
		case '-':
			self.HandleComment(b, self.wspace)
		case '/':
			self.HandleCStyleComment(b, self.wspace)

		// Variadic and Concatenation Handling
		case '.': self.HandleVariadics(b, self.wspace)
		
		// Absorptions
		case '>': self.HandleAbsorb(b, '=', TokenGt, TokenGtEq)
		case '<': self.HandleAbsorb(b, '=', TokenLt, TokenLtEq)
		case '=': self.HandleAbsorb(b, '=', TokenEq, TokenIsEq)
		case ':': self.HandleAbsorb(b, ':', TokenColon, TokenLabel)
		case '!', '~': self.HandleAbsorb(b, '=', TokenNot, TokenNotEq)

		default:
			err := self.HandleSymbol(b)

			if err != TokErrNone {
				self.err.Error(ErrorGeneral, string(err), string(b))
			}
		}
	}

	self.err.Dump()

	return self.middleware.PostTokenize(self.tokens), self.err.ShouldStop()
}

// Should the tokenizer stop?
func (self *Tokenizer) Stop(index int) bool {
	return index >= len(self.tokenizing)	
}

// Peek at current byte
func (self *Tokenizer) Peek() (error, byte) {
	if self.Stop(self.index) {
		return errors.New("Peeking at a byte that doesnt exist"), 'F'
	}

	b := self.tokenizing[self.index]

	return nil, b
}

// Peek at arbitrary byte
func (self *Tokenizer) PeekSome(amt int) (error, byte) {
	if self.Stop(self.index + amt) {
		return errors.New("Peeking at a byte that doesnt exist"), 'F'
	}

	b := self.tokenizing[self.index + amt]

	return nil, b
}

// Consume and return the current byte
func (self *Tokenizer) Consume() (error, byte) {
	if self.Stop(self.index) {
		return errors.New("Consuming a byte that doesnt exist"), 'F'
	}

	b := self.tokenizing[self.index]
	
	self.index++
	self.position.column++

	if b == '\n' {
		self.position.column = 0
		self.position.line++
	}

	return nil, b
}

// Consume an arbitrary amount
func (self *Tokenizer) ConsumeAmount(amt int) {
	for i := 0; i <= amt; i++ {
		self.Consume()
	}
}

// Append a token created prior
func (self *Tokenizer) Append(t *Token) {
	self.tokens = append(self.tokens, *t)
	self.last = &self.tokens[len(self.tokens) - 1]
}

// Debug: Dump Tokens
func (self *Tokenizer) DumpTokens() {
	for _, v := range self.tokens {
		v.Print()
	}
}
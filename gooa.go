package gooa

import (
	"errors"
	"os"
)

type Gooa struct {
	compiler Compiler
	parser Parser
	tokenizer Tokenizer

	err ErrorHandler
	middleware *MiddlewareHandler
}

func NewGooa() *Gooa {
	gooa := &Gooa{
		compiler: Compiler{},
		parser: Parser{},
		tokenizer: Tokenizer{},
	
		err: &BaseErrorHandler{},
		middleware: &MiddlewareHandler{},
	}

	gooa.compiler.middleware = gooa.middleware
	gooa.parser.middleware = gooa.middleware
	gooa.tokenizer.middleware = gooa.middleware

	return gooa 
}

func (self *Gooa) Compile(s []byte) (string, bool) {
	toks, stop := self.tokenizer.Tokenize(s, &self.err)

	if stop {
		self.err.Error(ErrorFatal, "Tokenization Failed")
		return "Failed to tokenize", true
	}

	ast, stop := self.parser.Parse(toks, &self.err)

	if stop {
		self.err.Error(ErrorFatal, "Parsing Failed")
		return "Failed to parse", true
	}

	cmp, stop := self.compiler.Compile(&ast, &self.err)

	if stop {
		self.err.Error(ErrorFatal, "Compilation Failed")
		return "Failed to compile", true
	}

	return cmp, false
}

func (self *Gooa) CompileFile(s string) (string, error) {
	f, err := os.ReadFile(s)

	if err != nil {
		return "Error Reading File", err
	}

	c, valid := self.Compile(f)

	if !valid {
		return "", errors.New(c)
	}

	return c, nil
}

func (self *Gooa) Use(m Middleware) {
	self.middleware.Use(m)
}
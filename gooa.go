package gooa

import (
	"errors"
	"os"
)

type Gooa struct {
	Compiler Compiler
	Parser Parser
	Tokenizer Tokenizer

	Err ErrorHandler
	Middleware *MiddlewareHandler
}

func NewGooa() *Gooa {
	gooa := &Gooa{
		Compiler: Compiler{},
		Parser: Parser{},
		Tokenizer: Tokenizer{},
	
		Err: &BaseErrorHandler{},
		Middleware: &MiddlewareHandler{},
	}

	gooa.Compiler.middleware = gooa.Middleware
	gooa.Parser.middleware = gooa.Middleware
	gooa.Tokenizer.middleware = gooa.Middleware

	return gooa 
}

func (self *Gooa) Compile(s []byte) (string, bool) {
	toks, stop := self.Tokenizer.Tokenize(s, &self.Err)

	if stop {
		self.Err.Error(ErrorFatal, "Tokenization Failed")
		return "Failed to tokenize", true
	}

	ast, stop := self.Parser.Parse(toks, &self.Err)

	if stop {
		self.Err.Error(ErrorFatal, "Parsing Failed")
		return "Failed to parse", true
	}

	cmp, stop := self.Compiler.Compile(&ast, &self.Err)

	if stop {
		self.Err.Error(ErrorFatal, "Compilation Failed")
		return "Failed to compile", true
	}

	return cmp, false
}

func (self *Gooa) CompileFile(s string) (string, error) {
	f, err := os.ReadFile(s)

	if err != nil {
		return "Error Reading File", err
	}

	c, errored := self.Compile(f)

	if errored {
		return "", errors.New(c)
	}

	return c, nil
}

func (self *Gooa) Use(m Middleware) {
	self.Middleware.Use(m)
}
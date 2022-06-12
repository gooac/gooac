package gooa

type MiddlewareHookFunc func(...interface{}) interface{}
type MiddlewareHook int

const (
	_ MiddlewareHook = iota

	MiddlewareHookPreTokenize
	MiddlewareHookPostTokenize
	MiddlewareHookPostParse
	MiddlewareHookPostCompile

	MiddlewareHookTokenizerHandleByte
	MiddlewareHookParserHandleToken
	MiddlewareHookCompileHandleNode
	MiddlewareHookCompileHandleRootNode
)

type Middleware interface {
	GetHooks() map[MiddlewareHook]MiddlewareHookFunc
}

type MiddlewareHookInternal struct {
	owner Middleware
	call  MiddlewareHookFunc
	hook  MiddlewareHook
}

type MiddlewareHandler struct {
	middleware []Middleware

	hooks map[MiddlewareHook][]MiddlewareHookInternal
}

func (self *MiddlewareHandler) Use(m Middleware) {
	if self.hooks == nil {
		self.hooks = map[MiddlewareHook][]MiddlewareHookInternal{}
	}

	for hook, fn := range m.GetHooks() {
		self.hooks[hook] = []MiddlewareHookInternal{{
			owner: m,
			call:  fn,
			hook:  hook,
		}}
	}

	(*self).hooks = self.hooks
}

func (self *MiddlewareHandler) CallHook(h MiddlewareHook, a ...interface{}) interface{} {
	if self == nil {
		return false
	}

	if (*self).hooks == nil {
		return false
	}

	m, val := self.hooks[h]

	if !val {
		return false
	}

	for _, v := range m {
		c := v.call(a...)
		if c != nil {
			return c
		}
	}

	return false
}

func (self *MiddlewareHandler) PreTokenize(str []byte) []byte {
	c := self.CallHook(MiddlewareHookPreTokenize, str)

	if c == false {
		return str
	}

	return c.([]byte)
}

func (self *MiddlewareHandler) PostTokenize(toks []Token) []Token {
	c := self.CallHook(MiddlewareHookPostTokenize, toks)

	if c == false {
		return toks
	}

	return c.([]Token)
}

func (self *MiddlewareHandler) PostParse(ast *AST) *AST {
	c := self.CallHook(MiddlewareHookPostParse, ast)

	if c == false {
		return ast
	}

	return c.(*AST)
}

func (self *MiddlewareHandler) PostCompile(str string) string {
	c := self.CallHook(MiddlewareHookPostCompile, str)

	if c == false {
		return str
	}

	return c.(string)
}

func (self *MiddlewareHandler) TokenizerHandleByte(tok *Tokenizer, b byte) bool {
	c := self.CallHook(MiddlewareHookTokenizerHandleByte, tok, b)

	return c.(bool)
}

func (self *MiddlewareHandler) ParserHandleToken(parser *Parser, tok Token) bool {
	c := self.CallHook(MiddlewareHookParserHandleToken, parser, tok)

	return c.(bool)
}

func (self *MiddlewareHandler) CompileHandleNode(cmp *Compiler, node *ASTNode) bool {
	c := self.CallHook(MiddlewareHookCompileHandleNode, cmp, node)

	return c.(bool)
}

func (self *MiddlewareHandler) CompileHandleRootNode(cmp *Compiler, node *ASTNode) bool {
	c := self.CallHook(MiddlewareHookCompileHandleRootNode, cmp, node)

	return c.(bool)
}
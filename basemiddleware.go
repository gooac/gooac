
package gooa

const (
	AttributeMiddlewareMustBeCall 	string = "Expected function call in function attribute, got %v"
	AttributeMiddlewareMustBeFn 	string = "Expected function definition after function attribute declaration, got '%v'"
) 

type attributeMiddleware struct {
	hooks map[MiddlewareHook]MiddlewareHookFunc
}

func (self attributeMiddleware) GetHooks() map[MiddlewareHook]MiddlewareHookFunc {
	return self.hooks
}

func AttributeMiddleware() attributeMiddleware {
	m := attributeMiddleware{}
	m.hooks = map[MiddlewareHook]MiddlewareHookFunc{}
	m.hooks[MiddlewareHookParserHandleToken] = func(i ...interface{}) interface{} {
		prs := i[0].(*Parser)
		tok := i[1].(Token)

		if tok.token == TokenAttr && prs.PeekSome(1).token == TokenLBrac {
			prs.Consume()
			prs.Consume()
			m.Handle(prs)
			
			return true
		}

		return false
	}

	return m
}

func (mw *attributeMiddleware) Handle(self *Parser) {
	nodes := []*ASTNode{}

	for {
		p := self.Peek()

		expr := self.ExpectExpression(p)

		if expr.nodetype != NodeCall && expr.nodetype != NodeMethodCall {
			self.err.Error(ErrorFatal, AttributeMiddlewareMustBeCall, expr.nodetype)
			return
		}

		nodes = append(nodes, expr)
	
		if self.Peek().token == TokenComma {
			self.Consume()
			continue
		} else {
			break
		}
	}

	t := self.Consume()
	if t.token != TokenRBrac {
		self.err.Error(ErrorFatal, ParserErrorExpectedSymbol, "]", t.value)
		return
	}

	if len(nodes) == 0 {
		self.err.Error(ErrorFatal, AttributeMiddlewareMustBeCall, "none")
		return
	}

	for k, v := range nodes {
		v.Dump(k*2)
		println("----")
	}

	n := self.Peek()
	var fn *ASTNode

	if n.token != TokenKeyword {
		self.err.Error(ErrorFatal, AttributeMiddlewareMustBeFn, n.value)
		return
	}

	if n.kwtype == KeywordFunction {
		self.Consume()
		fn = self.HandleFunction(n)
	} else if n.kwtype == KeywordLocal {
		self.Consume()
		self.Consume()
		fn = self.HandleFunction(n)

		fn.nodetype = NodeLocalFunction
	}

	if fn == nil || (fn.nodetype != NodeFunction && fn.nodetype != NodeLocalFunction) {
		self.err.Error(ErrorFatal, AttributeMiddlewareMustBeFn, fn.nodetype)
		return
	}

	self.ast.Add(fn)

	for _, v := range nodes {
		v.body[1].body = append([]*ASTNode{
			fn.body[0],
		}, v.body[1].body...)
		self.ast.Add(&ASTNode{
			nodetype: NodeVariableAssignment,
			body: []*ASTNode{
				{
					nodetype: NodeVariableNameList,
					body: []*ASTNode{
						fn.body[0],
					},
				},
				{
					nodetype: NodeVariableValList,
					body: []*ASTNode{
						v,
					},
				},
			},
		})
	}
}
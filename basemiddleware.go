
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

		if tok.Token == TokenAttr && prs.PeekSome(1).Token == TokenLBrac {
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

		if expr.Nodetype != NodeCall && expr.Nodetype != NodeMethodCall {
			self.err.Error(ErrorFatal, AttributeMiddlewareMustBeCall, expr.Nodetype)
			return
		}

		nodes = append(nodes, expr)
	
		if self.Peek().Token == TokenComma {
			self.Consume()
			continue
		} else {
			break
		}
	}

	t := self.Consume()
	if t.Token != TokenRBrac {
		self.err.Error(ErrorFatal, ParserErrorExpectedSymbol, "]", t.Value)
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

	if n.Token != TokenKeyword {
		self.err.Error(ErrorFatal, AttributeMiddlewareMustBeFn, n.Value)
		return
	}

	if n.Keyword == KeywordFunction {
		self.Consume()
		fn = self.HandleFunction(n)
	} else if n.Keyword == KeywordLocal {
		self.Consume()
		self.Consume()
		fn = self.HandleFunction(n)

		fn.Nodetype = NodeLocalFunction
	}

	if fn == nil || (fn.Nodetype != NodeFunction && fn.Nodetype != NodeLocalFunction) {
		self.err.Error(ErrorFatal, AttributeMiddlewareMustBeFn, fn.Nodetype)
		return
	}

	self.ast.Add(fn)

	for _, v := range nodes {
		v.Body[1].Body = append([]*ASTNode{
			fn.Body[0],
		}, v.Body[1].Body...)
		self.ast.Add(&ASTNode{
			Nodetype: NodeVariableAssignment,
			Body: []*ASTNode{
				{
					Nodetype: NodeVariableNameList,
					Body: []*ASTNode{
						fn.Body[0],
					},
				},
				{
					Nodetype: NodeVariableValList,
					Body: []*ASTNode{
						v,
					},
				},
			},
		})
	}
}
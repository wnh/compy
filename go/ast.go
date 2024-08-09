package main

type Node interface {
	isNode()
	Codegen(thing *CodegenModule)
}
type node struct{}

type AstModule struct {
	node
	Name       *AstIdent
	Statements []AstStatement
}

type AstStatement interface {
	Node
	isStatement()
	ForwardDecl(thing *CodegenModule)
}
type AstExpr interface {
	Node
}

type AstConstAssign struct {
	node
	Ident string
	Type  *AstType
	Value AstExpr
}

type AstIdent struct {
	Name string
}

type AstType struct {
	node
	Name *AstIdent
}

type AstIntLitExpr struct {
	node
	Value int
}
type AstStringLitExpr struct {
	node
	Value string
}

type AstFnDecl struct {
	node
	Name       *AstIdent
	ReturnType *AstType
	Params     []*AstParam // Obviously not
	Body       *AstBlock
}

type AstBlock struct {
	Body []AstStatement
}

type AstParam struct {
	Name *AstIdent
	Type *AstType
}

type AstFnCall struct {
	Name *AstIdent
	Args []AstExpr
}

func (n *AstConstAssign) isNode()   {}
func (s *AstIntLitExpr) isNode()    {}
func (s *AstStringLitExpr) isNode() {}
func (s *AstFnDecl) isNode()        {}
func (s *AstFnCall) isNode()        {}
func (s *AstIdent) isNode()         {}

func (s *AstConstAssign) isStatement() {}
func (s *AstFnDecl) isStatement()      {}
func (s *AstFnCall) isStatement()      {}

func (s *AstIntLitExpr) isExpr()    {}
func (s *AstStringLitExpr) isExpr() {}
func (s *AstIdent) isExpr()         {}

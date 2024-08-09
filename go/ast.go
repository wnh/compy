package main

type Node interface {
	isNode()
	Codegen(thing *CodegenModule)
}
type node struct{}

type AstModule struct {
	node
	Name       string
	Statements []AstStatement
}

type AstStatement interface {
	Node
	isStatement()
}
type AstExpr interface {
	Node
	isExpr()
}

type AstConstAssign struct {
	node
	Ident string
	Type  *AstType
	Value AstExpr
}

type AstType struct {
	node
	Name string
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
	Name       string
	ReturnType *AstType
	Params     []*AstParam // Obviously not
	Body       *AstBlock
}

type AstBlock struct {
	Body []AstStatement
}

type AstParam struct {
	Name string
	Type *AstType
}

func (n *AstConstAssign) isNode()   {}
func (s *AstIntLitExpr) isNode()    {}
func (s *AstStringLitExpr) isNode() {}
func (s *AstFnDecl) isNode()        {}

func (s *AstConstAssign) isStatement() {}
func (s *AstFnDecl) isStatement()      {}

func (s *AstIntLitExpr) isExpr()    {}
func (s *AstStringLitExpr) isExpr() {}

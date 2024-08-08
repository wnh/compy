
package main

type Node struct {}

type AstModule struct {
	Node
	Name string
	Statements []AstStatement
}

type AstStatement interface {
	isStatement()
}
type AstExpr interface {}

type AstConstAssign struct {
	Ident string
	Type *AstType
	Value AstExpr
}

func (AstConstAssign) isStatement() {}

type AstType struct {
	Name string
}

type AstIntLitExpr struct {
	Value int
}
type AstStringLitExpr struct {
	Value string
}

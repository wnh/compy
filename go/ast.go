
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


type AstConstAssign struct {
	Ident string
	Value int
}
func (AstConstAssign) isStatement() {}

package main

import (
	"fmt"
	"strings"
)

type CodegenModule struct {
	Code strings.Builder
}

func (c *CodegenModule) Write(s string) error {
	_, _ = c.Code.WriteString(" ")
	_, err := c.Code.WriteString(s)
	_, _ = c.Code.WriteString(" ")
	return err
}

func (c *CodegenModule) Nl() {
	_, _ = c.Code.WriteString("\n")
}

func (c *CodegenModule) Writef(format string, a ...any) error {
	return c.Write(fmt.Sprintf(format, a...))
}

func (c *CodegenModule) WriteRuntime() {
	c.Write("#include <stdio.h>\n")
	c.Write("typedef char* string;")
	c.Nl()
}

func (n *AstModule) Codegen(cg *CodegenModule) {
	cg.WriteRuntime()
	cg.Writef("/* Module: %s */", n.Name)
	cg.Nl()
	for _, stmt := range n.Statements {
		stmt.ForwardDecl(cg)
		cg.Write(";\n")
	}
	cg.Nl()
	for _, stmt := range n.Statements {
		stmt.Codegen(cg)
		cg.Nl()
	}

	cg.Nl()
}

func (n *AstConstAssign) Codegen(cg *CodegenModule) {
	cg.Write("const")
	n.Type.Codegen(cg)
	cg.Write(n.Ident)
	cg.Write("=")
	n.Value.Codegen(cg)
}
func (n *AstConstAssign) ForwardDecl(cg *CodegenModule) {}

func (n *AstType) Codegen(cg *CodegenModule) {
	cg.Write(n.Name.Name)
}

func (n *AstIntLitExpr) Codegen(cg *CodegenModule) {
	cg.Writef("%d", n.Value)
}

func (n *AstStringLitExpr) Codegen(cg *CodegenModule) {
	cg.Writef("\"%s\"", n.Value)
}

func (n *AstFnDecl) Codegen(cg *CodegenModule) {
	n.ReturnType.Codegen(cg)
	cg.Write(n.Name.Name + "(")
	paramCount := len(n.Params)
	for i, p := range n.Params {
		p.Type.Codegen(cg)
		cg.Write(p.Name.Name)
		if i != paramCount-1 {
			cg.Write(",")
		}
	}
	cg.Write(")")
	n.Body.Codegen(cg)
}
func (n *AstFnDecl) ForwardDecl(cg *CodegenModule) {
	n.ReturnType.Codegen(cg)
	cg.Write(n.Name.Name + "(")
	for i, p := range n.Params {
		if i != 0 {
			cg.Write(",")
		}
		p.Type.Codegen(cg)
	}
	cg.Write(")")
}

func (n *AstBlock) Codegen(cg *CodegenModule) {
	cg.Write("{\n")
	for _, s := range n.Body {
		s.Codegen(cg)
		cg.Write(";\n")
	}
	cg.Write("\n}")
}

func (n *AstFnCall) Codegen(cg *CodegenModule) {
	cg.Write(n.Name.Name + "(")
	for i, arg := range n.Args {
		if i != 0 {
			cg.Write(",")
		}
		arg.Codegen(cg)
	}
	cg.Write(")")
}

func (n *AstFnCall) ForwardDecl(cg *CodegenModule) {}

func (n *AstIdent) Codegen(cg *CodegenModule) {
	cg.Write(n.Name)
}

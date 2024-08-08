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
	c.Write("typedef char* string;")
	c.Nl()
}

func (n *AstModule) Codegen(cg *CodegenModule) {
	cg.WriteRuntime()
	cg.Writef("/* Module: %s */", n.Name)
	cg.Nl()
	for _, stmt := range n.Statements {
		stmt.Codegen(cg)
		cg.Nl()
	}

	cg.Nl()
	cg.Writef("int main() { return 0; }")
}

func (n *AstConstAssign) Codegen(cg *CodegenModule) {
	cg.Write("const")
	n.Type.Codegen(cg)
	cg.Write(n.Ident)
	cg.Write("=")
	n.Value.Codegen(cg)
	cg.Write(";")
}

func (n *AstType) Codegen(cg *CodegenModule) {
	cg.Write(n.Name)
}

func (n *AstIntLitExpr) Codegen(cg *CodegenModule) {
	cg.Writef("%d", n.Value)
}

func (n *AstStringLitExpr) Codegen(cg *CodegenModule) {
	cg.Writef("\"%s\"", n.Value)
}

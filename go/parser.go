package main

import (
	"fmt"
	"strconv"
)

type Parser struct {
	lexer   Lexer
	tok     Token
	nextTok Token

	filename string
}

type ParseError struct {
	Msg      string
	Filename string
	Line     uint
	Col      uint
}

func NewParser(input string, filename string) Parser {
	lex := NewLexer(input)
	t0 := lex.Next()
	t1 := lex.Next()
	p := Parser{lexer: lex, tok: t0, nextTok: t1, filename: filename}
	return p
}

func (p *Parser) nextToken() {
	p.tok = p.nextTok
	p.nextTok = p.lexer.Next()
}

func (p *Parser) expectv(expected TokenKind) (string, error) {
	//fmt.Printf("expectv: %v, %v, %v\n", p, expected, p.tok.Kind)
	if p.tok.Kind != expected {
		return "", p.parseErrorExp(expected)
	}
	val := p.tok.Text
	p.nextToken()
	return val, nil
}

func (p *Parser) expect(expected TokenKind) error {
	_, err := p.expectv(expected)
	return err
}

func (p *Parser) peek() TokenKind {
	return p.tok.Kind
}
func (p *Parser) peekIs(expected TokenKind) bool {
	return p.tok.Kind == expected
}

func (p *Parser) parseError() error {
	return p.parseErrorMsg(fmt.Sprintf("unexpected token: %v", p.tok.Kind))
}
func (p *Parser) parseErrorExp(expect TokenKind) error {
	return p.parseErrorMsg(fmt.Sprintf("expected token: %v got: %v", expect, p.tok.Kind))
}
func (p *Parser) parseErrorMsg(msg string) error {
	return ParseError{
		Msg:      msg,
		Filename: p.filename,
		Line:     p.tok.Line,
		Col:      p.tok.Col,
	}
}
func (e ParseError) Error() string {
	return fmt.Sprintf("%s:%d:%d  %s", e.Filename, e.Line, e.Col, e.Msg)
}

func (p *Parser) ParseModule() (*AstModule, error) {
	if err := p.expect(TokModule); err != nil {
		return nil, err
	}
	modName, err := p.expectv(TokIdent)
	if err != nil {
		return nil, err
	}
	if err := p.expect(TokSemi); err != nil {
		return nil, err
	}
	mod := AstModule{Name: modName}
	for {
		if p.peekIs(TokEof) {
			break
		}
		st, err := p.ParseStatement()
		if err != nil {
			return nil, err
		}
		mod.Statements = append(mod.Statements, st)
	}
	return &mod, nil
}

func (p *Parser) ParseStatement() (AstStatement, error) {
	switch p.peek() {
	case TokLet:
		return p.ParseConstAssign()
	case TokFn:
		return p.ParseFnDecl()
	}
	return nil, p.parseError()
}

func (p *Parser) ParseConstAssign() (*AstConstAssign, error) {
	if err := p.expect(TokLet); err != nil {
		return nil, err
	}
	constName, err := p.expectv(TokIdent)
	if err != nil {
		return nil, err
	}
	if err := p.expect(TokColon); err != nil {
		return nil, err
	}
	type_, err := p.ParseType()
	if err != nil {
		return nil, err
	}
	if err := p.expect(TokAssign); err != nil {
		return nil, err
	}
	valExpr, err := p.ParseExpr()
	if err != nil {
		return nil, err
	}
	if err := p.expect(TokSemi); err != nil {
		return nil, err
	}
	return &AstConstAssign{Ident: constName, Type: type_, Value: valExpr}, nil
}

func (p *Parser) ParseExpr() (AstExpr, error) {
	if p.peekIs(TokInt) {
		return p.ParseIntLitExpr()
	} else if p.peekIs(TokString) {
		return p.ParseStringLitExpr()
	}
	return nil, p.parseError()
}

func (p *Parser) ParseIntLitExpr() (*AstIntLitExpr, error) {
	intText, err := p.expectv(TokInt)
	if err != nil {
		return nil, err
	}
	intVal, err := strconv.Atoi(intText)
	if err != nil {
		return nil, err
	}
	return &AstIntLitExpr{Value: intVal}, nil
}

func (p *Parser) ParseStringLitExpr() (*AstStringLitExpr, error) {
	text, err := p.expectv(TokString)
	if err != nil {
		return nil, err
	}
	return &AstStringLitExpr{Value: text}, nil
}

func (p *Parser) ParseType() (*AstType, error) {
	text, err := p.expectv(TokIdent)
	if err != nil {
		return nil, err
	}
	return &AstType{Name: text}, nil
}

func (p *Parser) ParseFnDecl() (*AstFnDecl, error) {
	if err := p.expect(TokFn); err != nil {
		return nil, err
	}
	fnName, err := p.expectv(TokIdent)
	if err != nil {
		return nil, err
	}
	if err := p.expect(TokLpar); err != nil {
		return nil, err
	}
	params := []*AstParam{}
	for {
		if p.peekIs(TokRpar) {
			break
		} else if p.peekIs(TokIdent) {
			param, err := p.ParseParam()
			if err != nil {
				return nil, err
			}
			params = append(params, param)
			// Consume the comma if it is there
			if p.peekIs(TokComma) {
				p.nextToken()
			}
		} else {
			fmt.Println(p.peek())
			return nil, p.parseError()
		}
	}
	//if i == 99 {
	//	panic("blew the loop")
	//}
	if err := p.expect(TokRpar); err != nil {
		return nil, err
	}
	if err := p.expect(TokColon); err != nil {
		return nil, err
	}
	returnType, err := p.ParseType()
	if err != nil {
		return nil, err
	}
	if err := p.expect(TokLbrace); err != nil {
		return nil, err
	}
	if err := p.expect(TokRbrace); err != nil {
		return nil, err
	}
	ret := AstFnDecl{
		Name:       fnName,
		ReturnType: returnType,
		Params:     params,
		Body:       []AstStatement{},
	}
	return &ret, nil
}

func (p *Parser) ParseParam() (*AstParam, error) {
	text, err := p.expectv(TokIdent)
	if err != nil {
		return nil, err
	}
	if err := p.expect(TokColon); err != nil {
		return nil, err
	}
	type_, err := p.ParseType()
	if err != nil {
		return nil, err
	}
	return &AstParam{Name: text, Type: type_}, nil
}

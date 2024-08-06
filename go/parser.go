package main

import (
	"fmt"
	"strconv"
)

type Parser struct {
	lexer   Lexer
	tok     Token
	nextTok Token
}

type ParseError struct {
	Msg      string
	Filename string
	Line     uint
	Col      uint
}

func NewParser(input string) Parser {
	lex := NewLexer(input)
	t0 := lex.Next()
	t1 := lex.Next()
	p := Parser{lexer: lex, tok: t0, nextTok: t1}
	return p
}

func (p *Parser) nextToken() {
	p.tok = p.nextTok
	p.nextTok = p.lexer.Next()
}

func (p *Parser) expectv(expected TokenKind) (string, error) {
	//fmt.Printf("expectv: %v, %v, %v\n", p, expected, p.tok.Kind)
	if p.tok.Kind != expected {
		return "", fmt.Errorf("parse error: expected: %v token, got %v", expected, p.tok)
	}
	val := p.tok.Text
	p.nextToken()
	return val, nil
}

func (p *Parser) expect(expected TokenKind) error {
	_, err := p.expectv(expected)
	return err
}

func (p *Parser) peekIs(expected TokenKind) bool {
	return p.tok.Kind == expected
}

func (p *Parser) parseError() error {
	return ParseError{
		Msg:      fmt.Sprintf("unexpected token: %v", p.tok.Kind),
		Filename: "todo_file.b",
		Line:     0,
		Col:      0,
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
		fmt.Printf("ll: %v\n", p.tok)
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
	if p.peekIs(TokLet) {
		return p.ParseConstAssign()
	}
	return nil, p.parseError()
}

func (p *Parser) ParseConstAssign() (AstConstAssign, error) {
	if err := p.expect(TokLet); err != nil {
		return AstConstAssign{}, err
	}
	constName, err := p.expectv(TokIdent)
	if err != nil {
		return AstConstAssign{}, err
	}
	if err := p.expect(TokAssign); err != nil {
		return AstConstAssign{}, err
	}
	intText, err := p.expectv(TokInt)
	if err != nil {
		return AstConstAssign{}, err
	}
	intVal, err := strconv.Atoi(intText)
	if err != nil {
		return AstConstAssign{}, err
	}
	if err := p.expect(TokSemi); err != nil {
		return AstConstAssign{}, err
	}
	return AstConstAssign{Ident: constName, Value: intVal}, nil
}

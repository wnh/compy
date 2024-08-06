package main

import "fmt"

type Parser struct {
	lexer   Lexer
	tok     Token
	nextTok Token
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

func (p *Parser) ParseModule() (*AstModule, error) {
	if err := p.expect(TokModule); err != nil {
		return nil, err
	}
	modName, err := p.expectv(TokIdent)
	//fmt.Println("ParseModule():", modName, err)
	if err != nil {
		return nil, err
	}
	return &AstModule{Name: modName}, nil
}

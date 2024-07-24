package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Lexer struct {
	i     int
	input string
	// f
	text string
	num  int
}

//go:generate stringer -type=Token
type Token int

const (
	TokErr Token = -1
	TokEof Token = iota
	TokInt
	TokIdent
	TokFn
	TokLpar
	TokRpar
	TokLbrace
	TokRbrace
	TokLsq
	TokRsq
	TokColon
	TokLet
	TokAssign
	TokString
	TokSemi
	TokModule
	TokGte
	TokGt
	TokIf
	TokReturn
	TokDiv
)

var keywords = map[string]Token{
	"fn": TokFn,
	"let": TokLet,
	"module": TokModule,
	"return": TokReturn,
	"if": TokIf,
}

func NewLexer(input string) Lexer {
	return Lexer{input: input, i: 0}
}

func (l *Lexer) IsEof() bool {
	return l.i >= len(l.input)
}

func (l *Lexer) char() byte {
	if l.IsEof() {
		return 0
	}
	return l.input[l.i]
}

func (l *Lexer) nextChar() byte {
	l.i += 1
	return l.char()
}

func (l *Lexer) Next() (Token, error) {
	l.text = ""
	l.num = 0
	//fmt.Printf("Input length: %d\n", len(l.input))
	if l.i == len(l.input) {
		return TokEof, nil
	}
	c := l.char()
	switch {
	case isWS(c):
		l.skipWS()
		return l.Next()
	case isAlpha(c):
		return l.TokenizeIdent()
	case isNum(c):
		return l.TokenizeInt()
	case c == '"':
		return l.TokenizeString()
	case c == '(':
		l.nextChar()
		return TokLpar, nil
	case c == ')':
		l.nextChar()
		return TokRpar, nil
	case c == '[':
		l.nextChar()
		return TokLsq, nil
	case c == ']':
		l.nextChar()
		return TokRsq, nil
	case c == '{':
		l.nextChar()
		return TokLbrace, nil
	case c == '}':
		l.nextChar()
		return TokRbrace, nil
	case c == ':':
		l.nextChar()
		return TokColon, nil
	case c == '=':
		l.nextChar()
		return TokAssign, nil
	case c == ';':
		l.nextChar()
		return TokSemi, nil
	case c == '>':
		if l.nextChar() == '=' {
			l.i += 1
			return TokGte, nil
		} else {
			return TokGt, nil
		}
	case c == '/':
		if l.nextChar() == '/' {
			l.nextChar()
			l.skipComment()
			return l.Next()
		} else {
			return TokDiv, nil
		}
	default:
		return TokErr, errors.New("this is not quite working yet")
	}
}
func isNum(c byte) bool {
	return c >= '0' && c <= '9'
}
func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func isWS(c byte) bool {
	return strings.ContainsAny(string(c), " \t\r\n")
}

func (l *Lexer) skipWS() {
	for isWS(l.char()) {
		l.i++
	}
}

func (l *Lexer) skipComment() {
	for l.nextChar() != '\n' {
	}
	l.nextChar() // Move over newline
}

func (l *Lexer) TokenizeString() (Token, error) {
	l.text=""
	for l.nextChar() != '"' {
		l.text += string(l.char())
	}
	l.nextChar() // Move over closing quote
	return TokString, nil
}

func (l *Lexer) TokenizeIdent() (Token, error) {
	l.text = ""
	for isAlpha(l.char()) || isNum(l.char()){
		l.text += string(l.char())
		l.nextChar()
	}
	if kw, ok := keywords[l.text]; ok {
		return kw, nil
	}
	return TokIdent, nil
}

func (l *Lexer) TokenizeInt() (Token, error) {
	//fmt.Println("Token int", l.i)
	l.text = ""
	for isNum(l.char()) {
		l.text = l.text + string(l.char())
		l.nextChar()
	}
	// TODO: parse text into num
	var err error
	l.num, err = strconv.Atoi(l.text)
	if err != nil {
		errors.New(fmt.Sprintf("unable to convert Int token (%#v)", l.text))
	}
	return TokInt, nil
}

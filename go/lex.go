package main

import (
	"fmt"
	"io"
	"strings"
)

type Lexer struct {
	input string
	reader *strings.Reader
	isEof bool
	current rune
}

//go:generate stringer -type=TokenKind
type TokenKind int
type Token struct {
	Kind TokenKind
	// the raw text for the toke, unless kind == TokErr then its the error text
	// TODO(wnh): keep those separate 
	Text string 
	Error error
}

const (
	TokErr TokenKind = -1
	TokEof TokenKind = iota
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

func MkToken(kind TokenKind, text string) Token {
	return Token{kind, text, nil}
}
func MkTokenErr(err error) Token {
	return Token{TokErr, "", err}
}

var keywords = map[string]TokenKind{
	"fn": TokFn,
	"let": TokLet,
	"module": TokModule,
	"return": TokReturn,
	"if": TokIf,
}

func NewLexer(input string) Lexer {
	l := Lexer{input: input, reader: strings.NewReader(input)}
	l.nextChar()
	return l
}

func (l *Lexer) nextChar() rune {
	r, runeLen, err :=  l.reader.ReadRune()
	_ = runeLen
	//fmt.Printf("nextChar(): %v %#v %#v %v \n", string(r), runeLen, err, err == io.EOF)
	if err == io.EOF {
		//fmt.Println("Im returning?")
		l.isEof = true
		l.current = 0
		return 0
	}
	if err != nil {
		 panic("bad rune read")
	}
	l.current = r
	return r
}

func (l *Lexer) char() rune {
	return l.current
}

func (l *Lexer) Next() Token {
	c := l.char()
	//fmt.Println("Next(): l.isEof", l.isEof)
	if isWS(c) {
		l.skipWS()
		c = l.char()
	}
	if l.isEof {
		return MkToken(TokEof, "")
	}
	//fmt.Printf("char: %#v\n", c)
	switch {
	case isAlpha(c):
		return l.TokenizeIdent()
	case isNum(c):
		return l.TokenizeInt()
	case c == '"':
		return l.TokenizeString()
	case c == '(':
		l.nextChar()
		return MkToken(TokLpar, "(")
	case c == ')':
		l.nextChar()
		return MkToken(TokRpar, ")")
 	case c == '[':
		l.nextChar()
		return MkToken(TokLsq, "")
	case c == ']':
		l.nextChar()
		return MkToken(TokRsq, "")
	case c == '{':
		l.nextChar()
		return MkToken(TokLbrace, "")
	case c == '}':
		l.nextChar()
		return MkToken(TokRbrace, "")
	case c == ':':
		l.nextChar()
		return MkToken(TokColon, "")
	case c == '=':
		l.nextChar()
		return MkToken(TokAssign, "")
	case c == ';':
		l.nextChar()
		return MkToken(TokSemi, "")
	case c == '>':
		if l.nextChar() == '=' {
			l.nextChar()
			return MkToken(TokGte, "")
		} else {
			return MkToken(TokGt, "")
		}
	case c == '/':
		if l.nextChar() == '/' {
			l.nextChar()
			l.skipComment()
			return l.Next()
		} else {
			return MkToken(TokDiv, "")
		}
	default:
		return MkTokenErr(fmt.Errorf("parse error: unknown character %#v", l.char()))
	}
}
func isNum(c rune) bool {
	return c >= '0' && c <= '9'
}
func isAlpha(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c  >= 'A' && c <= 'Z')
}

func isWS(c rune) bool {
	//fmt.Println("isWS():", string(c))
	return strings.ContainsAny(string(c), " \t\r\n")
}

func (l *Lexer) skipWS() {
	//fmt.Println("skipWs()")
	for isWS(l.char()) {
		l.nextChar()
	}
}

func (l *Lexer) skipComment() {
	for l.nextChar() != '\n' {}
	l.nextChar() // Move over newline
}

func (l *Lexer) TokenizeString() Token {
	txt := ""
	for l.nextChar() != '"' {
		txt += string(l.char())
	}
	l.nextChar() // Move over closing quote
	return MkToken(TokString, txt)
}

func (l *Lexer) TokenizeIdent() Token {
	val := ""
	for isAlpha(l.char()) || isNum(l.char()){
		val += string(l.char())
		l.nextChar()
	}
	if kwKind, ok := keywords[val]; ok {
		return MkToken(kwKind, val)
	}
	return MkToken(TokIdent, val)
}

func (l *Lexer) TokenizeInt() Token {
	//fmt.Println("Token int", l.i)
	val := ""
	for isNum(l.char()) {
		val = val + string(l.char())
		l.nextChar()
	}
	return MkToken(TokInt, val)
}

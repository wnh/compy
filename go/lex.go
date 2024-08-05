package main

import (
	"fmt"
	"strconv"
	"io"
	"strings"
)

type Lexer struct {
	input string
	reader *strings.Reader
	isEof bool
	current rune
	// f
	textValue string
	numValue  int
	Error error
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
	l.textValue = ""
	l.numValue = 0
	c := l.char()
	//fmt.Println("Next(): l.isEof", l.isEof)
	if isWS(c) {
		l.skipWS()
		c = l.char()
	}
	if l.isEof {
		return TokEof
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
		return TokLpar
	case c == ')':
		l.nextChar()
		return TokRpar
	case c == '[':
		l.nextChar()
		return TokLsq
	case c == ']':
		l.nextChar()
		return TokRsq
	case c == '{':
		l.nextChar()
		return TokLbrace
	case c == '}':
		l.nextChar()
		return TokRbrace
	case c == ':':
		l.nextChar()
		return TokColon
	case c == '=':
		l.nextChar()
		return TokAssign
	case c == ';':
		l.nextChar()
		return TokSemi
	case c == '>':
		if l.nextChar() == '=' {
			l.nextChar()
			return TokGte
		} else {
			return TokGt
		}
	case c == '/':
		if l.nextChar() == '/' {
			l.nextChar()
			l.skipComment()
			return l.Next()
		} else {
			return TokDiv
		}
	default:
		l.Error =  fmt.Errorf("this is not quite working yet")
		l.textValue = l.Error.Error()
		return TokErr 
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
	l.textValue=""
	for l.nextChar() != '"' {
		l.textValue += string(l.char())
	}
	l.nextChar() // Move over closing quote
	return TokString
}

func (l *Lexer) TokenizeIdent() Token {
	l.textValue = ""
	for isAlpha(l.char()) || isNum(l.char()){
		l.textValue += string(l.char())
		l.nextChar()
	}
	if kw, ok := keywords[l.textValue]; ok {
		return kw
	}
	return TokIdent
}

func (l *Lexer) TokenizeInt() Token {
	//fmt.Println("Token int", l.i)
	l.textValue = ""
	for isNum(l.char()) {
		l.textValue = l.textValue + string(l.char())
		l.nextChar()
	}
	// TODO: parse text into num
	var err error
	l.numValue, err = strconv.Atoi(l.textValue)
	if err != nil {
		l.Error = fmt.Errorf("unable to convert Int token (%#v)", l.textValue)
		l.textValue = l.Error.Error()
		return TokErr
	}
	return TokInt
}

package main

import (
	"fmt"
	"io"
	"strings"
)

type Lexer struct {
	input   string
	reader  io.RuneReader
	isEof   bool
	current rune

	// Where we are in the RuneReader
	line uint
	col  uint
	// Location of the token being currently parsed
	startLine uint
	startCol  uint
}

//go:generate stringer -type=TokenKind
type TokenKind int
type Token struct {
	Kind TokenKind
	// the raw text for the toke, unless kind == TokErr then its the error text
	// TODO(wnh): keep those separate
	Text  string
	Error error
	Line  uint
	Col   uint
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
	TokComma
	TokModule
	TokGte
	TokGt
	TokIf
	TokReturn
	TokDiv
)

func (l *Lexer) MkToken(kind TokenKind, text string) Token {
	return Token{kind, text, nil, l.startLine, l.startCol}
}
func (l *Lexer) MkTokenErr(err error) Token {
	return l.MkToken(TokErr, "")
}

var keywords = map[string]TokenKind{
	"fn":     TokFn,
	"let":    TokLet,
	"module": TokModule,
	"return": TokReturn,
	"if":     TokIf,
}

func NewLexer(input string) Lexer {
	l := Lexer{input: input, reader: strings.NewReader(input), line: 1, col: 0}
	l.nextChar()
	return l
}

func (l *Lexer) nextChar() rune {
	r, runeLen, err := l.reader.ReadRune()
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
	if r == '\n' {
		l.line++
		l.col = 0
	} else {
		l.col++
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
		return l.MkToken(TokEof, "")
	}
	l.startLine = l.line
	l.startCol = l.col

	switch {
	case isAlpha(c) || c == '_':
		return l.TokenizeIdent()
	case isNum(c):
		return l.TokenizeInt()
	case c == '"':
		return l.TokenizeString()
	case c == '(':
		l.nextChar()
		return l.MkToken(TokLpar, "(")
	case c == ')':
		l.nextChar()
		return l.MkToken(TokRpar, ")")
	case c == '[':
		l.nextChar()
		return l.MkToken(TokLsq, "")
	case c == ']':
		l.nextChar()
		return l.MkToken(TokRsq, "")
	case c == '{':
		l.nextChar()
		return l.MkToken(TokLbrace, "")
	case c == '}':
		l.nextChar()
		return l.MkToken(TokRbrace, "")
	case c == ':':
		l.nextChar()
		return l.MkToken(TokColon, "")
	case c == '=':
		l.nextChar()
		return l.MkToken(TokAssign, "")
	case c == ';':
		l.nextChar()
		return l.MkToken(TokSemi, "")
	case c == ',':
		l.nextChar()
		return l.MkToken(TokComma, "")
	case c == '>':
		if l.nextChar() == '=' {
			l.nextChar()
			return l.MkToken(TokGte, "")
		} else {
			return l.MkToken(TokGt, "")
		}
	case c == '/':
		if l.nextChar() == '/' {
			l.nextChar()
			l.skipComment()
			return l.Next()
		} else {
			return l.MkToken(TokDiv, "")
		}
	default:
		return l.MkTokenErr(fmt.Errorf("parse error: unknown character %#v", l.char()))
	}
}
func isNum(c rune) bool {
	return c >= '0' && c <= '9'
}
func isAlpha(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
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
	for l.nextChar() != '\n' {
	}
	l.nextChar() // Move over newline
}

func (l *Lexer) TokenizeString() Token {
	txt := ""
	for l.nextChar() != '"' {
		txt += string(l.char())
	}
	l.nextChar() // Move over closing quote
	return l.MkToken(TokString, txt)
}

func (l *Lexer) TokenizeIdent() Token {
	val := ""
	for isAlpha(l.char()) || isNum(l.char()) || l.char() == '_' {
		val += string(l.char())
		l.nextChar()
	}
	if kwKind, ok := keywords[val]; ok {
		return l.MkToken(kwKind, val)
	}
	return l.MkToken(TokIdent, val)
}

func (l *Lexer) TokenizeInt() Token {
	//fmt.Println("Token int", l.i)
	val := ""
	for isNum(l.char()) {
		val = val + string(l.char())
		l.nextChar()
	}
	return l.MkToken(TokInt, val)
}

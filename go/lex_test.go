
package main

import (
	"testing"
)

func TestNextToken(t *testing.T) {
	type testcase struct { n string; input string; tokens []Token }
	tests := []testcase{
		testcase{ n: "Integer", input: "123", tokens: []Token{TokInt} },
		testcase{ n: "Trailing space", input: "123 ", tokens: []Token{TokInt} },
		testcase{ n: "Leading space", input: "   123", tokens: []Token{TokInt} },
		testcase{ n: "Leading+Trailing space", input: "   123", tokens: []Token{TokInt} },
		testcase{ n: "Two Ints", input: "123 456", tokens: []Token{TokInt, TokInt} },
		testcase{ n: "Ident", input: "foo", tokens: []Token{TokIdent} },
		testcase{ n: "keyword and Ident", input: "fn foo", tokens: []Token{TokFn, TokIdent} },
		testcase{ n: "Fn", input: "fn foo(a : int)", tokens: []Token{TokFn, TokIdent, TokLpar, TokIdent, TokColon, TokIdent, TokRpar} },
	}

	for _, tc := range tests {
		t.Run(tc.n, func (t *testing.T) {
			l := NewLexer(tc.input)
			for _, expected := range tc.tokens {
				got := l.Next()
				if got != expected {
					t.Errorf("Expected %v got %v", expected, got)
				}
			}
			lastTok := l.Next()
			if lastTok != TokEof {
				t.Errorf("Expected EOF got %v", lastTok)
			}
		})
	}
}


package main

import (
	"testing"
)

func TestNextToken(t *testing.T) {
	type testcase struct { n string; input string; tokens []TokenKind }
	tests := []testcase{
		testcase{ n: "Integer", input: "123", tokens: []TokenKind{TokInt} },
		testcase{ n: "Trailing space", input: "123 ", tokens: []TokenKind{TokInt} },
		testcase{ n: "Leading space", input: "   123", tokens: []TokenKind{TokInt} },
		testcase{ n: "Leading+Trailing space", input: "   123", tokens: []TokenKind{TokInt} },
		testcase{ n: "Two Ints", input: "123 456", tokens: []TokenKind{TokInt, TokInt} },
		testcase{ n: "Ident", input: "foo", tokens: []TokenKind{TokIdent} },
		testcase{ n: "Ident with underscore", input: "foo_bar", tokens: []TokenKind{TokIdent} },
		testcase{ n: "Ident with leading underscore", input: "_foobar", tokens: []TokenKind{TokIdent} },
		testcase{ n: "keyword and Ident", input: "fn foo", tokens: []TokenKind{TokFn, TokIdent} },
		testcase{ n: "Fn", input: "fn foo(a : int)", tokens: []TokenKind{TokFn, TokIdent, TokLpar, TokIdent, TokColon, TokIdent, TokRpar} },
		testcase{ n: "Two char tokens", input: "> = >= >==", tokens: []TokenKind{TokGt, TokAssign, TokGte, TokGte, TokAssign } },
	}

	for _, tc := range tests {
		t.Run(tc.n, func (t *testing.T) {
			l := NewLexer(tc.input)
			for _, expected := range tc.tokens {
				got := l.Next()
				if got.Kind != expected {
					t.Errorf("Expected %v got %v", expected, got)
				}
			}
			lastTok := l.Next()
			if lastTok.Kind != TokEof {
				t.Errorf("Expected EOF got %v", lastTok)
			}
		})
	}
}



package main

import (
	"testing"
)

func TestNextToken(t *testing.T) {
	type testcase struct { input string; tokens []Token }
	tests := []testcase{
		testcase{ input: "123 ", tokens: []Token{TokInt} },
		testcase{ input: "123", tokens: []Token{TokInt} },
		testcase{ input: "123 456", tokens: []Token{TokInt, TokInt} },
		testcase{ input: "foo", tokens: []Token{TokIdent} },
		testcase{ input: "fn foo", tokens: []Token{TokFn, TokIdent} },
		testcase{ input: "fn foo(a : int)", tokens: []Token{TokFn, TokIdent, TokLpar, TokIdent, TokColon, TokIdent, TokRpar} },
	}

	for _, tc := range tests {
		t.Run(tc.input, func (t *testing.T) {
			l := NewLexer(tc.input)
			for _, tok := range tc.tokens {
				tt, err := l.Next()
				//fmt.Printf("token %v\n", tt)
				if err != nil {
					t.Fail()
				}
				if tt != tok {
					t.Fail()
				}
			}
			tt, err := l.Next()
			if err != nil {
				t.Fail()
			}
			if tt != TokEof {
				t.Fail()
			}
		})
	}
}
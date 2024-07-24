package main

import (
	"fmt"
	"os"
)


func main() {
	//fmt.Println("hey")
	lex := NewLexer(`
	module main;

	fn main(args : []string): int {
		let thing = print("hey  there");
		let out = if len(args) >= 2 {
			"too many args"
		} else {
			"something else"
		}
		{
			let x = 123; // random block
		}
		return 0;
	}
	`)
	for range 100 {
		t, err := lex.Next()
		if err != nil {
			fmt.Printf("token error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("got token: %v: %#v\n", t, lex.text)
		if t == TokEof {
			break
		}
	}
}

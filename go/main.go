package main

import (
	"fmt"
	//"os"
)

func main() {
	//fmt.Println("hey")
	txt := `module main;

	fn main(args: []string): int 
        {
		let thing: string = "hey  there";
		let out = if len(args) >= 2 {
			"too many args"
		} else {
			"something else"
		}

		{
			let x: int = 123; // random block
		}
		return 0;
	}
	`
	lex := NewLexer(txt)

	fmt.Println("Tokens for an example program")
	for i:=0; i<100; i++  {
		tok := lex.Next()
		fmt.Printf("  %v (%#v) \n", tok.Kind, tok.Text)
		if tok.Kind == TokEof {
			fmt.Println("Ended Successfully")
			break
		} else if tok.Kind == TokErr {
			fmt.Printf("token error: %v\n", tok.Error)
			break
		}
	}

}

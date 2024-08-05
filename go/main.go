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
	_ = txt

	//lex := NewLexer(" 123x ")
	//lex := NewLexer(" foo  bar")
	lex := NewLexer("foo ")
	//lex := NewLexer(txt)

	for i:=0; i<100; i++  {
		//fmt.Println("TOP:", i);
		tok := lex.Next()
		fmt.Printf("  got: %v (%#v) \n", tok, lex.textValue)
		if tok == TokEof {
			fmt.Println("Ended successfully?")
			break
		} else if tok == TokErr {
			fmt.Printf("token error: %v\n", lex.Error)
			break
		}
		//fmt.Printf("%v\n",  tok)
	}

}

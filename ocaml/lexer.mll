{ open Parser }

let digit = ['0'-'9']
let alpha = ['a'-'z''A'-'Z']
let alphanum = alpha | digit
let ident = alpha alphanum+
let ws = [' ' '\n' '\t' '\r']

rule token = parse 
     | ws+ { token lexbuf }
     | digit+ as n { INT (int_of_string n) }
     | '{' { LBRACE }
     | '}' { RBRACE }
     | '(' { LPAREN }
     | ')' { RPAREN }
     | '+' { PLUS }
     | '-' { MINUS }
     | "return" { RETURN }
     | "fun" { FUN }
     | ';' { SEMI }
     (* Make sure this is at the bottom *)
     | ident as i { IDENT i }
     | eof { EOF }
{ open Parser }

let digit = ['0'-'9']
let alpha = ['a'-'z''A'-'Z']
let alphanum = alpha | digit
let ident = alpha alphanum+
let ws = [' ' '\n' '\t' '\r']

rule token = parse 
     | ws+ { token lexbuf }
     | digit+ as i { INT (int_of_string i) }
     | ident { IDENT }
     | '=' { ASSIGN }
     | "==" { EQ }
     | "!=" { NEQ }
     | '<' { LT }
     | "<=" { LTE }
     | '>' { GT }
     | ">=" { GTE }
     | '+' { PLUS }
     | '-' { MINUS }
     | '*' { STAR }
     | '/' { DIV }
     | '&' { AMP }
     | '{' { LBRACE }
     | '}' { RBRACE }
     | '(' { LPAREN }
     | ')' { RPAREN }
     | "else" { ELSE }
     | "if" { IF }
     | "while" { WHILE }
     | "for" { FOR }
     | "return" { RETURN }
     | ';' { SEMI }
     | eof { EOF }
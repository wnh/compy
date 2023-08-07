%{
open Ast
%}
%token <int> INT

%start <Ast.block> main
%token ELSE IF WHILE FOR RETURN SEMI LBRACE RBRACE LPAREN RPAREN
%token IDENT STAR  LT MINUS GT LTE GTE MUL PLUS EQ ASSIGN AMP NEQ DIV
%token EOF 

%right PLUS MINUS
%right STAR DIV

%%

main:
  | block EOF { $1 }
;

block:
  | LBRACE stmt+ RBRACE { $2 }
; 

stmt:
  | simple_stmt SEMI { $1 }
;

simple_stmt:
  | INT { Number $1 }
  | RETURN { Return }
;
(*
  | assign { FixmeStatement "simple:assign" }
  | expr { FixmeStatement "simple:expr" }
*)

(* 
  | block { FixmeStatement "block" }
  | if_stmt { FixmeStatement "if" }
  | for_stmt { FixmeStatement "for" }
  | while_stmt { FixmeStatement "while" }
  ;
if_stmt:
  | IF expr stmt else_stmt?  {}
;
else_stmt:
  | ELSE stmt
    {}

for_stmt:
| FOR LPAREN simple_stmt SEMI simple_stmt SEMI simple_stmt RPAREN stmt {}
;

while_stmt: WHILE LPAREN expr RPAREN stmt {} ;



assign: ident ASSIGN expr {};

expr:
  | unary  { -1 }
  | deref  { -1 }
  | var_ref  { -1 }
  | INT  { $1 }
  | expr binop expr  { -1 }
  | LPAREN expr RPAREN { $2 }
;
binop:
  | EQ  { EQ  }
  | NEQ { NEQ }
  | LT { LT }
  | LTE { LTE }
  | GT { GT }
  | GTE { GTE }
  | PLUS { PLUS }
  | MINUS { MINUS }
  | MUL { MUL }
  | DIV { DIV }
;

var_ref: ident {}; 

ident: IDENT {} ;

unary:
  | uneg {}
  | addr_of {}
  | deref {}
  | uplus {}

uneg:  MINUS expr {};
addr_of:  AMP expr {};
deref:  STAR expr {};
uplus:  PLUS expr {};

*)

%{
open Ast
%}

%start <Ast.program> program
%token <int> INT
%token <string> IDENT
%token FUN LPAREN RPAREN RETURN LBRACE RBRACE SEMI PLUS MINUS MULT DIV
%token EOF 

%left PLUS MINUS
%left MULT DIV
%nonassoc NEGATE

%%

program:
  | func_def EOF { MainFunc $1 }
;

func_def:
  | FUN IDENT LPAREN RPAREN block { FuncDef {name = $2; code = $5} }
;

block:
  | LBRACE statement* RBRACE { $2 }
;

statement:
  | RETURN expr SEMI { Return $2 }
;

expr:
  | expr PLUS expr { BinOp (Plus, $1, $3) }
  | expr MINUS expr { BinOp (Minus, $1, $3) }
  | expr MULT expr { BinOp (Mult, $1, $3) }
  | expr DIV expr { BinOp (Div, $1, $3) }
  | INT { Integer $1 }
  | MINUS expr %prec NEGATE { Neg $2 }
  | PLUS expr %prec NEGATE { $2 }
  | LPAREN expr RPAREN { $2 }
;


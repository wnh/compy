%{
open Ast
%}

%start <Ast.program> program
%token <int> INT
%token <string> IDENT
%token FUN LPAREN RPAREN RETURN LBRACE RBRACE SEMI
%token EOF 

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
  | INT { Integer $1 }
;


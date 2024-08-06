%{
open Ast
%}

%start <Ast.program> program
%token <int> INT
%token <string> IDENT
%token FUN LPAREN RPAREN RETURN LBRACE RBRACE SEMI PLUS MINUS MULT DIV
%token EQ NEQ GT GTE LT LTE ASSIGN IF ELSE
%token EOF

%left EQ NEQ GT GTE LT LTE
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
  | IDENT ASSIGN expr SEMI { Assign ($1, $3) }
  | RETURN expr SEMI { Return $2 }
  | expr SEMI { ExprStmt $1 }
  | IF LPAREN cond=expr LPAREN then_blk=expr
    { IfStmt { cond; then_blk; else_blk=EmptyExpr } }
  | SEMI { EmptyStmt }
;

expr:
  | expr PLUS expr { BinOp (Plus, $1, $3) }
  | expr MINUS expr { BinOp (Minus, $1, $3) }
  | expr MULT expr { BinOp (Mult, $1, $3) }
  | expr DIV expr { BinOp (Div, $1, $3) }
  | expr EQ expr { BinOp (Eq, $1, $3) }
  | expr NEQ expr { BinOp (Neq, $1, $3) }
  | expr GT expr { BinOp (Gt, $1, $3) }
  | expr GTE expr { BinOp (Gte, $1, $3) }
  | expr LT expr { BinOp (Lt, $1, $3) }
  | expr LTE expr { BinOp (Lte, $1, $3) }
  | block { BlockExpr $1 }
  | INT { Integer $1 }
  | IDENT { VarRef $1 }
  | MINUS expr %prec NEGATE { Neg $2 }
  | PLUS expr %prec NEGATE { $2 }
  | LPAREN expr RPAREN { $2 }
;


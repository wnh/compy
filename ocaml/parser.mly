%token <int> INT
%token ELSE IF WHILE FOR RETURN SEMI LBRACE RBRACE LPAREN RPAREN
%token IDENT STAR  LT MINUS GT LTE GTE MUL PLUS EQ ASSIGN AMP NEQ DIV
%token EOF 

%start <Ast.block> main

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
  | simple_stmt SEMI { FixmeStatement }
  | return SEMI { Return }
  | block { FixmeStatement }
  | if_stmt { FixmeStatement }
  | for_stmt { FixmeStatement }
  | while_stmt { Ast.FixmeStatement }
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

simple_stmt:
  | assign { -1 }
  | expr { $1 }
  | empty_stmt { $1 }
;

empty_stmt: SEMI { 0 };

assign: ident ASSIGN expr {};
return: RETURN expr {};

expr:
  | unary  { -1 }
  | deref  { -1 }
  | var_ref  { -1 }
  | INT  { $1 }
  | expr binop expr  { -1 }
  | LPAREN expr RPAREN { $2 }
;
binop:
  |  EQ  {} 
  |  NEQ  {} 
  |  LT  {} 
  |  LTE  {} 
  |  GT  {} 
  |  GTE  {} 
  |  PLUS  {} 
  |  MINUS  {} 
  |  MUL  {} 
  |  DIV {}
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

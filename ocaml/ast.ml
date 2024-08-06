
open Sexplib.Std

type program = MainFunc of func_def
[@@deriving sexp_of]

and func_def = FuncDef of { name : ident; code : block}
[@@deriving sexp_of]

and ident = string
[@@deriving sexp_of]

and block = statement list
[@@deriving sexp_of]

and statement =
  | Assign of (ident * expr)
  | Return of expr
  | ExprStmt of expr
  | IfStmt of { cond: expr; then_blk: expr; else_blk: expr }
  | EmptyStmt
[@@deriving sexp_of]

and expr =
  | BinOp of (operation * expr * expr)
  | BlockExpr of block
  | Integer of int
  | Neg of expr
  | VarRef of string
  | EmptyExpr
[@@deriving sexp_of]

and operation = Plus | Minus | Mult | Div | Eq | Neq | Gt | Gte | Lt | Lte
[@@deriving sexp_of]

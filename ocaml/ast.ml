
open Sexplib.Std

type program = MainFunc of func_def
[@@deriving sexp_of]

and func_def = FuncDef of { name : ident; code : block}
[@@deriving sexp_of]

and ident = string
[@@deriving sexp_of]

and block = statement list
[@@deriving sexp_of]

and statement = Return of expr
[@@deriving sexp_of]

and expr = Integer of int 
[@@deriving sexp_of]

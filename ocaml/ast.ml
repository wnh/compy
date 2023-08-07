
open Sexplib.Std

type block =  statement list
[@@deriving sexp_of]

and statement = Number of int | Return
[@@deriving sexp_of]

and binop =
  |  EQ
  |  NEQ
  |  LT
  |  LTE
  |  GT
  |  GTE
  |  PLUS
  |  MINUS
  |  MUL
  |  DIV
[@@deriving sexp_of]

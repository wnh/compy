
open Sexplib.Std

type block =  statement list
[@@deriving sexp_of]

and statement = Return | FixmeStatement
[@@deriving sexp_of]
             

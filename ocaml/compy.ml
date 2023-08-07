
let _ =
  print_endline "Starting?";
  let lexbuf = Lexing.from_string "{ return }" in
  let ast = Parser.main Lexer.token lexbuf in 
  print_endline (Sexplib.Sexp.to_string (Ast.sexp_of_block ast));
  print_endline "Ending?";


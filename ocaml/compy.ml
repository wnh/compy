
let fmt_position (lexbuf: Lexing.lexbuf) =
  let pos = lexbuf.lex_curr_p in
  Printf.sprintf "%s:%d:%d" pos.pos_fname
    pos.pos_lnum (pos.pos_cnum - pos.pos_bol + 1)

let parse input_string =
  let lexbuf = Lexing.from_string input_string in
  (* let ast = Parser.main Lexer.token lexbuf in *)
  let ast =
    try Parser.main Lexer.token lexbuf with
    | Parser.Error ->
       Printf.printf "%s -> %s There was an problem\n" input_string (fmt_position lexbuf);
       (* Printf.fprintf stderr "%a: syntax error\n" print_position lexbuf; *)
       exit (-1)
  in
  Printf.printf "%s -> %s\n" input_string (Sexplib.Sexp.to_string (Ast.sexp_of_block ast))

let _ =
  parse "{ 12; }";
  parse "{ 12; 13; }";
  parse "{21;return;}";
  parse "{ return; }";

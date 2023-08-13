(*
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

*)

let usage_message = "compy_c [-o outfile] <file1> [<file2>] ..."
let files = ref []
let outfile = ref ""
let verbose = ref false
let speclist =
  [
    ("-verbose", Arg.Set verbose, "Output debug information");
    ("-o", Arg.Set_string outfile, "Set output file name");
  ]
let add_src filename = files := filename :: !files

let fail msg =
  print_endline msg; 
  exit 1

let get_args = 
  let () = Arg.parse speclist add_src usage_message in
  if List.length !files != 1  then 
    fail "Only a single src file for now...."
  else
    ((List.nth !files 0), !outfile)

let _ =
  let (_src_fname, out_fname) = get_args in
  let out_ch = open_out out_fname in
  (* let src_ch = open_in src_fname in *)
  Printf.fprintf out_ch "export function w $main() {\n";
  Printf.fprintf out_ch "@start\n";
  Printf.fprintf out_ch "        ret 0\n";
  Printf.fprintf out_ch "}\n";


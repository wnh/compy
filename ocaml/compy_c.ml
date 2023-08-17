let fmt_position (lexbuf: Lexing.lexbuf) =
  let pos = lexbuf.lex_curr_p in
  Printf.sprintf "%s:%d:%d" pos.pos_fname
    pos.pos_lnum (pos.pos_cnum - pos.pos_bol + 1)

let parse filename input =
  let lexbuf = Lexing.from_channel input in
  Lexing.set_filename lexbuf filename;
  try Parser.program Lexer.token lexbuf with
  | Parser.Error ->
     Printf.printf "parse_error: %s\n" (fmt_position lexbuf);
     exit 127


module CompyArgs = struct
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
end

module CodeGen = struct
  open Ast
  let emit out str = Printf.fprintf out "%s\n" str


  let gen_expr out expr = match expr with
    | Integer i -> emit out (string_of_int i)

  let gen_stmt out stmt = match stmt with
    | Return e ->
       Printf.fprintf out "ret "; gen_expr out e 

  let gen_block out stmts =
    ignore (List.map (fun s -> gen_stmt out s) stmts )

  let gen_func_def out funcdef = match funcdef with
    | FuncDef {name; code} -> 
       Printf.fprintf out "export function w $%s() {\n"  name;
       emit out "@start";
       gen_block out code;
       emit out "}"

  let generate out prg = match prg with
    | MainFunc func -> gen_func_def out func
end

let _ =
  let (src_fname, out_fname) = CompyArgs.get_args in
  let input_ch = open_in src_fname in
  let out_ch = open_out out_fname in
  let ast = parse src_fname input_ch in
  CodeGen.generate out_ch ast


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
  type context =
    { out: out_channel;
      local: int ref;
    }
  let emit ctx str = Printf.fprintf ctx.out "%s\n" str
  let next_local ctx = ctx.local := !(ctx.local) + 1; !(ctx.local)

  let rec gen_bin_op ctx op left right =
    let ln = gen_expr ctx left in
    let rn = gen_expr ctx right in
    let inst = match op with
      | Plus -> "add"
      | Minus -> "sub"
      | Mult -> "mul"
      | Div -> "div"
    in
    let n = next_local ctx in
    Printf.fprintf ctx.out "\t%%.%d =w %s %%.%d, %%.%d\n" n inst ln rn;
    n

  and gen_expr ctx expr = match expr with
    | Integer i ->
       let n = next_local ctx in
       Printf.fprintf ctx.out "\t%%.%d =w add 0, %d\n" n i;
       n
    | BinOp (op, left, right) -> gen_bin_op ctx op left right

  let gen_stmt ctx stmt = match stmt with
    | Return e ->
       let r = gen_expr ctx e in
       ignore (Printf.fprintf ctx.out "\tret %%.%d\n" r)

  let gen_block ctx stmts =
    ignore (List.map (fun s -> gen_stmt ctx s) stmts )

  let gen_func_def ctx funcdef = match funcdef with
    | FuncDef {name; code} -> 
       Printf.fprintf ctx.out "export function w $%s() {\n"  name;
       emit ctx "@start";
       gen_block ctx code;
       emit ctx "}"

  let generate out prg =
    let ctx = { out; local = ref 0 } in
    match prg with
    | MainFunc func -> gen_func_def ctx func
end

let _ =
  let (src_fname, out_fname) = CompyArgs.get_args in
  let input_ch = open_in src_fname in
  let out_ch = open_out out_fname in
  let ast = parse src_fname input_ch in
  CodeGen.generate out_ch ast


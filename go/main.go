package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	fmt.Printf("%#v\n", os.Args)
	if len(os.Args) != 2 {
		fmt.Println("usage: compy <filename>")
		os.Exit(1)
	}
	filename := os.Args[1]
	file, err := os.OpenFile(filename, os.O_RDONLY, 0644)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
	content, err := io.ReadAll(file)
	if err != nil {
		fmt.Printf("error reading file \"%s\": %v\n", filename, err)
		os.Exit(1)
	}
	fmt.Println("======= Src File content =======")
	fmt.Println(string(content))

	parser := NewParser(string(content), filename)
	mod, err := parser.ParseModule()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("======= Ast =======")
	fmt.Printf("%#v\n", mod)

	codeMod := &CodegenModule{}
	mod.Codegen(codeMod)
	fmt.Println("======= Module Output =======")
	fmt.Println(codeMod.Code.String())
}

package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

func main() {
	fmt.Printf("%#v\n", os.Args)
	if len(os.Args) != 2 {
		fmt.Println("usage: compy <filename>")
		os.Exit(1)
	}
	filename := os.Args[1]
	srcbase := path.Base(filename)
	if !strings.HasSuffix(srcbase, ".b") {
		fmt.Printf("Error: bad filename pattern '%s', need it to end in '.b'", filename)
		os.Exit(1)
	}
	basename, _ := strings.CutSuffix(srcbase, ".b")

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
	cCode := codeMod.Code.String()
	fmt.Println(cCode)

	f, err := os.CreateTemp("", basename+"_*.c")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("C source: %s\n", f.Name())
	defer os.Remove(f.Name()) // clean up

	if _, err := f.Write([]byte(cCode)); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
	compileCmd := exec.Command("cc", "-o", "example", f.Name())
	fmt.Println("xx", compileCmd)
	ccOut, err := compileCmd.CombinedOutput()
	fmt.Println("======= CC Output =======")
	fmt.Println(string(ccOut))
	if err != nil {
		fmt.Printf("Error running CC: %s\n", err)
		os.Exit(1)
	}
	fmt.Println("======== DONE ========")
}

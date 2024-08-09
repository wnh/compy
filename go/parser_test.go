package main

import (
	"testing"
)

func TestNewParser(t *testing.T) {
	_ = NewParser("", "<filename>")
}

func TestParseEmptyModule(t *testing.T) {
	//t.Skip("need to refactor first")
	mods := []struct{ m, n string }{
		{"module main;", "main"},
		{"module other;", "other"},
		{"module more_stuff;", "more_stuff"},
	}
	for _, modTest := range mods {
		p := NewParser(modTest.m, "<filename>")
		//t.Logf("%#v, %v, %v\n", p, p.tok, p.nextTok);
		mod, err := p.ParseModule()

		t.Logf("ERR: %v", err)
		t.Logf("MOD: %#v", mod)
		if err != nil || mod == nil {
			t.Fatal()
		}
		if mod.Name.Name != modTest.n {
			t.Errorf("bad module name: expected %#v got %#v", modTest.n, mod.Name)
		}
	}
}

func TestParseModuleWithConsts(t *testing.T) {
	txt := `
		module test;
		
		let foo: int = 12;
		let bar: string = "more";
	`
	p := NewParser(txt, "<filename>")
	mod, err := p.ParseModule()
	t.Logf("%v", err)
	t.Logf("%+v", mod)
	if err != nil || mod == nil {
		t.Fatal()
	}
	if mod.Name.Name != "test" {
		t.Errorf("bad module name: expected %#v got %#v", "test", mod.Name)
	}
	if len(mod.Statements) != 2 {
		t.Fatal("wrong number of statements")
	}
}

func TestParseFailures(t *testing.T) {
	badCases := []string{
		"",
		"module foo; let a = 12;",
		"module foo; let a: int = ",
	}
	var p Parser
	for _, mod := range badCases {
		p = NewParser(mod, "<filename")
		_, err := p.ParseModule()
		if err == nil {
			t.Errorf("expected failure parsing: %#v", mod)
		}
	}
}

func TestParseEmptyFunction(t *testing.T) {
	txt := `
		module test;

		fn main(): int {}
		fn bar(arg: string): int {}
		fn baz(arg: string,): int {}
		fn more(arg: string, another: int): string {}
	`
	p := NewParser(txt, "<filename>")
	mod, err := p.ParseModule()
	t.Logf("ERR: %v", err)
	t.Logf("Module: %+v", mod)
	if err != nil || mod == nil {
		t.Fatal()
	}
	if len(mod.Statements) != 4 {
		t.Fatal("wrong number of statements")
	}
}

func TestParseFunction(t *testing.T) {
	txt := `
		module test;

		fn main(): int {
			let x: string = "thing";
		}
	`
	p := NewParser(txt, "<filename>")
	mod, err := p.ParseModule()
	t.Logf("ERR: %v", err)
	t.Logf("Module: %+v", mod)
	if err != nil || mod == nil {
		t.Fatal()
	}
	if len(mod.Statements) != 1 {
		t.Fatal("wrong number of statements")
	}
}

func TestParseFunctionCalls(t *testing.T) {
	txt := `
		module test;

		fn main(): int {
			let x: string = "thing";
			println("more");
		}
	`
	p := NewParser(txt, "<filename>")
	mod, err := p.ParseModule()
	t.Logf("ERR: %v", err)
	t.Logf("Module: %+v", mod)
	if err != nil || mod == nil {
		t.Fatal()
	}
	if len(mod.Statements) != 1 {
		t.Fatal("wrong number of statements")
	}
}

func TestParseVarRefs(t *testing.T) {
	txt := `
		module test;

		fn main(): int {
			let x: string = "thing";
			println("more", x);
		}
	`
	p := NewParser(txt, "<filename>")
	mod, err := p.ParseModule()
	t.Logf("ERR: %v", err)
	t.Logf("Module: %+v", mod)
	if err != nil || mod == nil {
		t.Fatal()
	}
	if len(mod.Statements) != 1 {
		t.Fatal("wrong number of statements")
	}
}

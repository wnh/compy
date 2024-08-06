package main

import (
	"testing"
)

func TestNewParser(t *testing.T) {
	_ = NewParser("")
}

func TestParseEmptyModule(t *testing.T) {
	//t.Skip("need to refactor first")
	mods := []struct{ m, n string }{
		{"module main;", "main"},
		{"module other;", "other"},
		{"module more_stuff;", "more_stuff"},
	}
	for _, modTest := range mods {
		p := NewParser(modTest.m)
		//t.Logf("%#v, %v, %v\n", p, p.tok, p.nextTok);
		mod, err := p.ParseModule()

		t.Logf("ERR: %v", err)
		t.Logf("MOD: %#v", mod)
		if err != nil || mod == nil {
			t.Fatal()
		}
		if mod.Name != modTest.n {
			t.Errorf("bad module name: expected %#v got %#v", modTest.n, mod.Name)
		}
	}
}

func TestParseModuleWithConsts(t *testing.T) {
	txt := `
		module test;
		
		let foo = 12;
		let bar = "more";
	`
	p := NewParser(txt)
	mod, err := p.ParseModule()
	t.Logf("%v", err)
	t.Logf("%+v", mod)
	if err != nil || mod == nil {
		t.Fatal()
	}
	if mod.Name != "test" {
		t.Errorf("bad module name: expected %#v got %#v", "test", mod.Name)
	}
	if len(mod.Statements) != 2 {
		t.Fatal("wrong number of statements")
	}
}

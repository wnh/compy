

package main

import (
	"testing"
)

func TestNewParser(t *testing.T) {
	_ = NewParser("")
}

func TestParseEmptyModule(t *testing.T) {
	//t.Skip("need to refactor first")
	mods := []struct{m, n string}{
		{"module main;", "main"},
		{"module other;", "other"},
		{"module more_stuff;", "more_stuff"},
	}
	for _, modTest := range mods {
		p := NewParser(modTest.m)
		//t.Logf("%#v, %v, %v\n", p, p.tok, p.nextTok);
		mod, err := p.ParseModule()
		if err != nil || mod == nil {
			t.Fail()
		}
		if mod.Name != modTest.n {
			t.Errorf("bad module name: expected %#v got %#v", modTest.n, mod.Name)
		}
	}
}



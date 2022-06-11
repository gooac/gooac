package gooa

import (
	"embed"
	"testing"
)

//go:embed tests/*
var compiletest_lua embed.FS

func TestCompiler(t *testing.T) {
	tok := Tokenizer{}
	
	data, _ := compiletest_lua.ReadFile("tests/test1.lua")
	
	toks, stop := tok.Tokenize(data, nil)
	
	if stop {
		t.Error("Refusing to parse")
		return
	}
	
	parser := Parser{}
	ast, stop := parser.Parse(toks, &tok.err)

	if stop {
		t.Error("Refusing to compile")
		return
	}

	compiler := Compiler{}
	code, stop := compiler.Compile(&ast, &parser.err)

	if stop {
		t.Error("Compilation failed")
		return
	}

	println("Successful Compilation!")
	println("----------[[START CODE]]---------")
	println(code)
	println("----------[[ END CODE ]]---------")
}
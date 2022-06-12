package gooa

import (
	"embed"
	"os"
	"testing"
)

//go:embed tests/*
var compiletest_lua embed.FS

func TestCompilerTest(t *testing.T) {
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

	println("Successful Compilation! (" + compiler.luv + ")")
	// println("----------[[START CODE]]---------")
	// println(code)
	// println("----------[[ END CODE ]]---------")
	// println("----------[[WRITE FILE]]---------")
	
	f, err := os.OpenFile("tests/testc.lua", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	defer f.Close()
	
	if err != nil {
		println("failed to write to file ", err)
	}
	
	f.WriteString("-- Output:\n")
	f.WriteString(code)
	// println("----------[[   DONE   ]]---------")
}

func TestCompilerSpec(t *testing.T) {
	tok := Tokenizer{}
	parser := Parser{}
	compiler := Compiler{}

	out := ""

	dir, err := compiletest_lua.ReadDir("tests/5.1")

	if err != nil {
		t.Error("Please download the 5.1 test suite from https://www.lua.org/tests/index.html and place its all the .lua files in the root in tests/5.1/")
		
		return
	}

	for _, v := range dir {
		println("\n -- TOKENIZING " + v.Name() + " -- ")
		str,_ := compiletest_lua.ReadFile("tests/5.1/" + v.Name())
		toks, stop := tok.Tokenize([]byte(str), nil)
		
		tok.err.Dump()

		if stop {
			println(" --FT")
			continue
		}

		ast, stop := parser.Parse(toks, &tok.err)

		tok.err.Dump()

		if stop {
			println(" --FP")
			break
		}

		cmp, err := compiler.Compile(&ast, &tok.err)

		if err {
			println(" --FC")
		} else {
			println(" --DONE COMPILING")
		}

		out += cmp + "\n\n\n"
	}

	f, err := os.OpenFile("tests/testc.lua", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	defer f.Close()
	
	if err != nil {
		println("failed to write to file ", err)
	}
	
	f.WriteString(out)
}
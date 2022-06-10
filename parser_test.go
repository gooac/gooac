package gooa

import (
	"embed"
	"testing"
	"fmt"
)

//go:embed tests/*
var parsetest_lua embed.FS

func Test51Dir(t *testing.T) {
	fmt.Print()
	tok := Tokenizer{}
	parser := Parser{}

	dir, _ := parsetest_lua.ReadDir("tests/5.1")

	for _, v := range dir {
		println("\n -- TOKENIZING " + v.Name() + " -- ")
		str,_ := parsetest_lua.ReadFile("tests/5.1/" + v.Name())
		toks, stop := tok.Tokenize([]byte(str), nil)
		
		tok.err.Dump()

		if stop {
			println(" --FT")
			continue
		} else {
			println(" --DONE TOKENIZING")
		}

		_, stop = parser.Parse(toks, &tok.err)

		tok.err.Dump()

		if stop {
			println(" --FP")
			break
		} else {
			println(" --DONE PARSING")
		}
	}
}

func TestParser(t *testing.T) {
	tok := Tokenizer{}
	parser := Parser{}

	data, _ := parsetest_lua.ReadFile("tests/test1.lua")

	toks, stop := tok.Tokenize(data, nil)

	if stop {
		println("Refusing to parse")
		return
	}

	// for _, v := range toks {
	// 	println()
	// 	v.Print()
	// }
	println()
	
	ast, stop := parser.Parse(toks, &tok.err)

	if stop {
		println("Refusing to compile")
		return
	}

	println("\n-- ASTDUMP --")
	// fmt.Printf("%v", ast)
	ast.Dump()
	println("-- ASTDUMP --")
}
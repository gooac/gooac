
package gooa

import (
	"embed"
	"testing"
)

//go:embed tests/*
var toktest_lua embed.FS

func TestTokenizer(t *testing.T) {
	// tok := Tokenizer{}
	// // data, _ := toktest_lua.ReadFile("tests/rustic.lua")

	// _, _ = tok.Tokenize(data, nil)

	// tok.DumpTokens()
}
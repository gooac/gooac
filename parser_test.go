
package gooa

import (
	// "fmt"
	"testing"
)

func TestParser(*testing.T) {
	tk := &Tokenizer{}

	goto_tokens := tk.Tokenize(`local a.b.c = 123`)

	// for _, v := range goto_tokens {
	// 	v.Print()
	// }

	parser := &Parser{}
	parser.Parse(goto_tokens)
}
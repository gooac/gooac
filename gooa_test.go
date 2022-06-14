package gooa

import (
	"testing"
	"embed"
)

//go:embed tests/*
var gooatest_lua embed.FS

func TestGooa(t *testing.T) {
	g := NewGooa()
	g.Use(AttributeMiddleware())

	f,_ := gooatest_lua.ReadFile("tests/cl_circles.lua")
	out, failed := g.Compile(f)
	
	if failed {
		println("Failed to compile! ")
		return
	}

	println(out)
}
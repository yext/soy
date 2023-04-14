package soymsg

import (
	"testing"

	"github.com/yext/soy/ast"
	"github.com/yext/soy/parse"
)

func TestSetPlaceholders(t *testing.T) {
	type test struct {
		node  *ast.MsgNode
		phstr string
	}

	var tests = []test{
		{newMsg("Hello world"), "Hello world"},

		// Data refs
		{newMsg("Hello {$name}"), "Hello {NAME}"},
		{newMsg("{$a}, {$b}, and {$c}"), "{A}, {B}, and {C}"},
		{newMsg("{$a} {$a}"), "{A} {A}"},
		{newMsg("{$a} {$b.a}"), "{A_1} {A_2}"},
		{newMsg("{$a.a}{$a.b.a}"), "{A_1}{A_2}"},

		// Command sequences
		{newMsg("hello{sp}world"), "hello world"},

		// HTML
		{newMsg("Click <a>here</a>"), "Click {START_LINK}here{END_LINK}"},
		{newMsg("<br><br/><br/>"), "{START_BREAK}{BREAK}{BREAK}"},
		{newMsg("<a href=foo>Click</a> <a href=bar>here</a >"),
			"{START_LINK_1}Click{END_LINK_1} {START_LINK_2}here{END_LINK_2}"},
		{newMsg("<p>P1</p><p>P2</p><p>P3</p>"),
			"{START_PARAGRAPH}P1{END_PARAGRAPH}{START_PARAGRAPH}P2{END_PARAGRAPH}{START_PARAGRAPH}P3{END_PARAGRAPH}"},

		// BUG: Data refs + HTML
		// {newMsg("<a href={$url}>Click</a>"), "{START_LINK}Click{END_LINK}"},

		// TODO: phname

		// TODO: investigate globals
		// {newMsg("{GLOBAL}"), "{GLOBAL}"},
		// {newMsg("{sub.global}"), "{GLOBAL}"},
	}

	for _, test := range tests {
		var actual = PlaceholderString(test.node)
		if actual != test.phstr {
			t.Errorf("(actual) %v != %v (expected)", actual, test.phstr)
		}
	}
}

func TestSetPluralVarName(t *testing.T) {
	type test struct {
		node    *ast.MsgNode
		varname string
	}

	var tests = []test{
		{newMsg("{plural $eggs}{case 1}one{default}other{/plural}"), "EGGS"},
		{newMsg("{plural $eggs}{case 1}one{default}{$eggs}{/plural}"), "EGGS_1"},
		{newMsg("{plural length($eggs)}{case 1}one{default}other{/plural}"), "NUM"},
	}

	for _, test := range tests {
		var actual = test.node.Body.Children()[0].(*ast.MsgPluralNode).VarName
		if actual != test.varname {
			t.Errorf("(actual) %v != %v (expected)", actual, test.varname)
		}
	}
}

func newMsg(msg string) *ast.MsgNode {
	// TODO: data.Map{"GLOBAL": data.Int(1), "sub.global": data.Int(2)})
	var sf, err = parse.SoyFile("", `{msg desc=""}`+msg+`{/msg}`)
	if err != nil {
		panic(err)
	}
	var msgnode = sf.Body[0].(*ast.MsgNode)
	SetPlaceholdersAndID(msgnode)
	return msgnode
}

func TestBaseName(t *testing.T) {
	type test struct {
		expr string
		ph   string
	}
	var tests = []test{
		{"$foo", "FOO"},
		{"$foo.boo", "BOO"},
		{"$foo.boo[0].zoo", "ZOO"},
		{"$foo.boo.0.zoo", "ZOO"},

		// parse.Expr doesn't accept undefined globals.
		// {"GLOBAL", "GLOBAL"},
		// {"sub.GLOBAL", "GLOBAL"},

		{"$foo[0]", "XXX"},
		{"$foo.boo[0]", "XXX"},
		{"$foo.boo.0", "XXX"},
		{"$foo + 1", "XXX"},
		{"'text'", "XXX"},
		{"max(1, 3)", "XXX"},
	}

	for _, test := range tests {
		var n, err = parse.Expr(test.expr)
		if err != nil {
			t.Error(err)
			return
		}

		var actual = genBasePlaceholderName(&ast.PrintNode{0, n, nil}, "XXX")
		if actual != test.ph {
			t.Errorf("(actual) %v != %v (expected)", actual, test.ph)
		}
	}
}

func TestToUpperUnderscore(t *testing.T) {
	var tests = []struct{ in, out string }{
		{"booFoo", "BOO_FOO"},
		{"_booFoo", "BOO_FOO"},
		{"booFoo_", "BOO_FOO"},
		{"BooFoo", "BOO_FOO"},
		{"boo_foo", "BOO_FOO"},
		{"BOO_FOO", "BOO_FOO"},
		{"__BOO__FOO__", "BOO_FOO"},
		{"Boo_Foo", "BOO_FOO"},
		{"boo8Foo", "BOO_8_FOO"},
		{"booFoo88", "BOO_FOO_88"},
		{"boo88_foo", "BOO_88_FOO"},
		{"_boo_8foo", "BOO_8_FOO"},
		{"boo_foo8", "BOO_FOO_8"},
		{"_BOO__8_FOO_", "BOO_8_FOO"},
	}
	for _, test := range tests {
		var actual = toUpperUnderscore(test.in)
		if actual != test.out {
			t.Errorf("(actual) %v != %v (expected)", actual, test.out)
		}
	}
}

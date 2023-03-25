package parser

import (
	"context"
	"encoding/json"
	"testing"
	"yal/lexer"
)

func TestParser(t *testing.T) {
	ctx := context.Background()
	l := lexer.NewLexer(ctx, `
	let variable: int = (4 * 2) + 5;
	variable = 4;
	variable = 5 * 2 * ( 5 + 3 );

	let test = 5;

	if ((x > 5) || (x < 2)) {
		let y = (5 + 1) - 2;
	} else {
		let h = NULL;
	}
	
	while (7 < x) {
		--x;
	}

	for (let x = 10; x < 10; ++x) {
		let str = "hello";
		let a = "test";
	}
	`)

	tokens, _ := l.Scan()
	parser := NewParser(ctx, tokens)
	tree := parser.Run()
	data, _ := json.MarshalIndent(tree, "", "  ")
	t.Logf("%v\n", string(data))
}

package parser

import (
	"encoding/json"
	"testing"
)

func TestLexer(t *testing.T) {
	l := NewLexer(`
	let variable = (4 * 2) + 5;
	variable = 4;
	variable = 5 * 2 * ( 5 + 3 );
	{
		if ((x > 5) || (x < 2)) {
			let y = (5 + 1) - 2;
		} else {
			let h = NULL;
		}

	}
	while (true) {
		--x;
	}

	for (let x = 10; x < 10; ++x) {
		let str = "hello";
	}
	`)

	tokens, _ := l.Scan()
	parser := NewParser(tokens)
	tree := parser.Run()
	data, _ := json.MarshalIndent(tree, "", "  ")

	t.Logf("%v\n", string(data))
}

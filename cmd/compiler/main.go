package main

import (
	"context"
	"encoding/json"
	"fmt"
	"yal/lexer"
	"yal/parser"
)

func main() {
	ctx := context.Background()
	l := lexer.NewLexer(ctx, `
fn main() : void {
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
	// comment
	/* comment */
}
fn test(a: int, b: uint, c: char) : uint {
	print("a b c", -1, --1, x + 5);
}
`)

	tokens, _ := l.Scan()
	parser := parser.NewParser(ctx, tokens)
	tree := parser.Run()
	data, _ := json.MarshalIndent(tree, "", "  ")
	fmt.Printf("%v\n", string(data))
}

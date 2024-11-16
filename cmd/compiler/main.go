package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"yal/lexer"
	"yal/parser"
)

func main() {
	ctx := context.Background()

	if len(os.Args) != 2 {
		panic("there must have 1 parameter, a file or - for stdin")
	}

	var f *os.File
	var data []byte
	var err error

	if os.Args[1] == "-" {
		f = os.Stdin
	} else {
		f, err = os.Open(os.Args[1])
		if err != nil {
			panic(err)
		}
	}

	data, err = io.ReadAll(f)
	if err != nil {
		panic(err)
	}

	l := lexer.NewLexer(ctx, string(data))

	tokens, err := l.Scan()
	if err != nil {
		panic(err)
	}

	yalParser := parser.NewParser(ctx, tokens)
	tree := yalParser.Run()

	// for _, stmt := range tree {
	// 	switch stmt.(type) {
	// 	case *parser.Block:
	// 	case *parser.DefineTypeStatement:
	// 	case *parser.FnArgs:
	// 	case *parser.FnDeclStmt:
	// 		t, _ := stmt.(*parser.FnDeclStmt)
	// 		fmt.Println(t.Type)
	// 	case *parser.ForLoop:
	// 	case *parser.IfExpr:
	// 	case *parser.StatementExpression:
	// 	case *parser.VarDeclExpression:
	// 	case *parser.WhileLoop:
	// 	default:
	// 		panic(fmt.Sprintf("unexpected parser.IStatement: %#v", stmt))
	// 	}
	// }

	data, err = json.MarshalIndent(tree, "", " ")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", string(data))
}

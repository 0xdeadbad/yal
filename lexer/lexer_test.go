package lexer_test

import (
	"encoding/json"
	"setlang/lexer"
	"setlang/parser"
	"testing"
)

func TestLexer(t *testing.T) {
	l := lexer.NewLexer(`
	{
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
		}
	}
	
	`)

	tokens, _ := l.Scan()
	parser := parser.NewParser(tokens)
	tree := parser.Run()
	data, _ := json.MarshalIndent(tree, "", "  ")
	t.Logf("%v\n", string(data))
}

/*func recurseExpression(expr IExpression, t *testing.T) string {
	switch expr := expr.(type) {
	case *Binary:
		left := recurseExpression(expr.Left, t)
		op := expr.Operator
		right := recurseExpression(expr.Right, t)
		return fmt.Sprintf("%+v %+v %+v", left, op.Lexeme, right)
	case *Unary:
		operator := expr.Operator.Lexeme
		right := recurseExpression(expr.Right, t)
		return fmt.Sprintf("%+v%+v", operator, right)
	case *Literal[string]:
		return fmt.Sprintf("%+v", expr.Value)
	case *Literal[bool]:
		return fmt.Sprintf("%+v", expr.Value)
	case *Literal[int64]:
		return fmt.Sprintf("%+v", expr.Value)
	case *Grouping:
		return fmt.Sprintf("( %+v )", recurseExpression(expr.Grouped, t))
	case *Variable:
		return fmt.Sprintf("%+v", expr.Name.Lexeme)
	case *Logical:
		left := recurseExpression(expr.Left, t)
		operator := expr.Operator.Lexeme
		right := recurseExpression(expr.Right, t)
		return fmt.Sprintf("%+v %+v %+v", left, operator, right)
	case *Assign:
		name := expr.Name.Lexeme
		exprr := recurseExpression(expr.Expr, t)
		return fmt.Sprintf("%+v = %+v", name, exprr)
	default:
		return fmt.Sprintf("%+v", expr)
	}
}

func recurseStatement(stmt IStatement, t *testing.T) string {
	switch stmt := stmt.(type) {
	case *StatementExpression:
		expr := recurseExpression(stmt.Expr, t)
		return fmt.Sprintf("%+v;", expr)
	case *VarDeclExpression:
		name := stmt.Name.Lexeme
		initializer := recurseExpression(stmt.Initializer, t)
		return fmt.Sprintf("let %+v = %+v", name, initializer)
	case *Block:
		list := []string{}
		for _, st := range stmt.Statements {
			list = append(list[:], recurseStatement(st, t))
		}
		res := strings.Join(list, " ")
		return fmt.Sprintf("{ %+v }", res)
	case *IfExpr:
		condition := recurseExpression(stmt.Condition, t)
		thenBranch := recurseStatement(stmt.ThenBranch, t)
		elseBranch := recurseStatement(stmt.ElseBranch, t)
		return fmt.Sprintf("If %+v { %+v } else { %+v }", condition, thenBranch, elseBranch)
	case *WhileLoop:
		condition := recurseExpression(stmt.Condition, t)
		body := recurseStatement(stmt.Body, t)
		return fmt.Sprintf("While %+v %+v", condition, body)
	case *ForLoop:
		initializer := recurseStatement(stmt.Initializer, t)
		condition := recurseExpression(stmt.Condition, t)
		apply := recurseExpression(stmt.Apply, t)
		body := recurseStatement(stmt.Body, t)
		return fmt.Sprintf("For (%+v; %+v; %+v) %+v", initializer, condition, apply, body)
	default:
		return fmt.Sprintf("Statement { %+v }", stmt)
	}
}*/

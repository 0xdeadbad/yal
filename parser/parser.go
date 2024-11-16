package parser

import (
	"context"
	"fmt"
	. "yal/lexer"
)

// Parser struct with methods
type Parser struct {
	Tokens  []Token
	current uint64
	ctx     context.Context
}

// Returns a new Parser instance with the providade Tokens list
func NewParser(ctx context.Context, tokens []Token) *Parser {
	return &Parser{
		Tokens:  tokens,
		current: 0,
		ctx:     ctx,
	}
}

// Processes the Tokens list parsing them and returning an AST
func (p *Parser) Run() []IStatement {
	stmts := []IStatement{}

	for !p.isEof() {
		stmts = append(stmts[:], p.declaration())
	}

	return stmts
}

func (p *Parser) declaration() IStatement {
	// TODO: perhaps we can use switches instead of 'if' blocks
	ok, v := p.match(DefineType, Let)
	if !ok {
		return p.statement()
	}

	switch v.TokenType {
	case DefineType:
		return p.defineTypeStatement()
	case Let:
		return p.varDeclaration()
	default:
		p.panicReason("It's not supposed to reach here\n")
		return nil
	}
}

func (p *Parser) statement() IStatement {
	// TODO: perhaps we can use switches instead of 'if' blocks

	ok, v := p.match(Fn, For, While, LeftBrace)
	if !ok {
		return p.expressionStatement()
	}

	switch v.TokenType {
	case Fn:
		return p.fnStatement()
	case For:
		return p.forStatement()
	case While:
		return p.whileStatement()
	case LeftBrace:
		return p.block()
	case If:
		return p.ifStatement()
	}

	panic("It's not supposed to reach here")
}

func (p *Parser) fnReturn() IExpression {
	p.consume(Return, "Expected return keyword (???)")

	return &FnReturn{
		Value: p.expression(),
	}
}

func (p *Parser) varDeclaration() IStatement {
	name := p.consume(Identifier, "Expect variable name.")

	var type_ann *Token
	var initializer IExpression = &Literal{Value: nil}
	if p.matchNT(Colon) {
		if p.peek().TokenType == Identifier {
			type_ann = p.advance()
		} else {
			type_ann = nil
		}
	}
	if p.matchNT(Equal) {
		initializer = p.expression()
	} else if p.matchNT(Comma) || p.peek().TokenType == RightParen {
		return &VarDeclExpression{
			Name: name,
			Initializer: &Literal{
				Value: nil,
			},
			Type: type_ann,
		}
	}

	p.consume(Semicolon, "Expect ';' after variable declaration.")

	return &VarDeclExpression{
		Name:        name,
		Initializer: initializer,
		Type:        type_ann,
	}
}

func (p *Parser) defineTypeStatement() IStatement {
	name := p.consume(Identifier, "Expected type name for type definition")
	p.consume(Equal, "Expected = after type definition name.")
	tokenType := p.consume(Identifier, "Expected a type after =")
	p.consume(Semicolon, "Expected ';' after type definition")

	return &DefineTypeStatement{
		Name: name,
		Type: tokenType,
	}
}

func (p *Parser) expressionStatement() IExpression {
	if p.peek().TokenType == If {
		p.consume(If, "This is not supposed to fail...")
		return p.ifStatement()
	}

	expr := p.expression()
	if p.peek().TokenType == RightBrace {
		return &FnReturn{
			Value: expr,
		}
	}

	p.consume(Semicolon, "Expect ';' after expression.")
	return &StatementExpression{
		Expr: expr,
	}
}

func (p *Parser) fnCall() IExpression {
	fnName := p.consume(Identifier, "there's not a valid identifier for the fn call")
	p.consume(LeftParen, "missing ( after fn identifier")
	fnArgs := FnCallArgs{}
	for p.peek().TokenType != RightParen {
		fnArgs = append(fnArgs, p.fnCallArg())
	}
	p.consume(RightParen, "missing ) after fn args")

	return &FnCall{
		Name: fnName,
		Args: fnArgs,
		Type: nil,
	}
}

func (p *Parser) fnCallArg() IExpression {
	p.match(Comma)
	return p.expression()
}

func (p *Parser) assignment() IExpression {
	expr := p.or()

	if p.matchNT(Equal) {
		p.previous()
		value := p.assignment()

		v, ok := expr.(*Variable)
		if ok {
			name := v.Name
			return &Assign{
				Name: name,
				Expr: value,
			}
		}

		p.panicReason("Error at line %d column %d\n", v.Line, v.Column)
	}

	return expr
}

func (p *Parser) expression() IExpression {
	// fmt.Println("curr ", p.Tokens[p.current].Type.String(), p.Tokens[p.current].Lexeme, p.peek().Type.String(), p.peek().Lexeme, p.peekNext().Type.String(), p.peekNext().Lexeme)
	if p.peek().TokenType == Identifier && p.peekNext().TokenType == LeftParen {
		return p.fnCall()
	}
	if p.peek().TokenType == Identifier && p.matchNextNT(Dec, Inc) {
		return p.unaryLeft()
	}
	if p.peek().TokenType == Return {
		return p.fnReturn()
	}

	return p.assignment()
}

func (p *Parser) equality() IExpression {
	expr := p.logic()

	ok, op := p.match(BangEqual, EqualEqual)

	for ok {
		operator := op
		right := p.logic()
		expr = &Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}

		ok, op = p.match(BangEqual, EqualEqual)
	}

	return expr
}

func (p *Parser) logic() IExpression {
	expr := p.comparison()

	for p.matchNT(DoubleAmpersand, DoublePipe) {
		operator := p.previous()
		right := p.comparison()
		expr = &Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr
}

func (p *Parser) comparison() IExpression {
	expr := p.term()

	for p.matchNT(Lesser, LesserEqual, Greater, GreaterEqual, EqualEqual, BangEqual) {
		operator := p.previous()
		right := p.term()
		expr = &Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr
}

func (p *Parser) term() IExpression {
	expr := p.factor()

	for p.matchNT(Minus, Plus) {
		operator := p.previous()
		right := p.factor()
		expr = &Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr
}

func (p *Parser) factor() IExpression {
	expr := p.unaryRight()

	for p.matchNT(Slash, Star) {
		operator := p.previous()
		right := p.unaryRight()
		expr = &Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr
}

func (p *Parser) unaryRight() IExpression {
	if p.matchNT(Bang, Minus, Inc, Dec) {
		operator := p.previous()
		right := p.unaryRight()
		return &UnaryRight{
			Operator: operator,
			Right:    right,
		}
	}

	return p.primary()
}

func (p *Parser) unaryLeft() IExpression {
	if p.matchNT(Inc, Dec) {
		operator := p.previous()
		left := p.unaryLeft()
		return &UnaryLeft{
			Operator: operator,
			Left:     left,
		}
	}

	return p.primary()
}

func (p *Parser) primary() IExpression {

	if p.matchNT(Number2, Number8, Number10, Number16, String, False, True, Null) {
		return &Literal{
			Value: p.previous(),
		}
	}

	if p.matchNT(Identifier) {
		return &Variable{
			Loc:         Loc{},
			IExpression: nil,
			Name:        p.previous(),
		}
	}

	if p.matchNT(LeftParen) {
		expr := p.expression()
		p.consume(RightParen, "")
		return &Grouping{
			Grouped: expr,
		}
	}

	p.panicReason("Error on primary(): Line %d Column %d\nToken: %+v\nPrevious: %+v\nNext: %+v\n", p.peek().Line, p.peek().Column, p.peek(), p.previous(), p.peekNext())

	return nil
}

func (p *Parser) block() IStatement {
	statements := []IStatement{}

	ok, _ := p.check(RightBrace)

	for !p.isEof() && !ok {
		statements = append(statements[:], p.declaration())
		ok, _ = p.check(RightBrace)
	}

	tk := p.consume(RightBrace, "Expect '}' after block.")

	return &Block{
		Loc: Loc{
			Line:   tk.Line,
			Column: tk.Column,
		},
		Statements: statements,
	}
}

func (p *Parser) ifStatement() IStatement {
	p.consume(LeftParen, "Expect '(' after 'if'.")
	condition := p.expression()
	p.consume(RightParen, "Expect ')' after if condition.")

	thenBranch := p.statement()
	var elseBranch IStatement = nil
	if p.matchNT(Else) {
		elseBranch = p.statement()
	}

	return &IfExpr{
		Condition:  condition,
		ThenBranch: thenBranch,
		ElseBranch: elseBranch,
	}
}

func (p *Parser) fnStatement() IStatement {
	fnName := p.consume(Identifier, "Expect 'fn' name.")
	p.consume(LeftParen, "Expect '(' after 'fn' name.")
	fnArgs := FnArgs{}
	for p.peek().TokenType != RightParen {
		fnArgs = append(fnArgs, p.varDeclaration())
	}
	p.consume(RightParen, "Expect ')' fn args.")

	var fnType *Token

	if p.matchNT(Colon) {
		fnType = p.consume(Identifier, "Expect identifier after : in fn decl")
	}

	fnBody := p.statement()

	return &FnDeclStmt{
		Name: fnName,
		Body: fnBody,
		Type: fnType,
		Args: &fnArgs,
	}
}

func (p *Parser) or() IExpression {
	expr := p.and()

	for p.matchNT(DoublePipe) {
		operator := p.previous()
		right := p.and()
		expr = &Logical{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr
}

func (p *Parser) and() IExpression {
	expr := p.equality()

	for p.matchNT(DoubleAmpersand) {
		operator := p.previous()
		right := p.equality()
		expr = &Logical{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr
}

func (p *Parser) whileStatement() IStatement {
	p.consume(LeftParen, "Expect '(' after 'while'.")
	condition := p.expression()
	p.consume(RightParen, "Expect ')' after condition.")
	body := p.statement()

	return &WhileLoop{
		Condition: condition,
		Body:      body,
	}
}

func (p *Parser) forStatement() IStatement {
	p.consume(LeftParen, "Expect '(' after 'for'.")

	var initializer IStatement
	if p.matchNT(Semicolon) {
		initializer = nil
	} else if p.matchNT(Let) {
		initializer = p.varDeclaration()
	} else {
		initializer = p.expressionStatement()
	}

	var condition IExpression = nil
	if !p.matchNT(Semicolon) {
		condition = p.expression()
	}
	p.consume(Semicolon, "Expect ';' after 'condition'.")

	var apply IExpression = nil
	if !p.checkNT(RightParen) {
		apply = p.expression()
	}
	p.consume(RightParen, "Expect ')' after condition.")
	body := p.statement()

	return &ForLoop{
		Initializer: initializer,
		Condition:   condition,
		Apply:       apply,
		Body:        body,
	}
}

// ----- Utility and help functions -----

func (p *Parser) panicReason(s string, args ...any) {
	panic(fmt.Sprintf(s, args...))
}

func (p *Parser) match(token_types ...TokenType) (bool, *Token) {
	for _, tkt := range token_types {
		if p.checkNT(tkt) {
			return true, p.advance()
		}
	}

	return false, nil
}

func (p *Parser) matchNT(token_types ...TokenType) bool {
	for _, tkt := range token_types {
		if p.checkNT(tkt) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) matchNext(token_types ...TokenType) (bool, *Token) {
	for _, tkt := range token_types {
		if p.checkNextNT(tkt) {
			return true, p.advance()
		}
	}

	return false, nil
}

func (p *Parser) matchNextNT(token_types ...TokenType) bool {
	for _, tkt := range token_types {
		if p.checkNextNT(tkt) {
			p.advance()
			return true
		}
	}

	return false
}

// func (p *Parser) matchPrevious(token_types ...TokenType) bool {
// 	for _, tkt := range token_types {
// 		if p.checkPrevious(tkt) {
// 			p.advance()
// 			return true
// 		}
// 	}

// 	return false
// }

func (p *Parser) check(token_type TokenType) (bool, *Token) {
	if p.isEof() || p.peek().TokenType != token_type {
		return false, nil
	}

	return true, p.peek()
}

func (p *Parser) checkNT(token_type TokenType) bool {
	if p.isEof() || p.peek().TokenType != token_type {
		return false
	}

	return true
}

func (p *Parser) checkNext(token_type TokenType) (bool, *Token) {
	if p.isEof() || p.peekNext().TokenType != token_type {
		return false, nil
	}

	return true, p.peekNext()
}

func (p *Parser) checkNextNT(token_type TokenType) bool {
	if p.isEof() || p.peekNext().TokenType != token_type {
		return false
	}

	return true
}

func (p *Parser) checkPrevious(token_type TokenType) (bool, *Token) {
	if tk := p.peekPrevious(); tk.TokenType == token_type {
		return true, tk
	}

	return false, nil
}

func (p *Parser) checkPreviousNT(token_type TokenType) bool {
	return p.peekPrevious().TokenType == token_type
}

func (p *Parser) advance() *Token {
	if !p.isEof() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isEof() bool {
	return p.peek().TokenType == Eof
}

func (p *Parser) peek() *Token {
	return &p.Tokens[p.current]
}

func (p *Parser) peekNext() *Token {
	if p.isEof() {
		return nil
	}
	return &p.Tokens[p.current+1]
}

func (p *Parser) peekPrevious() *Token {
	return &p.Tokens[p.current-1]
}

func (p *Parser) previous() *Token {
	return &p.Tokens[p.current-1]
}

func (p *Parser) consume(token_type TokenType, message string) *Token {
	if ok, _ := p.check(token_type); ok {
		return p.advance()
	}
	return nil
}

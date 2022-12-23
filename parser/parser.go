package parser

import (
	"fmt"
	"strconv"
)

// Parser struct with methods
type Parser struct {
	tokens  []Token
	current uint64
}

// Returns a new Parser instance with the providade tokens list
func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

// Processes the tokens list parsing them and returning an AST
func (p *Parser) Run() []IStatement {
	stmts := []IStatement{}

	for !p.isEof() {
		stmts = append(stmts[:], p.declaration())
	}

	return stmts
}

func (p *Parser) statement() IStatement {
	if p.match(For) {
		return p.forStatement()
	}
	if p.match(While) {
		return p.whileStatement()
	}
	if p.match(If) {
		return p.ifStatement()
	}
	if p.match(LeftBrace) {
		return &Block{Statements: p.block()}
	}

	return p.expressionStatement()
}

func (p *Parser) declaration() IStatement {
	if p.match(Let) {
		return p.varDeclaration()
	}

	return p.statement()
}

func (p *Parser) varDeclaration() IStatement {
	name := p.consume(Identifier, "Expect variable name.")

	var initializer IExpression
	if p.match(Equal) {
		initializer = p.expression()
	} else {
		initializer = &Literal[interface{}]{Value: nil}
	}

	p.consume(Semicolon, "Expect ';' after variable declaration.")

	return &VarDeclExpression{Name: name, Initializer: initializer}
}

func (p *Parser) expressionStatement() IStatement {
	expr := p.expression()
	p.consume(Semicolon, "Expect ';' after expression.")
	return &StatementExpression{
		Expr: expr,
	}
}

func (p *Parser) assignment() IExpression {
	expr := p.or()

	if p.match(Equal) {
		_ = p.previous()
		value := p.assignment()

		v, ok := expr.(*Variable)
		if ok {
			name := v.Name
			return &Assign{Name: name, Expr: value}
		}

		panic("error on assigment()")
	}

	return expr
}

func (p *Parser) expression() IExpression {
	return p.assignment()
}

func (p *Parser) equality() IExpression {
	expr := p.logic()

	for p.match(BangEqual, EqualEqual) {
		operator := p.previous()
		right := p.logic()
		expr = &Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) logic() IExpression {
	expr := p.comparison()

	for p.match(DoubleAmpersand, DoublePipe) {
		operator := p.previous()
		right := p.comparison()
		expr = &Binary{expr, operator, right}
	}

	return expr
}

func (p *Parser) comparison() IExpression {
	expr := p.term()

	for p.match(Lesser, LesserEqual, Greater, GreaterEqual, EqualEqual, BangEqual) {
		operator := p.previous()
		right := p.term()
		expr = &Binary{expr, operator, right}
	}

	return expr
}

func (p *Parser) term() IExpression {
	expr := p.factor()

	for p.match(Minus, Plus) {
		operator := p.previous()
		right := p.factor()
		expr = &Binary{expr, operator, right}
	}

	return expr
}

func (p *Parser) factor() IExpression {
	expr := p.unary()

	for p.match(Slash, Star) {
		operator := p.previous()
		right := p.unary()
		expr = &Binary{expr, operator, right}
	}

	return expr
}

func (p *Parser) unary() IExpression {
	if p.match(Bang, Minus, Inc, Dec) {
		operator := p.previous()
		right := p.unary()
		return &Unary{Operator: operator, Right: right}
	}

	return p.primary()
}

func (p *Parser) leftUnary() IExpression {
	if p.matchNext(Inc, Dec) {
		operator := p.peekNext()
		right := p.unary()
		p.advance()
		return &Unary{Operator: operator, Right: right}
	}

	return p.primary()
}

func (p *Parser) primary() IExpression {

	if p.match(False) {
		return &Literal[bool]{
			Value: false,
		}
	}
	if p.match(True) {
		return &Literal[bool]{
			Value: true,
		}
	}
	if p.match(Null) {
		return &Literal[interface{}]{
			Value: nil,
		}
	}

	if p.match(Number10, String) {
		switch p.previous().Type {
		case String:
			return &Literal[string]{
				Value: p.previous().Lexeme,
			}
		case Number10:
			{
				i, _ := strconv.Atoi(p.previous().Lexeme)
				return &Literal[int64]{
					Value: int64(i),
				}
			}
		}
	}

	if p.match(Identifier) {
		return &Variable{Name: p.previous()}
	}

	if p.match(LeftParen) {
		expr := p.expression()
		p.consume(RightParen, "")
		return &Grouping{
			Grouped: expr,
		}
	}

	panic("Error on primary()")
}

func (p *Parser) match(token_types ...TokenType) bool {
	for _, tkt := range token_types {
		if p.check(tkt) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) matchNext(token_types ...TokenType) bool {
	for _, tkt := range token_types {
		if p.checkNext(tkt) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) matchPrevious(token_types ...TokenType) bool {
	for _, tkt := range token_types {
		if p.checkPrevious(tkt) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) check(token_type TokenType) bool {
	if p.isEof() {
		return false
	}

	return p.peek().Type == token_type
}

func (p *Parser) checkNext(token_type TokenType) bool {
	if p.isEof() {
		return false
	}

	return p.peekNext().Type == token_type
}

func (p *Parser) checkPrevious(token_type TokenType) bool {

	return p.peekPrevious().Type == token_type
}

func (p *Parser) advance() Token {
	if !p.isEof() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isEof() bool {
	return p.peek().Type == Eof
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) peekNext() Token {
	if !p.isEof() {
		return p.tokens[len(p.tokens)-1]
	}
	return p.tokens[p.current+1]
}

func (p *Parser) peekPrevious() Token {
	return p.tokens[p.current-1]
}

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}

func (p *Parser) consume(token_type TokenType, message string) Token {
	if p.check(token_type) {
		return p.advance()
	}

	panic(fmt.Sprintf("Problem on consume(%+v)", token_type))
}

func (p *Parser) block() []IStatement {
	statements := []IStatement{}

	for !p.isEof() && !p.check(RightBrace) {
		statements = append(statements[:], p.declaration())
	}

	p.consume(RightBrace, "Expect '}' after block.")

	return statements
}

func (p *Parser) ifStatement() IStatement {
	p.consume(LeftParen, "Expect '(' after 'if'.")
	condition := p.expression()
	p.consume(RightParen, "Expect ')' after if condition.")

	thenBranch := p.statement()
	var elseBranch IStatement = nil
	if p.match(Else) {
		elseBranch = p.statement()
	}

	return &IfExpr{Condition: condition, ThenBranch: thenBranch, ElseBranch: elseBranch}
}

func (p *Parser) or() IExpression {
	expr := p.and()

	for p.match(DoublePipe) {
		operator := p.previous()
		right := p.and()
		expr = &Logical{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) and() IExpression {
	expr := p.equality()

	for p.match(DoubleAmpersand) {
		operator := p.previous()
		right := p.equality()
		expr = &Logical{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) whileStatement() IStatement {
	p.consume(LeftParen, "Expect '(' after 'while'.")
	condition := p.expression()
	p.consume(RightParen, "Expect ')' after condition.")
	body := p.statement()

	return &WhileLoop{Condition: condition, Body: body}
}

func (p *Parser) forStatement() IStatement {
	p.consume(LeftParen, "Expect '(' after 'for'.")

	var initializer IStatement
	if p.match(Semicolon) {
		initializer = nil
	} else if p.match(Let) {
		initializer = p.varDeclaration()
	} else {
		initializer = p.expressionStatement()
	}

	var condition IExpression = nil
	if !p.match(Semicolon) {
		condition = p.expression()
	}
	p.consume(Semicolon, "Expect ';' after 'condition'.")

	var apply IExpression = nil
	if !p.check(RightParen) {
		apply = p.expression()
	}
	p.consume(RightParen, "Expect ')' after condition.")
	body := p.statement()

	return &ForLoop{Initializer: initializer, Condition: condition, Apply: apply, Body: body}
}

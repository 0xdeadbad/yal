package parser

import (
	. "yal/lexer"
)

type iExpression interface{}
type iStatement interface {
	iExpression
}
type IExpression = interface {
	iStatement | iExpression
}
type IStatement = iStatement
type TypeT = interface {
	IStatement | IExpression
}

type Loc struct {
	Line   uint64
	Column uint64
}

type Binary struct {
	Loc
	IExpression
	Left     IExpression
	Operator *Token
	Right    IExpression
}

func (b *Binary) exprNode() IExpression {
	return nil
}
func (b *Binary) GetType() any {
	return nil
}
func (b *Binary) TokenLoc() (uint64, uint64) {
	return b.Operator.Line, b.Operator.Column
}

type UnaryRight struct {
	Loc
	IExpression
	Operator *Token
	Right    IExpression
}

func (b *UnaryRight) exprNode() IExpression {
	return nil
}
func (b *UnaryRight) GetType() any {
	return nil
}

type UnaryLeft struct {
	Loc
	IExpression
	Operator *Token
	Left     IExpression
}

func (b *UnaryLeft) exprNode() IExpression {
	return nil
}
func (b *UnaryLeft) GetType() any {
	return nil
}

type Literal struct {
	Loc
	IExpression
	Value *Token
}

func (b *Literal) exprNode() IExpression {
	return nil
}
func (b *Literal) GetType() any {
	return nil
}

type FnReturn struct {
	Loc
	IStatement
	Value IExpression
}

func (b *FnReturn) stmtNode() {}
func (b *FnReturn) exprNode() IExpression {
	return nil
}
func (b *FnReturn) GetType() any {
	return nil
}

type Grouping struct {
	Loc
	IExpression
	Grouped IExpression
}

func (b *Grouping) exprNode() IExpression {
	return nil
}
func (b *Grouping) GetType() any {
	return nil
}

type Variable struct {
	Loc
	IExpression
	Name *Token
}

func (b *Variable) exprNode() IExpression {
	return nil
}
func (b *Variable) GetType() any {
	return nil
}

type Assign struct {
	Loc
	IExpression
	Name *Token
	Expr IExpression
}

func (b *Assign) exprNode() IExpression {
	return nil
}
func (b *Assign) GetType() any {
	return nil
}

type StatementExpression struct {
	Loc
	IExpression
	Expr IExpression
}

func (s *StatementExpression) stmtNode() {}
func (s *StatementExpression) exprNode() IExpression {
	return nil
}
func (b *StatementExpression) GetType() any {
	return nil
}

type VarDeclExpression struct {
	Loc
	IExpression
	Name        *Token
	Initializer IExpression
	Type        *Token
}

func (b *VarDeclExpression) stmtNode() {}
func (s *VarDeclExpression) exprNode() IExpression {
	return nil
}
func (b *VarDeclExpression) GetType() any {
	return nil
}

type DefineTypeStatement struct {
	Loc
	IStatement
	Name *Token
	Type *Token
}

func (b *DefineTypeStatement) stmtNode() {}

type Block struct {
	Loc
	IStatement
	Statements []IStatement
}

func (s *Block) stmtNode() {}
func (s *Block) exprNode() IExpression {
	return nil
}
func (s *Block) GetType() any {
	return nil
}

type IfExpr struct {
	Loc
	IStatement
	Condition  IExpression
	ThenBranch IStatement
	ElseBranch IStatement
}

func (b *IfExpr) stmtNode() {}
func (b *IfExpr) exprNode() IExpression {
	return nil
}
func (b *IfExpr) GetType() any {
	return nil
}

type FnDeclStmt struct {
	Loc
	IStatement
	Name *Token
	Type *Token
	Args *FnArgs
	Body IStatement
}

func (b *FnDeclStmt) stmtNode() {}

type FnArgs []IStatement

func (b *FnArgs) stmtNode() {}

type FnCallArgs []IExpression

type FnCall struct {
	Loc
	IExpression
	Name *Token
	Args FnCallArgs
	Type *Token
}

func (b *FnCall) exprNode() IExpression {
	return nil
}
func (b *FnCall) GetType() any {
	return nil
}

type Logical struct {
	Loc
	IExpression
	Left     IExpression
	Operator *Token
	Right    IExpression
}

func (b *Logical) exprNode() IExpression {
	return nil
}
func (b *Logical) GetType() any {
	return nil
}

type WhileLoop struct {
	Loc
	IStatement
	Condition IExpression
	Body      IStatement
}

func (b *WhileLoop) stmtNode() {}

type ForLoop struct {
	Loc
	IStatement
	Initializer IStatement
	Condition   IExpression
	Apply       IExpression
	Body        IStatement
}

func (b *ForLoop) stmtNode() {}

package parser

type IStatement interface {
	stmtNode()
}

type IExpression interface {
	exprNode()
}

type IPrimary interface {
	primNode()
}

type IUnary interface {
	unaryNode()
}

type IFactor interface {
	factorNode()
}

type ITerm interface {
	termNode()
}

type Binary struct {
	Left     IExpression
	Operator Token
	Right    IExpression
}

func (b *Binary) exprNode() {}

type Unary struct {
	Operator Token
	Right    IExpression
}

func (b *Unary) exprNode() {}

type LeftUnary struct {
	Operator Token
	Left     IExpression
}

func (b *LeftUnary) exprNode() {}

type Literal[T any] struct {
	Value T
}

func (b *Literal[T]) exprNode() {}

type Grouping struct {
	Grouped IExpression
}

func (b *Grouping) exprNode() {}

type Variable struct {
	Name Token
}

func (b *Variable) exprNode() {}

type Assign struct {
	Name Token
	Expr IExpression
}

func (b *Assign) exprNode() {}
func (b *Assign) stmtNode() {}

type StatementExpression struct {
	Expr IExpression
}

func (s *StatementExpression) stmtNode() {}

type VarDeclExpression struct {
	Name        Token
	Initializer IExpression
}

func (s *VarDeclExpression) exprNode() {}
func (b *VarDeclExpression) stmtNode() {}

type Block struct {
	Statements []IStatement
}

func (s *Block) stmtNode() {}

type IfExpr struct {
	Condition  IExpression
	ThenBranch IStatement
	ElseBranch IStatement
}

func (b *IfExpr) stmtNode() {}

type Logical struct {
	Left     IExpression
	Operator Token
	Right    IExpression
}

func (b *Logical) exprNode() {}

type WhileLoop struct {
	Condition IExpression
	Body      IStatement
}

func (b *WhileLoop) stmtNode() {}

type ForLoop struct {
	Initializer IStatement
	Condition   IExpression
	Apply       IExpression
	Body        IStatement
}

func (b *ForLoop) stmtNode() {}

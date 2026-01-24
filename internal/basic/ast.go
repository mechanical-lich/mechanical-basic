package basic

// Node is the base interface for all AST nodes
type Node interface {
	node()
	Position() (line, column int)
}

// Statement represents an executable statement
type Statement interface {
	Node
	statement()
}

// Expression represents a value-producing expression
type Expression interface {
	Node
	expression()
}

// -----------------------------------------------------------------------------
// Base position tracking
// -----------------------------------------------------------------------------

// Pos holds line and column information for error reporting
type Pos struct {
	Line   int
	Column int
}

func (p Pos) Position() (line, column int) {
	return p.Line, p.Column
}

// -----------------------------------------------------------------------------
// Program (root node)
// -----------------------------------------------------------------------------

// Program is the root node containing all statements
type Program struct {
	Statements []Statement
}

func (p *Program) node() {}
func (p *Program) Position() (line, column int) {
	if len(p.Statements) > 0 {
		return p.Statements[0].Position()
	}
	return 1, 1
}

// -----------------------------------------------------------------------------
// Statements
// -----------------------------------------------------------------------------

// LetStatement represents: LET x = expr
type LetStatement struct {
	Pos
	Name  string
	Value Expression
}

func (s *LetStatement) node()      {}
func (s *LetStatement) statement() {}

// AssignStatement represents: x = expr, x += 1, x++, x--
type AssignStatement struct {
	Pos
	Name     string
	Operator TokenType  // TOKEN_EQ, TOKEN_PLUS_EQ, TOKEN_MINUS_EQ, TOKEN_PLUS_PLUS, TOKEN_MINUS_MINUS
	Value    Expression // nil for ++ and --
}

func (s *AssignStatement) node()      {}
func (s *AssignStatement) statement() {}

// IfStatement represents: IF cond THEN ... [ELSEIF cond THEN ...] [ELSE ...] ENDIF
type IfStatement struct {
	Pos
	Condition     Expression
	ThenBlock     []Statement
	ElseIfClauses []ElseIfClause
	ElseBlock     []Statement // nil if no ELSE
}

func (s *IfStatement) node()      {}
func (s *IfStatement) statement() {}

// ElseIfClause represents a single ELSEIF branch
type ElseIfClause struct {
	Pos
	Condition Expression
	Block     []Statement
}

// ForStatement represents: FOR i = start TO end ... NEXT i
type ForStatement struct {
	Pos
	Variable string
	Start    Expression
	End      Expression
	Body     []Statement
}

func (s *ForStatement) node()      {}
func (s *ForStatement) statement() {}

// BreakStatement represents: BREAK
type BreakStatement struct {
	Pos
}

func (s *BreakStatement) node()      {}
func (s *BreakStatement) statement() {}

// FunctionStatement represents: FUNCTION name(params): ... ENDFUNCTION
type FunctionStatement struct {
	Pos
	Name   string
	Params []string
	Body   []Statement
}

func (s *FunctionStatement) node()      {}
func (s *FunctionStatement) statement() {}

// ReturnStatement represents: RETURN expr
type ReturnStatement struct {
	Pos
	Value Expression // nil for bare RETURN
}

func (s *ReturnStatement) node()      {}
func (s *ReturnStatement) statement() {}

// PrintStatement represents: PRINT expr
type PrintStatement struct {
	Pos
	Value Expression
}

func (s *PrintStatement) node()      {}
func (s *PrintStatement) statement() {}

// ExpressionStatement wraps an expression used as a statement (e.g., function call)
type ExpressionStatement struct {
	Pos
	Expr Expression
}

func (s *ExpressionStatement) node()      {}
func (s *ExpressionStatement) statement() {}

// -----------------------------------------------------------------------------
// Expressions
// -----------------------------------------------------------------------------

// IntLiteral represents an integer literal: 42
type IntLiteral struct {
	Pos
	Value int64
}

func (e *IntLiteral) node()       {}
func (e *IntLiteral) expression() {}

// FloatLiteral represents a float literal: 3.14
type FloatLiteral struct {
	Pos
	Value float64
}

func (e *FloatLiteral) node()       {}
func (e *FloatLiteral) expression() {}

// StringLiteral represents a string literal: "hello"
type StringLiteral struct {
	Pos
	Value string
}

func (e *StringLiteral) node()       {}
func (e *StringLiteral) expression() {}

// BoolLiteral represents: true, false
type BoolLiteral struct {
	Pos
	Value bool
}

func (e *BoolLiteral) node()       {}
func (e *BoolLiteral) expression() {}

// Identifier represents a variable reference: x, foo
type Identifier struct {
	Pos
	Name string
}

func (e *Identifier) node()       {}
func (e *Identifier) expression() {}

// BinaryExpr represents: left op right (e.g., x + y, a < b, x AND y)
type BinaryExpr struct {
	Pos
	Left     Expression
	Operator TokenType
	Right    Expression
}

func (e *BinaryExpr) node()       {}
func (e *BinaryExpr) expression() {}

// UnaryExpr represents: op expr (e.g., NOT x, -5)
type UnaryExpr struct {
	Pos
	Operator TokenType
	Operand  Expression
}

func (e *UnaryExpr) node()       {}
func (e *UnaryExpr) expression() {}

// CallExpr represents a function call: pow(2, 3), getX()
type CallExpr struct {
	Pos
	Name string
	Args []Expression
}

func (e *CallExpr) node()       {}
func (e *CallExpr) expression() {}

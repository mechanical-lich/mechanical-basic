package basic

import (
	"testing"
)

// TestNodeInterfaces verifies all node types implement the correct interfaces
func TestNodeInterfaces(t *testing.T) {
	// Statements must implement Statement (and therefore Node)
	var _ Statement = &LetStatement{}
	var _ Statement = &AssignStatement{}
	var _ Statement = &IfStatement{}
	var _ Statement = &ForStatement{}
	var _ Statement = &BreakStatement{}
	var _ Statement = &FunctionStatement{}
	var _ Statement = &ReturnStatement{}
	var _ Statement = &PrintStatement{}
	var _ Statement = &ExpressionStatement{}

	// Expressions must implement Expression (and therefore Node)
	var _ Expression = &IntLiteral{}
	var _ Expression = &FloatLiteral{}
	var _ Expression = &StringLiteral{}
	var _ Expression = &BoolLiteral{}
	var _ Expression = &Identifier{}
	var _ Expression = &BinaryExpr{}
	var _ Expression = &UnaryExpr{}
	var _ Expression = &CallExpr{}

	// Program is a Node
	var _ Node = &Program{}
}

func TestPositionTracking(t *testing.T) {
	pos := Pos{Line: 10, Column: 5}

	line, col := pos.Position()
	if line != 10 {
		t.Errorf("expected line 10, got %d", line)
	}
	if col != 5 {
		t.Errorf("expected column 5, got %d", col)
	}
}

func TestLetStatement(t *testing.T) {
	stmt := &LetStatement{
		Pos:   Pos{Line: 1, Column: 1},
		Name:  "x",
		Value: &IntLiteral{Pos: Pos{Line: 1, Column: 9}, Value: 42},
	}

	if stmt.Name != "x" {
		t.Errorf("expected name 'x', got %q", stmt.Name)
	}

	line, col := stmt.Position()
	if line != 1 || col != 1 {
		t.Errorf("expected position (1, 1), got (%d, %d)", line, col)
	}

	intLit, ok := stmt.Value.(*IntLiteral)
	if !ok {
		t.Fatal("expected IntLiteral")
	}
	if intLit.Value != 42 {
		t.Errorf("expected value 42, got %d", intLit.Value)
	}
}

func TestAssignStatement(t *testing.T) {
	// x = 5
	stmt := &AssignStatement{
		Pos:      Pos{Line: 1, Column: 1},
		Name:     "x",
		Operator: TOKEN_EQ,
		Value:    &IntLiteral{Pos: Pos{Line: 1, Column: 5}, Value: 5},
	}

	if stmt.Name != "x" {
		t.Errorf("expected name 'x', got %q", stmt.Name)
	}
	if stmt.Operator != TOKEN_EQ {
		t.Errorf("expected TOKEN_EQ, got %s", stmt.Operator)
	}

	// x++
	incStmt := &AssignStatement{
		Pos:      Pos{Line: 2, Column: 1},
		Name:     "x",
		Operator: TOKEN_PLUS_PLUS,
		Value:    nil,
	}

	if incStmt.Value != nil {
		t.Error("expected nil Value for increment")
	}
}

func TestIfStatement(t *testing.T) {
	stmt := &IfStatement{
		Pos: Pos{Line: 1, Column: 1},
		Condition: &BinaryExpr{
			Pos:      Pos{Line: 1, Column: 4},
			Left:     &Identifier{Pos: Pos{Line: 1, Column: 4}, Name: "x"},
			Operator: TOKEN_GT,
			Right:    &IntLiteral{Pos: Pos{Line: 1, Column: 8}, Value: 5},
		},
		ThenBlock: []Statement{
			&PrintStatement{
				Pos:   Pos{Line: 2, Column: 5},
				Value: &StringLiteral{Pos: Pos{Line: 2, Column: 11}, Value: "big"},
			},
		},
		ElseIfClauses: []ElseIfClause{
			{
				Pos: Pos{Line: 3, Column: 1},
				Condition: &BinaryExpr{
					Pos:      Pos{Line: 3, Column: 8},
					Left:     &Identifier{Pos: Pos{Line: 3, Column: 8}, Name: "x"},
					Operator: TOKEN_LT,
					Right:    &IntLiteral{Pos: Pos{Line: 3, Column: 12}, Value: 0},
				},
				Block: []Statement{
					&PrintStatement{
						Pos:   Pos{Line: 4, Column: 5},
						Value: &StringLiteral{Pos: Pos{Line: 4, Column: 11}, Value: "negative"},
					},
				},
			},
		},
		ElseBlock: []Statement{
			&PrintStatement{
				Pos:   Pos{Line: 6, Column: 5},
				Value: &StringLiteral{Pos: Pos{Line: 6, Column: 11}, Value: "small"},
			},
		},
	}

	if len(stmt.ThenBlock) != 1 {
		t.Errorf("expected 1 statement in ThenBlock, got %d", len(stmt.ThenBlock))
	}
	if len(stmt.ElseIfClauses) != 1 {
		t.Errorf("expected 1 ElseIfClause, got %d", len(stmt.ElseIfClauses))
	}
	if len(stmt.ElseBlock) != 1 {
		t.Errorf("expected 1 statement in ElseBlock, got %d", len(stmt.ElseBlock))
	}
}

func TestForStatement(t *testing.T) {
	stmt := &ForStatement{
		Pos:      Pos{Line: 1, Column: 1},
		Variable: "i",
		Start:    &IntLiteral{Pos: Pos{Line: 1, Column: 9}, Value: 1},
		End:      &IntLiteral{Pos: Pos{Line: 1, Column: 14}, Value: 10},
		Body: []Statement{
			&PrintStatement{
				Pos:   Pos{Line: 2, Column: 5},
				Value: &Identifier{Pos: Pos{Line: 2, Column: 11}, Name: "i"},
			},
		},
	}

	if stmt.Variable != "i" {
		t.Errorf("expected variable 'i', got %q", stmt.Variable)
	}

	startLit := stmt.Start.(*IntLiteral)
	if startLit.Value != 1 {
		t.Errorf("expected start 1, got %d", startLit.Value)
	}

	endLit := stmt.End.(*IntLiteral)
	if endLit.Value != 10 {
		t.Errorf("expected end 10, got %d", endLit.Value)
	}
}

func TestFunctionStatement(t *testing.T) {
	stmt := &FunctionStatement{
		Pos:    Pos{Line: 1, Column: 1},
		Name:   "add",
		Params: []string{"x", "y"},
		Body: []Statement{
			&ReturnStatement{
				Pos: Pos{Line: 2, Column: 5},
				Value: &BinaryExpr{
					Pos:      Pos{Line: 2, Column: 12},
					Left:     &Identifier{Pos: Pos{Line: 2, Column: 12}, Name: "x"},
					Operator: TOKEN_PLUS,
					Right:    &Identifier{Pos: Pos{Line: 2, Column: 16}, Name: "y"},
				},
			},
		},
	}

	if stmt.Name != "add" {
		t.Errorf("expected name 'add', got %q", stmt.Name)
	}
	if len(stmt.Params) != 2 {
		t.Errorf("expected 2 params, got %d", len(stmt.Params))
	}
	if stmt.Params[0] != "x" || stmt.Params[1] != "y" {
		t.Errorf("expected params [x, y], got %v", stmt.Params)
	}
}

func TestBinaryExpr(t *testing.T) {
	expr := &BinaryExpr{
		Pos:      Pos{Line: 1, Column: 1},
		Left:     &IntLiteral{Pos: Pos{Line: 1, Column: 1}, Value: 2},
		Operator: TOKEN_PLUS,
		Right:    &IntLiteral{Pos: Pos{Line: 1, Column: 5}, Value: 3},
	}

	if expr.Operator != TOKEN_PLUS {
		t.Errorf("expected TOKEN_PLUS, got %s", expr.Operator)
	}

	left := expr.Left.(*IntLiteral)
	right := expr.Right.(*IntLiteral)
	if left.Value != 2 || right.Value != 3 {
		t.Errorf("expected 2 + 3, got %d + %d", left.Value, right.Value)
	}
}

func TestUnaryExpr(t *testing.T) {
	// NOT true
	expr := &UnaryExpr{
		Pos:      Pos{Line: 1, Column: 1},
		Operator: TOKEN_NOT,
		Operand:  &BoolLiteral{Pos: Pos{Line: 1, Column: 5}, Value: true},
	}

	if expr.Operator != TOKEN_NOT {
		t.Errorf("expected TOKEN_NOT, got %s", expr.Operator)
	}

	operand := expr.Operand.(*BoolLiteral)
	if !operand.Value {
		t.Error("expected true")
	}

	// -5
	negExpr := &UnaryExpr{
		Pos:      Pos{Line: 1, Column: 1},
		Operator: TOKEN_MINUS,
		Operand:  &IntLiteral{Pos: Pos{Line: 1, Column: 2}, Value: 5},
	}

	if negExpr.Operator != TOKEN_MINUS {
		t.Errorf("expected TOKEN_MINUS, got %s", negExpr.Operator)
	}
}

func TestCallExpr(t *testing.T) {
	// pow(2, 3)
	expr := &CallExpr{
		Pos:  Pos{Line: 1, Column: 1},
		Name: "pow",
		Args: []Expression{
			&IntLiteral{Pos: Pos{Line: 1, Column: 5}, Value: 2},
			&IntLiteral{Pos: Pos{Line: 1, Column: 8}, Value: 3},
		},
	}

	if expr.Name != "pow" {
		t.Errorf("expected name 'pow', got %q", expr.Name)
	}
	if len(expr.Args) != 2 {
		t.Errorf("expected 2 args, got %d", len(expr.Args))
	}

	// getX()
	noArgsExpr := &CallExpr{
		Pos:  Pos{Line: 1, Column: 1},
		Name: "getX",
		Args: []Expression{},
	}

	if noArgsExpr.Name != "getX" {
		t.Errorf("expected name 'getX', got %q", noArgsExpr.Name)
	}
	if len(noArgsExpr.Args) != 0 {
		t.Errorf("expected 0 args, got %d", len(noArgsExpr.Args))
	}
}

func TestLiterals(t *testing.T) {
	intLit := &IntLiteral{Pos: Pos{Line: 1, Column: 1}, Value: 42}
	if intLit.Value != 42 {
		t.Errorf("expected 42, got %d", intLit.Value)
	}

	floatLit := &FloatLiteral{Pos: Pos{Line: 1, Column: 1}, Value: 3.14}
	if floatLit.Value != 3.14 {
		t.Errorf("expected 3.14, got %f", floatLit.Value)
	}

	strLit := &StringLiteral{Pos: Pos{Line: 1, Column: 1}, Value: "hello"}
	if strLit.Value != "hello" {
		t.Errorf("expected 'hello', got %q", strLit.Value)
	}

	trueLit := &BoolLiteral{Pos: Pos{Line: 1, Column: 1}, Value: true}
	if !trueLit.Value {
		t.Error("expected true")
	}

	falseLit := &BoolLiteral{Pos: Pos{Line: 1, Column: 1}, Value: false}
	if falseLit.Value {
		t.Error("expected false")
	}
}

func TestProgram(t *testing.T) {
	prog := &Program{
		Statements: []Statement{
			&LetStatement{
				Pos:   Pos{Line: 1, Column: 1},
				Name:  "x",
				Value: &IntLiteral{Pos: Pos{Line: 1, Column: 9}, Value: 5},
			},
			&PrintStatement{
				Pos:   Pos{Line: 2, Column: 1},
				Value: &Identifier{Pos: Pos{Line: 2, Column: 7}, Name: "x"},
			},
		},
	}

	if len(prog.Statements) != 2 {
		t.Errorf("expected 2 statements, got %d", len(prog.Statements))
	}

	// Program position comes from first statement
	line, col := prog.Position()
	if line != 1 || col != 1 {
		t.Errorf("expected position (1, 1), got (%d, %d)", line, col)
	}

	// Empty program
	emptyProg := &Program{Statements: []Statement{}}
	line, col = emptyProg.Position()
	if line != 1 || col != 1 {
		t.Errorf("expected default position (1, 1), got (%d, %d)", line, col)
	}
}

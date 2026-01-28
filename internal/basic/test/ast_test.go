package basic

import (
	"testing"

	"github.com/mechanical-lich/mechanical-basic/internal/basic"
)

// TestNodeInterfaces verifies all node types implement the correct interfaces
func TestNodeInterfaces(t *testing.T) {
	// Statements must implement Statement (and therefore Node)
	var _ basic.Statement = &basic.LetStatement{}
	var _ basic.Statement = &basic.AssignStatement{}
	var _ basic.Statement = &basic.IfStatement{}
	var _ basic.Statement = &basic.ForStatement{}
	var _ basic.Statement = &basic.BreakStatement{}
	var _ basic.Statement = &basic.FunctionStatement{}
	var _ basic.Statement = &basic.ReturnStatement{}
	var _ basic.Statement = &basic.PrintStatement{}
	var _ basic.Statement = &basic.ExpressionStatement{}

	// Expressions must implement Expression (and therefore Node)
	var _ basic.Expression = &basic.IntLiteral{}
	var _ basic.Expression = &basic.FloatLiteral{}
	var _ basic.Expression = &basic.StringLiteral{}
	var _ basic.Expression = &basic.BoolLiteral{}
	var _ basic.Expression = &basic.Identifier{}
	var _ basic.Expression = &basic.BinaryExpr{}
	var _ basic.Expression = &basic.UnaryExpr{}
	var _ basic.Expression = &basic.CallExpr{}

	// Program is a Node
	var _ basic.Node = &basic.Program{}
}

func TestPositionTracking(t *testing.T) {
	pos := basic.Pos{Line: 10, Column: 5}

	line, col := pos.Position()
	if line != 10 {
		t.Errorf("expected line 10, got %d", line)
	}
	if col != 5 {
		t.Errorf("expected column 5, got %d", col)
	}
}

func TestLetStatement(t *testing.T) {
	stmt := &basic.LetStatement{
		Pos:   basic.Pos{Line: 1, Column: 1},
		Name:  "x",
		Value: &basic.IntLiteral{Pos: basic.Pos{Line: 1, Column: 9}, Value: 42},
	}

	if stmt.Name != "x" {
		t.Errorf("expected name 'x', got %q", stmt.Name)
	}

	line, col := stmt.Position()
	if line != 1 || col != 1 {
		t.Errorf("expected position (1, 1), got (%d, %d)", line, col)
	}

	intLit, ok := stmt.Value.(*basic.IntLiteral)
	if !ok {
		t.Fatal("expected IntLiteral")
	}
	if intLit.Value != 42 {
		t.Errorf("expected value 42, got %d", intLit.Value)
	}
}

func TestAssignStatement(t *testing.T) {
	// x = 5
	stmt := &basic.AssignStatement{
		Pos:      basic.Pos{Line: 1, Column: 1},
		Name:     "x",
		Operator: basic.TOKEN_EQ,
		Value:    &basic.IntLiteral{Pos: basic.Pos{Line: 1, Column: 5}, Value: 5},
	}

	if stmt.Name != "x" {
		t.Errorf("expected name 'x', got %q", stmt.Name)
	}
	if stmt.Operator != basic.TOKEN_EQ {
		t.Errorf("expected TOKEN_EQ, got %s", stmt.Operator)
	}

	// x++
	incStmt := &basic.AssignStatement{
		Pos:      basic.Pos{Line: 2, Column: 1},
		Name:     "x",
		Operator: basic.TOKEN_PLUS_PLUS,
		Value:    nil,
	}

	if incStmt.Value != nil {
		t.Error("expected nil Value for increment")
	}
}

func TestIfStatement(t *testing.T) {
	stmt := &basic.IfStatement{
		Pos: basic.Pos{Line: 1, Column: 1},
		Condition: &basic.BinaryExpr{
			Pos:      basic.Pos{Line: 1, Column: 4},
			Left:     &basic.Identifier{Pos: basic.Pos{Line: 1, Column: 4}, Name: "x"},
			Operator: basic.TOKEN_GT,
			Right:    &basic.IntLiteral{Pos: basic.Pos{Line: 1, Column: 8}, Value: 5},
		},
		ThenBlock: []basic.Statement{
			&basic.PrintStatement{
				Pos:   basic.Pos{Line: 2, Column: 5},
				Value: &basic.StringLiteral{Pos: basic.Pos{Line: 2, Column: 11}, Value: "big"},
			},
		},
		ElseIfClauses: []basic.ElseIfClause{
			{
				Pos: basic.Pos{Line: 3, Column: 1},
				Condition: &basic.BinaryExpr{
					Pos:      basic.Pos{Line: 3, Column: 8},
					Left:     &basic.Identifier{Pos: basic.Pos{Line: 3, Column: 8}, Name: "x"},
					Operator: basic.TOKEN_LT,
					Right:    &basic.IntLiteral{Pos: basic.Pos{Line: 3, Column: 12}, Value: 0},
				},
				Block: []basic.Statement{
					&basic.PrintStatement{
						Pos:   basic.Pos{Line: 4, Column: 5},
						Value: &basic.StringLiteral{Pos: basic.Pos{Line: 4, Column: 11}, Value: "negative"},
					},
				},
			},
		},
		ElseBlock: []basic.Statement{
			&basic.PrintStatement{
				Pos:   basic.Pos{Line: 6, Column: 5},
				Value: &basic.StringLiteral{Pos: basic.Pos{Line: 6, Column: 11}, Value: "small"},
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
	stmt := &basic.ForStatement{
		Pos:      basic.Pos{Line: 1, Column: 1},
		Variable: "i",
		Start:    &basic.IntLiteral{Pos: basic.Pos{Line: 1, Column: 9}, Value: 1},
		End:      &basic.IntLiteral{Pos: basic.Pos{Line: 1, Column: 14}, Value: 10},
		Body: []basic.Statement{
			&basic.PrintStatement{
				Pos:   basic.Pos{Line: 2, Column: 5},
				Value: &basic.Identifier{Pos: basic.Pos{Line: 2, Column: 11}, Name: "i"},
			},
		},
	}

	if stmt.Variable != "i" {
		t.Errorf("expected variable 'i', got %q", stmt.Variable)
	}

	startLit := stmt.Start.(*basic.IntLiteral)
	if startLit.Value != 1 {
		t.Errorf("expected start 1, got %d", startLit.Value)
	}

	endLit := stmt.End.(*basic.IntLiteral)
	if endLit.Value != 10 {
		t.Errorf("expected end 10, got %d", endLit.Value)
	}
}

func TestFunctionStatement(t *testing.T) {
	stmt := &basic.FunctionStatement{
		Pos:    basic.Pos{Line: 1, Column: 1},
		Name:   "add",
		Params: []string{"x", "y"},
		Body: []basic.Statement{
			&basic.ReturnStatement{
				Pos: basic.Pos{Line: 2, Column: 5},
				Value: &basic.BinaryExpr{
					Pos:      basic.Pos{Line: 2, Column: 12},
					Left:     &basic.Identifier{Pos: basic.Pos{Line: 2, Column: 12}, Name: "x"},
					Operator: basic.TOKEN_PLUS,
					Right:    &basic.Identifier{Pos: basic.Pos{Line: 2, Column: 16}, Name: "y"},
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
	expr := &basic.BinaryExpr{
		Pos:      basic.Pos{Line: 1, Column: 1},
		Left:     &basic.IntLiteral{Pos: basic.Pos{Line: 1, Column: 1}, Value: 2},
		Operator: basic.TOKEN_PLUS,
		Right:    &basic.IntLiteral{Pos: basic.Pos{Line: 1, Column: 5}, Value: 3},
	}

	if expr.Operator != basic.TOKEN_PLUS {
		t.Errorf("expected TOKEN_PLUS, got %s", expr.Operator)
	}

	left := expr.Left.(*basic.IntLiteral)
	right := expr.Right.(*basic.IntLiteral)
	if left.Value != 2 || right.Value != 3 {
		t.Errorf("expected 2 + 3, got %d + %d", left.Value, right.Value)
	}
}

func TestUnaryExpr(t *testing.T) {
	// NOT true
	expr := &basic.UnaryExpr{
		Pos:      basic.Pos{Line: 1, Column: 1},
		Operator: basic.TOKEN_NOT,
		Operand:  &basic.BoolLiteral{Pos: basic.Pos{Line: 1, Column: 5}, Value: true},
	}

	if expr.Operator != basic.TOKEN_NOT {
		t.Errorf("expected TOKEN_NOT, got %s", expr.Operator)
	}

	operand := expr.Operand.(*basic.BoolLiteral)
	if !operand.Value {
		t.Error("expected true")
	}

	// -5
	negExpr := &basic.UnaryExpr{
		Pos:      basic.Pos{Line: 1, Column: 1},
		Operator: basic.TOKEN_MINUS,
		Operand:  &basic.IntLiteral{Pos: basic.Pos{Line: 1, Column: 2}, Value: 5},
	}

	if negExpr.Operator != basic.TOKEN_MINUS {
		t.Errorf("expected TOKEN_MINUS, got %s", negExpr.Operator)
	}
}

func TestCallExpr(t *testing.T) {
	// pow(2, 3)
	expr := &basic.CallExpr{
		Pos:  basic.Pos{Line: 1, Column: 1},
		Name: "pow",
		Args: []basic.Expression{
			&basic.IntLiteral{Pos: basic.Pos{Line: 1, Column: 5}, Value: 2},
			&basic.IntLiteral{Pos: basic.Pos{Line: 1, Column: 8}, Value: 3},
		},
	}

	if expr.Name != "pow" {
		t.Errorf("expected name 'pow', got %q", expr.Name)
	}
	if len(expr.Args) != 2 {
		t.Errorf("expected 2 args, got %d", len(expr.Args))
	}

	// getX()
	noArgsExpr := &basic.CallExpr{
		Pos:  basic.Pos{Line: 1, Column: 1},
		Name: "getX",
		Args: []basic.Expression{},
	}

	if noArgsExpr.Name != "getX" {
		t.Errorf("expected name 'getX', got %q", noArgsExpr.Name)
	}
	if len(noArgsExpr.Args) != 0 {
		t.Errorf("expected 0 args, got %d", len(noArgsExpr.Args))
	}
}

func TestLiterals(t *testing.T) {
	intLit := &basic.IntLiteral{Pos: basic.Pos{Line: 1, Column: 1}, Value: 42}
	if intLit.Value != 42 {
		t.Errorf("expected 42, got %d", intLit.Value)
	}

	floatLit := &basic.FloatLiteral{Pos: basic.Pos{Line: 1, Column: 1}, Value: 3.14}
	if floatLit.Value != 3.14 {
		t.Errorf("expected 3.14, got %f", floatLit.Value)
	}

	strLit := &basic.StringLiteral{Pos: basic.Pos{Line: 1, Column: 1}, Value: "hello"}
	if strLit.Value != "hello" {
		t.Errorf("expected 'hello', got %q", strLit.Value)
	}

	trueLit := &basic.BoolLiteral{Pos: basic.Pos{Line: 1, Column: 1}, Value: true}
	if !trueLit.Value {
		t.Error("expected true")
	}

	falseLit := &basic.BoolLiteral{Pos: basic.Pos{Line: 1, Column: 1}, Value: false}
	if falseLit.Value {
		t.Error("expected false")
	}
}

func TestProgram(t *testing.T) {
	prog := &basic.Program{
		Statements: []basic.Statement{
			&basic.LetStatement{
				Pos:   basic.Pos{Line: 1, Column: 1},
				Name:  "x",
				Value: &basic.IntLiteral{Pos: basic.Pos{Line: 1, Column: 9}, Value: 5},
			},
			&basic.PrintStatement{
				Pos:   basic.Pos{Line: 2, Column: 1},
				Value: &basic.Identifier{Pos: basic.Pos{Line: 2, Column: 7}, Name: "x"},
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
	emptyProg := &basic.Program{Statements: []basic.Statement{}}
	line, col = emptyProg.Position()
	if line != 1 || col != 1 {
		t.Errorf("expected default position (1, 1), got (%d, %d)", line, col)
	}
}

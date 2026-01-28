package basic

import (
	"testing"

	"github.com/mechanical-lich/mechanical-basic/internal/basic"
)

func parseCode(t *testing.T, code string) *basic.Program {
	tokens, err := basic.Tokenize(code)
	if err != nil {
		t.Fatalf("tokenize error: %v", err)
	}
	prog, err := basic.Parse(tokens)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	return prog
}

func TestParseLetStatement(t *testing.T) {
	prog := parseCode(t, "LET x = 5")

	if len(prog.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(prog.Statements))
	}

	let, ok := prog.Statements[0].(*basic.LetStatement)
	if !ok {
		t.Fatalf("expected LetStatement, got %T", prog.Statements[0])
	}

	if let.Name != "x" {
		t.Errorf("expected name 'x', got %q", let.Name)
	}

	intLit, ok := let.Value.(*basic.IntLiteral)
	if !ok {
		t.Fatalf("expected IntLiteral, got %T", let.Value)
	}
	if intLit.Value != 5 {
		t.Errorf("expected value 5, got %d", intLit.Value)
	}
}

func TestParseLetWithExpression(t *testing.T) {
	prog := parseCode(t, "LET x = 2 + 3 * 4")

	let := prog.Statements[0].(*basic.LetStatement)

	// Should be 2 + (3 * 4) due to precedence
	binExpr, ok := let.Value.(*basic.BinaryExpr)
	if !ok {
		t.Fatalf("expected BinaryExpr, got %T", let.Value)
	}

	if binExpr.Operator != basic.TOKEN_PLUS {
		t.Errorf("expected PLUS at top level, got %s", binExpr.Operator)
	}

	// Left should be 2
	leftInt := binExpr.Left.(*basic.IntLiteral)
	if leftInt.Value != 2 {
		t.Errorf("expected left 2, got %d", leftInt.Value)
	}

	// Right should be 3 * 4
	rightBin := binExpr.Right.(*basic.BinaryExpr)
	if rightBin.Operator != basic.TOKEN_STAR {
		t.Errorf("expected STAR on right, got %s", rightBin.Operator)
	}
}

func TestParseAssignment(t *testing.T) {
	prog := parseCode(t, "x = 10")

	assign, ok := prog.Statements[0].(*basic.AssignStatement)
	if !ok {
		t.Fatalf("expected AssignStatement, got %T", prog.Statements[0])
	}

	if assign.Name != "x" {
		t.Errorf("expected name 'x', got %q", assign.Name)
	}
	if assign.Operator != basic.TOKEN_EQ {
		t.Errorf("expected TOKEN_EQ, got %s", assign.Operator)
	}
}

func TestParseCompoundAssignment(t *testing.T) {
	tests := []struct {
		code string
		op   basic.TokenType
	}{
		{"x += 1", basic.TOKEN_PLUS_EQ},
		{"x -= 1", basic.TOKEN_MINUS_EQ},
	}

	for _, tt := range tests {
		prog := parseCode(t, tt.code)
		assign := prog.Statements[0].(*basic.AssignStatement)
		if assign.Operator != tt.op {
			t.Errorf("%s: expected %s, got %s", tt.code, tt.op, assign.Operator)
		}
	}
}

func TestParseIncDec(t *testing.T) {
	prog := parseCode(t, "x++")
	assign := prog.Statements[0].(*basic.AssignStatement)
	if assign.Operator != basic.TOKEN_PLUS_PLUS {
		t.Errorf("expected TOKEN_PLUS_PLUS, got %s", assign.Operator)
	}
	if assign.Value != nil {
		t.Error("expected nil Value for ++")
	}

	prog = parseCode(t, "x--")
	assign = prog.Statements[0].(*basic.AssignStatement)
	if assign.Operator != basic.TOKEN_MINUS_MINUS {
		t.Errorf("expected TOKEN_MINUS_MINUS, got %s", assign.Operator)
	}
}

func TestParsePrint(t *testing.T) {
	prog := parseCode(t, `print "Hello"`)

	print, ok := prog.Statements[0].(*basic.PrintStatement)
	if !ok {
		t.Fatalf("expected PrintStatement, got %T", prog.Statements[0])
	}

	str, ok := print.Value.(*basic.StringLiteral)
	if !ok {
		t.Fatalf("expected StringLiteral, got %T", print.Value)
	}
	if str.Value != "Hello" {
		t.Errorf("expected 'Hello', got %q", str.Value)
	}
}

func TestParseIfThenEndif(t *testing.T) {
	code := `if x > 5 then
    print "big"
endif`
	prog := parseCode(t, code)

	ifStmt, ok := prog.Statements[0].(*basic.IfStatement)
	if !ok {
		t.Fatalf("expected IfStatement, got %T", prog.Statements[0])
	}

	// Check condition
	binExpr, ok := ifStmt.Condition.(*basic.BinaryExpr)
	if !ok {
		t.Fatalf("expected BinaryExpr condition, got %T", ifStmt.Condition)
	}
	if binExpr.Operator != basic.TOKEN_GT {
		t.Errorf("expected GT, got %s", binExpr.Operator)
	}

	// Check then block
	if len(ifStmt.ThenBlock) != 1 {
		t.Errorf("expected 1 statement in ThenBlock, got %d", len(ifStmt.ThenBlock))
	}

	// No else
	if len(ifStmt.ElseBlock) != 0 {
		t.Errorf("expected no ElseBlock, got %d", len(ifStmt.ElseBlock))
	}
}

func TestParseIfElse(t *testing.T) {
	code := `if x > 5 then
    print "big"
else
    print "small"
endif`
	prog := parseCode(t, code)

	ifStmt := prog.Statements[0].(*basic.IfStatement)

	if len(ifStmt.ThenBlock) != 1 {
		t.Errorf("expected 1 statement in ThenBlock, got %d", len(ifStmt.ThenBlock))
	}
	if len(ifStmt.ElseBlock) != 1 {
		t.Errorf("expected 1 statement in ElseBlock, got %d", len(ifStmt.ElseBlock))
	}
}

func TestParseIfElseIf(t *testing.T) {
	code := `if x > 5 then
    print "big"
elseif x < 0 then
    print "negative"
else
    print "small"
endif`
	prog := parseCode(t, code)

	ifStmt := prog.Statements[0].(*basic.IfStatement)

	if len(ifStmt.ElseIfClauses) != 1 {
		t.Errorf("expected 1 ElseIfClause, got %d", len(ifStmt.ElseIfClauses))
	}

	elseIf := ifStmt.ElseIfClauses[0]
	binExpr := elseIf.Condition.(*basic.BinaryExpr)
	if binExpr.Operator != basic.TOKEN_LT {
		t.Errorf("expected LT in elseif, got %s", binExpr.Operator)
	}
}

func TestParseForLoop(t *testing.T) {
	code := `for i = 1 to 10
    print i
next i`
	prog := parseCode(t, code)

	forStmt, ok := prog.Statements[0].(*basic.ForStatement)
	if !ok {
		t.Fatalf("expected ForStatement, got %T", prog.Statements[0])
	}

	if forStmt.Variable != "i" {
		t.Errorf("expected variable 'i', got %q", forStmt.Variable)
	}

	start := forStmt.Start.(*basic.IntLiteral)
	if start.Value != 1 {
		t.Errorf("expected start 1, got %d", start.Value)
	}

	end := forStmt.End.(*basic.IntLiteral)
	if end.Value != 10 {
		t.Errorf("expected end 10, got %d", end.Value)
	}

	if len(forStmt.Body) != 1 {
		t.Errorf("expected 1 statement in body, got %d", len(forStmt.Body))
	}
}

func TestParseForWithBreak(t *testing.T) {
	code := `for i = 1 to 10
    if i = 5 then
        break
    endif
next i`
	prog := parseCode(t, code)

	forStmt := prog.Statements[0].(*basic.ForStatement)
	if len(forStmt.Body) != 1 {
		t.Fatalf("expected 1 statement in for body, got %d", len(forStmt.Body))
	}

	ifStmt := forStmt.Body[0].(*basic.IfStatement)
	breakStmt, ok := ifStmt.ThenBlock[0].(*basic.BreakStatement)
	if !ok {
		t.Fatalf("expected BreakStatement, got %T", ifStmt.ThenBlock[0])
	}
	_ = breakStmt
}

func TestParseFunction(t *testing.T) {
	code := `function add(x, y):
    let z = x + y
    return z
endfunction`
	prog := parseCode(t, code)

	fn, ok := prog.Statements[0].(*basic.FunctionStatement)
	if !ok {
		t.Fatalf("expected FunctionStatement, got %T", prog.Statements[0])
	}

	if fn.Name != "add" {
		t.Errorf("expected name 'add', got %q", fn.Name)
	}

	if len(fn.Params) != 2 || fn.Params[0] != "x" || fn.Params[1] != "y" {
		t.Errorf("expected params [x, y], got %v", fn.Params)
	}

	if len(fn.Body) != 2 {
		t.Errorf("expected 2 statements in body, got %d", len(fn.Body))
	}

	// Check return
	ret, ok := fn.Body[1].(*basic.ReturnStatement)
	if !ok {
		t.Fatalf("expected ReturnStatement, got %T", fn.Body[1])
	}

	ident := ret.Value.(*basic.Identifier)
	if ident.Name != "z" {
		t.Errorf("expected return 'z', got %q", ident.Name)
	}
}

func TestParseFunctionNoParams(t *testing.T) {
	code := `function greet():
    print "Hello"
endfunction`
	prog := parseCode(t, code)

	fn := prog.Statements[0].(*basic.FunctionStatement)
	if len(fn.Params) != 0 {
		t.Errorf("expected 0 params, got %d", len(fn.Params))
	}
}

func TestParseFunctionCall(t *testing.T) {
	code := `let x = add(1, 2)`
	prog := parseCode(t, code)

	let := prog.Statements[0].(*basic.LetStatement)
	call, ok := let.Value.(*basic.CallExpr)
	if !ok {
		t.Fatalf("expected CallExpr, got %T", let.Value)
	}

	if call.Name != "add" {
		t.Errorf("expected name 'add', got %q", call.Name)
	}
	if len(call.Args) != 2 {
		t.Errorf("expected 2 args, got %d", len(call.Args))
	}
}

func TestParseFunctionCallNoArgs(t *testing.T) {
	code := `let x = getX()`
	prog := parseCode(t, code)

	let := prog.Statements[0].(*basic.LetStatement)
	call := let.Value.(*basic.CallExpr)

	if call.Name != "getX" {
		t.Errorf("expected name 'getX', got %q", call.Name)
	}
	if len(call.Args) != 0 {
		t.Errorf("expected 0 args, got %d", len(call.Args))
	}
}

func TestParseFunctionCallAsStatement(t *testing.T) {
	code := `doSomething(1, 2)`
	prog := parseCode(t, code)

	exprStmt, ok := prog.Statements[0].(*basic.ExpressionStatement)
	if !ok {
		t.Fatalf("expected ExpressionStatement, got %T", prog.Statements[0])
	}

	call := exprStmt.Expr.(*basic.CallExpr)
	if call.Name != "doSomething" {
		t.Errorf("expected name 'doSomething', got %q", call.Name)
	}
}

func TestParseUnaryMinus(t *testing.T) {
	prog := parseCode(t, "let x = -5")

	let := prog.Statements[0].(*basic.LetStatement)
	unary, ok := let.Value.(*basic.UnaryExpr)
	if !ok {
		t.Fatalf("expected UnaryExpr, got %T", let.Value)
	}

	if unary.Operator != basic.TOKEN_MINUS {
		t.Errorf("expected MINUS, got %s", unary.Operator)
	}

	intLit := unary.Operand.(*basic.IntLiteral)
	if intLit.Value != 5 {
		t.Errorf("expected 5, got %d", intLit.Value)
	}
}

func TestParseUnaryNot(t *testing.T) {
	prog := parseCode(t, "let x = not true")

	let := prog.Statements[0].(*basic.LetStatement)
	unary := let.Value.(*basic.UnaryExpr)

	if unary.Operator != basic.TOKEN_NOT {
		t.Errorf("expected NOT, got %s", unary.Operator)
	}

	boolLit := unary.Operand.(*basic.BoolLiteral)
	if !boolLit.Value {
		t.Error("expected true")
	}
}

func TestParseLogicalOperators(t *testing.T) {
	prog := parseCode(t, "let x = a and b or c")

	let := prog.Statements[0].(*basic.LetStatement)

	// OR has lower precedence, so it should be at the top
	orExpr, ok := let.Value.(*basic.BinaryExpr)
	if !ok {
		t.Fatalf("expected BinaryExpr, got %T", let.Value)
	}
	if orExpr.Operator != basic.TOKEN_OR {
		t.Errorf("expected OR at top, got %s", orExpr.Operator)
	}

	// Left should be (a AND b)
	andExpr := orExpr.Left.(*basic.BinaryExpr)
	if andExpr.Operator != basic.TOKEN_AND {
		t.Errorf("expected AND on left, got %s", andExpr.Operator)
	}
}

func TestParseComparisonOperators(t *testing.T) {
	tests := []struct {
		code string
		op   basic.TokenType
	}{
		{"x < 5", basic.TOKEN_LT},
		{"x > 5", basic.TOKEN_GT},
		{"x <= 5", basic.TOKEN_LTE},
		{"x >= 5", basic.TOKEN_GTE},
		{"x = 5", basic.TOKEN_EQ},
		{"x <> 5", basic.TOKEN_NEQ},
		{"x != 5", basic.TOKEN_NEQ},
	}

	for _, tt := range tests {
		prog := parseCode(t, "let y = "+tt.code)
		let := prog.Statements[0].(*basic.LetStatement)
		binExpr := let.Value.(*basic.BinaryExpr)
		if binExpr.Operator != tt.op {
			t.Errorf("%s: expected %s, got %s", tt.code, tt.op, binExpr.Operator)
		}
	}
}

func TestParseParentheses(t *testing.T) {
	prog := parseCode(t, "let x = (2 + 3) * 4")

	let := prog.Statements[0].(*basic.LetStatement)

	// Top should be multiplication
	mulExpr := let.Value.(*basic.BinaryExpr)
	if mulExpr.Operator != basic.TOKEN_STAR {
		t.Errorf("expected STAR at top, got %s", mulExpr.Operator)
	}

	// Left should be addition (from parens)
	addExpr := mulExpr.Left.(*basic.BinaryExpr)
	if addExpr.Operator != basic.TOKEN_PLUS {
		t.Errorf("expected PLUS on left, got %s", addExpr.Operator)
	}
}

func TestParseLiterals(t *testing.T) {
	tests := []struct {
		code     string
		expected interface{}
	}{
		{"let x = 42", 42},
		{"let x = 3.14", 3.14},
		{`let x = "hello"`, "hello"},
		{"let x = true", true},
		{"let x = false", false},
	}

	for _, tt := range tests {
		prog := parseCode(t, tt.code)
		let := prog.Statements[0].(*basic.LetStatement)

		switch expected := tt.expected.(type) {
		case int:
			lit := let.Value.(*basic.IntLiteral)
			if lit.Value != expected {
				t.Errorf("%s: expected %d, got %d", tt.code, expected, lit.Value)
			}
		case float64:
			lit := let.Value.(*basic.FloatLiteral)
			if lit.Value != expected {
				t.Errorf("%s: expected %f, got %f", tt.code, expected, lit.Value)
			}
		case string:
			lit := let.Value.(*basic.StringLiteral)
			if lit.Value != expected {
				t.Errorf("%s: expected %q, got %q", tt.code, expected, lit.Value)
			}
		case bool:
			lit := let.Value.(*basic.BoolLiteral)
			if lit.Value != expected {
				t.Errorf("%s: expected %v, got %v", tt.code, expected, lit.Value)
			}
		}
	}
}

func TestParseMultipleStatements(t *testing.T) {
	code := `let x = 1
let y = 2
print x + y`
	prog := parseCode(t, code)

	if len(prog.Statements) != 3 {
		t.Errorf("expected 3 statements, got %d", len(prog.Statements))
	}
}

func TestParseError(t *testing.T) {
	tests := []string{
		"let = 5",       // missing identifier
		"let x 5",       // missing =
		"if x > 5",      // missing THEN
		"if x > 5 then", // missing ENDIF
		"for i = 1",     // missing TO
		"function",      // missing name
		"function foo",  // missing parens
	}

	for _, code := range tests {
		tokens, err := basic.Tokenize(code)
		if err != nil {
			continue // tokenize error is fine
		}
		_, err = basic.Parse(tokens)
		if err == nil {
			t.Errorf("expected parse error for: %s", code)
		}
	}
}

package basic

import (
	"fmt"
	"strings"
	"testing"

	"github.com/mechanical-lich/mechanical-basic/internal/basic"
)

func newTestInterpreter() (*basic.Interpreter, *[]interface{}) {
	interp := basic.NewInterpreter()
	var output []interface{}
	interp.SetPrintFunc(func(v interface{}) {
		output = append(output, v)
	})
	return interp, &output
}

func TestInterpretLetStatement(t *testing.T) {
	interp, output := newTestInterpreter()
	err := interp.Interpret(`
let x = 5
print x
`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(*output) != 1 || (*output)[0] != (5) {
		t.Errorf("expected [5], got %v", *output)
	}
}

func TestInterpretArithmetic(t *testing.T) {
	tests := []struct {
		code     string
		expected interface{}
	}{
		{"let x = 2 + 3\nprint x", (5)},
		{"let x = 10 - 4\nprint x", (6)},
		{"let x = 3 * 4\nprint x", (12)},
		{"let x = 15 / 3\nprint x", (5)},
		{"let x = 2 + 3 * 4\nprint x", (14)},             // Precedence
		{"let x = (2 + 3) * 4\nprint x", (20)},           // Parentheses
		{"let x = 10 / 3\nprint x", (3)},                 // Integer division
		{"let x = 10.0 / 3\nprint x", float64(10.0 / 3)}, // Float division
	}

	for _, tt := range tests {
		interp, output := newTestInterpreter()
		err := interp.Interpret(tt.code)
		if err != nil {
			t.Errorf("%s: unexpected error: %v", tt.code, err)
			continue
		}
		if len(*output) != 1 || (*output)[0] != tt.expected {
			t.Errorf("%s: expected %v, got %v", tt.code, tt.expected, *output)
		}
	}
}

func TestInterpretStringConcat(t *testing.T) {
	interp, output := newTestInterpreter()
	err := interp.Interpret(`
let x = "Hello" + " " + "World"
print x
`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if (*output)[0] != "Hello World" {
		t.Errorf("expected 'Hello World', got %v", (*output)[0])
	}
}

func TestInterpretStringWithNumber(t *testing.T) {
	interp, output := newTestInterpreter()
	err := interp.Interpret(`
let x = "Value: " + 42
print x
`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if (*output)[0] != "Value: 42" {
		t.Errorf("expected 'Value: 42', got %v", (*output)[0])
	}
}

func TestInterpretAssignment(t *testing.T) {
	interp, output := newTestInterpreter()
	err := interp.Interpret(`
let x = 5
x = 10
print x
`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if (*output)[0] != (10) {
		t.Errorf("expected 10, got %v", (*output)[0])
	}
}

func TestInterpretCompoundAssignment(t *testing.T) {
	interp, output := newTestInterpreter()
	err := interp.Interpret(`
let x = 10
x += 5
print x
x -= 3
print x
`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if (*output)[0] != (15) {
		t.Errorf("expected 15, got %v", (*output)[0])
	}
	if (*output)[1] != (12) {
		t.Errorf("expected 12, got %v", (*output)[1])
	}
}

func TestInterpretIncDec(t *testing.T) {
	interp, output := newTestInterpreter()
	err := interp.Interpret(`
let x = 10
x++
print x
x--
x--
print x
`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if (*output)[0] != (11) {
		t.Errorf("expected 11, got %v", (*output)[0])
	}
	if (*output)[1] != (9) {
		t.Errorf("expected 9, got %v", (*output)[1])
	}
}

func TestInterpretUnaryMinus(t *testing.T) {
	interp, output := newTestInterpreter()
	err := interp.Interpret(`
let x = -5
print x
let y = -3.14
print y
`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if (*output)[0] != (-5) {
		t.Errorf("expected -5, got %v", (*output)[0])
	}
	if (*output)[1] != float64(-3.14) {
		t.Errorf("expected -3.14, got %v", (*output)[1])
	}
}

func TestInterpretComparison(t *testing.T) {
	tests := []struct {
		code     string
		expected bool
	}{
		{"print 5 > 3", true},
		{"print 3 > 5", false},
		{"print 5 < 3", false},
		{"print 3 < 5", true},
		{"print 5 >= 5", true},
		{"print 5 <= 5", true},
		{"print 5 = 5", true},
		{"print 5 <> 3", true},
		{"print 5 != 3", true},
	}

	for _, tt := range tests {
		interp, output := newTestInterpreter()
		err := interp.Interpret(tt.code)
		if err != nil {
			t.Errorf("%s: unexpected error: %v", tt.code, err)
			continue
		}
		if (*output)[0] != tt.expected {
			t.Errorf("%s: expected %v, got %v", tt.code, tt.expected, (*output)[0])
		}
	}
}

func TestInterpretLogical(t *testing.T) {
	tests := []struct {
		code     string
		expected bool
	}{
		{"print true and true", true},
		{"print true and false", false},
		{"print true or false", true},
		{"print false or false", false},
		{"print not true", false},
		{"print not false", true},
	}

	for _, tt := range tests {
		interp, output := newTestInterpreter()
		err := interp.Interpret(tt.code)
		if err != nil {
			t.Errorf("%s: unexpected error: %v", tt.code, err)
			continue
		}
		if (*output)[0] != tt.expected {
			t.Errorf("%s: expected %v, got %v", tt.code, tt.expected, (*output)[0])
		}
	}
}

func TestInterpretIfThen(t *testing.T) {
	interp, output := newTestInterpreter()
	err := interp.Interpret(`
let x = 10
if x > 5 then
    print "big"
endif
`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(*output) != 1 || (*output)[0] != "big" {
		t.Errorf("expected ['big'], got %v", *output)
	}
}

func TestInterpretIfElse(t *testing.T) {
	interp, output := newTestInterpreter()
	err := interp.Interpret(`
let x = 3
if x > 5 then
    print "big"
else
    print "small"
endif
`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if (*output)[0] != "small" {
		t.Errorf("expected 'small', got %v", (*output)[0])
	}
}

func TestInterpretIfElseIf(t *testing.T) {
	interp, output := newTestInterpreter()
	err := interp.Interpret(`
let x = -5
if x > 5 then
    print "big"
elseif x < 0 then
    print "negative"
else
    print "small"
endif
`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if (*output)[0] != "negative" {
		t.Errorf("expected 'negative', got %v", (*output)[0])
	}
}

func TestInterpretForLoop(t *testing.T) {
	interp, output := newTestInterpreter()
	err := interp.Interpret(`
let sum = 0
for i = 1 to 5
    sum += i
next i
print sum
`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if (*output)[0] != (15) {
		t.Errorf("expected 15, got %v", (*output)[0])
	}
}

func TestInterpretForLoopScopeDoesNotLeak(t *testing.T) {
	interp, _ := newTestInterpreter()
	err := interp.Interpret(`
for i = 1 to 3
    print i
next i
print i
`)
	if err == nil {
		t.Error("expected error for undefined variable 'i' outside loop")
	}
	if !strings.Contains(err.Error(), "undefined variable") {
		t.Errorf("expected 'undefined variable' error, got: %v", err)
	}
}

func TestInterpretForLoopBreak(t *testing.T) {
	interp, output := newTestInterpreter()
	err := interp.Interpret(`
for i = 1 to 10
    if i = 5 then
        break
    endif
    print i
next i
`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(*output) != 4 {
		t.Errorf("expected 4 outputs, got %d: %v", len(*output), *output)
	}
}

func TestInterpretInfiniteLoopProtection(t *testing.T) {
	interp, _ := newTestInterpreter()
	interp.SetMaxIterations(100)

	err := interp.Interpret(`
for i = 1 to 1000
    print i
next i
`)
	if err == nil {
		t.Error("expected error for infinite loop protection")
	}
	if !strings.Contains(err.Error(), "maximum iterations") {
		t.Errorf("expected 'maximum iterations' error, got: %v", err)
	}
}

func TestInterpretFunction(t *testing.T) {
	interp, output := newTestInterpreter()
	err := interp.Interpret(`
function add(a, b):
    return a + b
endfunction

let result = add(3, 4)
print result
`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if (*output)[0] != (7) {
		t.Errorf("expected 7, got %v", (*output)[0])
	}
}

func TestInterpretFunctionWithLocals(t *testing.T) {
	interp, output := newTestInterpreter()
	err := interp.Interpret(`
function square(n):
    let result = n * n
    return result
endfunction

print square(5)
`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if (*output)[0] != (25) {
		t.Errorf("expected 25, got %v", (*output)[0])
	}
}

func TestInterpretRecursiveFunction(t *testing.T) {
	interp, output := newTestInterpreter()
	err := interp.Interpret(`
function factorial(n):
    if n <= 1 then
        return 1
    endif
    return n * factorial(n - 1)
endfunction

print factorial(5)
`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if (*output)[0] != (120) {
		t.Errorf("expected 120, got %v", (*output)[0])
	}
}

func TestInterpretExternalFunction(t *testing.T) {
	interp, output := newTestInterpreter()

	// Register external function
	interp.RegisterFunction("getX", func(args ...interface{}) (interface{}, error) {
		return (42), nil
	})

	err := interp.Interpret(`
let x = getX()
print x
`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if (*output)[0] != (42) {
		t.Errorf("expected 42, got %v", (*output)[0])
	}
}

func TestInterpretExternalFunctionWithArgs(t *testing.T) {
	interp, output := newTestInterpreter()

	interp.RegisterFunction("pow", func(args ...interface{}) (interface{}, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("pow requires 2 arguments")
		}
		base, ok1 := args[0].(int)
		exp, ok2 := args[1].(int)
		if !ok1 || !ok2 {
			return nil, fmt.Errorf("pow requires integer arguments")
		}
		result := (1)
		for i := (0); i < exp; i++ {
			result *= base
		}
		return result, nil
	})

	err := interp.Interpret(`
let x = pow(2, 8)
print x
`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if (*output)[0] != (256) {
		t.Errorf("expected 256, got %v", (*output)[0])
	}
}

func TestInterpretExternalFunctionAsStatement(t *testing.T) {
	interp, _ := newTestInterpreter()

	called := false
	interp.RegisterFunction("doSomething", func(args ...interface{}) (interface{}, error) {
		called = true
		return nil, nil
	})

	err := interp.Interpret(`doSomething()`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Error("expected doSomething to be called")
	}
}

func TestInterpretUndefinedVariable(t *testing.T) {
	interp, _ := newTestInterpreter()
	err := interp.Interpret(`print x`)
	if err == nil {
		t.Error("expected error for undefined variable")
	}
}

func TestInterpretUndefinedFunction(t *testing.T) {
	interp, _ := newTestInterpreter()
	err := interp.Interpret(`let x = unknown()`)
	if err == nil {
		t.Error("expected error for undefined function")
	}
}

func TestInterpretDivisionByZero(t *testing.T) {
	interp, _ := newTestInterpreter()
	err := interp.Interpret(`let x = 10 / 0`)
	if err == nil {
		t.Error("expected error for division by zero")
	}
}

func TestInterpretASTCaching(t *testing.T) {
	interp, output := newTestInterpreter()

	code := `print 42`

	// First run
	err := interp.Interpret(code)
	if err != nil {
		t.Fatalf("first run error: %v", err)
	}

	// Second run should use cached AST
	err = interp.Interpret(code)
	if err != nil {
		t.Fatalf("second run error: %v", err)
	}

	if len(*output) != 2 {
		t.Errorf("expected 2 outputs, got %d", len(*output))
	}
}

func TestInterpretValidate(t *testing.T) {
	interp := basic.NewInterpreter()

	// Valid code
	err := interp.Validate(`let x = 5`)
	if err != nil {
		t.Errorf("expected valid code, got error: %v", err)
	}

	// Invalid code
	err = interp.Validate(`let = 5`)
	if err == nil {
		t.Error("expected error for invalid code")
	}
}

func TestInterpretBooleanLiterals(t *testing.T) {
	interp, output := newTestInterpreter()
	err := interp.Interpret(`
let x = true
let y = false
print x
print y
`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if (*output)[0] != true {
		t.Errorf("expected true, got %v", (*output)[0])
	}
	if (*output)[1] != false {
		t.Errorf("expected false, got %v", (*output)[1])
	}
}

func TestInterpretCaseInsensitiveVariables(t *testing.T) {
	interp, output := newTestInterpreter()
	err := interp.Interpret(`
let X = 5
print x
`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if (*output)[0] != (5) {
		t.Errorf("expected 5, got %v", (*output)[0])
	}
}

func TestInterpretCaseInsensitiveFunctions(t *testing.T) {
	interp, output := newTestInterpreter()

	interp.RegisterFunction("GetValue", func(args ...interface{}) (interface{}, error) {
		return (99), nil
	})

	err := interp.Interpret(`print getvalue()`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if (*output)[0] != (99) {
		t.Errorf("expected 99, got %v", (*output)[0])
	}
}

func TestInterpretNestedBlocks(t *testing.T) {
	interp, output := newTestInterpreter()
	err := interp.Interpret(`
for i = 1 to 3
    for j = 1 to 2
        print i * 10 + j
    next j
next i
`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []int{11, 12, 21, 22, 31, 32}
	if len(*output) != len(expected) {
		t.Fatalf("expected %d outputs, got %d", len(expected), len(*output))
	}
	for i, exp := range expected {
		if (*output)[i] != exp {
			t.Errorf("output[%d]: expected %d, got %v", i, exp, (*output)[i])
		}
	}
}

// -----------------------------------------------------------------------------
// Load/Call Tests
// -----------------------------------------------------------------------------

func TestLoadAndCall(t *testing.T) {
	interp, output := newTestInterpreter()

	code := `
function greet(name):
    print "Hello, " + name
endfunction

function add(a, b):
    return a + b
endfunction
`
	err := interp.Load(code)
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}

	// Call greet
	_, err = interp.Call("greet", "World")
	if err != nil {
		t.Fatalf("Call greet error: %v", err)
	}
	if len(*output) != 1 || (*output)[0] != "Hello, World" {
		t.Errorf("expected 'Hello, World', got %v", *output)
	}

	// Call add
	result, err := interp.Call("add", (3), (4))
	if err != nil {
		t.Fatalf("Call add error: %v", err)
	}
	if result != (7) {
		t.Errorf("expected 7, got %v", result)
	}
}

func TestLoadExecutesTopLevel(t *testing.T) {
	interp, output := newTestInterpreter()

	code := `
print "Top-level executed"

function test():
    print "Function called"
endfunction
`
	err := interp.Load(code)
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}

	// Load should have executed top-level print
	if len(*output) != 1 || (*output)[0] != "Top-level executed" {
		t.Errorf("Load did not execute top-level code, output: %v", *output)
	}

	// Now call the function
	_, err = interp.Call("test")
	if err != nil {
		t.Fatalf("Call error: %v", err)
	}
	if len(*output) != 2 || (*output)[1] != "Function called" {
		t.Errorf("expected 'Function called', got %v", *output)
	}
}

func TestCallVariablesDoNotPersist(t *testing.T) {
	interp, _ := newTestInterpreter()

	code := `
function setAndGet():
    let x = 42
    print x
    return x
endfunction

function tryGetX():
    print x
endfunction
`
	err := interp.Load(code)
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}

	// First call sets x
	result, err := interp.Call("setAndGet")
	if err != nil {
		t.Fatalf("Call setAndGet error: %v", err)
	}
	if result != (42) {
		t.Errorf("expected 42, got %v", result)
	}

	// Second call should NOT see x (fresh scope)
	_, err = interp.Call("tryGetX")
	if err == nil {
		t.Error("expected error for undefined variable x")
	}
	if !strings.Contains(err.Error(), "undefined variable") {
		t.Errorf("expected 'undefined variable' error, got: %v", err)
	}
}

func TestTopLevelVariablesPersist(t *testing.T) {
	interp := basic.NewInterpreter()

	code := `
counter = 0

function increment():
    counter = counter + 1
    return counter
endfunction

function getCounter():
    return counter
endfunction
`
	err := interp.Load(code)
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}

	// First increment
	result, err := interp.Call("increment")
	if err != nil {
		t.Fatalf("Call increment error: %v", err)
	}
	if result != 1 {
		t.Errorf("expected 1, got %v (type %T)", result, result)
	}

	// Second increment - should see updated counter
	result, err = interp.Call("increment")
	if err != nil {
		t.Fatalf("Call increment error: %v", err)
	}
	if result != 2 {
		t.Errorf("expected 2, got %v (type %T)", result, result)
	}

	// Get counter should return 2
	result, err = interp.Call("getCounter")
	if err != nil {
		t.Fatalf("Call getCounter error: %v", err)
	}
	if result != 2 {
		t.Errorf("expected 2, got %v (type %T)", result, result)
	}
}

func TestCallUndefinedFunction(t *testing.T) {
	interp := basic.NewInterpreter()

	err := interp.Load(`
function exists():
    return 1
endfunction
`)
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}

	_, err = interp.Call("doesNotExist")
	if err == nil {
		t.Error("expected error for undefined function")
	}
}

func TestCallWrongArgumentCount(t *testing.T) {
	interp := basic.NewInterpreter()

	err := interp.Load(`
function add(a, b):
    return a + b
endfunction
`)
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}

	_, err = interp.Call("add", (1)) // Missing second arg
	if err == nil {
		t.Error("expected error for wrong argument count")
	}
}

func TestHasFunction(t *testing.T) {
	interp := basic.NewInterpreter()

	err := interp.Load(`
function init():
    return 1
endfunction

function update():
    return 2
endfunction
`)
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}

	if !interp.HasFunction("init") {
		t.Error("expected HasFunction('init') to be true")
	}
	if !interp.HasFunction("update") {
		t.Error("expected HasFunction('update') to be true")
	}
	if interp.HasFunction("destroy") {
		t.Error("expected HasFunction('destroy') to be false")
	}
}

func TestCallWithExternalFunctions(t *testing.T) {
	interp := basic.NewInterpreter()

	var velocityX, velocityY float64
	interp.RegisterFunction("setVelocity", func(args ...interface{}) (interface{}, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("setVelocity requires 2 arguments")
		}
		velocityX = toFloat(args[0])
		velocityY = toFloat(args[1])
		return nil, nil
	})

	err := interp.Load(`
function hit(forceX, forceY):
    setVelocity(forceX, forceY)
endfunction
`)
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}

	_, err = interp.Call("hit", 5.0, -3.0)
	if err != nil {
		t.Fatalf("Call error: %v", err)
	}

	if velocityX != 5.0 {
		t.Errorf("expected velocityX 5.0, got %v", velocityX)
	}
	if velocityY != -3.0 {
		t.Errorf("expected velocityY -3.0, got %v", velocityY)
	}
}

func toFloat(v interface{}) float64 {
	switch val := v.(type) {
	case float64:
		return val
	case int:
		return float64(val)
	default:
		return 0
	}
}

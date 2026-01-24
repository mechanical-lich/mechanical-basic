package basic

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

// MaxIterations is the default limit for loop iterations to prevent infinite loops
const MaxIterations = 100000

// ExternalFunc is the signature for registered external functions
type ExternalFunc func(args ...interface{}) (interface{}, error)

// PrintFunc is the signature for custom print handlers
type PrintFunc func(value interface{})

// Interpreter executes MechanicalBasic programs
type Interpreter struct {
	// External functions registered by the host application
	externalFuncs map[string]ExternalFunc

	// User-defined functions from the script
	userFuncs map[string]*FunctionStatement

	// Variable scopes (stack for function calls)
	scopes []map[string]interface{}

	// AST cache keyed by code hash
	astCache map[string]*Program

	// Configuration
	maxIterations int       // Max loop iterations (infinite loop protection)
	printFunc     PrintFunc // Custom print handler (defaults to fmt.Println)

	// Execution state
	iterationCount int  // Current iteration count for loop protection
	breakFlag      bool // Set when BREAK is encountered
	returnFlag     bool // Set when RETURN is encountered
	returnValue    interface{}
}

// NewInterpreter creates a new interpreter instance
func NewInterpreter() *Interpreter {
	return &Interpreter{
		externalFuncs: make(map[string]ExternalFunc),
		userFuncs:     make(map[string]*FunctionStatement),
		scopes:        []map[string]interface{}{make(map[string]interface{})},
		astCache:      make(map[string]*Program),
		maxIterations: MaxIterations,
		printFunc:     func(v interface{}) { fmt.Println(v) },
	}
}

// RegisterFunction registers an external function that can be called from scripts
func (i *Interpreter) RegisterFunction(name string, function ExternalFunc) {
	i.externalFuncs[strings.ToLower(name)] = function
}

// SetMaxIterations sets the maximum loop iterations allowed
func (i *Interpreter) SetMaxIterations(max int) {
	i.maxIterations = max
}

// SetPrintFunc sets a custom handler for PRINT statements
func (i *Interpreter) SetPrintFunc(fn PrintFunc) {
	i.printFunc = fn
}

// Interpret executes the given code string
func (i *Interpreter) Interpret(code string) error {
	prog, err := i.getOrParseProgram(code)
	if err != nil {
		return err
	}

	return i.executeProgram(prog)
}

// Load parses the code and registers function definitions without executing top-level code
func (i *Interpreter) Load(code string) error {
	prog, err := i.getOrParseProgram(code)
	if err != nil {
		return err
	}

	// Reset user functions and collect from this program
	i.userFuncs = make(map[string]*FunctionStatement)
	for _, stmt := range prog.Statements {
		if fn, ok := stmt.(*FunctionStatement); ok {
			i.userFuncs[strings.ToLower(fn.Name)] = fn
		}
	}

	return nil
}

// Call invokes a script-defined function by name with the provided arguments
// Each call starts with a fresh scope - variables do not persist between calls
func (i *Interpreter) Call(funcName string, args ...interface{}) (interface{}, error) {
	name := strings.ToLower(funcName)

	fn, ok := i.userFuncs[name]
	if !ok {
		return nil, fmt.Errorf("undefined function: %s", funcName)
	}

	if len(args) != len(fn.Params) {
		return nil, fmt.Errorf("function %s expects %d arguments, got %d", funcName, len(fn.Params), len(args))
	}

	// Reset execution state for this call
	i.iterationCount = 0
	i.breakFlag = false
	i.returnFlag = false
	i.returnValue = nil

	// Start with fresh scope (no variable persistence)
	i.scopes = []map[string]interface{}{make(map[string]interface{})}

	// Bind parameters to the scope
	for idx, param := range fn.Params {
		i.currentScope()[strings.ToLower(param)] = args[idx]
	}

	// Execute function body
	if err := i.executeBlock(fn.Body); err != nil {
		return nil, err
	}

	return i.returnValue, nil
}

// HasFunction checks if a function with the given name exists
func (i *Interpreter) HasFunction(funcName string) bool {
	_, ok := i.userFuncs[strings.ToLower(funcName)]
	return ok
}

// Validate checks the given code for syntax errors without executing it
func (i *Interpreter) Validate(code string) error {
	_, err := i.getOrParseProgram(code)
	return err
}

// getOrParseProgram returns a cached AST or parses and caches the code
func (i *Interpreter) getOrParseProgram(code string) (*Program, error) {
	hash := i.hashCode(code)

	if prog, ok := i.astCache[hash]; ok {
		return prog, nil
	}

	tokens, err := Tokenize(code)
	if err != nil {
		return nil, err
	}

	prog, err := Parse(tokens)
	if err != nil {
		return nil, err
	}

	i.astCache[hash] = prog
	return prog, nil
}

func (i *Interpreter) hashCode(code string) string {
	h := sha256.Sum256([]byte(code))
	return fmt.Sprintf("%x", h[:8])
}

// executeProgram runs the program
func (i *Interpreter) executeProgram(prog *Program) error {
	// Reset execution state
	i.iterationCount = 0
	i.breakFlag = false
	i.returnFlag = false
	i.returnValue = nil
	i.userFuncs = make(map[string]*FunctionStatement)

	// First pass: collect function definitions
	for _, stmt := range prog.Statements {
		if fn, ok := stmt.(*FunctionStatement); ok {
			i.userFuncs[strings.ToLower(fn.Name)] = fn
		}
	}

	// Second pass: execute top-level statements
	for _, stmt := range prog.Statements {
		if _, ok := stmt.(*FunctionStatement); ok {
			continue // Skip function definitions
		}

		if err := i.executeStatement(stmt); err != nil {
			return err
		}

		if i.returnFlag {
			break
		}
	}

	return nil
}

// -----------------------------------------------------------------------------
// Statement Execution
// -----------------------------------------------------------------------------

func (i *Interpreter) executeStatement(stmt Statement) error {
	switch s := stmt.(type) {
	case *LetStatement:
		return i.executeLetStatement(s)
	case *AssignStatement:
		return i.executeAssignStatement(s)
	case *IfStatement:
		return i.executeIfStatement(s)
	case *ForStatement:
		return i.executeForStatement(s)
	case *BreakStatement:
		i.breakFlag = true
		return nil
	case *ReturnStatement:
		return i.executeReturnStatement(s)
	case *PrintStatement:
		return i.executePrintStatement(s)
	case *ExpressionStatement:
		_, err := i.evaluateExpression(s.Expr)
		return err
	case *FunctionStatement:
		// Already collected in first pass
		return nil
	default:
		return fmt.Errorf("unknown statement type: %T", stmt)
	}
}

func (i *Interpreter) executeLetStatement(stmt *LetStatement) error {
	value, err := i.evaluateExpression(stmt.Value)
	if err != nil {
		return err
	}

	// LET always creates/overwrites in current scope
	i.currentScope()[strings.ToLower(stmt.Name)] = value
	return nil
}

func (i *Interpreter) executeAssignStatement(stmt *AssignStatement) error {
	name := strings.ToLower(stmt.Name)

	switch stmt.Operator {
	case TOKEN_PLUS_PLUS:
		val, err := i.getVariable(name)
		if err != nil {
			return err
		}
		newVal, err := i.addValues(val, 1)
		if err != nil {
			return i.runtimeError(stmt, "cannot increment %T", val)
		}
		i.setVariable(name, newVal)

	case TOKEN_MINUS_MINUS:
		val, err := i.getVariable(name)
		if err != nil {
			return err
		}
		newVal, err := i.subtractValues(val, 1)
		if err != nil {
			return i.runtimeError(stmt, "cannot decrement %T", val)
		}
		i.setVariable(name, newVal)

	case TOKEN_PLUS_EQ:
		val, err := i.getVariable(name)
		if err != nil {
			return err
		}
		addend, err := i.evaluateExpression(stmt.Value)
		if err != nil {
			return err
		}
		newVal, err := i.addValues(val, addend)
		if err != nil {
			return i.runtimeError(stmt, "cannot add %T to %T", addend, val)
		}
		i.setVariable(name, newVal)

	case TOKEN_MINUS_EQ:
		val, err := i.getVariable(name)
		if err != nil {
			return err
		}
		subtrahend, err := i.evaluateExpression(stmt.Value)
		if err != nil {
			return err
		}
		newVal, err := i.subtractValues(val, subtrahend)
		if err != nil {
			return i.runtimeError(stmt, "cannot subtract %T from %T", subtrahend, val)
		}
		i.setVariable(name, newVal)

	case TOKEN_EQ:
		value, err := i.evaluateExpression(stmt.Value)
		if err != nil {
			return err
		}
		i.setVariable(name, value)

	default:
		return i.runtimeError(stmt, "unknown assignment operator: %s", stmt.Operator)
	}

	return nil
}

func (i *Interpreter) executeIfStatement(stmt *IfStatement) error {
	cond, err := i.evaluateExpression(stmt.Condition)
	if err != nil {
		return err
	}

	if i.isTruthy(cond) {
		return i.executeBlock(stmt.ThenBlock)
	}

	// Check elseif clauses
	for _, elseIf := range stmt.ElseIfClauses {
		cond, err := i.evaluateExpression(elseIf.Condition)
		if err != nil {
			return err
		}
		if i.isTruthy(cond) {
			return i.executeBlock(elseIf.Block)
		}
	}

	// Execute else block if present
	if len(stmt.ElseBlock) > 0 {
		return i.executeBlock(stmt.ElseBlock)
	}

	return nil
}

func (i *Interpreter) executeForStatement(stmt *ForStatement) error {
	start, err := i.evaluateExpression(stmt.Start)
	if err != nil {
		return err
	}

	end, err := i.evaluateExpression(stmt.End)
	if err != nil {
		return err
	}

	startInt, ok := i.toInt(start)
	if !ok {
		return i.runtimeError(stmt, "FOR start value must be numeric")
	}

	endInt, ok := i.toInt(end)
	if !ok {
		return i.runtimeError(stmt, "FOR end value must be numeric")
	}

	// Create a new scope for the loop variable (doesn't leak)
	i.pushScope()
	defer i.popScope()

	varName := strings.ToLower(stmt.Variable)

	for j := startInt; j <= endInt; j++ {
		// Check infinite loop protection
		i.iterationCount++
		if i.iterationCount > i.maxIterations {
			return i.runtimeError(stmt, "maximum iterations exceeded (%d)", i.maxIterations)
		}

		i.currentScope()[varName] = j

		if err := i.executeBlock(stmt.Body); err != nil {
			return err
		}

		if i.breakFlag {
			i.breakFlag = false
			break
		}

		if i.returnFlag {
			break
		}
	}

	return nil
}

func (i *Interpreter) executeReturnStatement(stmt *ReturnStatement) error {
	if stmt.Value != nil {
		val, err := i.evaluateExpression(stmt.Value)
		if err != nil {
			return err
		}
		i.returnValue = val
	}
	i.returnFlag = true
	return nil
}

func (i *Interpreter) executePrintStatement(stmt *PrintStatement) error {
	val, err := i.evaluateExpression(stmt.Value)
	if err != nil {
		return err
	}
	i.printFunc(val)
	return nil
}

func (i *Interpreter) executeBlock(statements []Statement) error {
	for _, stmt := range statements {
		if err := i.executeStatement(stmt); err != nil {
			return err
		}
		if i.breakFlag || i.returnFlag {
			break
		}
	}
	return nil
}

// -----------------------------------------------------------------------------
// Expression Evaluation
// -----------------------------------------------------------------------------

func (i *Interpreter) evaluateExpression(expr Expression) (interface{}, error) {
	switch e := expr.(type) {
	case *IntLiteral:
		return e.Value, nil
	case *FloatLiteral:
		return e.Value, nil
	case *StringLiteral:
		return e.Value, nil
	case *BoolLiteral:
		return e.Value, nil
	case *Identifier:
		return i.getVariable(strings.ToLower(e.Name))
	case *BinaryExpr:
		return i.evaluateBinaryExpr(e)
	case *UnaryExpr:
		return i.evaluateUnaryExpr(e)
	case *CallExpr:
		return i.evaluateCallExpr(e)
	default:
		return nil, fmt.Errorf("unknown expression type: %T", expr)
	}
}

func (i *Interpreter) evaluateBinaryExpr(expr *BinaryExpr) (interface{}, error) {
	left, err := i.evaluateExpression(expr.Left)
	if err != nil {
		return nil, err
	}

	right, err := i.evaluateExpression(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator {
	// Arithmetic
	case TOKEN_PLUS:
		return i.addValues(left, right)
	case TOKEN_MINUS:
		return i.subtractValues(left, right)
	case TOKEN_STAR:
		return i.multiplyValues(left, right)
	case TOKEN_SLASH:
		return i.divideValues(left, right)

	// Comparison
	case TOKEN_EQ:
		return i.equalValues(left, right), nil
	case TOKEN_NEQ:
		return !i.equalValues(left, right), nil
	case TOKEN_LT:
		return i.compareValues(left, right) < 0, nil
	case TOKEN_GT:
		return i.compareValues(left, right) > 0, nil
	case TOKEN_LTE:
		return i.compareValues(left, right) <= 0, nil
	case TOKEN_GTE:
		return i.compareValues(left, right) >= 0, nil

	// Logical
	case TOKEN_AND:
		return i.isTruthy(left) && i.isTruthy(right), nil
	case TOKEN_OR:
		return i.isTruthy(left) || i.isTruthy(right), nil

	default:
		return nil, i.runtimeError(expr, "unknown binary operator: %s", expr.Operator)
	}
}

func (i *Interpreter) evaluateUnaryExpr(expr *UnaryExpr) (interface{}, error) {
	operand, err := i.evaluateExpression(expr.Operand)
	if err != nil {
		return nil, err
	}

	switch expr.Operator {
	case TOKEN_MINUS:
		switch v := operand.(type) {
		case int:
			return -v, nil
		case float64:
			return -v, nil
		default:
			return nil, i.runtimeError(expr, "cannot negate %T", operand)
		}

	case TOKEN_NOT:
		return !i.isTruthy(operand), nil

	default:
		return nil, i.runtimeError(expr, "unknown unary operator: %s", expr.Operator)
	}
}

func (i *Interpreter) evaluateCallExpr(expr *CallExpr) (interface{}, error) {
	name := strings.ToLower(expr.Name)

	// Evaluate arguments
	args := make([]interface{}, len(expr.Args))
	for idx, argExpr := range expr.Args {
		val, err := i.evaluateExpression(argExpr)
		if err != nil {
			return nil, err
		}
		args[idx] = val
	}

	// Check external functions first
	if fn, ok := i.externalFuncs[name]; ok {
		return fn(args...)
	}

	// Check user-defined functions
	if fn, ok := i.userFuncs[name]; ok {
		return i.callUserFunction(fn, args)
	}

	return nil, i.runtimeError(expr, "undefined function: %s", expr.Name)
}

func (i *Interpreter) callUserFunction(fn *FunctionStatement, args []interface{}) (interface{}, error) {
	if len(args) != len(fn.Params) {
		return nil, fmt.Errorf("function %s expects %d arguments, got %d", fn.Name, len(fn.Params), len(args))
	}

	// Push new scope for function
	i.pushScope()
	defer i.popScope()

	// Bind parameters
	for idx, param := range fn.Params {
		i.currentScope()[strings.ToLower(param)] = args[idx]
	}

	// Save and restore return state
	oldReturnFlag := i.returnFlag
	oldReturnValue := i.returnValue
	i.returnFlag = false
	i.returnValue = nil

	// Execute function body
	if err := i.executeBlock(fn.Body); err != nil {
		return nil, err
	}

	result := i.returnValue

	// Restore return state
	i.returnFlag = oldReturnFlag
	i.returnValue = oldReturnValue

	return result, nil
}

// -----------------------------------------------------------------------------
// Value Operations
// -----------------------------------------------------------------------------

func (i *Interpreter) addValues(left, right interface{}) (interface{}, error) {
	// String concatenation
	if ls, ok := left.(string); ok {
		return ls + i.toString(right), nil
	}
	if _, ok := right.(string); ok {
		return i.toString(left) + i.toString(right), nil
	}

	// Numeric addition
	lf, lok := i.toFloat64(left)
	rf, rok := i.toFloat64(right)
	if !lok || !rok {
		return nil, fmt.Errorf("cannot add %T and %T", left, right)
	}

	// If both are ints, return int
	if li, ok := left.(int); ok {
		if ri, ok := right.(int); ok {
			return li + ri, nil
		}
	}

	return lf + rf, nil
}

func (i *Interpreter) subtractValues(left, right interface{}) (interface{}, error) {
	lf, lok := i.toFloat64(left)
	rf, rok := i.toFloat64(right)
	if !lok || !rok {
		return nil, fmt.Errorf("cannot subtract %T from %T", right, left)
	}

	if li, ok := left.(int); ok {
		if ri, ok := right.(int); ok {
			return li - ri, nil
		}
	}

	return lf - rf, nil
}

func (i *Interpreter) multiplyValues(left, right interface{}) (interface{}, error) {
	lf, lok := i.toFloat64(left)
	rf, rok := i.toFloat64(right)
	if !lok || !rok {
		return nil, fmt.Errorf("cannot multiply %T and %T", left, right)
	}

	if li, ok := left.(int); ok {
		if ri, ok := right.(int); ok {
			return li * ri, nil
		}
	}

	return lf * rf, nil
}

func (i *Interpreter) divideValues(left, right interface{}) (interface{}, error) {
	lf, lok := i.toFloat64(left)
	rf, rok := i.toFloat64(right)
	if !lok || !rok {
		return nil, fmt.Errorf("cannot divide %T by %T", left, right)
	}

	if rf == 0 {
		return nil, fmt.Errorf("division by zero")
	}

	if li, ok := left.(int); ok {
		if ri, ok := right.(int); ok {
			return li / ri, nil
		}
	}

	return lf / rf, nil
}

func (i *Interpreter) equalValues(left, right interface{}) bool {
	// Type-aware comparison
	switch lv := left.(type) {
	case int:
		if rv, ok := right.(int); ok {
			return lv == rv
		}
		if rv, ok := right.(float64); ok {
			return float64(lv) == rv
		}
	case float64:
		if rv, ok := right.(float64); ok {
			return lv == rv
		}
		if rv, ok := right.(int); ok {
			return lv == float64(rv)
		}
	case string:
		if rv, ok := right.(string); ok {
			return lv == rv
		}
	case bool:
		if rv, ok := right.(bool); ok {
			return lv == rv
		}
	}
	return false
}

func (i *Interpreter) compareValues(left, right interface{}) int {
	lf, lok := i.toFloat64(left)
	rf, rok := i.toFloat64(right)

	if lok && rok {
		if lf < rf {
			return -1
		}
		if lf > rf {
			return 1
		}
		return 0
	}

	// String comparison
	ls := i.toString(left)
	rs := i.toString(right)
	if ls < rs {
		return -1
	}
	if ls > rs {
		return 1
	}
	return 0
}

// -----------------------------------------------------------------------------
// Type Helpers
// -----------------------------------------------------------------------------

func (i *Interpreter) isTruthy(val interface{}) bool {
	switch v := val.(type) {
	case nil:
		return false
	case bool:
		return v
	case int:
		return v != 0
	case float64:
		return v != 0
	case string:
		return v != ""
	default:
		return true
	}
}

func (i *Interpreter) toFloat64(val interface{}) (float64, bool) {
	switch v := val.(type) {
	case int:
		return float64(v), true
	case float64:
		return v, true
	default:
		return 0, false
	}
}

func (i *Interpreter) toInt(val interface{}) (int, bool) {
	switch v := val.(type) {
	case int:
		return v, true
	case float64:
		return int(v), true
	default:
		return 0, false
	}
}

func (i *Interpreter) toString(val interface{}) string {
	switch v := val.(type) {
	case string:
		return v
	case int:
		return fmt.Sprintf("%d", v)
	case float64:
		return fmt.Sprintf("%g", v)
	case bool:
		if v {
			return "true"
		}
		return "false"
	case nil:
		return ""
	default:
		return fmt.Sprintf("%v", v)
	}
}

// -----------------------------------------------------------------------------
// Scope Management
// -----------------------------------------------------------------------------

func (i *Interpreter) currentScope() map[string]interface{} {
	return i.scopes[len(i.scopes)-1]
}

func (i *Interpreter) pushScope() {
	i.scopes = append(i.scopes, make(map[string]interface{}))
}

func (i *Interpreter) popScope() {
	if len(i.scopes) > 1 {
		i.scopes = i.scopes[:len(i.scopes)-1]
	}
}

func (i *Interpreter) getVariable(name string) (interface{}, error) {
	// Search from innermost scope outward
	for j := len(i.scopes) - 1; j >= 0; j-- {
		if val, ok := i.scopes[j][name]; ok {
			return val, nil
		}
	}
	return nil, fmt.Errorf("undefined variable: %s", name)
}

func (i *Interpreter) setVariable(name string, value interface{}) {
	// Find existing variable in any scope, or create in current scope
	for j := len(i.scopes) - 1; j >= 0; j-- {
		if _, ok := i.scopes[j][name]; ok {
			i.scopes[j][name] = value
			return
		}
	}
	// Create in current scope if not found
	i.currentScope()[name] = value
}

// -----------------------------------------------------------------------------
// Error Helpers
// -----------------------------------------------------------------------------

func (i *Interpreter) runtimeError(node Node, format string, args ...interface{}) error {
	line, col := node.Position()
	msg := fmt.Sprintf(format, args...)
	return fmt.Errorf("runtime error at line %d, column %d: %s", line, col, msg)
}

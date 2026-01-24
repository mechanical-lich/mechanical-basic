package basic

import (
	"github.com/mechanical-lich/mechanical-basic/internal/basic"
	mathlib "github.com/mechanical-lich/mechanical-basic/internal/math_lib"
)

type MechBasic struct {
	interpreter *basic.Interpreter
}

func NewMechanicalBasic() *MechBasic {
	mb := &MechBasic{
		interpreter: basic.NewInterpreter(),
	}

	// Register built-in math functions
	mb.RegisterMathLibrary()

	return mb
}

func (mb *MechBasic) RegisterFunc(name string, function func(args ...any) (any, error)) {
	mb.interpreter.RegisterFunction(name, function)
}

func (mb *MechBasic) Run(code string) error {
	return mb.interpreter.Interpret(code)
}

// Load parses the script and registers function definitions without executing top-level code
func (mb *MechBasic) Load(code string) error {
	return mb.interpreter.Load(code)
}

// Call invokes a script-defined function by name with the provided arguments
// Each call starts with a fresh scope - variables do not persist between calls
func (mb *MechBasic) Call(funcName string, args ...any) (any, error) {
	return mb.interpreter.Call(funcName, args...)
}

// HasFunction checks if a function with the given name exists in the loaded script
func (mb *MechBasic) HasFunction(funcName string) bool {
	return mb.interpreter.HasFunction(funcName)
}

func (mb *MechBasic) RegisterMathLibrary() {
	mb.interpreter.RegisterFunction("pow", mathlib.Pow)
	mb.interpreter.RegisterFunction("abs", mathlib.Abs)
	mb.interpreter.RegisterFunction("atn", mathlib.Atn)
	mb.interpreter.RegisterFunction("cos", mathlib.Cos)
	mb.interpreter.RegisterFunction("exp", mathlib.Exp)
	mb.interpreter.RegisterFunction("int", mathlib.Int)
	mb.interpreter.RegisterFunction("log", mathlib.Log)
	mb.interpreter.RegisterFunction("rnd", mathlib.Rnd)
	mb.interpreter.RegisterFunction("sin", mathlib.Sin)
	mb.interpreter.RegisterFunction("tan", mathlib.Tan)
	mb.interpreter.RegisterFunction("sqr", mathlib.Sqr)
}

func (mb *MechBasic) SetPrintFunc(fn func(value any)) {
	mb.interpreter.SetPrintFunc(fn)
}

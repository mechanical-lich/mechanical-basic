package basic

import (
	"testing"
)

var benchSmallProgram = `
LET x = 1
LET y = 2
PRINT x + y
`

var benchLoopProgram = `
LET sum = 0
FOR i = 1 TO 100
    sum = sum + i
NEXT
PRINT sum
`

var benchFuncProgram = `
FUNCTION add(a, b)
    RETURN a + b
ENDFUNCTION

LET s = 0
FOR i = 1 TO 100
    s = s + add(i, i)
NEXT
PRINT s
`

var benchExternalProgram = `
FOR i = 1 TO 100
    noop(i)
NEXT
`

func BenchmarkInterpret_ParseAndExec(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		interp := NewInterpreter()
		if err := interp.Interpret(benchSmallProgram); err != nil {
			b.Fatalf("interpret error: %v", err)
		}
	}
}

func BenchmarkInterpret_Cached(b *testing.B) {
	b.ReportAllocs()
	interp := NewInterpreter()
	// warm cache
	if err := interp.Interpret(benchSmallProgram); err != nil {
		b.Fatalf("warm cache error: %v", err)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if err := interp.Interpret(benchSmallProgram); err != nil {
			b.Fatalf("interpret error: %v", err)
		}
	}
}

func BenchmarkInterpret_Loop(b *testing.B) {
	b.ReportAllocs()
	interp := NewInterpreter()
	for n := 0; n < b.N; n++ {
		if err := interp.Interpret(benchLoopProgram); err != nil {
			b.Fatalf("interpret error: %v", err)
		}
	}
}

func BenchmarkFunctionCalls(b *testing.B) {
	b.ReportAllocs()
	interp := NewInterpreter()
	for n := 0; n < b.N; n++ {
		if err := interp.Interpret(benchFuncProgram); err != nil {
			b.Fatalf("interpret error: %v", err)
		}
	}
}

func BenchmarkExternalCalls(b *testing.B) {
	b.ReportAllocs()
	interp := NewInterpreter()
	interp.RegisterFunction("noop", func(args ...interface{}) (interface{}, error) {
		return nil, nil
	})
	for n := 0; n < b.N; n++ {
		if err := interp.Interpret(benchExternalProgram); err != nil {
			b.Fatalf("interpret error: %v", err)
		}
	}
}

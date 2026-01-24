package mathlib

import (
	"fmt"
	"math"
	"math/rand"

	basic "github.com/mechanical-lich/mechanical-basic/pkg/functions"
)

// Pow returns base raised to exponent power
func Pow(args ...interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("pow requires 2 arguments")
	}

	a, err := basic.EnsureFloat(args[0])
	if err != nil {
		return nil, fmt.Errorf("pow: first argument must be numeric: %v", err)
	}

	b, err := basic.EnsureFloat(args[1])
	if err != nil {
		return nil, fmt.Errorf("pow: second argument must be numeric: %v", err)
	}

	return math.Pow(a, b), nil
}

// Abs returns the absolute value of a number
func Abs(args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("abs requires 1 argument")
	}

	val, err := basic.EnsureFloat(args[0])
	if err != nil {
		return nil, fmt.Errorf("abs: argument must be numeric: %v", err)
	}

	return math.Abs(val), nil
}

// Atn returns the arctangent of a number in radians
func Atn(args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("atn requires 1 argument")
	}

	val, err := basic.EnsureFloat(args[0])
	if err != nil {
		return nil, fmt.Errorf("atn: argument must be numeric: %v", err)
	}

	return math.Atan(val), nil
}

// Cos returns the cosine of an angle in radians
func Cos(args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("cos requires 1 argument")
	}

	val, err := basic.EnsureFloat(args[0])
	if err != nil {
		return nil, fmt.Errorf("cos: argument must be numeric: %v", err)
	}

	return math.Cos(val), nil
}

// Exp returns e raised to the power of the argument
func Exp(args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("exp requires 1 argument")
	}

	val, err := basic.EnsureFloat(args[0])
	if err != nil {
		return nil, fmt.Errorf("exp: argument must be numeric: %v", err)
	}

	return math.Exp(val), nil
}

// Int returns the integer part (floor) of a number
func Int(args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("int requires 1 argument")
	}

	val, err := basic.EnsureFloat(args[0])
	if err != nil {
		return nil, fmt.Errorf("int: argument must be numeric: %v", err)
	}

	return int(math.Floor(val)), nil
}

// Log returns the natural logarithm of a number
func Log(args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("log requires 1 argument")
	}

	val, err := basic.EnsureFloat(args[0])
	if err != nil {
		return nil, fmt.Errorf("log: argument must be numeric: %v", err)
	}

	if val <= 0 {
		return nil, fmt.Errorf("log: argument must be positive")
	}

	return math.Log(val), nil
}

// Rnd returns a random number between 0 and 1 if no argument,
// or between 0 and the specified value if an argument is provided
func Rnd(args ...interface{}) (interface{}, error) {
	if len(args) == 0 {
		return rand.Float64(), nil
	}

	if len(args) == 1 {
		max, err := basic.EnsureFloat(args[0])
		if err != nil {
			return nil, fmt.Errorf("rnd: argument must be numeric: %v", err)
		}
		return rand.Float64() * max, nil
	}

	return nil, fmt.Errorf("rnd requires 0 or 1 argument")
}

// Sin returns the sine of an angle in radians
func Sin(args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("sin requires 1 argument")
	}

	val, err := basic.EnsureFloat(args[0])
	if err != nil {
		return nil, fmt.Errorf("sin: argument must be numeric: %v", err)
	}

	return math.Sin(val), nil
}

// Sqr returns the square root of a number
func Sqr(args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("sqr requires 1 argument")
	}

	val, err := basic.EnsureFloat(args[0])
	if err != nil {
		return nil, fmt.Errorf("sqr: argument must be numeric: %v", err)
	}

	if val < 0 {
		return nil, fmt.Errorf("sqr: argument must be non-negative")
	}

	return math.Sqrt(val), nil
}

// Tan returns the tangent of an angle in radians
func Tan(args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("tan requires 1 argument")
	}

	val, err := basic.EnsureFloat(args[0])
	if err != nil {
		return nil, fmt.Errorf("tan: argument must be numeric: %v", err)
	}

	return math.Tan(val), nil
}

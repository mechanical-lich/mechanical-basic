package mathlib

import (
	"math"
	"testing"
)

func TestPow(t *testing.T) {
	result, err := Pow(2.0, 3.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != 8.0 {
		t.Errorf("expected 8, got %v", result)
	}
}

func TestAbs(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected float64
	}{
		{-5.0, 5.0},
		{5.0, 5.0},
		{-3, 3.0},
		{0, 0.0},
	}

	for _, tt := range tests {
		result, err := Abs(tt.input)
		if err != nil {
			t.Errorf("Abs(%v): unexpected error: %v", tt.input, err)
			continue
		}
		if result != tt.expected {
			t.Errorf("Abs(%v): expected %v, got %v", tt.input, tt.expected, result)
		}
	}
}

func TestAtn(t *testing.T) {
	result, err := Atn(1.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := math.Atan(1.0)
	if result != expected {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestCos(t *testing.T) {
	result, err := Cos(0.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != 1.0 {
		t.Errorf("expected 1, got %v", result)
	}

	result, err = Cos(math.Pi)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if math.Abs(result.(float64)-(-1.0)) > 0.0001 {
		t.Errorf("expected -1, got %v", result)
	}
}

func TestExp(t *testing.T) {
	result, err := Exp(1.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := math.E
	if math.Abs(result.(float64)-expected) > 0.0001 {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestInt(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected int
	}{
		{3.14, 3},
		{3.99, 3},
		{-2.5, -3},
		{5.0, 5},
	}

	for _, tt := range tests {
		result, err := Int(tt.input)
		if err != nil {
			t.Errorf("Int(%v): unexpected error: %v", tt.input, err)
			continue
		}
		if result != tt.expected {
			t.Errorf("Int(%v): expected %v, got %v", tt.input, tt.expected, result)
		}
	}
}

func TestLog(t *testing.T) {
	result, err := Log(math.E)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if math.Abs(result.(float64)-1.0) > 0.0001 {
		t.Errorf("expected 1, got %v", result)
	}

	_, err = Log(-1.0)
	if err == nil {
		t.Error("expected error for negative input")
	}

	_, err = Log(0.0)
	if err == nil {
		t.Error("expected error for zero input")
	}
}

func TestRnd(t *testing.T) {
	result, err := Rnd()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	val := result.(float64)
	if val < 0 || val >= 1 {
		t.Errorf("expected value in [0, 1), got %v", val)
	}

	result, err = Rnd(10.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	val = result.(float64)
	if val < 0 || val >= 10 {
		t.Errorf("expected value in [0, 10), got %v", val)
	}
}

func TestSin(t *testing.T) {
	result, err := Sin(0.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != 0.0 {
		t.Errorf("expected 0, got %v", result)
	}

	result, err = Sin(math.Pi / 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if math.Abs(result.(float64)-1.0) > 0.0001 {
		t.Errorf("expected 1, got %v", result)
	}
}

func TestSqr(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected float64
	}{
		{4.0, 2.0},
		{9.0, 3.0},
		{0.0, 0.0},
		{2.0, math.Sqrt(2)},
	}

	for _, tt := range tests {
		result, err := Sqr(tt.input)
		if err != nil {
			t.Errorf("Sqr(%v): unexpected error: %v", tt.input, err)
			continue
		}
		if math.Abs(result.(float64)-tt.expected) > 0.0001 {
			t.Errorf("Sqr(%v): expected %v, got %v", tt.input, tt.expected, result)
		}
	}

	_, err := Sqr(-1.0)
	if err == nil {
		t.Error("expected error for negative input")
	}
}

func TestTan(t *testing.T) {
	result, err := Tan(0.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != 0.0 {
		t.Errorf("expected 0, got %v", result)
	}

	result, err = Tan(math.Pi / 4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if math.Abs(result.(float64)-1.0) > 0.0001 {
		t.Errorf("expected 1, got %v", result)
	}
}

func TestErrorCases(t *testing.T) {
	_, err := Pow(1.0)
	if err == nil {
		t.Error("Pow: expected error for wrong number of arguments")
	}

	_, err = Abs()
	if err == nil {
		t.Error("Abs: expected error for missing argument")
	}

	_, err = Rnd(1.0, 2.0)
	if err == nil {
		t.Error("Rnd: expected error for too many arguments")
	}

	_, err = Sin("hello")
	if err == nil {
		t.Error("Sin: expected error for non-numeric argument")
	}
}

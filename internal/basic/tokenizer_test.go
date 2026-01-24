package basic

import (
	"testing"
)

func TestTokenizeBasicAssignment(t *testing.T) {
	input := "LET x = 5"
	tokens, err := Tokenize(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []struct {
		typ   TokenType
		value string
	}{
		{TOKEN_LET, "LET"},
		{TOKEN_IDENTIFIER, "x"},
		{TOKEN_EQ, "="},
		{TOKEN_INT, "5"},
		{TOKEN_EOF, ""},
	}

	if len(tokens) != len(expected) {
		t.Fatalf("expected %d tokens, got %d", len(expected), len(tokens))
	}

	for i, exp := range expected {
		if tokens[i].Type != exp.typ {
			t.Errorf("token %d: expected type %s, got %s", i, exp.typ, tokens[i].Type)
		}
		if tokens[i].Value != exp.value {
			t.Errorf("token %d: expected value %q, got %q", i, exp.value, tokens[i].Value)
		}
	}
}

func TestTokenizeNumbers(t *testing.T) {
	input := "42 3.14 0.5"
	tokens, err := Tokenize(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []struct {
		typ   TokenType
		value string
	}{
		{TOKEN_INT, "42"},
		{TOKEN_FLOAT, "3.14"},
		{TOKEN_FLOAT, "0.5"},
		{TOKEN_EOF, ""},
	}

	if len(tokens) != len(expected) {
		t.Fatalf("expected %d tokens, got %d", len(expected), len(tokens))
	}

	for i, exp := range expected {
		if tokens[i].Type != exp.typ {
			t.Errorf("token %d: expected type %s, got %s", i, exp.typ, tokens[i].Type)
		}
		if tokens[i].Value != exp.value {
			t.Errorf("token %d: expected value %q, got %q", i, exp.value, tokens[i].Value)
		}
	}
}

func TestTokenizeString(t *testing.T) {
	input := `"Hello World"`
	tokens, err := Tokenize(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(tokens) != 2 { // STRING + EOF
		t.Fatalf("expected 2 tokens, got %d", len(tokens))
	}

	if tokens[0].Type != TOKEN_STRING {
		t.Errorf("expected STRING, got %s", tokens[0].Type)
	}
	if tokens[0].Value != "Hello World" {
		t.Errorf("expected 'Hello World', got %q", tokens[0].Value)
	}
}

func TestTokenizeStringWithEscapes(t *testing.T) {
	input := `"He said \"Hello\""`
	tokens, err := Tokenize(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(tokens) != 2 {
		t.Fatalf("expected 2 tokens, got %d", len(tokens))
	}

	if tokens[0].Value != `He said "Hello"` {
		t.Errorf("expected 'He said \"Hello\"', got %q", tokens[0].Value)
	}
}

func TestTokenizeOperators(t *testing.T) {
	input := "+ - * / = < > <= >= <> != += -= ++ --"
	tokens, err := Tokenize(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []TokenType{
		TOKEN_PLUS, TOKEN_MINUS, TOKEN_STAR, TOKEN_SLASH, TOKEN_EQ,
		TOKEN_LT, TOKEN_GT, TOKEN_LTE, TOKEN_GTE, TOKEN_NEQ, TOKEN_NEQ,
		TOKEN_PLUS_EQ, TOKEN_MINUS_EQ, TOKEN_PLUS_PLUS, TOKEN_MINUS_MINUS,
		TOKEN_EOF,
	}

	if len(tokens) != len(expected) {
		t.Fatalf("expected %d tokens, got %d", len(expected), len(tokens))
	}

	for i, exp := range expected {
		if tokens[i].Type != exp {
			t.Errorf("token %d: expected type %s, got %s", i, exp, tokens[i].Type)
		}
	}
}

func TestTokenizeKeywords(t *testing.T) {
	input := "if then else elseif endif for to next break function endfunction return print and or not let true false"
	tokens, err := Tokenize(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []TokenType{
		TOKEN_IF, TOKEN_THEN, TOKEN_ELSE, TOKEN_ELSEIF, TOKEN_ENDIF,
		TOKEN_FOR, TOKEN_TO, TOKEN_NEXT, TOKEN_BREAK,
		TOKEN_FUNCTION, TOKEN_ENDFUNCTION, TOKEN_RETURN, TOKEN_PRINT,
		TOKEN_AND, TOKEN_OR, TOKEN_NOT, TOKEN_LET, TOKEN_TRUE, TOKEN_FALSE,
		TOKEN_EOF,
	}

	if len(tokens) != len(expected) {
		t.Fatalf("expected %d tokens, got %d", len(expected), len(tokens))
	}

	for i, exp := range expected {
		if tokens[i].Type != exp {
			t.Errorf("token %d: expected type %s, got %s", i, exp, tokens[i].Type)
		}
	}
}

func TestTokenizeKeywordsCaseInsensitive(t *testing.T) {
	input := "IF Then ELSE"
	tokens, err := Tokenize(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []TokenType{TOKEN_IF, TOKEN_THEN, TOKEN_ELSE, TOKEN_EOF}

	for i, exp := range expected {
		if tokens[i].Type != exp {
			t.Errorf("token %d: expected type %s, got %s", i, exp, tokens[i].Type)
		}
	}
}

func TestTokenizeIfStatement(t *testing.T) {
	input := `if x > 5 then
    print "big"
endif`
	tokens, err := Tokenize(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []struct {
		typ   TokenType
		value string
	}{
		{TOKEN_IF, "if"},
		{TOKEN_IDENTIFIER, "x"},
		{TOKEN_GT, ">"},
		{TOKEN_INT, "5"},
		{TOKEN_THEN, "then"},
		{TOKEN_NEWLINE, "\\n"},
		{TOKEN_PRINT, "print"},
		{TOKEN_STRING, "big"},
		{TOKEN_NEWLINE, "\\n"},
		{TOKEN_ENDIF, "endif"},
		{TOKEN_EOF, ""},
	}

	if len(tokens) != len(expected) {
		t.Fatalf("expected %d tokens, got %d", len(expected), len(tokens))
	}

	for i, exp := range expected {
		if tokens[i].Type != exp.typ {
			t.Errorf("token %d: expected type %s, got %s", i, exp.typ, tokens[i].Type)
		}
	}
}

func TestTokenizeForLoop(t *testing.T) {
	input := `for i = 1 to 10
    print i
next i`
	tokens, err := Tokenize(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []TokenType{
		TOKEN_FOR, TOKEN_IDENTIFIER, TOKEN_EQ, TOKEN_INT, TOKEN_TO, TOKEN_INT,
		TOKEN_NEWLINE,
		TOKEN_PRINT, TOKEN_IDENTIFIER,
		TOKEN_NEWLINE,
		TOKEN_NEXT, TOKEN_IDENTIFIER,
		TOKEN_EOF,
	}

	if len(tokens) != len(expected) {
		t.Fatalf("expected %d tokens, got %d", len(expected), len(tokens))
	}

	for i, exp := range expected {
		if tokens[i].Type != exp {
			t.Errorf("token %d: expected type %s, got %s", i, exp, tokens[i].Type)
		}
	}
}

func TestTokenizeFunctionCall(t *testing.T) {
	input := "pow(2, 3)"
	tokens, err := Tokenize(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []TokenType{
		TOKEN_IDENTIFIER, TOKEN_LPAREN, TOKEN_INT, TOKEN_COMMA, TOKEN_INT, TOKEN_RPAREN,
		TOKEN_EOF,
	}

	if len(tokens) != len(expected) {
		t.Fatalf("expected %d tokens, got %d", len(expected), len(tokens))
	}

	for i, exp := range expected {
		if tokens[i].Type != exp {
			t.Errorf("token %d: expected type %s, got %s", i, exp, tokens[i].Type)
		}
	}
}

func TestTokenizeFunction(t *testing.T) {
	input := `function add(x, y):
    let z = x + y
    return z
endfunction`
	tokens, err := Tokenize(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []TokenType{
		TOKEN_FUNCTION, TOKEN_IDENTIFIER, TOKEN_LPAREN, TOKEN_IDENTIFIER, TOKEN_COMMA, TOKEN_IDENTIFIER, TOKEN_RPAREN, TOKEN_COLON,
		TOKEN_NEWLINE,
		TOKEN_LET, TOKEN_IDENTIFIER, TOKEN_EQ, TOKEN_IDENTIFIER, TOKEN_PLUS, TOKEN_IDENTIFIER,
		TOKEN_NEWLINE,
		TOKEN_RETURN, TOKEN_IDENTIFIER,
		TOKEN_NEWLINE,
		TOKEN_ENDFUNCTION,
		TOKEN_EOF,
	}

	if len(tokens) != len(expected) {
		t.Fatalf("expected %d tokens, got %d", len(expected), len(tokens))
	}

	for i, exp := range expected {
		if tokens[i].Type != exp {
			t.Errorf("token %d: expected type %s, got %s", i, exp, tokens[i].Type)
		}
	}
}

func TestTokenizeCommentSkipped(t *testing.T) {
	input := `x = 5 # this is a comment
y = 10`
	tokens, err := Tokenize(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Comments should be skipped
	expected := []TokenType{
		TOKEN_IDENTIFIER, TOKEN_EQ, TOKEN_INT,
		TOKEN_NEWLINE,
		TOKEN_IDENTIFIER, TOKEN_EQ, TOKEN_INT,
		TOKEN_EOF,
	}

	if len(tokens) != len(expected) {
		t.Fatalf("expected %d tokens, got %d", len(expected), len(tokens))
	}

	for i, exp := range expected {
		if tokens[i].Type != exp {
			t.Errorf("token %d: expected type %s, got %s", i, exp, tokens[i].Type)
		}
	}
}

func TestTokenizeLineNumbers(t *testing.T) {
	input := `x = 1
y = 2
z = 3`
	tokens, err := Tokenize(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check that z is on line 3
	for _, tok := range tokens {
		if tok.Value == "z" {
			if tok.Line != 3 {
				t.Errorf("expected 'z' on line 3, got line %d", tok.Line)
			}
		}
	}
}

func TestTokenizeUnterminatedString(t *testing.T) {
	input := `"Hello`
	_, err := Tokenize(input)
	if err == nil {
		t.Error("expected error for unterminated string")
	}
}

func TestTokenizeUnexpectedCharacter(t *testing.T) {
	input := "x = @"
	_, err := Tokenize(input)
	if err == nil {
		t.Error("expected error for unexpected character")
	}
}

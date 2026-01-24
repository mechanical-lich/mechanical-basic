package basic

import (
	"fmt"
	"strings"
	"unicode"
)

// Tokenizer holds the state for lexical analysis
type Tokenizer struct {
	input    string
	pos      int // Current position in input
	line     int // Current line number (1-indexed)
	column   int // Current column number (1-indexed)
	start    int // Start position of current token
	startCol int // Start column of current token
}

// NewTokenizer creates a new tokenizer for the given input
func NewTokenizer(input string) *Tokenizer {
	return &Tokenizer{
		input:  input,
		pos:    0,
		line:   1,
		column: 1,
	}
}

// Tokenize converts an input string into a slice of tokens
func Tokenize(input string) ([]Token, error) {
	t := NewTokenizer(input)
	return t.ScanAll()
}

// ScanAll scans all tokens from the input
func (t *Tokenizer) ScanAll() ([]Token, error) {
	var tokens []Token

	for {
		tok, err := t.NextToken()
		if err != nil {
			return nil, err
		}

		// Skip comments - they're not needed for parsing
		if tok.Type == TOKEN_COMMENT {
			continue
		}

		tokens = append(tokens, tok)

		if tok.Type == TOKEN_EOF {
			break
		}
	}

	return tokens, nil
}

// NextToken scans and returns the next token
func (t *Tokenizer) NextToken() (Token, error) {
	t.skipWhitespace()

	t.start = t.pos
	t.startCol = t.column

	if t.isAtEnd() {
		return t.makeToken(TOKEN_EOF, ""), nil
	}

	ch := t.advance()

	// Comments
	if ch == '#' {
		return t.scanComment(), nil
	}

	// Newlines
	if ch == '\n' {
		tok := t.makeToken(TOKEN_NEWLINE, "\\n")
		t.line++
		t.column = 1
		return tok, nil
	}

	// Strings
	if ch == '"' {
		return t.scanString()
	}

	// Numbers
	if unicode.IsDigit(ch) {
		return t.scanNumber(), nil
	}

	// Identifiers and keywords
	if unicode.IsLetter(ch) || ch == '_' {
		return t.scanIdentifier(), nil
	}

	// Operators and delimiters
	switch ch {
	case '(':
		return t.makeToken(TOKEN_LPAREN, "("), nil
	case ')':
		return t.makeToken(TOKEN_RPAREN, ")"), nil
	case ',':
		return t.makeToken(TOKEN_COMMA, ","), nil
	case ':':
		return t.makeToken(TOKEN_COLON, ":"), nil
	case '*':
		return t.makeToken(TOKEN_STAR, "*"), nil
	case '/':
		return t.makeToken(TOKEN_SLASH, "/"), nil
	case '+':
		if t.match('+') {
			return t.makeToken(TOKEN_PLUS_PLUS, "++"), nil
		}
		if t.match('=') {
			return t.makeToken(TOKEN_PLUS_EQ, "+="), nil
		}
		return t.makeToken(TOKEN_PLUS, "+"), nil
	case '-':
		if t.match('-') {
			return t.makeToken(TOKEN_MINUS_MINUS, "--"), nil
		}
		if t.match('=') {
			return t.makeToken(TOKEN_MINUS_EQ, "-="), nil
		}
		return t.makeToken(TOKEN_MINUS, "-"), nil
	case '=':
		return t.makeToken(TOKEN_EQ, "="), nil
	case '<':
		if t.match('=') {
			return t.makeToken(TOKEN_LTE, "<="), nil
		}
		if t.match('>') {
			return t.makeToken(TOKEN_NEQ, "<>"), nil
		}
		return t.makeToken(TOKEN_LT, "<"), nil
	case '>':
		if t.match('=') {
			return t.makeToken(TOKEN_GTE, ">="), nil
		}
		return t.makeToken(TOKEN_GT, ">"), nil
	case '!':
		if t.match('=') {
			return t.makeToken(TOKEN_NEQ, "!="), nil
		}
		return Token{}, t.error("unexpected character '!'")
	}

	return Token{}, t.error(fmt.Sprintf("unexpected character '%c'", ch))
}

// scanComment consumes a comment until end of line
func (t *Tokenizer) scanComment() Token {
	for !t.isAtEnd() && t.peek() != '\n' {
		t.advance()
	}
	value := t.input[t.start:t.pos]
	return t.makeToken(TOKEN_COMMENT, value)
}

// scanString scans a string literal with escape sequence handling
func (t *Tokenizer) scanString() (Token, error) {
	var builder strings.Builder
	startLine := t.line

	for !t.isAtEnd() {
		ch := t.peek()

		if ch == '\n' {
			return Token{}, t.error("unterminated string")
		}

		if ch == '"' {
			t.advance() // consume closing quote
			return Token{
				Type:   TOKEN_STRING,
				Value:  builder.String(),
				Line:   startLine,
				Column: t.startCol,
			}, nil
		}

		if ch == '\\' {
			t.advance() // consume backslash
			if t.isAtEnd() {
				return Token{}, t.error("unterminated string escape")
			}

			escaped := t.advance()
			switch escaped {
			case '"':
				builder.WriteRune('"')
			case '\\':
				builder.WriteRune('\\')
			case 'n':
				builder.WriteRune('\n')
			case 't':
				builder.WriteRune('\t')
			case 'r':
				builder.WriteRune('\r')
			default:
				// For unknown escapes, just include the character literally
				builder.WriteRune(escaped)
			}
		} else {
			builder.WriteRune(t.advance())
		}
	}

	return Token{}, t.error("unterminated string")
}

// scanNumber scans an integer or float literal
func (t *Tokenizer) scanNumber() Token {
	for !t.isAtEnd() && unicode.IsDigit(t.peek()) {
		t.advance()
	}

	// Check for decimal point
	isFloat := false
	if !t.isAtEnd() && t.peek() == '.' {
		// Look ahead to make sure there's a digit after the dot
		if t.pos+1 < len(t.input) && unicode.IsDigit(rune(t.input[t.pos+1])) {
			isFloat = true
			t.advance() // consume '.'
			for !t.isAtEnd() && unicode.IsDigit(t.peek()) {
				t.advance()
			}
		}
	}

	value := t.input[t.start:t.pos]
	if isFloat {
		return t.makeToken(TOKEN_FLOAT, value)
	}
	return t.makeToken(TOKEN_INT, value)
}

// scanIdentifier scans an identifier or keyword
func (t *Tokenizer) scanIdentifier() Token {
	for !t.isAtEnd() && (unicode.IsLetter(t.peek()) || unicode.IsDigit(t.peek()) || t.peek() == '_') {
		t.advance()
	}

	value := t.input[t.start:t.pos]
	lower := strings.ToLower(value)

	// Check if it's a keyword
	tokenType := LookupKeyword(lower)

	return Token{
		Type:   tokenType,
		Value:  value,
		Line:   t.line,
		Column: t.startCol,
	}
}

// Helper methods

func (t *Tokenizer) isAtEnd() bool {
	return t.pos >= len(t.input)
}

func (t *Tokenizer) peek() rune {
	if t.isAtEnd() {
		return 0
	}
	return rune(t.input[t.pos])
}

func (t *Tokenizer) advance() rune {
	ch := rune(t.input[t.pos])
	t.pos++
	t.column++
	return ch
}

func (t *Tokenizer) match(expected rune) bool {
	if t.isAtEnd() {
		return false
	}
	if rune(t.input[t.pos]) != expected {
		return false
	}
	t.pos++
	t.column++
	return true
}

func (t *Tokenizer) skipWhitespace() {
	for !t.isAtEnd() {
		ch := t.peek()
		switch ch {
		case ' ', '\t', '\r':
			t.advance()
		default:
			return
		}
	}
}

func (t *Tokenizer) makeToken(tokenType TokenType, value string) Token {
	return Token{
		Type:   tokenType,
		Value:  value,
		Line:   t.line,
		Column: t.startCol,
	}
}

func (t *Tokenizer) error(message string) error {
	return fmt.Errorf("line %d, column %d: %s", t.line, t.startCol, message)
}

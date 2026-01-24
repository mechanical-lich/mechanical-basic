package basic

// TokenType represents the type of a token
type TokenType int

const (
	// Special tokens
	TOKEN_EOF TokenType = iota
	TOKEN_NEWLINE
	TOKEN_COMMENT

	// Literals
	TOKEN_IDENTIFIER
	TOKEN_INT
	TOKEN_FLOAT
	TOKEN_STRING
	TOKEN_TRUE
	TOKEN_FALSE

	// Keywords
	TOKEN_LET
	TOKEN_IF
	TOKEN_THEN
	TOKEN_ELSE
	TOKEN_ELSEIF
	TOKEN_ENDIF
	TOKEN_FOR
	TOKEN_TO
	TOKEN_NEXT
	TOKEN_BREAK
	TOKEN_FUNCTION
	TOKEN_ENDFUNCTION
	TOKEN_RETURN
	TOKEN_PRINT
	TOKEN_AND
	TOKEN_OR
	TOKEN_NOT

	// Operators
	TOKEN_PLUS        // +
	TOKEN_MINUS       // -
	TOKEN_STAR        // *
	TOKEN_SLASH       // /
	TOKEN_EQ          // =
	TOKEN_NEQ         // <> or !=
	TOKEN_LT          // <
	TOKEN_GT          // >
	TOKEN_LTE         // <=
	TOKEN_GTE         // >=
	TOKEN_PLUS_EQ     // +=
	TOKEN_MINUS_EQ    // -=
	TOKEN_PLUS_PLUS   // ++
	TOKEN_MINUS_MINUS // --

	// Delimiters
	TOKEN_LPAREN // (
	TOKEN_RPAREN // )
	TOKEN_COMMA  // ,
	TOKEN_COLON  // :
)

// Token represents a lexical token with its type, value, and position
type Token struct {
	Type   TokenType
	Value  string // The literal value of the token
	Line   int    // Line number (1-indexed)
	Column int    // Column number (1-indexed)
}

// String returns a human-readable name for the token type
func (t TokenType) String() string {
	names := map[TokenType]string{
		TOKEN_EOF:         "EOF",
		TOKEN_NEWLINE:     "NEWLINE",
		TOKEN_COMMENT:     "COMMENT",
		TOKEN_IDENTIFIER:  "IDENTIFIER",
		TOKEN_INT:         "INT",
		TOKEN_FLOAT:       "FLOAT",
		TOKEN_STRING:      "STRING",
		TOKEN_TRUE:        "TRUE",
		TOKEN_FALSE:       "FALSE",
		TOKEN_LET:         "LET",
		TOKEN_IF:          "IF",
		TOKEN_THEN:        "THEN",
		TOKEN_ELSE:        "ELSE",
		TOKEN_ELSEIF:      "ELSEIF",
		TOKEN_ENDIF:       "ENDIF",
		TOKEN_FOR:         "FOR",
		TOKEN_TO:          "TO",
		TOKEN_NEXT:        "NEXT",
		TOKEN_BREAK:       "BREAK",
		TOKEN_FUNCTION:    "FUNCTION",
		TOKEN_ENDFUNCTION: "ENDFUNCTION",
		TOKEN_RETURN:      "RETURN",
		TOKEN_PRINT:       "PRINT",
		TOKEN_AND:         "AND",
		TOKEN_OR:          "OR",
		TOKEN_NOT:         "NOT",
		TOKEN_PLUS:        "PLUS",
		TOKEN_MINUS:       "MINUS",
		TOKEN_STAR:        "STAR",
		TOKEN_SLASH:       "SLASH",
		TOKEN_EQ:          "EQ",
		TOKEN_NEQ:         "NEQ",
		TOKEN_LT:          "LT",
		TOKEN_GT:          "GT",
		TOKEN_LTE:         "LTE",
		TOKEN_GTE:         "GTE",
		TOKEN_PLUS_EQ:     "PLUS_EQ",
		TOKEN_MINUS_EQ:    "MINUS_EQ",
		TOKEN_PLUS_PLUS:   "PLUS_PLUS",
		TOKEN_MINUS_MINUS: "MINUS_MINUS",
		TOKEN_LPAREN:      "LPAREN",
		TOKEN_RPAREN:      "RPAREN",
		TOKEN_COMMA:       "COMMA",
		TOKEN_COLON:       "COLON",
	}
	if name, ok := names[t]; ok {
		return name
	}
	return "UNKNOWN"
}

// keywords maps keyword strings to their token types (case-insensitive)
var keywords = map[string]TokenType{
	"let":         TOKEN_LET,
	"if":          TOKEN_IF,
	"then":        TOKEN_THEN,
	"else":        TOKEN_ELSE,
	"elseif":      TOKEN_ELSEIF,
	"endif":       TOKEN_ENDIF,
	"for":         TOKEN_FOR,
	"to":          TOKEN_TO,
	"next":        TOKEN_NEXT,
	"break":       TOKEN_BREAK,
	"function":    TOKEN_FUNCTION,
	"endfunction": TOKEN_ENDFUNCTION,
	"return":      TOKEN_RETURN,
	"print":       TOKEN_PRINT,
	"and":         TOKEN_AND,
	"or":          TOKEN_OR,
	"not":         TOKEN_NOT,
	"true":        TOKEN_TRUE,
	"false":       TOKEN_FALSE,
}

// LookupKeyword checks if an identifier is a keyword and returns the appropriate token type
func LookupKeyword(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return TOKEN_IDENTIFIER
}

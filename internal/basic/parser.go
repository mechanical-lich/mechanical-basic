package basic

import (
	"fmt"
	"strconv"
)

// Parser converts tokens into an AST
type Parser struct {
	tokens  []Token
	pos     int
	current Token
}

// NewParser creates a new parser for the given tokens
func NewParser(tokens []Token) *Parser {
	p := &Parser{
		tokens: tokens,
		pos:    0,
	}
	if len(tokens) > 0 {
		p.current = tokens[0]
	}
	return p
}

// Parse parses the tokens into a Program AST
func Parse(tokens []Token) (*Program, error) {
	p := NewParser(tokens)
	return p.ParseProgram()
}

// ParseProgram parses the entire program
func (p *Parser) ParseProgram() (*Program, error) {
	program := &Program{
		Statements: []Statement{},
	}

	for !p.isAtEnd() {
		// Skip newlines between statements
		p.skipNewlines()
		if p.isAtEnd() {
			break
		}

		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
	}

	return program, nil
}

// parseStatement parses a single statement
func (p *Parser) parseStatement() (Statement, error) {
	switch p.current.Type {
	case TOKEN_LET:
		return p.parseLetStatement()
	case TOKEN_IF:
		return p.parseIfStatement()
	case TOKEN_FOR:
		return p.parseForStatement()
	case TOKEN_BREAK:
		return p.parseBreakStatement()
	case TOKEN_FUNCTION:
		return p.parseFunctionStatement()
	case TOKEN_RETURN:
		return p.parseReturnStatement()
	case TOKEN_PRINT:
		return p.parsePrintStatement()
	case TOKEN_IDENTIFIER:
		return p.parseIdentifierStatement()
	default:
		return nil, p.error("unexpected token %s", p.current.Type)
	}
}

// parseLetStatement parses: LET name = expr
func (p *Parser) parseLetStatement() (*LetStatement, error) {
	stmt := &LetStatement{
		Pos: Pos{Line: p.current.Line, Column: p.current.Column},
	}

	p.advance() // consume LET

	if p.current.Type != TOKEN_IDENTIFIER {
		return nil, p.error("expected identifier after LET")
	}
	stmt.Name = p.current.Value
	p.advance()

	if p.current.Type != TOKEN_EQ {
		return nil, p.error("expected '=' after variable name")
	}
	p.advance()

	expr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	stmt.Value = expr

	p.consumeNewlineOrEOF()
	return stmt, nil
}

// parseIdentifierStatement parses assignment or expression statement starting with identifier
func (p *Parser) parseIdentifierStatement() (Statement, error) {
	pos := Pos{Line: p.current.Line, Column: p.current.Column}
	name := p.current.Value
	p.advance()

	// Check for assignment operators
	switch p.current.Type {
	case TOKEN_EQ, TOKEN_PLUS_EQ, TOKEN_MINUS_EQ:
		op := p.current.Type
		p.advance()
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		p.consumeNewlineOrEOF()
		return &AssignStatement{Pos: pos, Name: name, Operator: op, Value: expr}, nil

	case TOKEN_PLUS_PLUS:
		p.advance()
		p.consumeNewlineOrEOF()
		return &AssignStatement{Pos: pos, Name: name, Operator: TOKEN_PLUS_PLUS, Value: nil}, nil

	case TOKEN_MINUS_MINUS:
		p.advance()
		p.consumeNewlineOrEOF()
		return &AssignStatement{Pos: pos, Name: name, Operator: TOKEN_MINUS_MINUS, Value: nil}, nil

	case TOKEN_LPAREN:
		// Function call as statement
		p.advance() // consume (
		args, err := p.parseArguments()
		if err != nil {
			return nil, err
		}
		p.consumeNewlineOrEOF()
		return &ExpressionStatement{
			Pos:  pos,
			Expr: &CallExpr{Pos: pos, Name: name, Args: args},
		}, nil

	default:
		return nil, p.error("expected assignment operator or function call after identifier")
	}
}

// parseIfStatement parses: IF cond THEN ... [ELSEIF cond THEN ...] [ELSE ...] ENDIF
func (p *Parser) parseIfStatement() (*IfStatement, error) {
	stmt := &IfStatement{
		Pos: Pos{Line: p.current.Line, Column: p.current.Column},
	}

	p.advance() // consume IF

	// Parse condition
	cond, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	stmt.Condition = cond

	// Expect THEN
	if p.current.Type != TOKEN_THEN {
		return nil, p.error("expected THEN after IF condition")
	}
	p.advance()
	p.consumeNewline()

	// Parse THEN block
	stmt.ThenBlock, err = p.parseBlock(TOKEN_ELSEIF, TOKEN_ELSE, TOKEN_ENDIF)
	if err != nil {
		return nil, err
	}

	// Parse ELSEIF clauses
	for p.current.Type == TOKEN_ELSEIF {
		elseIfPos := Pos{Line: p.current.Line, Column: p.current.Column}
		p.advance() // consume ELSEIF

		elseIfCond, err := p.parseExpression()
		if err != nil {
			return nil, err
		}

		if p.current.Type != TOKEN_THEN {
			return nil, p.error("expected THEN after ELSEIF condition")
		}
		p.advance()
		p.consumeNewline()

		elseIfBlock, err := p.parseBlock(TOKEN_ELSEIF, TOKEN_ELSE, TOKEN_ENDIF)
		if err != nil {
			return nil, err
		}

		stmt.ElseIfClauses = append(stmt.ElseIfClauses, ElseIfClause{
			Pos:       elseIfPos,
			Condition: elseIfCond,
			Block:     elseIfBlock,
		})
	}

	// Parse ELSE block
	if p.current.Type == TOKEN_ELSE {
		p.advance() // consume ELSE
		p.consumeNewline()

		stmt.ElseBlock, err = p.parseBlock(TOKEN_ENDIF)
		if err != nil {
			return nil, err
		}
	}

	// Expect ENDIF
	if p.current.Type != TOKEN_ENDIF {
		return nil, p.error("expected ENDIF")
	}
	p.advance()
	p.consumeNewlineOrEOF()

	return stmt, nil
}

// parseForStatement parses: FOR var = start TO end ... NEXT var
func (p *Parser) parseForStatement() (*ForStatement, error) {
	stmt := &ForStatement{
		Pos: Pos{Line: p.current.Line, Column: p.current.Column},
	}

	p.advance() // consume FOR

	if p.current.Type != TOKEN_IDENTIFIER {
		return nil, p.error("expected identifier after FOR")
	}
	stmt.Variable = p.current.Value
	p.advance()

	if p.current.Type != TOKEN_EQ {
		return nil, p.error("expected '=' after loop variable")
	}
	p.advance()

	start, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	stmt.Start = start

	if p.current.Type != TOKEN_TO {
		return nil, p.error("expected TO in FOR loop")
	}
	p.advance()

	end, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	stmt.End = end

	p.consumeNewline()

	// Parse body
	stmt.Body, err = p.parseBlock(TOKEN_NEXT)
	if err != nil {
		return nil, err
	}

	// Expect NEXT
	if p.current.Type != TOKEN_NEXT {
		return nil, p.error("expected NEXT")
	}
	p.advance()

	// Optional variable name after NEXT
	if p.current.Type == TOKEN_IDENTIFIER {
		if p.current.Value != stmt.Variable {
			return nil, p.error("NEXT variable '%s' doesn't match FOR variable '%s'", p.current.Value, stmt.Variable)
		}
		p.advance()
	}

	p.consumeNewlineOrEOF()
	return stmt, nil
}

// parseBreakStatement parses: BREAK
func (p *Parser) parseBreakStatement() (*BreakStatement, error) {
	stmt := &BreakStatement{
		Pos: Pos{Line: p.current.Line, Column: p.current.Column},
	}
	p.advance()
	p.consumeNewlineOrEOF()
	return stmt, nil
}

// parseFunctionStatement parses: FUNCTION name(params): ... ENDFUNCTION
func (p *Parser) parseFunctionStatement() (*FunctionStatement, error) {
	stmt := &FunctionStatement{
		Pos: Pos{Line: p.current.Line, Column: p.current.Column},
	}

	p.advance() // consume FUNCTION

	if p.current.Type != TOKEN_IDENTIFIER {
		return nil, p.error("expected function name")
	}
	stmt.Name = p.current.Value
	p.advance()

	if p.current.Type != TOKEN_LPAREN {
		return nil, p.error("expected '(' after function name")
	}
	p.advance()

	// Parse parameters
	stmt.Params = []string{}
	for p.current.Type != TOKEN_RPAREN {
		if p.current.Type != TOKEN_IDENTIFIER {
			return nil, p.error("expected parameter name")
		}
		stmt.Params = append(stmt.Params, p.current.Value)
		p.advance()

		if p.current.Type == TOKEN_COMMA {
			p.advance()
		} else if p.current.Type != TOKEN_RPAREN {
			return nil, p.error("expected ',' or ')' in parameter list")
		}
	}
	p.advance() // consume )

	// Optional colon
	if p.current.Type == TOKEN_COLON {
		p.advance()
	}
	p.consumeNewline()

	// Parse body
	var err error
	stmt.Body, err = p.parseBlock(TOKEN_ENDFUNCTION)
	if err != nil {
		return nil, err
	}

	if p.current.Type != TOKEN_ENDFUNCTION {
		return nil, p.error("expected ENDFUNCTION")
	}
	p.advance()
	p.consumeNewlineOrEOF()

	return stmt, nil
}

// parseReturnStatement parses: RETURN [expr]
func (p *Parser) parseReturnStatement() (*ReturnStatement, error) {
	stmt := &ReturnStatement{
		Pos: Pos{Line: p.current.Line, Column: p.current.Column},
	}
	p.advance() // consume RETURN

	// Check if there's an expression (not newline or EOF)
	if p.current.Type != TOKEN_NEWLINE && p.current.Type != TOKEN_EOF {
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		stmt.Value = expr
	}

	p.consumeNewlineOrEOF()
	return stmt, nil
}

// parsePrintStatement parses: PRINT expr
func (p *Parser) parsePrintStatement() (*PrintStatement, error) {
	stmt := &PrintStatement{
		Pos: Pos{Line: p.current.Line, Column: p.current.Column},
	}
	p.advance() // consume PRINT

	expr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	stmt.Value = expr

	p.consumeNewlineOrEOF()
	return stmt, nil
}

// parseBlock parses statements until one of the terminator tokens is found
func (p *Parser) parseBlock(terminators ...TokenType) ([]Statement, error) {
	var statements []Statement

	for !p.isAtEnd() {
		p.skipNewlines()
		if p.isAtEnd() {
			break
		}

		// Check for terminators
		for _, t := range terminators {
			if p.current.Type == t {
				return statements, nil
			}
		}

		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		if stmt != nil {
			statements = append(statements, stmt)
		}
	}

	return statements, nil
}

// -----------------------------------------------------------------------------
// Expression Parsing (Pratt parser / precedence climbing)
// -----------------------------------------------------------------------------

type precedence int

const (
	precNone       precedence = iota
	precOr                    // OR
	precAnd                   // AND
	precEquality              // = <> !=
	precComparison            // < > <= >=
	precTerm                  // + -
	precFactor                // * /
	precUnary                 // NOT -
	precCall                  // ()
)

func (p *Parser) parseExpression() (Expression, error) {
	return p.parsePrecedence(precOr)
}

func (p *Parser) parsePrecedence(minPrec precedence) (Expression, error) {
	left, err := p.parseUnary()
	if err != nil {
		return nil, err
	}

	for {
		prec := p.getInfixPrecedence()
		if prec < minPrec {
			break
		}

		op := p.current.Type
		p.advance()

		right, err := p.parsePrecedence(prec + 1)
		if err != nil {
			return nil, err
		}

		line, col := left.Position()
		left = &BinaryExpr{
			Pos:      Pos{Line: line, Column: col},
			Left:     left,
			Operator: op,
			Right:    right,
		}
	}

	return left, nil
}

func (p *Parser) getInfixPrecedence() precedence {
	switch p.current.Type {
	case TOKEN_OR:
		return precOr
	case TOKEN_AND:
		return precAnd
	case TOKEN_EQ, TOKEN_NEQ:
		return precEquality
	case TOKEN_LT, TOKEN_GT, TOKEN_LTE, TOKEN_GTE:
		return precComparison
	case TOKEN_PLUS, TOKEN_MINUS:
		return precTerm
	case TOKEN_STAR, TOKEN_SLASH:
		return precFactor
	default:
		return precNone
	}
}

func (p *Parser) parseUnary() (Expression, error) {
	if p.current.Type == TOKEN_NOT || p.current.Type == TOKEN_MINUS {
		pos := Pos{Line: p.current.Line, Column: p.current.Column}
		op := p.current.Type
		p.advance()

		operand, err := p.parseUnary()
		if err != nil {
			return nil, err
		}

		return &UnaryExpr{Pos: pos, Operator: op, Operand: operand}, nil
	}

	return p.parseCall()
}

func (p *Parser) parseCall() (Expression, error) {
	expr, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}

	// Check for function call
	if ident, ok := expr.(*Identifier); ok && p.current.Type == TOKEN_LPAREN {
		pos := ident.Pos
		p.advance() // consume (

		args, err := p.parseArguments()
		if err != nil {
			return nil, err
		}

		return &CallExpr{Pos: pos, Name: ident.Name, Args: args}, nil
	}

	return expr, nil
}

func (p *Parser) parseArguments() ([]Expression, error) {
	args := []Expression{}

	if p.current.Type == TOKEN_RPAREN {
		p.advance()
		return args, nil
	}

	for {
		arg, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		args = append(args, arg)

		if p.current.Type == TOKEN_COMMA {
			p.advance()
		} else {
			break
		}
	}

	if p.current.Type != TOKEN_RPAREN {
		return nil, p.error("expected ')' after arguments")
	}
	p.advance()

	return args, nil
}

func (p *Parser) parsePrimary() (Expression, error) {
	pos := Pos{Line: p.current.Line, Column: p.current.Column}

	switch p.current.Type {
	case TOKEN_INT:
		value, err := strconv.ParseInt(p.current.Value, 10, 64)
		if err != nil {
			return nil, p.error("invalid integer: %s", p.current.Value)
		}
		p.advance()
		return &IntLiteral{Pos: pos, Value: value}, nil

	case TOKEN_FLOAT:
		value, err := strconv.ParseFloat(p.current.Value, 64)
		if err != nil {
			return nil, p.error("invalid float: %s", p.current.Value)
		}
		p.advance()
		return &FloatLiteral{Pos: pos, Value: value}, nil

	case TOKEN_STRING:
		value := p.current.Value
		p.advance()
		return &StringLiteral{Pos: pos, Value: value}, nil

	case TOKEN_TRUE:
		p.advance()
		return &BoolLiteral{Pos: pos, Value: true}, nil

	case TOKEN_FALSE:
		p.advance()
		return &BoolLiteral{Pos: pos, Value: false}, nil

	case TOKEN_IDENTIFIER:
		name := p.current.Value
		p.advance()
		return &Identifier{Pos: pos, Name: name}, nil

	case TOKEN_LPAREN:
		p.advance() // consume (
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		if p.current.Type != TOKEN_RPAREN {
			return nil, p.error("expected ')' after expression")
		}
		p.advance()
		return expr, nil

	default:
		return nil, p.error("unexpected token in expression: %s", p.current.Type)
	}
}

// -----------------------------------------------------------------------------
// Helper methods
// -----------------------------------------------------------------------------

func (p *Parser) isAtEnd() bool {
	return p.current.Type == TOKEN_EOF
}

func (p *Parser) advance() {
	p.pos++
	if p.pos < len(p.tokens) {
		p.current = p.tokens[p.pos]
	}
}

func (p *Parser) skipNewlines() {
	for p.current.Type == TOKEN_NEWLINE {
		p.advance()
	}
}

func (p *Parser) consumeNewline() {
	if p.current.Type == TOKEN_NEWLINE {
		p.advance()
	}
}

func (p *Parser) consumeNewlineOrEOF() {
	if p.current.Type == TOKEN_NEWLINE {
		p.advance()
	}
	// EOF is also acceptable at end of statement
}

func (p *Parser) error(format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
	return fmt.Errorf("line %d, column %d: %s", p.current.Line, p.current.Column, msg)
}

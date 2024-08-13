// This is the parser for the Clear programming language
// It uses Vaughan Pratt's "top "down operator precedence" parser defined here: https://dl.acm.org/doi/pdf/10.1145/512927.512931

package parser

import (
	"fmt"
	"strconv"

	"github.com/ajtroup1/clearv2/ast"
	"github.com/ajtroup1/clearv2/lexer"
	"github.com/ajtroup1/clearv2/token"
)

// Iota of precedences representing their integer 'powers'
const (
	_           int = iota
	LOWEST          // Lowest precedence level, used as a base
	EQUALS          // Precedence level for '==' and '!='
	LESSGREATER     // Precedence level for '<' and '>'
	SUM             // Precedence level for '+' and '-'
	PRODUCT         // Precedence level for '*' and '/'
	PREFIX          // Precedence level for prefix operators like '-X' or '!X'
	CALL            // Precedence level for function calls like 'myFunction(X)'
)

// Maps tokens to their corresponding precedence levels
var precedences = map[token.TokenType]int{ // Precedence table
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
}

type (
	// prefixParseFn is a function that parses expressions with a prefix operator.
	// For example, in the expression "-5", the "-" is a prefix operator.
	prefixParseFn func() ast.Expression

	// infixParseFn is a function that parses expressions with an infix operator.
	// For example, in the expression "2 + 3", the "+" is an infix operator.
	infixParseFn func(ast.Expression) ast.Expression
)

type Parser struct {
	l         *lexer.Lexer // lexer that supplies the tokens
	curToken  token.Token  // The current token being examined
	peekToken token.Token  // The token being compared to the currToken, or the next token to be examined
	errors    []string     // List of errors accrued when parsing the source code

	prefixParseFns map[token.TokenType]prefixParseFn // Registered prefix parsing functions
	infixParseFns  map[token.TokenType]infixParseFn  // Registered infix parsing functions
}

// Associates a token type with a prefix parse function
func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

// Associates a token type with an infix parse function
func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

// Instantiates a new instances of Parser given a lexer containing a stream of tokens from the source code
func New(l *lexer.Lexer) *Parser {
	// Instantiate parser object
	p := &Parser{l: l, errors: []string{}}

	// Register all prefix parsing functions
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)

	// Register all infix parsing functions
	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()
	return p
}

// Returns the list of errors accrued when parsing
func (p *Parser) Errors() []string {
	return p.errors
}
func (p *Parser) nextToken() {
	// 'consume' method
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// Parses the entire program and returns the root node of the AST
func (p *Parser) ParseProgram() *ast.Program {
	// Returns a list of statements given tokens
	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	for !p.curTokenIs(token.EOF) { // Loop until the end of input
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program

}

// Parses and identifier and returns it as an expression node
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

// Evaluates which type of statement to parse based on the current token
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	// Unless explicitly defined as LET or RETURN, most everything is an expression
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	// let x = 5
	stmt := &ast.LetStatement{Token: p.curToken} // Let token
	// Identifier (x, y ...) follows let keyword
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	// "=" follows the identifier
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}
	p.nextToken()
	// And any expression follows the "="
	stmt.Value = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken} // Return token
	p.nextToken()
	stmt.ReturnValue = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) { // As long as the next token doesn't end the statement
		p.nextToken()
	}
	return stmt
}

// Parses an expression as a statement
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST) // Start parsing with the lowest precedence

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// Parses an expression given a precedence
// The heart of the Pratt parset
func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type] // Lookup prefixParseFn for current token type
	if prefix == nil {                          // If there isn't one, this situation is unaccounted for
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}

	leftExp := prefix() // Parse the prefix expression

	// Continue parsing expressions as long as they have a higher precedence and it isn't the end of the line
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}
	return leftExp
}

// Parses an integer literal and returns it as an expression node
func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}             // Instantiates a literal value for the currToken
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64) // Uses strconv to parse from string to int64
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.Value = value
	return lit
}

// Parses an expression with a prefix operator: "!", "-"
func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,         // The prefix operator token
		Operator: p.curToken.Literal, // The prefix operator itself
	}

	// Advance to parse the expression that follows the prefix operator
	p.nextToken()

	// Parse the expression after the prefix operator
	expression.Right = p.parseExpression(PREFIX)

	return expression
}

// Parses functions with an infix operator: "+", "*", "=="...
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,         // The infix operator token
		Operator: p.curToken.Literal, // The infix operator itself
		Left:     left,               // The expression to the left of the infix operator
	}

	// Retreive the precedence of the infix operator
	precedence := p.curPrecedence()

	// Advance to the expression that follows the infix operator
	p.nextToken()

	// Parse the expression after the infix operator
	expression.Right = p.parseExpression(precedence)

	return expression
}

// Parses a boolean literal: "true", "false"
func (p *Parser) parseBoolean() ast.Expression {
	// Create a boolean node with the token's value
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

// Parses an expression encased in parentheses
func (p *Parser) parseGroupedExpression() ast.Expression {
	// Advance past open parenthesis
	p.nextToken()

	// Parse the expression inside the parentheses
	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return exp
}

// Parses an if expression: "if (condition) {x}" and returns an expression
func (p *Parser) parseIfExpression() ast.Expression {
	// Instantiate if expression token
	expression := &ast.IfExpression{Token: p.curToken}
	// Must receive a condition encased within parentheses: "if (x < y)"

	// "("
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	p.nextToken()
	// "x < y"
	expression.Condition = p.parseExpression(LOWEST)
	// ")"
	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	// Check for required consequence
	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	// Optionally parses an else expression, or alternative
	if p.peekTokenIs(token.ELSE) { // If the if expression contains an else
		p.nextToken()
		// Must contain a left brace to encase alternative
		if !p.expectPeek(token.LBRACE) {
			return nil
		}
		// Assign the statement to the alternative of the if expression
		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

// Parses a block statement: "{x}", "{add(5, 7) * 2}", ...
func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	// Instantiate block statement token
	block := &ast.BlockStatement{Token: p.curToken}
	// Initialize the list of statements contained in the block
	block.Statements = []ast.Statement{}
	p.nextToken()
	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) { // As long as the token isn't the end of the block "}" or the end of the file (illegal)
		// Parse the statement and add it to the list of statements in the block
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}
	return block
}

// Parses a function literal expression
func (p *Parser) parseFunctionLiteral() ast.Expression {
	// Instantiate the function object
	lit := &ast.FunctionLiteral{Token: p.curToken}

	// Parse the parameters, which are encased in parentheses and separated by commas
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	lit.Parameters = p.parseFunctionParameters()
	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	// Parse the function body, which is just a block statement
	lit.Body = p.parseBlockStatement()
	return lit
}

// Parses the parameter list as a slice of identifier for a function literal
func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	// Check if the parameter list is empty (right paren immedietely follows left paren: "fn()")
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		// If so, return the empty slice
		return identifiers
	}
	p.nextToken()
	// Instantiate first parameter as an identifier and add it to the slice
	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)
	for p.peekTokenIs(token.COMMA) { // Continue to parse params checking if there is another listed ahead
		// Consume ident and comma
		p.nextToken()
		p.nextToken()
		// Instantiate next param
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
	}
	// Must conclude param list with right paren
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return identifiers
}

// Parses the call to a defined function
func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	// Instantiate a call expression with a given function
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	// Parse the arguments list
	exp.Arguments = p.parseCallArguments()
	return exp
}

// Parses the list of function call arguments and returns them as a slice of expression
// Works similarly to parseFunctionParameters() above
func (p *Parser) parseCallArguments() []ast.Expression {
	// Instantiate the slice
	args := []ast.Expression{}
	// Arguments list must be encased in parentheses
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		// If not, return the empty slice
		return args
	}
	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))
	for p.peekTokenIs(token.COMMA) { // Continue through comma separated list and parse the individual arguments
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return args
}

// Check for if the CURRENT token matches the sent token type (param)
func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

// Check for if the PEEK token matches the sent token type (param)
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// Checks for if the PEEK token matches the sent token type (param) and advances if it does
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t) // Record error if the token type doesn't match
		return false
	}
}

// Returns the precedence of the peek token type. Defaults to LOWEST if it doesn't have one
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

// Returns the precedence of the current token type. Defaults to LOWEST if it doesn't have one
func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

// Returns an error msg if the next token doesn't match the send token type (param)
func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

// Records an error message if no prefix parse function is found for the current token type
func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

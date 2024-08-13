// Lexer for the Clear programming language
// Converts the given source code into a stream of tokens to be analyzed by the parser
// This is a basic and common implementation of a lexer used in many languages
package lexer

import "github.com/ajtroup1/clearv2/token"

// Lexer struct contains the data necessary for lexical analysis
// input: The entire source code to be tokenized
// position: Current position in the input string
// readPosition: Next position to read in the input string
// ch: Current character being examined
type Lexer struct {
	input        string // The entire source code
	position     int    // Current position in the input string
	readPosition int    // Next position to read in the input string
	ch           byte   // Current character under examination
}

// Creates a new Lexer instance with the given source code
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar() // Initialize the first character
	return l
}

// Reads the next character from the input string and updates the lexer state
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) { // Check if the end of input is reached
		l.ch = 0 // Null character indicating end of input
	} else {
		l.ch = l.input[l.readPosition] // Read the current character
	}
	l.position = l.readPosition // Update the current position
	l.readPosition += 1 // Move to the next character
}

// Returns the next token from the input stream
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace() // Skip any whitespace characters

	// Tokenize based on the current character
	switch l.ch {
	case '=':
		if l.peekChar() == '=' { // Check for comparison "=="
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.ASSIGN, l.ch) // Single '='
		}
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		if l.peekChar() == '=' { // Check for counter-comparison "!="
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.NOT_EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.BANG, l.ch) // Single '!'
		}
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF // End of file
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier() // Read an identifier
			tok.Type = token.LookupIdent(tok.Literal) // Lookup identifier token type
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT // Integer literal
			tok.Literal = l.readNumber() // Read the number
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch) // Illegal character
		}
	}

	l.readChar() // Read the next character
	return tok
}

// Creates a new token of the specified type with a given character
func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

// Skips any whitespace characters (spaces, tabs, newlines, etc.) in the input
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar() // Move to the next character
	}
}

// Reads an identifier from the input
// An identifier is a sequence of letters and underscores
func (l *Lexer) readIdentifier() string {
	position := l.position // Start position of the identifier
	for isLetter(l.ch) {
		l.readChar() // Move to the next character
	}
	return l.input[position:l.position] // Return the identifier
}

// Determines if the current character is a valid letter or underscore for identifiers
// This function can be adjusted to match the identifier rules of your language
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// Reads a sequence of digits from the input
func (l *Lexer) readNumber() string {
	position := l.position // Start position of the number
	for isDigit(l.ch) {
		l.readChar() // Move to the next character
	}
	return l.input[position:l.position] // Return the number
}

// Determines if the given character is a digit
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// Peeks at the next character in the input without advancing the read position
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0 // End of input
	} else {
		return l.input[l.readPosition] // Return the next character
	}
}

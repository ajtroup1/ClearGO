// Defines the token types accounted for in the Clear programming language
package token

// Represents the type of token in string format
type TokenType string

// Represents a single token object in the Clear programming language
// Tokens have a type (keyword, operator, ...) and a literal value associated with it (+, 5, x, ...)
type Token struct {
	Type    TokenType
	Literal string
}

// Constants for various token types used in the Clear language
const (
	ILLEGAL = "ILLEGAL" // Represents an unrecognized token
	EOF     = "EOF"     // End of file

	// Identifiers and literals
	IDENT = "IDENT" // General identifier (e.g., variable names, function names)
	INT   = "INT"   // Integer literal (e.g., 12345)

	// Operators
	ASSIGN   = "="  // Assignment operator
	EQ       = "==" // Equality operator
	NOT_EQ   = "!=" // Not-equal operator
	PLUS     = "+"  // Addition operator
	MINUS    = "-"  // Subtraction operator
	BANG     = "!"  // Logical negation (not) operator
	ASTERISK = "*"  // Multiplication operator
	SLASH    = "/"  // Division operator
	LT       = "<"  // Less-than operator
	GT       = ">"  // Greater-than operator

	// Delimiters
	COMMA     = "," // Comma separator
	SEMICOLON = ";" // Semicolon separator
	LPAREN    = "(" // Left parenthesis
	RPAREN    = ")" // Right parenthesis
	LBRACE    = "{" // Left brace (beginning of a block)
	RBRACE    = "}" // Right brace (end of a block)

	// Keywords
	FUNCTION = "FUNCTION" // Function keyword (e.g., function definitions)
	LET      = "LET"      // Let keyword (variable declarations)
	TRUE     = "TRUE"     // Boolean literal true
	FALSE    = "FALSE"    // Boolean literal false
	IF       = "IF"       // If keyword (conditional statements)
	ELSE     = "ELSE"     // Else keyword (alternative conditional branches)
	RETURN   = "RETURN"   // Return keyword (function return statements)
)

// Keyword map for reserved words in Clear
var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

// Check for if the given identifier exists as a reserved word in Clear
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		// If it is, return the corresponding token type
		return tok
	}
	// If not, it must just be an identifier
	return IDENT
}

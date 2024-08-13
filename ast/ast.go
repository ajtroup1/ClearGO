// Defines the AST as well as its nodes
// Contains definitions for every node accounted for in Clear
package ast

import (
	"bytes"
	"strings"

	"github.com/ajtroup1/clearv2/token"
)

// High-level node structure that serves as the foundation for all nodes in Clear
// Absolutely ALL nodes in Clear must implement the TokenLiteral() and String() methods since they all implement Node
type Node interface {
	TokenLiteral() string // Returns the literal value of the given node. Used extensively and necessary for all nodes
	String() string       // Simple method that returns a string representation of the given node
}

// Node containing a statement. Statements are evaulted lines such as "let x = 5", "return x"...
// Clear code is a Program node made up of a slice of these statements
type Statement interface {
	Node
	statementNode() // Marker method used to distinguish statements from expressions. Implement this if the type is a statement
}

// Node containing an expression. An expression is a stream of tokens waiting to be evaluated such as "1 + 2", "x = true"...
type Expression interface {
	Node
	expressionNode() // Marker method used to distinguish statements from expressions. Implement this if the type is an expression
}

// Represents the entire program. The "root" node of the AST
type Program struct {
	Statements []Statement // A Clear program is just a slice of statements
}

// Returns the first token's literal value (as long as it contains at least one statement)
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

// Returns the string representation of the entire program
// Concatentates the string representation of all the program's statements
func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString((s.String()))
	}

	return out.String()
}

// List of statements & expressions accounted for in Clear's AST
// ALL statements & expressions must implement the TokenLiteral() and String() methods

// LET statement
type LetStatement struct {
	Token token.Token // The token.LET token
	Name  *Identifier // Name of the identifier: "x", "foobar"...
	Value Expression  // Value stored in the variable: "let x = 5", 5 is the value
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

func (ls *LetStatement) String() string {
	// let x = 5;
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ") // "let "
	out.WriteString(ls.Name.String())        // "x"
	out.WriteString(" = ")                   // " = "

	if ls.Value != nil {
		out.WriteString(ls.Value.String()) // "5"
	}

	out.WriteString(";") // ";"

	return out.String()
}

// The identifier for a let statement / variable: "x", "foobar"
// Identifiers are treated as expressions because they represent values that can be evaluated.
type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string      // Actual string of the name of the ident
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

// Return statement
type ReturnStatement struct {
	Token       token.Token // the token.RETURN token
	ReturnValue Expression  // Value being returned (to the right of "return"): "0", "x"...
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }

func (rs *ReturnStatement) String() string {
	// return x + 5;
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ") // "return "

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String()) // "x + 5"
	}

	out.WriteString(";") // ";"

	return out.String()
}

// Represents a statement consisting of a single expression
type ExpressionStatement struct {
	Token      token.Token // The first token of the expression
	Expression Expression  // The expression itself
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// Represents an integer value
// Integer literals are considered expressions because they represent values that can be evaluated in arithmetic operations OR assigned to variables.
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

// Represents ant prefix expression. In Clear, these are only "!" and "-"
type PrefixExpression struct {
	Token    token.Token // The prefix token: "!", "-"
	Operator string      // The actual operator representaion
	Right    Expression  // Expression to the right of the operator: In "!myFunction()" the "myFunction()" would be 'Right'
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	// Groups the prefix operator with its operand using parentheses
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

// Represents infix expression. These are most commmon expressions: "1 + 2", "x * 2.5"...
type InfixExpression struct {
	// EX. "1 + 2"
	Token    token.Token // Token represents the operator token in the infix expression: "+", "*"...
	Left     Expression  // The left 'value' of the expression: "1"
	Operator string      // The operator in the expression: "+"
	Right    Expression  // The right 'value' of the expression: "2"
}

func (oe *InfixExpression) expressionNode()      {}
func (oe *InfixExpression) TokenLiteral() string { return oe.Token.Literal }
func (oe *InfixExpression) String() string {
	// Groups expression elements together using parentheses
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(oe.Left.String())
	out.WriteString(" " + oe.Operator + " ")
	out.WriteString(oe.Right.String())
	out.WriteString(")")
	return out.String()
}

// Represents a boolean value: true, false
type Boolean struct {
	Token token.Token // The token.TRUE or token.FALSE token
	Value bool        // The GO value of the given token
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }

// Represents an if expression
// If expressions contain an if token, a condition to be rendered, something that happens if it renders true, and optionally an alternative for if it renders false
type IfExpression struct {
	Token       token.Token     // The 'if' token
	Condition   Expression      // What is being evaluated. Can be any expression
	Consequence *BlockStatement // What happens if the condition is true
	Alternative *BlockStatement // What happens if the condition is false
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())
	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}
	return out.String()
}

// Represents a block statement, which is just a series a statements
// Like in if else possibly containing a list of statements to execute depending on a result
type BlockStatement struct {
	Token      token.Token // the { token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	for _, s := range bs.Statements {
		// Output each statement present in the slice
		out.WriteString(s.String())
	}
	return out.String()
}

// Represents a function literal, which is an expression
// Comprised of "fn" keyword, list of params enclosed in parentheses and separated by commas, and a body enclosed in braces
// EX. let myFunction = fn(x, y) { return x + y; }
type FunctionLiteral struct {
	Token      token.Token // The 'fn' token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}
	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())
	return out.String()
}

// Represents a call to a defined function
// Contains a function identifier and a list of function arguments encased in parentheses and separated by commas
type CallExpression struct {
	Token     token.Token  // The '(' token
	Function  Expression   // Identifier or FunctionLiteral
	Arguments []Expression // Arguments being passed in as function params
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out bytes.Buffer
	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}
	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	return out.String()
}

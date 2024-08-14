package evaluator

import (
	"github.com/ajtroup1/clearv2/ast"
	"github.com/ajtroup1/clearv2/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

// The core evaluation function. Traverses the AST from the ast.Program down
// Evaluates the given type of node and returns it as the corresponding evaluated value
func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	// Statements
	// Initially evalute the entire program recursively
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)
	}

	// Unrecognized
	return nil
}

// Receives a list of statements and returns them one by one
func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object
	// Initially, evalute the entire slice of statements in the program
	for _, statement := range stmts {
		result = Eval(statement)
	}
	return result
}

// Converts native boolean to our boolean object
func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

// Applies the native prefix operator to the operand of the right expression
func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return NULL
	}
}

// Evaluates the native bang prefix operator to the right expression operand
func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

// Evaluates the native negaitve prefix operator to the right expression operand
func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return NULL
	}
	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

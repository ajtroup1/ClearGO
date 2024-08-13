package object

import "fmt"

// String representation of the object's type. Similar to TokenType in token
type ObjectType string

const (
	INTEGER_OBJ = "INTEGER"
	BOOLEAN_OBJ = "BOOLEAN"
	NULL_OBJ    = "NULL"
)

// When evaluating input source code, data is parsed into the respective node. That node is then turned into a Object.Integer, for example
type Object interface {
	Type() ObjectType
	Inspect() string
}

// Represents integers, taking ast.IntegerLiteral
type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

// Represents booleans, taking ast.Boolean
type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

// Represents a null value. Doesn't wrap any data, but represents the absence of a value
type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }

package ast

import (
	"testing"

	"github.com/ajtroup1/clearv2/token"
)

const (
	Red    = "\033[31m"
	Yellow = "\033[33m"
	Green  = "\033[32m"
	Reset  = "\033[0m"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}
	expected := "let myVar = anotherVar;"
	actual := program.String()
	if actual != expected {
		t.Errorf(Red+"program.String() wrong. expected=%q, got=%q"+Reset, expected, actual)
	} else {
		t.Logf(Green+"program.String() is correct. got=%q"+Reset, actual)
	}
}

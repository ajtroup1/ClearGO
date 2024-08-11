package parser

import (
	"fmt"
	"testing"

	"github.com/ajtroup1/clearv2/ast"
	"github.com/ajtroup1/clearv2/lexer"
)

const (
	Red    = "\033[31m"
	Yellow = "\033[33m"
	Green  = "\033[32m"
	Reset  = "\033[0m"
)

func TestLetStatements(t *testing.T) {
	input := `
		let x = 5;
		let y = 10;
		let foobar = 838383;
		`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf(Red + "ParseProgram() returned nil" + Reset)
	}
	if len(program.Statements) != 3 {
		t.Fatalf(Red+"program.Statements does not contain 3 statements. got %d"+Reset,
			len(program.Statements))
	}
	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}
	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		} else {
			t.Logf(Green+"Test passed for let statement: %s"+Reset, tt.expectedIdentifier)
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf(Red+"s.TokenLiteral not 'let'. got=%q"+Reset, s.TokenLiteral())
		return false
	}
	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf(Red+"s not *ast.LetStatement. got=%T"+Reset, s)
		return false
	}
	if letStmt.Name.Value != name {
		t.Errorf(Red+"letStmt.Name.Value not '%s'. got=%s"+Reset, name, letStmt.Name.Value)
		return false
	}
	if letStmt.Name.TokenLiteral() != name {
		t.Errorf(Red+"s.Name not '%s'. got=%s"+Reset, name, letStmt.Name)
		return false
	}
	return true
}

func TestReturnStatements(t *testing.T) {
	input := `
	return 5;
	return 10;
	return 99999;
	`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf(Red+"expected 3 statements, got=%d"+Reset, program.Statements)
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf(Red+"stmt not *ast.returnStatement. got=%T"+Reset, stmt)
			continue
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf(Red+"returnStmt.TokenLiteral not 'return', got=%q"+Reset,
				returnStmt.TokenLiteral())
		} else {
			t.Logf(Green+"Test passed for return statement with value: %s"+Reset, returnStmt.TokenLiteral())
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf(Red+"program has not enough statements. got=%d"+Reset,
			len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf(Red+"program.Statements[0] is not ast.ExpressionStatement. got=%T"+Reset,
			program.Statements[0])
	}
	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf(Red+"exp not *ast.Identifier. got=%T"+Reset, stmt.Expression)
	}
	if ident.Value != "foobar" {
		t.Errorf(Red+"ident.Value not %s. got=%s"+Reset, "foobar", ident.Value)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Errorf(Red+"ident.TokenLiteral not %s. got=%s"+Reset, "foobar",
			ident.TokenLiteral())
	} else {
		t.Logf(Green+"Test passed for identifier: %s"+Reset, ident.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf(Red+"program has not enough statements. got=%d"+Reset,
			len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement) // type assertion. checks the type and possibly casts it to *ast.Exp... ok holds whether this was successful
	if !ok {
		t.Fatalf(Red+"program.Statements[0] is not ast.ExpressionStatement. got=%T"+Reset,
			program.Statements[0])
	}
	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf(Red+"exp not *ast.IntegerLiteral. got=%T"+Reset, stmt.Expression)
	}
	if literal.Value != 5 {
		t.Errorf(Red+"literal.Value not %d. got=%d"+Reset, 5, literal.Value)
	}
	if literal.TokenLiteral() != "5" {
		t.Errorf(Red+"literal.TokenLiteral not %s. got=%s"+Reset, Yellow+"5"+Reset,
			literal.TokenLiteral())
	} else {
		t.Logf(Green+"Test passed for integer literal: %s"+Reset, literal.TokenLiteral())
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
	}
	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		if len(program.Statements) != 1 {
			t.Fatalf(Red+"program.Statements does not contain %d statements. got=%d\n"+Reset,
				1, len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf(Red+"program.Statements[0] is not ast.ExpressionStatement. got=%T"+Reset,
				program.Statements[0])
		}
		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf(Red+"stmt is not ast.PrefixExpression. got=%T"+Reset, stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Fatalf(Red+"exp.Operator is not '%s'. got=%s"+Reset,
				tt.operator, exp.Operator)
		}
		if !testIntegerLiteral(t, exp.Right, tt.integerValue) {
			return
		} else {
			t.Logf(Green+"Test passed for prefix operator: %s"+Reset, exp.Operator)
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
	}
	passCount := 0
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		actual := program.String()
		if actual != tt.expected {
			t.Errorf(Red+"expected=%q, got=%q"+Reset, tt.expected, actual)
		} else {
			passCount++
		}
	}
	if passCount > 0 {
		t.Logf(Green+"%d/%d infix operator tests passed"+Reset, passCount, len(tests))
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()

	if len(errors) == 0 {
		return
	}

	t.Errorf(Yellow+"parser encountered %d errors"+Reset, len(errors))
	for _, msg := range errors {
		t.Errorf(Red+"parser error: %q"+Reset, msg)
	}
	t.FailNow()
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}
	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}
	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d. got=%s", value,
			integ.TokenLiteral())
		return false
	}
	return true
}

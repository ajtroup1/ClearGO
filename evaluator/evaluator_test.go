package evaluator

import (
	"testing"

	"github.com/ajtroup1/clearv2/lexer"
	"github.com/ajtroup1/clearv2/object"
	"github.com/ajtroup1/clearv2/parser"
)

const (
	Red    = "\033[31m"
	Yellow = "\033[33m"
	Green  = "\033[32m"
	Reset  = "\033[0m"
)

func logTestResult(t *testing.T, passed bool, testName string) {
	if passed {
		t.Logf(Green+"%s passed"+Reset, testName)
	} else {
		t.Logf(Red+"%s failed"+Reset, testName)
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	return Eval(program)
}

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
	}

	passed := true
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		if !testIntegerObject(t, evaluated, tt.expected) {
			passed = false
		}
	}

	logTestResult(t, passed, "TestEvalIntegerExpression")
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d",
			result.Value, expected)
		return false
	}
	return true
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t",
			result.Value, expected)
		return false
	}
	return true
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	passed := true
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		if !testBooleanObject(t, evaluated, tt.expected) {
			passed = false
		}
	}

	if passed {
		logTestResult(t, true, "TestBangOperator")
	} else {
		logTestResult(t, false, "TestBangOperator")
	}
}

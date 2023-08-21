package eval

import (
	"testing"

	"github.com/Arch-4ng3l/Monkey/lexer"
	"github.com/Arch-4ng3l/Monkey/object"
	"github.com/Arch-4ng3l/Monkey/parser"
)

func TestFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"let f = func(x) { x; }; f(1);", 1},
		{"let f = func(x) { x * x; }; f(2);", 4},
		{"let f = func(x) { x / x; }; f(1);", 1},
		{"let f = func(x) { x * x; }; f(f(2));", 16},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIngegerObject(t, evaluated, tt.expected)

	}
}
func TestFuncObject(t *testing.T) {
	input := "func(x) { x + 2; };"
	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("Object is not A Function got %T", evaluated)
	}

	if len(fn.Params) != 1 {
		t.Fatalf("Wrong Amount of Params got %d", len(fn.Params))
	}

	expectedBody := "(x + 2)"

	if fn.Body.String() != expectedBody {
		t.Fatalf("Wrong Function Body %s", fn.Body.String())
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"let x = 10; x;", 10},
		{"let x = 10; let y = 20; y;", 20},
		{"let x = 10 * 10; x;", 100},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIngegerObject(t, evaluated, tt.expected)

	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"return 10; 9", 10},
		{"9; return 10; 9", 10},
		{"9; return 2 * 5; 9", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIngegerObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 5 } else { 20 }", 20},
		{"if (10) { 20 }", 20},
		{"if (false) { 10 }", nil},
		{"if (10 == 10) { 10 }", 10},
		{"if (10 != 10) { 20 } else { 10 }", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		i, ok := tt.expected.(int)
		t.Log(ok)
		if ok {
			testIngegerObject(t, evaluated, i)
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object isnt null got %T", obj)
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
		{"!2", false},
		{"!!false", false},
		{"!!true", true},
		{"!!2", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBoolObject(t, evaluated, tt.expected, tt.input)
	}
}

func TestEvalBoolExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 > 2", false},
		{"2 == 1", false},
		{"2 != 1", true},
		{"1 != 2", true},
		{"1 < 2", true},
		{"true == true", true},
		{"true != false", true},
		{"true == false", false},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBoolObject(t, evaluated, tt.expected, tt.input)
	}
}

func testBoolObject(t *testing.T, obj object.Object, expected bool, input string) bool {
	res, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("Objects is not Boolean got %T", obj)
		return false
	}

	if res.Value != expected {
		t.Errorf(input)
		t.Errorf("Object has wrong value want %t got %t", expected, res.Value)
		return false
	}
	return true
}

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"5", 5},
		{"10", 10},
		{"-10", -10},
		{"-69", -69},
		{"-69 + 10", -59},
		{"10 * 3", 30},
		{"10 * 5 - 10 * 5", 0},
		{"10 / 5", 2},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIngegerObject(t, evaluated, tt.expected)
	}

}

func testEval(input string) object.Object {
	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	program := p.ParseProgram()
	env := object.NewEnv()
	return Eval(program, env)
}

func testIngegerObject(t *testing.T, obj object.Object, expected int) bool {
	res, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("Objects is not Integer got %T", obj)
		return false
	}

	if res.Value != expected {
		t.Errorf("Object has wrong value want %d got %d", expected, res.Value)
		return false
	}
	return true
}

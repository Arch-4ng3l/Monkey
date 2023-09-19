package vm

import (
	"fmt"
	"testing"

	"github.com/Arch-4ng3l/Monkey/ast"
	"github.com/Arch-4ng3l/Monkey/compiler"
	"github.com/Arch-4ng3l/Monkey/lexer"
	"github.com/Arch-4ng3l/Monkey/object"
	"github.com/Arch-4ng3l/Monkey/parser"
)

type vmTestCase struct {
	input    string
	expected interface{}
}

func parse(input string) *ast.Program {
	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	return p.ParseProgram()
}

func TestFunctionCalls(t *testing.T) {
	tests := []vmTestCase{
		{`var f = func() { 5 + 10 }; f();`, 15},
		{`var fo = func() { return 20; };var f = func() { 5 + fo() }; f();`, 25},
		{`var f = func() {}; f();`, Null},
		{`var f = func() {return 1;}; var x = func() { return f; } x()();`, 1},
		{`var f = func() {var f = 10; f};f()`, 10},
	}
	runVmTest(t, tests)
}

func TestArrayExpression(t *testing.T) {
	tests := []vmTestCase{
		{"[]", []int{}},
		{"[][0]", Null},
		{"[1, 2, 3]", []int{1, 2, 3}},
		{"[1 + 2, 2 * 3]", []int{3, 6}},
		{"[1 + 2, 2 * 3][0]", 3},
		{"[1 + 2, 2 * 3][-1]", Null},
		{"[1 + 2, 2 * 3][-1]", Null},
		{"var x = []; x[0]", Null},
		{`["apple", "banana", "cherry"]`, []string{"apple", "banana", "cherry"}},
		{`["Hello", " ", "world"][2]`, "world"},
		{`[[1, 2], [3, 4, 5], [6]]`, [][]int{{1, 2}, {3, 4, 5}, {6}}},
		{`[[1, 2], ["a", "b", "c"], [true, false]]`, [][]interface{}{{1, 2}, {"a", "b", "c"}, {true, false}}},
	}
	runVmTest(t, tests)

}

func TestStringExpressions(t *testing.T) {
	tests := []vmTestCase{
		{`"monkey"`, "monkey"},
		{`"mon" + "key"`, "monkey"},
		{`var x = "mon" + "key"; x[1]`, "o"},
	}
	runVmTest(t, tests)
}

func TestGlobalVariables(t *testing.T) {
	tests := []vmTestCase{
		{"var o = 1; o", 1},
		{"var o = 1; var i = 2; o+i", 3},
		{"var a = 5; var b = 3; a + b", 8},
		{"var x = 10; var y = 20; x + y", 30},
		{"var p = 7; var q = 12; var r = 3; p + q + r", 22},
		{"var a = 10; var b = 3; a - b", 7},
		{"var x = 15; var y = 7; x - y", 8},
		{"var p = 20; var q = 5; var r = 3; p - q - r", 12},
		{"var a = 6; var b = 4; a * b", 24},
		{"var x = 8; var y = 3; x * y", 24},
		{"var p = 5; var q = 2; var r = 3; p * q * r", 30},
		{"var x = 30; var y = 6; x / y", 5},
		{"var a = 15; var b = 3; a / b", 5},
		{"var p = 20; var q = 2; var r = 4; p / q / r", 2},
		{"var a = 4; var b = 3; var c = 2; a + b * c", 10},
		{"var x = 10; var y = 3; var z = 2; x - y / z", 9},
		{"var p = 6; var q = 2; var r = 3; p * q - r", 9},
		{"var a = 15; var b = 20; var z = if (a * 2 < b) { 100 } else { a + b }; z", 35},
		{"var x = 10; var y = 5; var z = if (x < y) { 20 } else { 30 }; z", 30},
		{"var p = 8; var q = 12; var r = 5; var z = if (p < q) { if (q > r) { 50 } else { 60 } } else { if (r < p) { 70 } else { 80 } }; z", 50},
	}

	runVmTest(t, tests)
}

func TestConditionals(t *testing.T) {
	tests := []vmTestCase{
		{"if (true) { 10 }", 10},
		{"if (true) { 10 } else { 20 }", 10},
		{"if (10 < 5) { 10 } else { 20 }", 20},
		{"if (10 > 5) { 10 } else { 20 }", 10},
		{"if (!false) { 10 } else { 20 }", 10},
		{"if (false) { 10 }", Null},
		{"if (1 > 2) { 10 }", Null},
		{"if (!true) { 10 }", Null},
		{"if (true) { 10 } else { 20 }", 10},
		{"if (false) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (10 == 10) { 10 } else { 20 }", 10},
		{"if (10 != 10) { 10 } else { 20 }", 20},
		{"if (!false) { 10 }", 10},
		{"if (!true) { 10 }", Null},
		{"if (2 * 3 == 6) { 10 } else { 20 }", 10},
		{"if (8 / 2 == 4) { 10 } else { 20 }", 10},
		{"if (true) { if (false) { 10 } else { 20 } }", 20},
		{"if (true) { if (true) { 10 } else { 20 } }", 10},
		{"if (false) { if (true) { 10 } else { 20 } }", Null},
		{"if (true) { if (1 > 2) { 10 } else { 20 } }", 20},
		{"if (true) { if (1 < 2) { 10 } else { if (false) { 5 } else { 15 } } }", 10},
	}

	runVmTest(t, tests)
}

func TestBoolArithmetic(t *testing.T) {
	tests := []vmTestCase{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"false != true", true},
		{"!false", true},
		{"!!false", false},
		{"(1 < 2) == true", true},
		{"(1 > 2) == true", false},
		{"(1 > 2) == (2 < 1)", true},
	}
	runVmTest(t, tests)
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []vmTestCase{
		{"1", 1},
		{"-1", -1},
		{"-2 + 2", 0},
		{"1 + 2", 3},
		{"2 - 1", 1},
		{"1 * 2 / 2 + 1", 2},
		{"2 * (2 + 1)", 6},
	}
	runVmTest(t, tests)
}

func testBoolObject(expected bool, actual object.Object) error {
	result, ok := actual.(*object.Boolean)
	if !ok {
		return fmt.Errorf("Object is not a Boolean got %T", actual)
	}
	if result.Value != expected {
		return fmt.Errorf("Object has wrong Value want %t got %t", expected, result.Value)
	}

	return nil
}

func testIntegerObject(expected int, actual object.Object) error {
	result, ok := actual.(*object.Integer)
	if !ok {
		return fmt.Errorf("Object is not an Integer got %T", actual)
	}
	if result.Value != expected {
		return fmt.Errorf("Object has wrong Value want %d got %d", expected, result.Value)
	}

	return nil
}

func testStringObject(expected string, actual object.Object) error {
	result, ok := actual.(*object.String)
	if !ok {
		return fmt.Errorf("Object is not an String got %T", actual)
	}
	if result.Value != expected {
		return fmt.Errorf("Object has wrong Value want %s got %s", expected, result.Value)
	}

	return nil
}

func runVmTest(t *testing.T, tests []vmTestCase) {
	t.Helper()

	for _, tt := range tests {
		program := parse(tt.input)
		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			t.Fatalf("%s", err)
		}
		vm := New(comp.Bytecode())

		err = vm.Run()
		if err != nil {
			t.Fatalf("%s", err)
		}
		stackElem := vm.LastPoppedStackElement()

		testExpectedObject(t, tt.expected, stackElem)
	}
}

func testExpectedObject(t *testing.T, expected interface{}, actual object.Object) {
	t.Helper()
	switch expected := expected.(type) {
	case *object.Null:
		if actual != Null {
			t.Errorf("object is not Null got %T", actual)
		}
	case int:
		err := testIntegerObject(expected, actual)
		if err != nil {
			t.Errorf("test Integer Object failed: %s", err)
		}
	case []int:
		arr, ok := actual.(*object.Array)
		if !ok {
			t.Errorf("object Is Not An Array")
			return
		}

		if len(arr.Elements) != len(expected) {
			return
		}
		for i, el := range expected {
			err := testIntegerObject(el, arr.Elements[i])
			if err != nil {
				t.Errorf("test Integer Object failed: %s", err)
			}
		}
	case bool:
		err := testBoolObject(expected, actual)
		if err != nil {
			t.Errorf("test Bool Object failed: %s", err)
		}
	case string:
		err := testStringObject(expected, actual)
		if err != nil {
			t.Errorf("test String Object failed: %s", err)
		}
	}
}

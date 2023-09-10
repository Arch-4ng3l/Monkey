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
	case int:
		err := testIntegerObject(expected, actual)
		if err != nil {
			t.Errorf("test Integer Object failed: %s", err)
		}
	case bool:
		err := testBoolObject(expected, actual)
		if err != nil {
			t.Errorf("test Bool Object failed: %s", err)
		}
	}
}

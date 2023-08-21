package parser

import (
	"fmt"
	"testing"

	"github.com/Arch-4ng3l/Monkey/ast"
	"github.com/Arch-4ng3l/Monkey/lexer"
)

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		intValue int64
	}{
		{"!5;", "!", 5},
		{"-125;", "-", 125},
	}

	for _, tt := range prefixTests {
		l := lexer.NewLexer(tt.input)
		p := NewParser(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("Wrong Number of Statements got %d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpresssionStatement)
		if !ok {
			t.Fatalf("Program.Statements[0] isnt an ExpressionStatement got %T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("Expression isnt a PrefixExpression got %T", stmt.Expression)
		}

		if exp.Operator != tt.operator {
			t.Fatalf("Wrong Operator got %s", exp.Operator)
		}

		if !testIntLiteral(t, exp.Right, tt.intValue) {
			return
		}

	}
}

func testIntLiteral(t *testing.T, il ast.Expression, value int64) bool {
	n, ok := il.(*ast.IntLiteral)
	if !ok {
		t.Errorf("Expression isnt an IntLiteral got %T", il)
		return false
	}
	if n.Value != value {
		t.Errorf("Wrong Value got %d", n.Value)
		return false
	}
	if n.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("Wrong Token Literal got %s", n.TokenLiteral())
		return false
	}

	return true
}

func TestIntExpression(t *testing.T) {
	input := "5;"

	l := lexer.NewLexer(input)
	p := NewParser(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Program has not enough statements got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpresssionStatement)

	if !ok {
		t.Fatalf("statement has wrong type got %T", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.IntLiteral)

	if !ok {
		t.Fatalf("expression not identifier got %T", stmt.Expression)
	}

	if literal.Value != 5 {
		t.Fatalf("wrong value in identifier got %d", literal.Value)
	}

	if literal.TokenLiteral() != "5" {
		t.Fatalf("wrong value in Token Literal got %s", literal.TokenLiteral())
	}
}

func TestIdentExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.NewLexer(input)
	p := NewParser(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Program has not enough statements got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpresssionStatement)

	if !ok {
		t.Fatalf("statement has wrong type got %T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Ident)

	if !ok {
		t.Fatalf("expression not identifier got %T", stmt.Expression)
	}

	if ident.Value != "foobar" {
		t.Fatalf("wrong value in identifier got %s", ident.Value)
	}

	if ident.TokenLiteral() != "foobar" {
		t.Fatalf("wrong value in Token Literal got %s", ident.TokenLiteral())
	}
}

func TestReturnStatements(t *testing.T) {
	input := `
		return 5;
		return 10;
		return 123;
	`
	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatal("ParseProgram returned nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements doesn't contain 3 statements")
	}
	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not Return Statement got %T", stmt)
			continue
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("Token Literal not return got %q", returnStmt.TokenLiteral())

		}
	}
}

func TestLetStatements(t *testing.T) {
	input := `
		let x = 5;
		let y = 10;
		let foobar = 123;
	`
	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatal("ParseProgram returned nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements doesn't contain 3 statements")
	}
	tests := []struct {
		expectedIdent string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}
	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdent) {
			return
		}
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()

	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error :%q", msg)
	}
	t.FailNow()
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {

	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let' got %q", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement got %T", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s' got %s", name, letStmt.Name.Value)
		return false
	}
	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiter() not '%s' got %s", name, letStmt.Name.TokenLiteral())
		return false
	}

	return true

}

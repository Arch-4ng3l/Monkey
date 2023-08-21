package ast

import (
	"testing"

	"github.com/Arch-4ng3l/Monkey/token"
)

func TestString(t *testing.T) {

	program := &Program{

		Statements: []Statement{

			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Ident{
					Token: token.Token{Type: token.IDENT, Literal: "var"},
					Value: "var",
				},

				Value: &Ident{
					Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	if program.String() != "let var = anotherVar;" {
		t.Errorf("program.String() wrong got %q", program.String())
	}

}

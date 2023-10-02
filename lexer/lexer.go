package lexer

import (
	"github.com/Arch-4ng3l/Monkey/token"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	char         byte
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()

	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.char = 0
	} else {
		l.char = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++

}

func (l *Lexer) readIdent() string {
	position := l.position
	for isLetter(l.char) || isDigit(l.char) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.char) {
		l.readChar()
	}

	return l.input[position:l.position]
}

func (l *Lexer) skipWhitespace() {
	for l.char == ' ' || l.char == '\t' || l.char == '\n' || l.char == '\r' {
		l.readChar()
	}
}

func (l *Lexer) peakChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]

}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.skipWhitespace()
	switch l.char {
	case '=':
		if l.peakChar() == '=' {
			literal := "=="
			l.readChar()

			tok = token.Token{Type: token.EQ, Literal: literal}
		} else {
			tok = newToken(token.ASSIGN, l.char)
		}

	case '!':
		if l.peakChar() == '=' {
			literal := "!="
			l.readChar()
			tok = token.Token{Type: token.NOT_EQ, Literal: literal}
		} else {
			tok = newToken(token.BANG, l.char)
		}
	case '+':
		if l.peakChar() == '=' {
			literal := "+="
			l.readChar()
			tok = token.Token{Type: token.PLUS_ASSIGN, Literal: literal}
		} else {
			tok = newToken(token.PLUS, l.char)
		}
	case ',':
		tok = newToken(token.COMMA, l.char)
	case ';':
		tok = newToken(token.SEMICOLON, l.char)
	case '(':
		tok = newToken(token.LPAREN, l.char)
	case ')':
		tok = newToken(token.RPAREN, l.char)
	case '{':
		tok = newToken(token.LBRACE, l.char)
	case '}':
		tok = newToken(token.RBRACE, l.char)
	case '-':
		if l.peakChar() == '=' {
			literal := "-="
			l.readChar()
			tok = token.Token{Type: token.MINUS_ASSIGN, Literal: literal}
		} else {
			tok = newToken(token.MINUS, l.char)
		}
	case '*':
		if l.peakChar() == '=' {
			literal := "*="
			l.readChar()
			tok = token.Token{Type: token.STAR_ASSIGN, Literal: literal}
		} else {
			tok = newToken(token.STAR, l.char)
		}
	case '/':
		if l.peakChar() == '=' {
			literal := "/="
			l.readChar()
			tok = token.Token{Type: token.SLASH_ASSIGN, Literal: literal}
		} else {
			tok = newToken(token.SLASH, l.char)
		}
	case '<':
		if l.peakChar() == '=' {
			literal := "<="
			l.readChar()
			tok = token.Token{Type: token.LT_EQ, Literal: literal}
		} else {
			tok = newToken(token.LT, l.char)
		}
	case '>':
		if l.peakChar() == '=' {
			literal := "<="
			l.readChar()
			tok = token.Token{Type: token.GT_EQ, Literal: literal}
		} else {
			tok = newToken(token.GT, l.char)
		}

	case '"':
		tok = token.Token{Type: token.STR, Literal: l.readStr()}

	case '[':
		tok = newToken(token.LBRACKET, l.char)

	case ']':
		tok = newToken(token.RBRACKET, l.char)
	case '^':
		tok = newToken(token.POWER, l.char)

	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.char) {
			tok.Literal = l.readIdent()
			tok.Type = token.LookUpIdent(tok.Literal)
			return tok
		} else if isDigit(l.char) {
			firstNum := l.readNumber()
			if l.char == '.' {
				tok.Type = token.FLOAT
				l.readChar()
				firstNum += "." + l.readNumber()
			} else {
				tok.Type = token.INT
			}
			tok.Literal = firstNum

			return tok

		} else {
			tok = newToken(token.ILLEGAL, l.char)
		}
	}
	l.readChar()
	return tok
}

func (l *Lexer) readStr() string {
	pos := l.position + 1
	for {
		l.readChar()
		if l.char == '"' || l.char == 0 {
			break
		}
	}

	return l.input[pos:l.position]
}

func isLetter(char byte) bool {
	return 'a' <= char && char <= 'z' || 'A' <= char && char <= 'Z' || char == '_'
}
func isDigit(char byte) bool {
	return '0' <= char && char <= '9'
}

func newToken(tokenType token.TokenType, char byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(char)}
}

package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENT = "IDENT"
	INT   = "INT"
	FLOAT = "FLOAT"
	STR   = "STR"

	ASSIGN       = "="
	PLUS         = "+"
	MINUS        = "-"
	STAR         = "*"
	SLASH        = "/"
	PLUS_ASSIGN  = "+="
	MINUS_ASSIGN = "-="
	STAR_ASSIGN  = "*="
	SLASH_ASSIGN = "/="

	BANG   = "!"
	EQ     = "=="
	NOT_EQ = "!="
	LT_EQ  = "<="
	GT_EQ  = ">="
	LT     = "<"
	GT     = ">"

	COMMA     = ","
	SEMICOLON = ";"

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	FUNCTION = "FUNCTION"
	LET      = "LET"
	RETURN   = "RETURN"
	IF       = "IF"
	ELSE     = "ELSE"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	FOR      = "FOR"
	WHILE    = "WHILE"
)

var keywords = map[string]TokenType{
	"func":   FUNCTION,
	"var":    LET,
	"return": RETURN,
	"if":     IF,
	"else":   ELSE,
	"true":   TRUE,
	"false":  FALSE,
	"for":    FOR,
	"while":  WHILE,
}

func LookUpIdent(input string) TokenType {
	if tok, ok := keywords[input]; ok {
		return tok
	}
	return IDENT
}

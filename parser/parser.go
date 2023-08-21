package parser

import (
	"fmt"
	"strconv"

	"github.com/Arch-4ng3l/Monkey/ast"
	"github.com/Arch-4ng3l/Monkey/lexer"
	"github.com/Arch-4ng3l/Monkey/token"
)

const (
	_ int = iota
	LOWEST
	RANGE
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
	INDEX
)

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.GT_EQ:    LESSGREATER,
	token.LT_EQ:    LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.STAR:     PRODUCT,
	token.LPAREN:   CALL,
	token.LBRACKET: INDEX,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
	errors    []string

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)

	p.registerPrefix(token.IDENT, p.parseIdent)
	p.registerPrefix(token.INT, p.parseIntLiteral)
	p.registerPrefix(token.FLOAT, p.parseFloatLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFnLiteral)
	p.registerPrefix(token.STR, p.parseStrLiteral)
	p.registerPrefix(token.LBRACKET, p.parseArrLiteral)
	p.registerPrefix(token.FOR, p.parseForLoop)
	p.registerPrefix(token.WHILE, p.parseWhileLoop)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registeInfix(token.PLUS, p.parseInfixExpression)
	p.registeInfix(token.MINUS, p.parseInfixExpression)
	p.registeInfix(token.SLASH, p.parseInfixExpression)
	p.registeInfix(token.STAR, p.parseInfixExpression)
	p.registeInfix(token.EQ, p.parseInfixExpression)
	p.registeInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registeInfix(token.LT, p.parseInfixExpression)
	p.registeInfix(token.GT, p.parseInfixExpression)
	p.registeInfix(token.LT_EQ, p.parseInfixExpression)
	p.registeInfix(token.GT_EQ, p.parseInfixExpression)
	p.registeInfix(token.LPAREN, p.parseCallExpression)

	p.registeInfix(token.PLUS_ASSIGN, p.parseInfixExpression)
	p.registeInfix(token.MINUS_ASSIGN, p.parseInfixExpression)
	p.registeInfix(token.SLASH_ASSIGN, p.parseInfixExpression)
	p.registeInfix(token.STAR_ASSIGN, p.parseInfixExpression)

	p.registeInfix(token.LBRACKET, p.parseIndexExpression)

	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) parseWhileLoop() ast.Expression {

	wl := &ast.WhileLoop{Token: p.curToken}

	if !p.peekTokenIs(token.LPAREN) {
		return nil
	}

	p.nextToken()

	wl.LoopCond = p.parseExpression(LOWEST)

	p.nextToken()

	wl.Body = p.parseBlockStatement()

	p.nextToken()

	return wl

}

func (p *Parser) parseForLoop() ast.Expression {
	fl := &ast.ForLoop{Token: p.curToken}

	if !p.peekTokenIs(token.LPAREN) {
		return nil
	}

	p.nextToken()
	p.nextToken()

	fl.LoopVar = p.parseLetStatement()

	p.nextToken()

	fl.LoopCond = p.parseExpression(LOWEST)

	p.nextToken()
	p.nextToken()

	fl.PostLoop = p.parseExpression(LOWEST)

	p.nextToken()
	p.nextToken()

	fl.Body = p.parseBlockStatement()

	p.nextToken()

	return fl

}

func (p *Parser) parseFloatLiteral() ast.Expression {

	fl := &ast.FloatLiteral{Token: p.curToken}

	n, err := strconv.ParseFloat(p.curToken.Literal, 64)

	if err != nil {
		msg := fmt.Sprintf("Token Literal isnt a Valid Float got %s", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	fl.Value = n

	return fl
}
func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{
		Token: p.curToken,
		Left:  left,
	}

	p.nextToken()

	exp.Index = p.parseExpression(LOWEST)

	if !p.expectedPeek(token.RBRACKET) {
		return nil
	}

	return exp
}

func (p *Parser) parseArrLiteral() ast.Expression {
	arr := &ast.ArrayLiteral{
		Token: p.curToken,
	}

	arr.Elements = p.parseExpressionList(token.RBRACKET)

	return arr
}

func (p *Parser) parseStrLiteral() ast.Expression {
	return &ast.StrLiteral{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{
		Token:    p.curToken,
		Function: function,
	}
	exp.Args = p.parseExpressionList(token.RPAREN)
	return exp
}

func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	list := []ast.Expression{}
	if p.peekTokenIs(end) {
		p.nextToken()
		return nil
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectedPeek(end) {
		return nil
	}

	return list
}

func (p *Parser) parseFnLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{
		Token: p.curToken,
	}

	if !p.expectedPeek(token.LPAREN) {
		return nil
	}
	lit.Params = p.parseFnParams()

	if !p.expectedPeek(token.LBRACE) {
		return nil
	}

	lit.Body = p.parseBlockStatement()

	return lit
}

func (p *Parser) parseFnParams() []*ast.Ident {
	idents := []*ast.Ident{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return idents
	}

	p.nextToken()

	ident := &ast.Ident{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
	idents = append(idents, ident)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		ident := &ast.Ident{
			Token: p.curToken,
			Value: p.curToken.Literal,
		}
		idents = append(idents, ident)
	}

	if !p.expectedPeek(token.RPAREN) {
		return nil
	}
	return idents
}
func (p *Parser) parseIfExpression() ast.Expression {

	expression := &ast.IfExpression{
		Token: p.curToken,
	}
	if !p.expectedPeek(token.LPAREN) {
		return nil
	}
	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectedPeek(token.RPAREN) {
		return nil
	}

	if !p.expectedPeek(token.LBRACE) {
		return nil
	}

	expression.If = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		if !p.expectedPeek(token.LBRACE) {
			return nil
		}

		expression.Else = p.parseBlockStatement()
	}

	return expression
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{
		Token: p.curToken,
	}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if !p.expectedPeek(token.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{
		Token: p.curToken,
		Value: p.curTokenIs(token.TRUE),
	}
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}
	precedences := p.curPrecedence()

	p.nextToken()
	expression.Right = p.parseExpression(precedences)

	return expression
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	exp := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	exp.Right = p.parseExpression(PREFIX)

	return exp
}

func (p *Parser) parseIntLiteral() ast.Expression {
	il := &ast.IntLiteral{Token: p.curToken}

	n, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("Token Literal isnt a Valid Integer got %s", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	il.Value = n

	return il
}

func (p *Parser) parseIdent() ast.Expression {

	var re *ast.ReasignExpression

	ident := &ast.Ident{Token: p.curToken, Value: p.curToken.Literal}

	switch p.peekToken.Type {
	case token.ASSIGN:
		re = &ast.ReasignExpression{
			Token:    p.curToken,
			Operator: "=",
		}
	case token.PLUS_ASSIGN:
		re = &ast.ReasignExpression{
			Token:    p.curToken,
			Operator: "+=",
		}
	case token.MINUS_ASSIGN:
		re = &ast.ReasignExpression{
			Token:    p.curToken,
			Operator: "-=",
		}
	case token.STAR_ASSIGN:
		re = &ast.ReasignExpression{
			Token:    p.curToken,
			Operator: "*=",
		}
	case token.SLASH_ASSIGN:
		re = &ast.ReasignExpression{
			Token:    p.curToken,
			Operator: "/=",
		}

	default:
		return ident
	}
	p.nextToken()
	p.nextToken()
	re.Var = ident
	re.Value = p.parseExpression(LOWEST)
	return re
}

func (p *Parser) registerPrefix(tokenType token.TokenType, f prefixParseFn) {
	p.prefixParseFns[tokenType] = f
}

func (p *Parser) registeInfix(tokenType token.TokenType, f infixParseFn) {
	p.infixParseFns[tokenType] = f
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t

}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectedPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectedPeek(token.IDENT) {
		return nil
	}
	stmt.Name = &ast.Ident{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectedPeek(token.ASSIGN) {
		return nil
	}
	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseExpressionStatement() *ast.ExpresssionStatement {
	stmt := &ast.ExpresssionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]

	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}

	leftExp := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt

}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()

	}
	return program
}

package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jesseduffield/OK/ast"
	"github.com/jesseduffield/OK/lexer"
	"github.com/jesseduffield/OK/token"
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

const (
	_ int = iota
	LOWEST
	ASSIGN
	LAZY            // lazy myFunc()
	ANDOR           // && or ||
	EQUALS          // ==
	LESSGREATER     // > or <
	SUM_AND_PRODUCT // *
	PREFIX          // -X or !X
	NEW             // new Person()
	CALL            // myFunction(X)
	MEMBERACCESS    // myStruct.foo
	INDEX           // array[index]
)

const MAX_IDENTIFIER_LENGTH = 8

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM_AND_PRODUCT,
	token.MINUS:    SUM_AND_PRODUCT,
	token.SLASH:    SUM_AND_PRODUCT,
	token.ASTERISK: SUM_AND_PRODUCT,
	token.LPAREN:   CALL,
	token.LBRACKET: INDEX,
	token.AND:      ANDOR,
	token.OR:       ANDOR,
	token.ASSIGN:   ASSIGN,
	token.PERIOD:   MEMBERACCESS,
	token.LAZY:     LAZY,
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.NULL, p.parseNull)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)
	p.registerPrefix(token.LBRACE, p.parseHashLiteral)
	p.registerPrefix(token.SWITCH, p.parseSwitchExpression)
	p.registerPrefix(token.NEW, p.parseStructInstantiation)
	p.registerPrefix(token.LAZY, p.parseLazyExpression)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.AND, p.parseLogicalInfixExpression)
	p.registerInfix(token.OR, p.parseLogicalInfixExpression)
	p.registerInfix(token.ASSIGN, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)
	p.registerInfix(token.PERIOD, p.parseMemberAccessExpression)

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) parseMemberAccessExpression(left ast.Expression) ast.Expression {
	exp := &ast.StructMemberAccessExpression{Left: left, Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	exp.MemberName = p.curToken.Literal

	return exp
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseLazyExpression() ast.Expression {
	expression := &ast.LazyExpression{
		Token: p.curToken,
	}

	p.nextToken()

	expression.Right = p.parseExpression(LAZY)

	return expression
}

func (p *Parser) parseIdentifier() ast.Expression {
	if !p.validateIdentifier(p.curToken.Literal) {
		return nil
	}

	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) validateIdentifier(identifier string) bool {
	if len(identifier) > MAX_IDENTIFIER_LENGTH {
		p.errors = append(p.errors, "Identifier must be at most eight characters long")
		return false
	}

	if strings.ToLower(identifier) != identifier {
		p.errors = append(p.errors, "Identifier must not contain uppercase characters")
		return false
	}

	if strings.Contains(identifier, "_") {
		p.errors = append(p.errors, "Identifier must not contain underscores")
		return false
	}

	return true
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	p.appendError(
		fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type),
	)
}

func (p *Parser) appendError(msg string) {
	p.errors = append(p.errors, msg)
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
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

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.STRUCT:
		return p.parseStruct()
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
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

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	if !p.validateIdentifier(p.curToken.Literal) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value

	return lit
}

func (p *Parser) parseNull() ast.Expression {
	return &ast.NullLiteral{Token: p.curToken}
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
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

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

// this is for the '&&' and '||' operators
func (p *Parser) parseLogicalInfixExpression(left ast.Expression) ast.Expression {
	exp := p.parseInfixExpression(left)
	castExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		return nil
	}

	for _, operand := range []struct {
		exp  ast.Expression
		side string
	}{{castExp.Left, "left"}, {castExp.Right, "right"}} {
		switch v := operand.exp.(type) {
		case *ast.Identifier:
		case *ast.InfixExpression:
			if v.Operator != "&&" && v.Operator != "||" {
				p.errors = append(p.errors, fmt.Sprintf("%s operand of logical expression must be a variable. Consider storing the %s operand in a variable", operand.side, operand.side))
				return nil
			}
		default:
			p.errors = append(p.errors, fmt.Sprintf("%s operand of logical expression must be a variable. Consider storing the %s operand in a variable", operand.side, operand.side))
			return nil
		}
	}

	return exp
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

func (p *Parser) parseSwitchExpression() ast.Expression {
	expression := &ast.SwitchExpression{Token: p.curToken}

	p.nextToken()

	expression.Subject = p.parseExpression(LOWEST)

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	p.nextToken()

	expression.Cases = p.parseSwitchCases()
	if expression.Cases == nil {
		return nil
	}

	if p.curTokenIs(token.DEFAULT) {
		if !p.expectPeek(token.COLON) {
			return nil
		}

		expression.Default = p.parseSwitchBlockStatement()
		if expression.Default == nil {
			return nil
		}
	}

	return expression
}

func (p *Parser) parseSwitchCases() []ast.SwitchCase {
	cases := []ast.SwitchCase{}

	for p.curTokenIs(token.CASE) {
		switchCase := ast.SwitchCase{}

		p.nextToken()

		switchCase.Value = p.parseExpression(LOWEST)

		if !p.expectPeek(token.COLON) {
			return nil
		}

		switchCase.Block = p.parseSwitchBlockStatement()
		if switchCase.Block == nil {
			return nil
		}

		cases = append(cases, switchCase)
	}

	return cases
}

func (p *Parser) parseSwitchBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	maxAllowedStatements := 1
	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) && !p.curTokenIs(token.DEFAULT) && !p.curTokenIs(token.CASE) {
		if len(block.Statements) >= maxAllowedStatements {
			p.appendError("switch blocks can only contain a single statement. If you want to include multiple statements, use a function call")
			return nil
		}
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
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

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	lit.Body = p.parseBlockStatement()

	return lit
}

func (p *Parser) parseStruct() *ast.Struct {
	str := &ast.Struct{Token: p.curToken}
	str.Methods = map[string]ast.StructMethod{}

	p.nextToken()

	str.Name = p.curToken.Literal

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	if p.peekTokenIs(token.PACK) {
		p.nextToken()
		p.nextToken()

		str.PrivacyAcknowledgement = p.parseStringLiteral().String()
	}

	for p.peekTokenIs(token.FIELD) {
		p.nextToken()
		fieldName := p.peekToken.Literal
		if !p.validateIdentifier(fieldName) {
			return nil
		}
		// no public struct fields for now
		str.Fields = append(str.Fields, ast.StructField{Name: fieldName, Public: false})
		p.nextToken()
	}

	for !p.peekTokenIs(token.RBRACE) {
		isPublic := false
		if p.peekTokenIs(token.PUBLIC) {
			p.nextToken()
			isPublic = true
			if p.peekTokenIs(token.FIELD) {
				p.appendError("public nac fields are not permitted")
				return nil
			}
		}

		methodName := p.peekToken.Literal
		if !p.validateIdentifier(methodName) {
			return nil
		}
		p.nextToken()
		p.nextToken()

		fn := p.parseFunctionLiteral().(*ast.FunctionLiteral)
		str.Methods[methodName] = ast.StructMethod{Public: isPublic, FunctionLiteral: fn}
	}

	p.nextToken()

	return str
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return identifiers
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseExpressionList(token.RPAREN)
	return exp
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.curToken}

	array.Elements = p.parseExpressionList(token.RBRACKET)

	return array
}

func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	list := []ast.Expression{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return list
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Token: p.curToken, Left: left}

	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	return exp
}

func (p *Parser) parseHashLiteral() ast.Expression {
	hash := &ast.HashLiteral{Token: p.curToken}
	hash.Pairs = make(map[ast.Expression]ast.Expression)

	for !p.peekTokenIs(token.RBRACE) {
		p.nextToken()
		key := p.parseExpression(LOWEST)

		if !p.expectPeek(token.COLON) {
			return nil
		}

		p.nextToken()
		value := p.parseExpression(LOWEST)

		hash.Pairs[key] = value

		if !p.peekTokenIs(token.RBRACE) && !p.expectPeek(token.COMMA) {
			return nil
		}
	}

	if !p.expectPeek(token.RBRACE) {
		return nil
	}

	return hash
}

func (p *Parser) parseStructInstantiation() ast.Expression {
	// typical line:
	// new Person(arg1, arg2)

	exp := &ast.StructInstantiation{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	exp.StructName = p.curToken.Literal
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	exp.Arguments = p.parseExpressionList(token.RPAREN)

	return exp
}

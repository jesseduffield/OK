package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jesseduffield/OK/ok/ast"
	"github.com/jesseduffield/OK/ok/lexer"
	"github.com/jesseduffield/OK/ok/token"
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

const MAX_IDENTIFIER_LENGTH = 8

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
	p.registerInfix(token.GTEQ, p.parseInfixExpression)
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
	p.validateIdentifier(p.curToken.Literal)

	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) validateIdentifier(identifier string) {
	if len(identifier) > MAX_IDENTIFIER_LENGTH {
		suggested := shortenedIdentifier(identifier)

		p.appendError(fmt.Sprintf(
			"Identifier must be at most eight characters long; consider using '%s' instead.\nSee https://github.com/jesseduffield/ok#familiarity-admits-brevity",
			suggested,
		))
	}

	if strings.ToLower(identifier) != identifier {
		p.appendError(fmt.Sprintf(
			"Identifier must not contain uppercase characters; consider using '%s' instead.\nSee https://github.com/jesseduffield/ok#familiarity-admits-brevity",
			strings.ToLower(identifier),
		))
	}

	if strings.Contains(identifier, "_") {
		p.appendError(fmt.Sprintf(
			"Identifier must not contain underscores; consider using '%s' instead.\nSee https://github.com/jesseduffield/ok#familiarity-admits-brevity",
			removeUnderscores(identifier),
		))
	}
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
	locatedMsg := fmt.Sprintf("%s (%s): %s", p.curToken.Location(), p.curToken.Literal, msg)
	p.errors = append(p.errors, locatedMsg)
}

func (p *Parser) appendErrorForExpression(msg string, exp ast.Expression) {
	locatedMsg := fmt.Sprintf("%s (%s): %s", exp.GetToken().Location(), exp.String(), msg)
	p.errors = append(p.errors, locatedMsg)
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
	case token.COMMENT:
		return p.parseCommentStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekSemiColon() {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken)
		return nil
	}
	leftExp := prefix()

	for !p.peekSemiColon() && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) peekSemiColon() bool {
	return p.peekTokenIs(token.SEMICOLON)
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	p.validateIdentifier(p.curToken.Literal)
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekSemiColon() {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.peekSemiColon() {
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
		p.appendError(fmt.Sprintf("'%s' is not a valid integer", p.curToken.Literal))
		return nil
	}

	lit.Value = value

	return lit
}

func (p *Parser) parseNull() ast.Expression {
	return &ast.NullLiteral{Token: p.curToken}
}

func (p *Parser) handleUnexpectedToken(t token.Token) {
	switch t.Literal {
	case ">", "<", "<=", "==", "!=":
		p.appendError(
			fmt.Sprintf(
				"Unexpected token '%s'. There is only one comparison operator: '>='.\nSee https://github.com/jesseduffield/ok#one-comparison-operator",
				t.Literal,
			),
		)
	default:
		p.appendError(fmt.Sprintf("Unexpected token '%s'", t.Literal))
	}
}

func (p *Parser) noPrefixParseFnError(t token.Token) {
	p.handleUnexpectedToken(t)
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

func (p *Parser) parseCommentStatement() *ast.CommentStatement {
	// TODO: handle when the two slashes aren't followed by a space
	text := strings.TrimPrefix(p.curToken.Literal, "// ")

	stmt := &ast.CommentStatement{
		Token: p.curToken,
		Text:  text,
	}

	return stmt
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
	}{{castExp.Left, "Left"}, {castExp.Right, "Right"}} {
		switch v := operand.exp.(type) {
		case *ast.Identifier:
		case *ast.InfixExpression:
			if v.Operator != "&&" && v.Operator != "||" {
				p.appendErrorForExpression(fmt.Sprintf("%s operand of logical expression must be a variable. Consider storing '%s' in a variable", operand.side, v.String()), v)
				return nil
			}
		default:
			p.appendErrorForExpression(fmt.Sprintf("%s operand of logical expression must be a variable. Consider storing '%s' in a variable", operand.side, v.String()), v)
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
	statementCount := 0
	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) && !p.curTokenIs(token.DEFAULT) && !p.curTokenIs(token.CASE) {
		if statementCount >= maxAllowedStatements {
			p.appendError(
				"switch blocks can only contain a single statement. If you want to include multiple statements, use a function call\nSee https://github.com/jesseduffield/ok#readable-switches",
			)
			return nil
		}
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
			if _, ok := stmt.(*ast.CommentStatement); !ok {
				statementCount += 1
			}
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
		p.nextToken()
		fieldName := p.curToken.Literal
		p.validateIdentifier(fieldName)
		// no public struct fields for now
		str.Fields = append(str.Fields, ast.StructField{Name: fieldName, Public: false})
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

		p.nextToken()
		methodName := p.curToken.Literal
		p.validateIdentifier(methodName)
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

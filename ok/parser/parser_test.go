package parser

import (
	"fmt"
	"strings"
	"testing"

	"github.com/jesseduffield/OK/ok/ast"
	"github.com/jesseduffield/OK/ok/lexer"
)

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d",
				len(program.Statements))
		}

		stmt := program.Statements[0]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}

		val := stmt.(*ast.LetStatement).Value
		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 5;", 5},
		{"return true;", true},
		{"return foobar;", "foobar"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d",
				len(program.Statements))
		}

		stmt := program.Statements[0]
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Fatalf("stmt not *ast.ReturnStatement. got=%T", stmt)
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Fatalf("returnStmt.TokenLiteral not 'return', got %q",
				returnStmt.TokenLiteral())
		}
		if testLiteralExpression(t, returnStmt.ReturnValue, tt.expectedValue) {
			return
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expression)
	}
	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %s. got=%s", "foobar", ident.Value)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "foobar",
			ident.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp not *ast.IntegerLiteral. got=%T", stmt.Expression)
	}
	if literal.Value != 5 {
		t.Errorf("literal.Value not %d. got=%d", 5, literal.Value)
	}
	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral not %s. got=%s", "5",
			literal.TokenLiteral())
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!foobar;", "!", "foobar"},
		{"-foobar;", "-", "foobar"},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s",
				tt.operator, exp.Operator)
		}
		if !testLiteralExpression(t, exp.Right, tt.value) {
			return
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 >= 5;", 5, ">=", 5},
		{"foobar + barfoo;", "foobar", "+", "barfoo"},
		{"foobar - barfoo;", "foobar", "-", "barfoo"},
		{"foobar * barfoo;", "foobar", "*", "barfoo"},
		{"foobar / barfoo;", "foobar", "/", "barfoo"},
		{"foobar >= barfoo;", "foobar", ">=", "barfoo"},
		{"true >= true", true, ">=", true},
		{"a && b", "a", "&&", "b"},
		{"a || b", "a", "||", "b"},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		if !testInfixExpression(t, stmt.Expression, tt.leftValue,
			tt.operator, tt.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"((a + b) / c)",
		},
		{
			"a + b * c + d / e - f",
			"(((((a + b) * c) + d) / e) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"(5 + 5) * 2 * (5 + 5)",
			"(((5 + 5) * 2) * (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((((a + b) + c) * d) / f) + g))",
		},
		{
			"a * [1, 2, 3, 4][b * c] * d",
			"((a * ([1, 2, 3, 4][(b * c)])) * d)",
		},
		{
			"add(a * b[2], b[1], 2 * [1, 2][1])",
			"add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))",
		},
		{
			"x[1][2][3]",
			"(((x[1])[2])[3])",
		},
		{
			"nil",
			"nil",
		},
		{
			"a && b && c",
			"((a && b) && c)",
		},
		{
			"a && b || c",
			"((a && b) || c)",
		},
		{
			"a || b || c",
			"((a || b) || c)",
		},
		{
			"x = 5 + 5",
			"(x = (5 + 5))",
		},
		{
			"x = lazy 5 + 5",
			"(x = lazy((5 + 5)))",
		},
		{
			"x = x + 1",
			"(x = (x + 1))",
		},
		{
			"fn() { r = r + 1 }();",
			"fn() { (r = (r + 1)) }()",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input           string
		expectedBoolean bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program has not enough statements. got=%d",
				len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		boolean, ok := stmt.Expression.(*ast.Boolean)
		if !ok {
			t.Fatalf("exp not *ast.Boolean. got=%T", stmt.Expression)
		}
		if boolean.Value != tt.expectedBoolean {
			t.Errorf("boolean.Value not %t. got=%t", tt.expectedBoolean,
				boolean.Value)
		}
	}
}

func TestIfExpression(t *testing.T) {
	input := `if (x >= y) { x }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T",
			stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", ">=", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n",
			len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if exp.Alternative != nil {
		t.Errorf("exp.Alternative.Statements was not nil. got=%+v", exp.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x >= y) { x } else { y }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", ">=", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n",
			len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if len(exp.Alternative.Statements) != 1 {
		t.Errorf("exp.Alternative.Statements does not contain 1 statements. got=%d\n",
			len(exp.Alternative.Statements))
	}

	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.Alternative.Statements[0])
	}

	if !testIdentifier(t, alternative.Expression, "y") {
		return
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x, y) { x + y; }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T",
			stmt.Expression)
	}

	if len(function.Parameters) != 2 {
		t.Fatalf("function literal parameters wrong. want 2, got=%d\n",
			len(function.Parameters))
	}

	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")

	if len(function.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements has not 1 statements. got=%d\n",
			len(function.Body.Statements))
	}

	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function body stmt is not ast.ExpressionStatement. got=%T",
			function.Body.Statements[0])
	}

	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "fn() {};", expectedParams: []string{}},
		{input: "fn(x) {};", expectedParams: []string{"x"}},
		{input: "fn(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		function := stmt.Expression.(*ast.FunctionLiteral)

		if len(function.Parameters) != len(tt.expectedParams) {
			t.Errorf("length parameters wrong. want %d, got=%d\n",
				len(tt.expectedParams), len(function.Parameters))
		}

		for i, ident := range tt.expectedParams {
			testLiteralExpression(t, function.Parameters[i], ident)
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T",
			stmt.Expression)
	}

	if !testIdentifier(t, exp.Function, "add") {
		return
	}

	if len(exp.Arguments) != 3 {
		t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
	}

	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}

func TestCallExpressionParameterParsing(t *testing.T) {
	tests := []struct {
		input         string
		expectedIdent string
		expectedArgs  []string
	}{
		{
			input:         "add();",
			expectedIdent: "add",
			expectedArgs:  []string{},
		},
		{
			input:         "add(1);",
			expectedIdent: "add",
			expectedArgs:  []string{"1"},
		},
		{
			input:         "add(1, 2 * 3, 4 + 5);",
			expectedIdent: "add",
			expectedArgs:  []string{"1", "(2 * 3)", "(4 + 5)"},
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		exp, ok := stmt.Expression.(*ast.CallExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T",
				stmt.Expression)
		}

		if !testIdentifier(t, exp.Function, tt.expectedIdent) {
			return
		}

		if len(exp.Arguments) != len(tt.expectedArgs) {
			t.Fatalf("wrong number of arguments. want=%d, got=%d",
				len(tt.expectedArgs), len(exp.Arguments))
		}

		for i, arg := range tt.expectedArgs {
			if exp.Arguments[i].String() != arg {
				t.Errorf("argument %d wrong. want=%q, got=%q", i,
					arg, exp.Arguments[i].String())
			}
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral() not '%s'. got=%s",
			name, letStmt.Name.TokenLiteral())
		return false
	}

	return true
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{},
	operator string, right interface{}) bool {

	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.InfixExpression. got=%T(%s)", exp, exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
		return false
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true
}

func testLiteralExpression(
	t *testing.T,
	exp ast.Expression,
	expected interface{},
) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d. got=%s", value,
			integ.TokenLiteral())
		return false
	}

	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got=%s", value,
			ident.TokenLiteral())
		return false
	}

	return true
}

func testComment(t *testing.T, exp ast.Statement, value string) bool {
	comment, ok := exp.(*ast.CommentStatement)
	if !ok {
		t.Errorf("exp not *ast.CommentStatement. got=%T", exp)
		return false
	}

	if comment.Text != value {
		t.Errorf("comment text not %s. got=%s", value, comment.Text)
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bo, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp not *ast.Boolean. got=%T", exp)
		return false
	}

	if bo.Value != value {
		t.Errorf("bo.Value not %t. got=%t", value, bo.Value)
		return false
	}

	if bo.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("bo.TokenLiteral not %t. got=%s",
			value, bo.TokenLiteral())
		return false
	}

	return true
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"hello world";`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	literal, ok := stmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("exp not *ast.StringLiteral. got=%T", stmt.Expression)
	}

	if literal.Value != "hello world" {
		t.Errorf("literal.Value not %q. got=%q", "hello world", literal.Value)
	}
}

func TestParsingArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	array, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("exp not ast.ArrayLiteral. got=%T", stmt.Expression)
	}

	if len(array.Elements) != 3 {
		t.Fatalf("len(array.Elements) not 3. got=%d", len(array.Elements))
	}

	testIntegerLiteral(t, array.Elements[0], 1)
	testInfixExpression(t, array.Elements[1], 2, "*", 2)
	testInfixExpression(t, array.Elements[2], 3, "+", 3)
}

func TestParsingIndexExpressions(t *testing.T) {
	input := "myarray[1 + 1]"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	indexExp, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("exp not *ast.IndexExpression. got=%T", stmt.Expression)
	}

	if !testIdentifier(t, indexExp.Left, "myarray") {
		return
	}

	if !testInfixExpression(t, indexExp.Index, 1, "+", 1) {
		return
	}
}

func TestParsingHashLiteralsStringKeys(t *testing.T) {
	input := `{"one": 1, "two": 2, "three": 3}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expression)
	}

	if len(hash.Pairs) != 3 {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}

	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral. got=%T", key)
		}

		expectedValue := expected[literal.String()]

		testIntegerLiteral(t, value, expectedValue)
	}
}

func TestParsingEmptyHashLiteral(t *testing.T) {
	input := "{}"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expression)
	}

	if len(hash.Pairs) != 0 {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}
}

func TestParsingHashLiteralsWithExpressions(t *testing.T) {
	input := `{"one": 0 + 1, "two": 10 - 8, "three": 15 / 5}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expression)
	}

	if len(hash.Pairs) != 3 {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}

	tests := map[string]func(ast.Expression){
		"one": func(e ast.Expression) {
			testInfixExpression(t, e, 0, "+", 1)
		},
		"two": func(e ast.Expression) {
			testInfixExpression(t, e, 10, "-", 8)
		},
		"three": func(e ast.Expression) {
			testInfixExpression(t, e, 15, "/", 5)
		},
	}

	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral. got=%T", key)
			continue
		}

		testFunc, ok := tests[literal.String()]
		if !ok {
			t.Errorf("No test function for key %q found", literal.String())
			continue
		}

		testFunc(value)
	}
}

func TestParsingSwitchStatement(t *testing.T) {
	input := `switch x >= y { case 1 + 5: x; case true: x; default: 9 }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.SwitchExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.SwitchExpression. got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, exp.Subject, "x", ">=", "y") {
		return
	}

	if len(exp.Cases) != 2 {
		t.Errorf("exp.Cases does not contain %d cases. got=%d\n",
			2, len(exp.Cases))
	}

	if !testInfixExpression(t, exp.Cases[0].Value, 1, "+", 5) {
		return
	}

	statement := exp.Cases[0].Block.Statements[0].(*ast.ExpressionStatement)
	if !testIdentifier(t, statement.Expression, "x") {
		return
	}

	if !testBooleanLiteral(t, exp.Cases[1].Value, true) {
		return
	}

	statement = exp.Cases[1].Block.Statements[0].(*ast.ExpressionStatement)
	if !testIdentifier(t, statement.Expression, "x") {
		return
	}

	statement = exp.Default.Statements[0].(*ast.ExpressionStatement)
	if !testIntegerLiteral(t, statement.Expression, 9) {
		return
	}
}

func TestParsingSwitchStatementWithComment(t *testing.T) {
	input := `
		switch x >= y {
			case 1 + 5:
				// comment
				x;
			case true:
				x;
			default:
				9
		}
	`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.SwitchExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.SwitchExpression. got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, exp.Subject, "x", ">=", "y") {
		return
	}

	if len(exp.Cases) != 2 {
		t.Errorf("exp.Cases does not contain %d cases. got=%d\n",
			2, len(exp.Cases))
	}

	if !testInfixExpression(t, exp.Cases[0].Value, 1, "+", 5) {
		return
	}

	statement := exp.Cases[0].Block.Statements[0]
	if !testComment(t, statement, "comment") {
		return
	}

	expStatement := exp.Cases[0].Block.Statements[1].(*ast.ExpressionStatement)
	if !testIdentifier(t, expStatement.Expression, "x") {
		return
	}

	if !testBooleanLiteral(t, exp.Cases[1].Value, true) {
		return
	}

	expStatement = exp.Cases[1].Block.Statements[0].(*ast.ExpressionStatement)
	if !testIdentifier(t, expStatement.Expression, "x") {
		return
	}

	expStatement = exp.Default.Statements[0].(*ast.ExpressionStatement)
	if !testIntegerLiteral(t, expStatement.Expression, 9) {
		return
	}
}

func TestParsingStructDefinition(t *testing.T) {
	input := `notaclass person { pack "test" field name field email public foo fn(selfish, a, b) { return 5 } bar fn(selfish) { return 3 } } notaclass other { field blah }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			2, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.Struct)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.Struct. got=%T",
			program.Statements[0])
	}

	str := stmt.String()
	// our methods are in a hashmap and we're not sorting by anything.
	// I'm too lazy to fix right now
	expected := `notaclass person {
	pack "test"

	field name
	field email

	public foo fn(selfish, a, b) { return 5; }
	bar fn(selfish) { return 3; }
}`
	alternativeExpected := `notaclass person {
	pack "test"

	field name
	field email

	bar fn(selfish) { return 3; }
	public foo fn(selfish, a, b) { return 5; }
}`
	if str != expected && str != alternativeExpected {
		t.Fatalf("unexpected struct got=\n%s\nexpected=\n%s\n", str, expected)
	}

	stmt, ok = program.Statements[1].(*ast.Struct)
	if !ok {
		t.Fatalf("program.Statements[1] is not ast.Struct. got=%T",
			program.Statements[1])
	}

	str = stmt.String()
	expected = `notaclass other {
	field blah
}`
	if str != expected {
		t.Fatalf("unexpected statement got=\n%s\nexpected=\n%s\n", str, expected)
	}
}

func TestParsingInvalidStructDefinition(t *testing.T) {
	input := `notaclass person { public field name }`
	expectedError := "line 1, column 20 (public): public nac fields are not permitted"

	l := lexer.New(input)
	p := New(l)
	p.ParseProgram()
	for _, err := range p.errors {
		if err == expectedError {
			return
		}
	}
	t.Fatalf("expected error: %s\nGot: %s", expectedError, strings.Join(p.errors, "\n"))
}

func TestParsingStructInstantiation(t *testing.T) {
	input := `let x = new person(a, b);`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	expectStatements(t, program.Statements, []string{
		`let x = new person(a, b);`,
	})
}

func TestParsingStructMemberAccess(t *testing.T) {
	input := `x.foo; x.bar(a,b);`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	expectStatements(t, program.Statements, []string{
		`x.foo`,
		`x.bar(a, b)`,
	})
}

func TestParsingInvalidSwitch(t *testing.T) {
	input := `switch x { case true: x; y; default: x; }`

	l := lexer.New(input)
	p := New(l)
	p.ParseProgram()
	expectedError := "line 1, column 26 (y): switch blocks can only contain a single statement. If you want to include multiple statements, use a function call\nSee https://github.com/jesseduffield/ok#readable-switches"
	for _, err := range p.errors {
		if err == expectedError {
			return
		}
	}
	t.Fatalf("expected error:\n%s\nActual errors:\n%s", expectedError, strings.Join(p.errors, "\n"))
}

func TestParsingInvalidExpressions(t *testing.T) {
	tests := []struct {
		input         string
		expectedError string
	}{
		{
			input: "a && b()",
			// TODO: use the column of the start of b() not the end.
			expectedError: "line 1, column 7 (b()): Right operand of logical expression must be a variable. Consider storing 'b()' in a variable",
		},
		{
			input:         "a() && b",
			expectedError: "line 1, column 2 (a()): Left operand of logical expression must be a variable. Consider storing 'a()' in a variable",
		},
		{
			input:         "a && true",
			expectedError: "line 1, column 6 (true): Right operand of logical expression must be a variable. Consider storing 'true' in a variable",
		},
		{
			input:         "REALLY_LONG_VARIABLE_NAME",
			expectedError: "line 1, column 1 (REALLY_LONG_VARIABLE_NAME): Identifier must be at most eight characters long; consider using 'rlvn' instead.\nSee https://github.com/jesseduffield/ok#familiarity-admits-brevity",
		},
		{
			input:         "a_b",
			expectedError: "line 1, column 1 (a_b): Identifier must not contain underscores; consider using 'ab' instead.\nSee https://github.com/jesseduffield/ok#familiarity-admits-brevity",
		},
		{
			input:         "abC",
			expectedError: "line 1, column 1 (abC): Identifier must not contain uppercase characters; consider using 'abc' instead.\nSee https://github.com/jesseduffield/ok#familiarity-admits-brevity",
		},
		{
			input:         "let really_long_variable_name = 5",
			expectedError: "line 1, column 5 (really_long_variable_name): Identifier must be at most eight characters long; consider using 'rlvn' instead.\nSee https://github.com/jesseduffield/ok#familiarity-admits-brevity",
		},
		{
			input:         "notaclass me { field really_long_variable_name }",
			expectedError: "line 1, column 22 (really_long_variable_name): Identifier must be at most eight characters long; consider using 'rlvn' instead.\nSee https://github.com/jesseduffield/ok#familiarity-admits-brevity",
		},
		{
			input:         "notaclass me { really_long_variable_name fn() { return 5 } }",
			expectedError: "line 1, column 16 (really_long_variable_name): Identifier must be at most eight characters long; consider using 'rlvn' instead.\nSee https://github.com/jesseduffield/ok#familiarity-admits-brevity",
		},
		{
			input:         "a < b",
			expectedError: "line 1, column 3 (<): Unexpected token '<'. There is only one comparison operator: '>='.\nSee https://github.com/jesseduffield/ok#one-comparison-operator",
		},
		{
			input:         "a ** b",
			expectedError: "line 1, column 4 (*): Unexpected token '*'",
		},
		{
			input:         "a * b\nb ** c",
			expectedError: "line 2, column 4 (*): Unexpected token '*'",
		},
		{
			input:         "999999999999999999999999999",
			expectedError: "line 1, column 1 (999999999999999999999999999): '999999999999999999999999999' is not a valid integer",
		},
	}

outer:
	for _, test := range tests {
		l := lexer.New(test.input)
		p := New(l)
		p.ParseProgram()
		for _, err := range p.errors {
			if err == test.expectedError {
				continue outer
			}
		}
		t.Fatalf("expected error:\n%s\nActual errors:\n%s", test.expectedError, strings.Join(p.errors, "\n"))
	}
}

func TestParsingLazyExpression(t *testing.T) {
	input := `let x = lazy 3 >= 4`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	expectStatements(t, program.Statements, []string{
		`let x = lazy((3 >= 4));`,
	})
}

func TestParsingCommentExpression(t *testing.T) {
	input := `let x = 3; // comment 1
	let y = 4;
	// comment 2
	let z = 4;`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	expectStatements(t, program.Statements, []string{
		`let x = 3;`,
		`// comment 1`,
		`let y = 4;`,
		`// comment 2`,
		`let z = 4;`,
	})
}

func expectStatements(t *testing.T, statements []ast.Statement, expected []string) {
	statementStrings := []string{}
	for _, statement := range statements {
		statementStrings = append(statementStrings, statement.String())
	}

	if len(statements) != len(expected) {
		t.Fatalf("program.Statements does not contain %d statements. got=%d:\n%s",
			len(expected), len(statements), strings.Join(statementStrings, "\n"))
	}

	for i, str := range expected {
		if statements[i].String() != str {
			t.Fatalf("unexpected statement got=\n%s\nexpected=\n%s\n",
				statements[i].String(), expected[i])
		}
	}
}

func TestShortenedIdentifier(t *testing.T) {
	scenarios := []struct {
		input    string
		expected string
	}{
		// already short enough
		{"aa", "aa"},
		// already short enough
		{"abcdefgh", "abcdefgh"},
		// removing underscores is sufficient to satisfy max length
		{"a_b_c_d_e_f_g", "abcdefg"},
		// removing underscores is not sufficient to satisfy max length,
		// so abbreivation is used
		{"really_long_variable_name", "rlvn"},
		// same with camelCase
		{"reallyLongVariableName", "rlvn"},
		// not abbreviating here; falling back to more generic approach
		{"reallylongvariablename", "rlynvbnm"},
		// removes vowels after first letter
		{"abcdfghi", "abcdfghi"},
		// starts removing consonants after all the vowels are gone
		{"aabaacaadaafaagaahaajaak", "abdfghjk"},
		// truncates if still too long after removing consonants
		{"alskdfhljkahsdfaoipequwaksjdhjfklajreopiwqhjkaf", "alkflksf"},
		// more examples
		{"realvalues", "rlvalues"},
		{"smallValue1", "sv1"},
		{"longishValue", "lngshvle"},
		{"another_one", "anthrone"},
		{"anotherOne", "anthrone"},
		{"oneTwoThree", "ott"},
	}

	for _, scenario := range scenarios {
		input, expected := scenario.input, scenario.expected
		output := shortenedIdentifier(input)
		if output != expected {
			t.Fatalf("shortenedIdentifier(%s) = %s, expected %s", input, output, expected)
		}
	}
}

package lexer

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/jesseduffield/OK/ok/token"
)

func TestNextToken(t *testing.T) {
	input := `let five = 5;
let ten = 10;

let add = fn(x, y) {
  x + y;
};

let result = add(five, ten);
!-/*5;
5 >= 10 >= 5;

if (5 >= 10) {
	return true;
} else {
	return false;
}

10 >= 10;
"foobar"
"foo bar"
[1, 2];
{"foo": "bar"}
switch true {
	case true:
		return true;
	case false:
		return false;
	default:
		return false;
}
foo && bar
foo || bar
10 != 12
lazy 3 >= 4
NO!
test // this is my comment
testb
// this is my other comment
arr1
a!a
a?a
?a
!a
<=
<
>
==
!=
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTERISK, "*"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.INT, "5"},
		{token.GTEQ, ">="},
		{token.INT, "10"},
		{token.GTEQ, ">="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.INT, "5"},
		{token.GTEQ, ">="},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.INT, "10"},
		{token.GTEQ, ">="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.STRING, "foobar"},
		{token.STRING, "foo bar"},
		{token.LBRACKET, "["},
		{token.INT, "1"},
		{token.COMMA, ","},
		{token.INT, "2"},
		{token.RBRACKET, "]"},
		{token.SEMICOLON, ";"},
		{token.LBRACE, "{"},
		{token.STRING, "foo"},
		{token.COLON, ":"},
		{token.STRING, "bar"},
		{token.RBRACE, "}"},
		{token.SWITCH, "switch"},
		{token.TRUE, "true"},
		{token.LBRACE, "{"},
		{token.CASE, "case"},
		{token.TRUE, "true"},
		{token.COLON, ":"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.CASE, "case"},
		{token.FALSE, "false"},
		{token.COLON, ":"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.DEFAULT, "default"},
		{token.COLON, ":"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.IDENT, "foo"},
		{token.AND, "&&"},
		{token.IDENT, "bar"},
		{token.IDENT, "foo"},
		{token.OR, "||"},
		{token.IDENT, "bar"},
		{token.INT, "10"},
		{token.ILLEGAL, "!="},
		{token.INT, "12"},
		{token.LAZY, "lazy"},
		{token.INT, "3"},
		{token.GTEQ, ">="},
		{token.INT, "4"},
		{token.NULL, "NO!"},
		{token.IDENT, "test"},
		{token.COMMENT, "// this is my comment"},
		{token.IDENT, "testb"},
		{token.COMMENT, "// this is my other comment"},
		{token.IDENT, "arr1"},
		{token.IDENT, "a!a"},
		{token.IDENT, "a?a"},
		{token.ILLEGAL, "?"},
		{token.IDENT, "a"},
		{token.BANG, "!"},
		{token.IDENT, "a"},
		{token.ILLEGAL, "<="},
		{token.ILLEGAL, "<"},
		{token.ILLEGAL, ">"},
		{token.ILLEGAL, "=="},
		{token.ILLEGAL, "!="},
		{token.EOF, ""},
	}

	l := New(input)

	tokens := []token.Token{}

	for i, tt := range tests {
		tok := l.NextToken()
		tokens = append(tokens, tok)

		if tok.Type != tt.expectedType {
			t.Log(spew.Sdump(tokens))
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Log(spew.Sdump(tokens))
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestLocationMarking(t *testing.T) {
	input := "let five = 5;\nfive >= 4;\n\"aa\" >= \"b\""

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
		expectedLine    int
		expectedColumn  int
	}{
		// this is stored as line zero but displayed as line 1 in the Location() method
		{token.LET, "let", 0, 1},
		{token.IDENT, "five", 0, 5},
		{token.ASSIGN, "=", 0, 10},
		{token.INT, "5", 0, 12},
		{token.SEMICOLON, ";", 0, 13},
		{token.IDENT, "five", 1, 1},
		{token.GTEQ, ">=", 1, 6},
		{token.INT, "4", 1, 9},
		{token.SEMICOLON, ";", 1, 10},
		{token.STRING, "aa", 2, 1},
		{token.GTEQ, ">=", 2, 6},
		{token.STRING, "b", 2, 9},
		{token.EOF, "", 2, 12},
	}

	l := New(input)

	tokens := []token.Token{}

	for i, tt := range tests {
		tok := l.NextToken()
		tokens = append(tokens, tok)

		if tok.Type != tt.expectedType {
			t.Log(spew.Sdump(tokens))
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Log(spew.Sdump(tokens))
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}

		if tok.Line != tt.expectedLine {
			t.Log(spew.Sdump(tokens))
			t.Fatalf("tests[%d] - line wrong. expected=%d, got=%d",
				i, tt.expectedLine, tok.Line)
		}

		if tok.Column != tt.expectedColumn {
			t.Log(spew.Sdump(tokens))
			t.Fatalf("tests[%d] - column wrong. expected=%d, got=%d",
				i, tt.expectedColumn, tok.Column)
		}
	}
}

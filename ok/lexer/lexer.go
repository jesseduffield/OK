// lexer/lexer.go
package lexer

import (
	"github.com/jesseduffield/OK/ok/token"
)

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination

	line   int
	column int
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = eofByte()
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
	l.column++
}

type mapNode struct {
	key       byte
	tokenType token.TokenType
	mapping   map[byte]mapNode
}

// The order matters here: if you have a token of two characters, you need to preceed
// it with a token of just its first character, even if that's just an illegal token.
// We could improve this algorithm but it's good enough for now.
var mapping = []struct{ key, tokenType token.TokenType }{
	{"!", token.BANG},
	{"!=", token.ILLEGAL},
	{"=", token.ASSIGN},
	{"==", token.ILLEGAL},
	{">", token.ILLEGAL},
	{">=", token.GTEQ},
	{"<", token.ILLEGAL},
	{"<=", token.ILLEGAL},
	{"+", token.PLUS},
	{"-", token.MINUS},
	{"&", token.ILLEGAL},
	{"&&", token.AND},
	{"|", token.ILLEGAL},
	{"||", token.OR},
	{"*", token.ASTERISK},
	{";", token.SEMICOLON},
	{",", token.COMMA},
	{"(", token.LPAREN},
	{")", token.RPAREN},
	{"{", token.LBRACE},
	{"}", token.RBRACE},
	{"[", token.LBRACKET},
	{"]", token.RBRACKET},
	{":", token.COLON},
	{".", token.PERIOD},
}

var tokenTree = generateTokenTree()

func generateTokenTree() map[byte]mapNode {
	result := map[byte]mapNode{}

	for _, node := range mapping {
		if len(node.key) > 2 {
			panic("only known tokens of length 1 and 2 are supported")
		}
		inner := result
		for _, ch := range []byte(node.key) {
			if _, ok := inner[ch]; !ok {
				inner[ch] = mapNode{key: ch, tokenType: node.tokenType, mapping: map[byte]mapNode{}}
			}
			inner = inner[ch].mapping
		}
	}

	return result
}

func (l *Lexer) ReadKnownToken(line int, column int) (token.Token, bool) {
	node, ok := tokenTree[l.ch]
	if !ok {
		return token.Token{}, false
	}

	// no registered tokens longer than the given token so we can return immediately
	if node.mapping == nil {
		return l.newToken(node.tokenType, l.ch), true
	}

	// if we're here then maybe we've got '=' but we want to see if the token is
	// actually '=='.
	char := l.ch
	nextChar := l.peekChar()
	if next, ok := node.mapping[nextChar]; ok {
		l.readChar()
		return l.newStringToken(next.tokenType, string(char)+string(nextChar), line, column), true
	}

	return l.newToken(node.tokenType, l.ch), true
}

func (l *Lexer) NextToken() token.Token {
	l.skipWhitespace()

	var tok token.Token
	tok.Column = l.column
	tok.Line = l.line

	switch l.ch {
	case '/':
		if l.peekChar() == '/' {
			tok.Literal = l.readComment()
			tok.Type = token.COMMENT
			return tok
		}

		tok = l.newToken(token.SLASH, l.ch)
	case eofByte():
		tok.Literal = ""
		tok.Type = token.EOF
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	default:
		if tok, ok := l.ReadKnownToken(l.line, l.column); ok {
			l.readChar()
			return tok
		}

		if isValidIdentifierStartChar(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		}

		if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		}

		tok = l.newToken(token.ILLEGAL, l.ch)
	}

	l.readChar()
	return tok
}

func (l *Lexer) newToken(tokenType token.TokenType, ch byte) token.Token {
	return l.newStringToken(tokenType, string(ch), l.line, l.column)
}

func (l *Lexer) newStringToken(
	tokenType token.TokenType,
	literal string,
	line int,
	column int,
) token.Token {
	return token.Token{Type: tokenType, Literal: literal, Line: line, Column: column}
}

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == eofByte() {
			break
		}
	}
	return l.input[position:l.position]
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	// the bang is here for the sake of the 'NO!' token,
	// and the ? is here for the sake of the 'ayok?' builtin function
	for isValidIdentifierStartChar(l.ch) || (l.ch > '0' && l.ch < '9') || l.ch == '!' || l.ch == '?' {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isValidIdentifierStartChar(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func eofByte() byte {
	return 0
}

func (l *Lexer) readComment() string {
	position := l.position
	for l.ch != '\n' && l.ch != eofByte() {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		if l.ch == '\n' {
			l.line++
			l.column = 0
		}
		l.readChar()
	}
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return eofByte()
	} else {
		return l.input[l.readPosition]
	}
}

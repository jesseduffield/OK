package parser

import "github.com/jesseduffield/OK/token"

const (
	_ int = iota
	LOWEST
	ASSIGN
	LAZY            // lazy myFunc()
	ANDOR           // && or ||
	COMPARISON      // >=
	SUM_AND_PRODUCT // *
	PREFIX          // -X or !X
	NEW             // new Person()
	CALL            // myFunction(X)
	MEMBERACCESS    // myStruct.foo
	INDEX           // array[index]
	COMMENT         // e.g. '// I acknowledge that I shouldn't use this private field'
)

var precedences = map[token.TokenType]int{
	token.GTEQ:     COMPARISON,
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
	token.COMMENT:  COMMENT,
}

package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
	// Identifiers + literals
	IDENT  = "IDENT" // add, foobar, x, y, ...
	INT    = "INT"
	STRING = "STRING"
	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	AND      = "&&"
	OR       = "||"

	GT = ">"

	EQ = "=="

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"
	LBRACKET  = "["
	RBRACKET  = "]"
	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
	SWITCH   = "SWITCH"
	CASE     = "CASE"
	DEFAULT  = "DEFAULT"
	NULL     = "NO!"
	LAZY     = "LAZY"

	// structs
	STRUCT = "STRUCT"
	PACK   = "PACK"
	FIELD  = "FIELD"
	PUBLIC = "PUBLIC"
	NEW    = "NEW"
	PERIOD = "PERIOD"
)

var keywords = map[string]TokenType{
	"fn":        FUNCTION,
	"let":       LET,
	"true":      TRUE,
	"false":     FALSE,
	"NO!":       NULL,
	"if":        IF,
	"else":      ELSE,
	"return":    RETURN,
	"switch":    SWITCH,
	"case":      CASE,
	"default":   DEFAULT,
	"notaclass": STRUCT,
	"pack":      PACK,
	"field":     FIELD,
	"public":    PUBLIC,
	"new":       NEW,
	"lazy":      LAZY,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

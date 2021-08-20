package ast

import (
	"bytes"
	"strings"

	"github.com/jesseduffield/OK/ok/token"
)

type Node interface {
	TokenLiteral() string
	GetToken() token.Token
	String() string
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

type Statement interface {
	Node
	statementNode()
}

func (self *Program) TokenLiteral() string {
	if len(self.Statements) > 0 {
		return self.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (self *Program) GetToken() token.Token {
	if len(self.Statements) > 0 {
		return self.Statements[0].GetToken()
	} else {
		return token.Token{} // pretty sure we're okay to do this
	}
}

func (self *Program) String() string {
	var out bytes.Buffer

	for _, s := range self.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (self *Identifier) expressionNode()       {}
func (self *Identifier) GetToken() token.Token { return self.Token }
func (self *Identifier) TokenLiteral() string  { return self.Token.Literal }
func (self *Identifier) String() string        { return self.Value }

type IntegerLiteral struct {
	Token token.Token // the token.IDENT token
	Value int64
}

func (self *IntegerLiteral) expressionNode()       {}
func (self *IntegerLiteral) GetToken() token.Token { return self.Token }
func (self *IntegerLiteral) TokenLiteral() string  { return self.Token.Literal }
func (self *IntegerLiteral) String() string        { return self.Token.Literal }

type NullLiteral struct {
	Token token.Token // the token.NULL token
}

func (self *NullLiteral) expressionNode()       {}
func (self *NullLiteral) GetToken() token.Token { return self.Token }
func (self *NullLiteral) TokenLiteral() string  { return self.Token.Literal }
func (self *NullLiteral) String() string        { return self.Token.Literal }

type StringLiteral struct {
	Token token.Token
	Value string
}

func (self *StringLiteral) expressionNode()       {}
func (self *StringLiteral) GetToken() token.Token { return self.Token }
func (self *StringLiteral) TokenLiteral() string  { return self.Token.Literal }
func (self *StringLiteral) String() string        { return self.Token.Literal }

type PrefixExpression struct {
	Token    token.Token // The prefix token, e.g. !
	Operator string
	Right    Expression
}

func (self *PrefixExpression) expressionNode()       {}
func (self *PrefixExpression) GetToken() token.Token { return self.Token }
func (self *PrefixExpression) TokenLiteral() string  { return self.Token.Literal }
func (self *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(self.Operator)
	out.WriteString(self.Right.String())
	out.WriteString(")")

	return out.String()
}

type InfixExpression struct {
	Token    token.Token // The operator token, e.g. +
	Left     Expression
	Operator string
	Right    Expression
}

func (self *InfixExpression) expressionNode()       {}
func (self *InfixExpression) GetToken() token.Token { return self.Token }
func (self *InfixExpression) TokenLiteral() string  { return self.Token.Literal }
func (self *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(self.Left.String())
	out.WriteString(" " + self.Operator + " ")
	out.WriteString(self.Right.String())
	out.WriteString(")")

	return out.String()
}

type Boolean struct {
	Token token.Token
	Value bool
}

func (self *Boolean) expressionNode()       {}
func (self *Boolean) GetToken() token.Token { return self.Token }
func (self *Boolean) TokenLiteral() string  { return self.Token.Literal }
func (self *Boolean) String() string        { return self.Token.Literal }

type IfExpression struct {
	Token       token.Token // The 'if' token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (self *IfExpression) expressionNode()       {}
func (self *IfExpression) GetToken() token.Token { return self.Token }
func (self *IfExpression) TokenLiteral() string  { return self.Token.Literal }
func (self *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(self.Condition.String())
	out.WriteString(" ")
	out.WriteString(self.Consequence.String())

	if self.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(self.Alternative.String())
	}

	return out.String()
}

type SwitchCase struct {
	Value Expression
	Block *BlockStatement
}

type SwitchExpression struct {
	Token   token.Token // The 'switch' token
	Subject Expression
	Cases   []SwitchCase
	Default *BlockStatement
}

func (self *SwitchExpression) expressionNode()       {}
func (self *SwitchExpression) GetToken() token.Token { return self.Token }
func (self *SwitchExpression) TokenLiteral() string  { return self.Token.Literal }
func (self *SwitchExpression) String() string {
	var out bytes.Buffer

	out.WriteString("switch")
	out.WriteString(" ")
	out.WriteString(self.Subject.String())
	out.WriteString(" {")
	for _, e := range self.Cases {
		out.WriteString(" case ")
		out.WriteString(e.Value.String())
		out.WriteString(": { ")
		out.WriteString(e.Block.String())
		out.WriteString(" }")
	}

	if self.Default != nil {
		out.WriteString(" default: ")
		out.WriteString(self.Default.String())
	}

	out.WriteString("}")

	return out.String()
}

type BlockStatement struct {
	Token      token.Token // the { token
	Statements []Statement
}

func (self *BlockStatement) statementNode()        {}
func (self *BlockStatement) GetToken() token.Token { return self.Token }
func (self *BlockStatement) TokenLiteral() string  { return self.Token.Literal }
func (self *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range self.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

type FunctionLiteral struct {
	Token      token.Token // The 'fn' token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (self *FunctionLiteral) expressionNode()       {}
func (self *FunctionLiteral) GetToken() token.Token { return self.Token }
func (self *FunctionLiteral) TokenLiteral() string  { return self.Token.Literal }
func (self *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range self.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(self.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") { ")
	out.WriteString(self.Body.String())
	out.WriteString(" }")

	return out.String()
}

type CallExpression struct {
	Token     token.Token // The '(' token
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
}

func (self *CallExpression) expressionNode()       {}
func (self *CallExpression) GetToken() token.Token { return self.Token }
func (self *CallExpression) TokenLiteral() string  { return self.Token.Literal }
func (self *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range self.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(self.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

type ArrayLiteral struct {
	Token    token.Token // the '[' token
	Elements []Expression
}

func (self *ArrayLiteral) expressionNode()       {}
func (self *ArrayLiteral) GetToken() token.Token { return self.Token }
func (self *ArrayLiteral) TokenLiteral() string  { return self.Token.Literal }
func (self *ArrayLiteral) String() string {
	var out bytes.Buffer

	elements := []string{}
	for _, el := range self.Elements {
		elements = append(elements, el.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

type IndexExpression struct {
	Token token.Token // The [ token
	Left  Expression
	Index Expression
}

func (self *IndexExpression) expressionNode()       {}
func (self *IndexExpression) GetToken() token.Token { return self.Token }
func (self *IndexExpression) TokenLiteral() string  { return self.Token.Literal }
func (self *IndexExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(self.Left.String())
	out.WriteString("[")
	out.WriteString(self.Index.String())
	out.WriteString("])")

	return out.String()
}

type HashLiteral struct {
	Token token.Token // the '{' token
	Pairs map[Expression]Expression
}

func (self *HashLiteral) expressionNode()       {}
func (self *HashLiteral) GetToken() token.Token { return self.Token }
func (self *HashLiteral) TokenLiteral() string  { return self.Token.Literal }
func (self *HashLiteral) String() string {
	var out bytes.Buffer

	pairs := []string{}
	for key, value := range self.Pairs {
		pairs = append(pairs, key.String()+":"+value.String())
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

type LetStatement struct {
	Token token.Token // the token.LET token
	Name  *Identifier
	Value Expression
}

func (self *LetStatement) statementNode()        {}
func (self *LetStatement) GetToken() token.Token { return self.Token }
func (self *LetStatement) TokenLiteral() string  { return self.Token.Literal }
func (self *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(self.TokenLiteral() + " ")
	out.WriteString(self.Name.String())
	out.WriteString(" = ")

	if self.Value != nil {
		out.WriteString(self.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

type ReturnStatement struct {
	Token       token.Token // the 'return' token
	ReturnValue Expression
}

func (self *ReturnStatement) statementNode()        {}
func (self *ReturnStatement) GetToken() token.Token { return self.Token }
func (self *ReturnStatement) TokenLiteral() string  { return self.Token.Literal }
func (self *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(self.TokenLiteral() + " ")

	if self.ReturnValue != nil {
		out.WriteString(self.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}

type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}

func (self *ExpressionStatement) statementNode()        {}
func (self *ExpressionStatement) GetToken() token.Token { return self.Token }
func (self *ExpressionStatement) TokenLiteral() string  { return self.Token.Literal }
func (self *ExpressionStatement) String() string {
	if self.Expression != nil {
		return self.Expression.String()
	}
	return ""
}

type LazyExpression struct {
	Token token.Token // The prefix token, i.e. 'lazy'
	Right Expression
}

func (self *LazyExpression) expressionNode()       {}
func (self *LazyExpression) GetToken() token.Token { return self.Token }
func (self *LazyExpression) TokenLiteral() string  { return self.Token.Literal }
func (self *LazyExpression) String() string {
	var out bytes.Buffer

	out.WriteString("lazy(")
	out.WriteString(self.Right.String())
	out.WriteString(")")

	return out.String()
}

type CommentStatement struct {
	Token token.Token
	Text  string
}

func (self *CommentStatement) statementNode()        {}
func (self *CommentStatement) GetToken() token.Token { return self.Token }
func (self *CommentStatement) TokenLiteral() string  { return self.Token.Literal }
func (self *CommentStatement) String() string {
	var out bytes.Buffer

	out.WriteString("// ")
	out.WriteString(self.Text)

	return out.String()
}

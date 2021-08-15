package ast

import (
	"bytes"
	"fmt"
	"monkey/token"
)

type StructField struct {
	Name   string
	Public bool
}

type StructMethod struct {
	FunctionLiteral *FunctionLiteral
	Public          bool
}

type Struct struct {
	Token token.Token // The 'struct' token
	Name  string

	// will be empty if no privacy acknowledgement is set
	PrivacyAcknowledgement string

	Fields  []StructField
	Methods map[string]StructMethod
}

func (self *Struct) statementNode()       {}
func (self *Struct) TokenLiteral() string { return self.Token.Literal }
func (self *Struct) String() string {
	var out bytes.Buffer

	out.WriteString("notAClass ")
	out.WriteString(self.Name)
	out.WriteString(" {\n")

	if self.PrivacyAcknowledgement != "" {
		out.WriteString("\tpack \"")
		out.WriteString(self.PrivacyAcknowledgement)
		out.WriteString("\"\n\n")
	}

	for _, f := range self.Fields {
		out.WriteString("\t")
		if f.Public {
			out.WriteString("public ")
		}
		out.WriteString("field ")
		out.WriteString(f.Name)
		out.WriteString("\n")
	}

	if len(self.Fields) > 0 && len(self.Methods) > 0 {
		out.WriteString("\n")
	}

	for methodName, method := range self.Methods {
		out.WriteString("\t")
		if method.Public {
			out.WriteString("public ")
		}
		out.WriteString(methodName)
		out.WriteString(" ")

		out.WriteString(method.FunctionLiteral.String())
		out.WriteString("\n")
	}
	out.WriteString("}")

	return out.String()
}

type StructInstantiation struct {
	Token      token.Token
	StructName string
	Arguments  []Expression
}

func (self *StructInstantiation) expressionNode()      {}
func (self *StructInstantiation) TokenLiteral() string { return self.Token.Literal }
func (self *StructInstantiation) String() string {
	var out bytes.Buffer

	out.WriteString("new ")
	out.WriteString(self.StructName)
	out.WriteString("(")
	for i, arg := range self.Arguments {
		out.WriteString(arg.String())
		if i < len(self.Arguments)-1 {
			out.WriteString(", ")
		}
	}
	out.WriteString(")")

	return out.String()
}

type StructMemberAccessExpression struct {
	Token      token.Token // The . token
	Left       Expression
	MemberName string
}

func (self *StructMemberAccessExpression) expressionNode()      {}
func (self *StructMemberAccessExpression) TokenLiteral() string { return self.Token.Literal }
func (self *StructMemberAccessExpression) String() string {
	return fmt.Sprintf("%s.%s", self.Left.String(), self.MemberName)
}

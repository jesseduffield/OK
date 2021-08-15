package object

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/jesseduffield/OK/ast"
)

type StructInstance struct {
	Fields map[string]Object
	Struct *ast.Struct
}

func (self *StructInstance) Type() ObjectType { return HASH_OBJ }
func (self *StructInstance) Inspect() string {
	var out bytes.Buffer

	out.WriteString(self.Struct.Name)

	pairs := []string{}
	for key, obj := range self.Fields {
		pairs = append(pairs, fmt.Sprintf("%s: %s", key, obj.Inspect()))
	}

	out.WriteString(": {")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

func (self *StructInstance) IsField(fieldName string) bool {
	for _, field := range self.Struct.Fields {
		if field.Name == fieldName {
			return true
		}
	}
	return false
}

func (self *StructInstance) IsMethod(methodName string) bool {
	for name := range self.Struct.Methods {
		if name == methodName {
			return true
		}
	}
	return false
}

// struct method instance object: it contains a pointer to the struct instance so that when evaluated in the context of the struct, it can access the fields of the struct

func (self *StructInstance) GetMethod(methodName string) Object {
	for name, method := range self.Struct.Methods {
		if name == methodName {
			return &Method{
				StructInstance: self,
				Name:           methodName,
				StructMethod:   method,
			}
		}
	}

	return NewError(fmt.Sprintf("No such method for struct %s: %s", self.Struct.Name, methodName))
}

func (self *StructInstance) GetFieldValue(fieldName string) Object {
	for name, field := range self.Fields {
		if name == fieldName {
			return field
		}
	}
	return NULL
}

// assumes we've already determined that the field exists
func (self *StructInstance) IsPublicField(fieldName string) bool {
	for _, field := range self.Struct.Fields {
		if field.Name == fieldName {
			return field.Public
		}
	}
	return false
}

// assumes we've already determined that the method exists
func (self *StructInstance) IsPublicMethod(methodName string) bool {
	for name, method := range self.Struct.Methods {
		if name == methodName {
			return method.Public
		}
	}
	return false
}

func (self *StructInstance) SetFieldValue(fieldName string, value Object) Object {
	if !self.IsField(fieldName) {
		return NewError(fmt.Sprintf("No such field for struct %s: %s", self.Struct.Name, fieldName))
	}

	self.Fields[fieldName] = value

	return value
}

func (self *StructInstance) EvolveInto(other *StructInstance) {
	self.Struct = other.Struct
	self.Fields = other.Fields
}

type Method struct {
	Name         string
	StructMethod ast.StructMethod

	// not sure if we need an env here, or if we do which env to pass.
	Env *Environment

	StructInstance *StructInstance
}

func (self *Method) Type() ObjectType { return METHOD_OBJ }
func (self *Method) Inspect() string {
	params := []string{}
	for _, p := range self.StructMethod.FunctionLiteral.Parameters {
		params = append(params, p.String())
	}

	return fmt.Sprintf("(%s) %s fn(%s) {\n\t%s\n}", self.StructInstance.Struct.Name, self.Name, strings.Join(params, ", "), self.StructMethod.FunctionLiteral.Body.String())
}

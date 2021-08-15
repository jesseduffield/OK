package object

import (
	"fmt"
	"monkey/ast"
)

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	vs := make(map[string]*ast.Struct)
	return &Environment{variableStore: s, structStore: vs, outer: nil}
}

// Hack for now that allows us to access all the defined structs without accessing any variables.
// A better solution would be to store that in a separate environment
func OnlyStructs(env *Environment) *Environment {
	if env == nil {
		return nil
	}
	s := make(map[string]Object)
	return &Environment{variableStore: s, structStore: env.structStore, outer: OnlyStructs(env.outer)}
}

type Environment struct {
	variableStore         map[string]Object
	structStore           map[string]*ast.Struct
	outer                 *Environment
	currentStructInstance *StructInstance
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.variableStore[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

// Set is for declaring variables and then assigning to them in the current environment. If the variable is declared in a parent environment, that variable will now be shadowed
func (e *Environment) Set(name string, val Object) Object {
	e.variableStore[name] = val
	return val
}

// assign expects the variable to be declared and findable in the environment chain somewhere. Wherever it's found in the chain we'll store the new value
func (e *Environment) Assign(name string, val Object) Object {
	// I need to actually assign to the env that currently has the key
	env := e
	for env != nil {
		_, ok := env.variableStore[name]
		if ok {
			env.variableStore[name] = val
			return val
		}
		env = env.outer
	}

	return NewError("%s has not been declared", name)
}

func (e *Environment) GetStruct(name string) (*ast.Struct, bool) {
	obj, ok := e.structStore[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.GetStruct(name)
	}
	return obj, ok
}

func (e *Environment) SetStruct(structDef *ast.Struct) *ast.Struct {
	e.structStore[structDef.Name] = structDef
	return structDef
}

func (e *Environment) SetCurrentStructInstance(structInstance *StructInstance) {
	e.currentStructInstance = structInstance
}

func (e *Environment) IsCurrentStructInstance(structInstance *StructInstance) bool {
	return e.currentStructInstance == structInstance || (e.outer != nil && e.outer.IsCurrentStructInstance(structInstance))
}

func (e *Environment) String() string {
	result := ""
	result += "Variables:\n"
	for name, obj := range e.variableStore {
		if obj == nil {
			result += fmt.Sprintf("%s: nil\n", name)
		} else {
			result += name + ": " + obj.Inspect() + "\n"
		}
	}
	result += "Structs:\n"
	for name, obj := range e.structStore {
		result += name + ": " + obj.String() + "\n"
	}
	if e.outer != nil {
		result += "Outer:\n"
		result += e.outer.String()
	}

	return result
}

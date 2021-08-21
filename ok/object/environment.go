package object

import (
	"fmt"
	"sync"

	"github.com/jesseduffield/OK/ok/ast"
)

type Environment struct {
	variableStore         map[string]Object
	structStore           map[string]*ast.Struct
	outer                 *Environment
	currentStructInstance *StructInstance
	acknowledgements      map[string]bool

	mutex sync.Mutex
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	vs := make(map[string]*ast.Struct)
	acknowledgements := make(map[string]bool)

	return &Environment{
		variableStore:    s,
		structStore:      vs,
		outer:            nil,
		acknowledgements: acknowledgements,
	}
}

// Hack for now that allows us to access all the defined structs without accessing any variables.
// A better solution would be to store that in a separate environment
func OnlyStructs(env *Environment) *Environment {
	if env == nil {
		return nil
	}
	s := make(map[string]Object)
	return &Environment{
		variableStore: s,
		structStore:   env.structStore,
		outer:         OnlyStructs(env.outer),
	}
}

func (e *Environment) Get(name string) (Object, bool) {
	current := e
	for current != nil {
		current.mutex.Lock()
		defer current.mutex.Unlock()

		obj, ok := current.variableStore[name]
		if ok {
			return obj, ok
		}

		current = current.outer
	}

	return NULL, false
}

// Set is for declaring variables and then assigning to them in the current environment. If the variable is declared in a parent environment, that variable will now be shadowed
func (e *Environment) Set(name string, val Object) Object {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	e.variableStore[name] = val
	return val
}

// assign expects the variable to be declared and findable in the environment chain somewhere. Wherever it's found in the chain we'll store the new value
func (e *Environment) Assign(name string, val Object) (Object, error) {
	// I need to actually assign to the current that currently has the key
	current := e
	for current != nil {
		current.mutex.Lock()
		defer current.mutex.Unlock()

		_, ok := current.variableStore[name]
		if ok {
			current.variableStore[name] = val
			return val, nil
		}
		current = current.outer
	}

	return NULL, fmt.Errorf("%s has not been declared", name)
}

func (e *Environment) GetStruct(name string) (*ast.Struct, bool) {
	current := e
	for current != nil {
		current.mutex.Lock()
		defer current.mutex.Unlock()

		obj, ok := e.structStore[name]
		if ok {
			return obj, ok
		}

		current = current.outer
	}

	return nil, false
}

func (e *Environment) SetStruct(structDef *ast.Struct) *ast.Struct {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	e.structStore[structDef.Name] = structDef
	return structDef
}

func (e *Environment) SetCurrentStructInstance(structInstance *StructInstance) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	e.currentStructInstance = structInstance
}

func (e *Environment) IsCurrentStructInstance(structInstance *StructInstance) bool {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	return e.currentStructInstance == structInstance ||
		(e.outer != nil && e.outer.IsCurrentStructInstance(structInstance))
}

func (e *Environment) String() string {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	result := ""
	result += "Variables:\n"
	for name, obj := range e.variableStore {
		if obj == nil {
			result += fmt.Sprintf("%s: NO!\n", name)
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

func (e *Environment) AddAcknowledgement(ack string) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	e.acknowledgements[ack] = true
}

func (e *Environment) AllowsPrivateAccess(str *ast.Struct) bool {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	return e.acknowledgements[str.PrivacyAcknowledgement]
}

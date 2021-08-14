package object

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, outer: nil}
}

type Environment struct {
	store map[string]Object
	outer *Environment
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

// Set is for declaring variables and then assigning to them in the current environment. If the variable is declared in a parent environment, that variable will now be shadowed
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

// assign expects the variable to be declared and findable in the environment chain somewhere. Wherever it's found in the chain we'll store the new value
func (e *Environment) Assign(name string, val Object) Object {
	// I need to actually assign to the env that currently has the key
	env := e
	for env != nil {
		_, ok := env.store[name]
		if ok {
			env.store[name] = val
			return val
		}
		env = env.outer
	}

	return NewError("%s has not been declared", name)
}

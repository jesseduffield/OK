package evaluator

import (
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/jesseduffield/OK/ok/object"
)

func (e *Evaluator) getBuiltins(out io.Writer) map[string]*object.Builtin {
	return map[string]*object.Builtin{
		"len": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return object.NewError("wrong number of arguments. got=%d, want=1",
						len(args))
				}

				switch arg := args[0].(type) {
				case *object.Array:
					return &object.Integer{Value: int64(len(arg.Elements))}
				case *object.String:
					return &object.Integer{Value: int64(len(arg.Value))}
				default:
					return object.NewError("argument to `len` not supported, got %s",
						args[0].Type())
				}
			},
		},
		"first": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return object.NewError("wrong number of arguments. got=%d, want=1",
						len(args))
				}
				if args[0].Type() != object.ARRAY_OBJ {
					return object.NewError("argument to `first` must be ARRAY, got %s",
						args[0].Type())
				}

				arr := args[0].(*object.Array)
				if len(arr.Elements) > 0 {
					return arr.Elements[0]
				}

				return object.NULL
			},
		},
		"last": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return object.NewError("wrong number of arguments. got=%d, want=1",
						len(args))
				}
				if args[0].Type() != object.ARRAY_OBJ {
					return object.NewError("argument to `last` must be ARRAY, got %s",
						args[0].Type())
				}

				arr := args[0].(*object.Array)
				length := len(arr.Elements)
				if length > 0 {
					return arr.Elements[length-1]
				}

				return object.NULL
			},
		},
		"rest": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return object.NewError("wrong number of arguments. got=%d, want=1",
						len(args))
				}
				if args[0].Type() != object.ARRAY_OBJ {
					return object.NewError("argument to `rest` must be ARRAY, got %s",
						args[0].Type())
				}

				arr := args[0].(*object.Array)
				length := len(arr.Elements)
				if length > 0 {
					newElements := make([]object.Object, length-1)
					copy(newElements, arr.Elements[1:length])
					return &object.Array{Elements: newElements}
				}

				return object.NULL
			},
		},
		"push": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return object.NewError("wrong number of arguments. got=%d, want=2",
						len(args))
				}
				if args[0].Type() != object.ARRAY_OBJ {
					return object.NewError("argument to `push` must be ARRAY, got %s",
						args[0].Type())
				}

				arr := args[0].(*object.Array)
				length := len(arr.Elements)

				newElements := make([]object.Object, length+1)
				copy(newElements, arr.Elements)
				newElements[length] = args[1]

				return &object.Array{Elements: newElements}
			},
		},
		"puts": {
			Fn: func(args ...object.Object) object.Object {
				for _, arg := range args {
					fmt.Fprintln(out, arg.Inspect())
				}

				return object.NULL
			},
		},
		"ayok?": {
			Fn: func(args ...object.Object) object.Object {
				return nativeBoolToBooleanObject(args[0] != object.NULL)
			},
		},
		"sleep": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return object.NewError("wrong number of arguments. got=%d, want=1", len(args))
				}
				if args[0].Type() != object.INTEGER_OBJ {
					return object.NewError("argument to `sleep` must be INTEGER, got %s", args[0].Type())
				}

				time.Sleep(time.Duration(args[0].(*object.Integer).Value) * time.Second)

				return object.NULL
			},
		},
		"map": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return object.NewError("wrong number of arguments. got=%d, want=2", len(args))
				}

				arr := args[0]
				if arr.Type() != object.ARRAY_OBJ {
					return object.NewError("First argument to `map` must be ARRAY, got %s", arr.Type())
				}

				fn := args[1]
				if fn.Type() != object.FUNCTION_OBJ {
					return object.NewError("Second argument to `map` must be FUNCTION, got %s", fn.Type())
				}

				arrObj := arr.(*object.Array)
				fnObj := fn.(*object.Function)
				if len(fnObj.Parameters) > 2 || len(fnObj.Parameters) < 1 {
					return object.NewError("Function must have 1 or 2 parameters, got %d", len(fnObj.Parameters))
				}

				result := &object.Array{Elements: make([]object.Object, len(arrObj.Elements))}
				waitGroup := &sync.WaitGroup{}
				waitGroup.Add(len(arrObj.Elements))
				for i, el := range arrObj.Elements {
					el := el
					i := i
					go func() {
						if len(fnObj.Parameters) == 1 {
							result.Elements[i] = e.applyUserFunction(fnObj, []object.Object{el})
						} else {
							result.Elements[i] = e.applyUserFunction(fnObj, []object.Object{el, &object.Integer{Value: int64(i)}})
						}

						waitGroup.Done()
					}()
				}

				waitGroup.Wait()

				return result
			},
		},
	}
}

package evaluator

import (
	"fmt"
	"io"
	"strings"

	"github.com/jesseduffield/OK/ok/ast"
	"github.com/jesseduffield/OK/ok/object"
)

type Evaluator struct {
	out io.Writer
}

func New(out io.Writer) *Evaluator {
	return &Evaluator{out: out}
}

func (e *Evaluator) Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {

	// Statements
	case *ast.Program:
		return e.evalProgram(node, env)

	case *ast.ExpressionStatement:
		return e.Eval(node.Expression, env)

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.PrefixExpression:
		return e.evalPrefixExpression(node.Operator, node.Right, env)

	case *ast.InfixExpression:
		switch node.Operator {
		case "=":
			return e.evalAssignmentExpression(node.Left, node.Right, env)
		case "&&":
			left := e.Eval(node.Left, env)
			if isError(left) {
				return left
			}

			if !isTruthy(left) {
				return nativeBoolToBooleanObject(false)
			}

			right := e.Eval(node.Right, env)
			if isError(right) {
				return right
			}

			return nativeBoolToBooleanObject(isTruthy(right))
		case "||":
			left := e.Eval(node.Left, env)
			if isError(left) {
				return left
			}

			if isTruthy(left) {
				return nativeBoolToBooleanObject(true)
			}

			right := e.Eval(node.Right, env)
			if isError(right) {
				return right
			}

			return nativeBoolToBooleanObject(isTruthy(right))
		}

		left := e.Eval(node.Left, env)
		if isError(left) {
			return left
		}

		right := e.Eval(node.Right, env)
		if isError(right) {
			return right
		}

		return e.evalInfixExpression(node.Operator, left, right)

	case *ast.BlockStatement:
		return e.evalBlockStatement(node, env)

	case *ast.IfExpression:
		return e.evalIfExpression(node, env)

	case *ast.SwitchExpression:
		return e.evalSwitchExpression(node, env)

	case *ast.ReturnStatement:
		val := e.Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	case *ast.LetStatement:
		val := e.Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)

	case *ast.Identifier:
		return e.evalIdentifier(node, env)

	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Env: env, Body: body}

	case *ast.CallExpression:
		function := e.Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := e.evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return e.applyFunction(function, args, env)

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	case *ast.ArrayLiteral:
		elements := e.evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}

	case *ast.IndexExpression:
		left := e.Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := e.Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return e.evalIndexExpression(left, index)

	case *ast.HashLiteral:
		return e.evalHashLiteral(node, env)

	case *ast.Struct:
		return e.evalStructDefinition(node, env)

	case *ast.StructInstantiation:
		return e.evalStructInstantiation(node, env)

	case *ast.StructMemberAccessExpression:
		return e.evalStructMemberAccess(node, env)

	case *ast.NullLiteral:
		return object.NULL

	case *ast.LazyExpression:
		return e.evalLazyExpression(node)

	case *ast.CommentStatement:
		return e.evalCommentStatement(node, env)

	case nil:
		// TODO: I'm not actually sure why this would ever be nil. Might need to investigate
		return object.NULL

	default:
		return object.NewError("Unknown expression type: %T", node)
	}

	return nil
}

func (e *Evaluator) evalCommentStatement(node *ast.CommentStatement, env *object.Environment) object.Object {
	acknowledgePrefix := "I acknowledge that "
	if strings.HasPrefix(node.Text, acknowledgePrefix) {
		text := strings.TrimPrefix(node.Text, acknowledgePrefix)
		env.AddAcknowledgement(text)
	}

	return object.NULL
}

func (e *Evaluator) evalLazyExpression(node *ast.LazyExpression) object.Object {
	return &object.LazyObject{
		Right: node.Right,
	}
}

// this is only for when used as an expression
func (e *Evaluator) evalStructMemberAccess(node *ast.StructMemberAccessExpression, env *object.Environment) object.Object {
	left := e.Eval(node.Left, env)
	structInstance, ok := left.(*object.StructInstance)
	if !ok {
		return object.NewError(fmt.Sprintf("`%s` is not a nac", node.Left.String()))
	}

	if structInstance.IsField(node.MemberName) {
		if !structInstance.IsPublicField(node.MemberName) && !env.IsCurrentStructInstance(structInstance) && !env.AllowsPrivateAccess(structInstance.Struct) {
			return object.NewError(fmt.Sprintf("`%s` is a private field on nac %s", node.MemberName, structInstance.Struct.Name))
		}
		return structInstance.GetFieldValue(node.MemberName)
	} else if structInstance.IsMethod(node.MemberName) {
		if !structInstance.IsPublicMethod(node.MemberName) && !env.IsCurrentStructInstance(structInstance) && !env.AllowsPrivateAccess(structInstance.Struct) {
			return object.NewError(fmt.Sprintf("`%s` is a private method on nac %s", node.MemberName, structInstance.Struct.Name))
		}
		return structInstance.GetMethod(node.MemberName)
	} else {
		return object.NewError(fmt.Sprintf("undefined field for nac %s: %s", structInstance.Struct.Name, node.MemberName))
	}
}

// TODO: support referring to global variables from within a struct
func (e *Evaluator) evalStructDefinition(structDef *ast.Struct, env *object.Environment) object.Object {
	env.SetStruct(structDef)
	return object.NULL
}

func (e *Evaluator) evalStructInstantiation(node *ast.StructInstantiation, env *object.Environment) object.Object {
	instance := &object.StructInstance{}
	instance.Fields = make(map[string]object.Object)
	// need to find the struct in our env
	tmp, ok := env.GetStruct(node.StructName)
	if !ok {
		return object.NewError(fmt.Sprintf("undefined nac %s", node.StructName))
	}

	instance.Struct = tmp

	// ignoring fields for now

	return instance
}

func (e *Evaluator) evalExpressions(
	exps []ast.Expression,
	env *object.Environment,
) []object.Object {
	var result []object.Object

	for _, exp := range exps {
		evaluated := e.Eval(exp, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func (e *Evaluator) evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = e.Eval(statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func (e *Evaluator) evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = e.Eval(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return object.TRUE
	}
	return object.FALSE
}

func (e *Evaluator) evalPrefixExpression(operator string, rightNode ast.Node, env *object.Environment) object.Object {
	switch operator {
	case "!":
		right := e.Eval(rightNode, env)
		if isError(right) {
			return right
		}
		return e.evalBangOperatorExpression(right)
	case "-":
		right := e.Eval(rightNode, env)
		if isError(right) {
			return right
		}
		return e.evalMinusPrefixOperatorExpression(right)
	case "lazy":
		return &object.LazyObject{Right: rightNode}
	default:
		return object.NewError("unknown operator: %s for %s", operator, rightNode.String())
	}
}

func (e *Evaluator) evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case object.TRUE:
		return object.FALSE
	case object.FALSE:
		return object.TRUE
	case object.NULL:
		return object.TRUE
	default:
		return object.FALSE
	}
}

func (e *Evaluator) evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return object.NewError("unknown operator: -%s", right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func (e *Evaluator) evalInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return e.evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return e.evalStringInfixExpression(operator, left, right)
	case operator == ">=":
		// for bools, nil, structs, hashes, and arrays, >= is true if and only if
		// == is true
		return nativeBoolToBooleanObject(left == right)
	case operator == "==":
		// this is allowed internally but illegal in the lexer
		return nativeBoolToBooleanObject(left == right)
	case left.Type() != right.Type():
		return object.NewError("type mismatch: %s %s %s",
			left.Type(), operator, right.Type())
	default:
		return object.NewError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func (e *Evaluator) evalAssignmentExpression(left ast.Expression, right ast.Expression, env *object.Environment) object.Object {
	val := e.Eval(right, env)
	if isError(val) {
		return val
	}

	switch v := left.(type) {
	case *ast.Identifier:
		return env.Assign(v.Value, val)
	case *ast.IndexExpression:
		// I can just evaluate the left entirely and that will leave me with an object
		key := e.Eval(v.Index, env)
		if isError(key) {
			return key
		}
		leftVal := e.Eval(v.Left, env)
		if isError(leftVal) {
			return leftVal
		}

		switch l := leftVal.(type) {
		case *object.Array:
			indexVal, ok := key.(*object.Integer)
			if !ok {
				return object.NewError("Index must be an integer")
			}
			if indexVal.Value < 0 {
				return object.NewError("Index must be positive")
			}
			if int(indexVal.Value) > len(l.Elements)-1 {
				return object.NewError(fmt.Sprintf("Index %d is out of bounds (array length %d)", indexVal.Value, len(l.Elements)))
			}
			l.Elements[indexVal.Value] = val
		case *object.Hash:
			hashKey, ok := key.(object.Hashable)
			if !ok {
				return object.NewError("Unusable as hash key: %s", key.Type())
			}

			l.Pairs[hashKey.HashKey()] = object.HashPair{Key: key, Value: val}
		case *object.Null:
			return object.NewError("Attempted index of NULL object")
		default:
			return object.NewError(fmt.Sprintf("`%s` is neither a hash nor array so you cannot index into it", v.Left.String()))
		}
	case *ast.StructMemberAccessExpression:
		leftVal := e.Eval(v.Left, env)
		if isError(leftVal) {
			return leftVal
		}

		structInstance, ok := leftVal.(*object.StructInstance)
		if !ok {
			return object.NewError(fmt.Sprintf("`%s` is not a nac instance", v.Left.String()))
		}

		if structInstance.IsMethod(v.MemberName) {
			return object.NewError(fmt.Sprintf("`%s` is a method, not a field, on nac %s. You cannot reassign it", v.MemberName, structInstance.Struct.Name))
		}
		if !structInstance.IsPublicField(v.MemberName) && !env.IsCurrentStructInstance(structInstance) && !env.AllowsPrivateAccess(structInstance.Struct) {
			return object.NewError(fmt.Sprintf("`%s` is a private field on nac %s", v.MemberName, structInstance.Struct.Name))
		}

		structInstance.SetFieldValue(v.MemberName, val)

	default:
		return object.NewError("LHS must be an identifier or index expression")
	}

	return val
}

func (e *Evaluator) evalStringInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	switch operator {
	case "+":
		leftVal := left.(*object.String).Value
		rightVal := right.(*object.String).Value
		return &object.String{Value: leftVal + rightVal}
	case ">=":
		leftVal := left.(*object.String).Value
		rightVal := right.(*object.String).Value
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	case "==":
		leftVal := left.(*object.String).Value
		rightVal := right.(*object.String).Value
		// this is allowed internally but illegal in the lexer
		return nativeBoolToBooleanObject(leftVal == rightVal)
	default:
		return object.NewError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func (e *Evaluator) evalIntegerInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	case "==":
		// this is allowed internally but illegal in the lexer
		return nativeBoolToBooleanObject(leftVal == rightVal)
	default:
		return object.NewError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func (e *Evaluator) evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := e.Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return e.Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return e.Eval(ie.Alternative, env)
	} else {
		return object.NULL
	}
}

func (e *Evaluator) evalSwitchExpression(se *ast.SwitchExpression, env *object.Environment) object.Object {
	subject := e.Eval(se.Subject, env)
	if isError(subject) {
		return subject
	}

	for _, c := range se.Cases {
		value := e.Eval(c.Value, env)
		if value.Type() != subject.Type() {
			return object.NewError("mismatched types in switch statement: %s %s",
				subject.Type(), value.Type())
		}
		test := e.evalInfixExpression("==", subject, value)
		if test == object.TRUE {
			return e.Eval(c.Block, env)
		}
	}

	if se.Default != nil {
		return e.Eval(se.Default, env)
	}

	return object.NULL
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case object.NULL:
		return false
	case object.TRUE:
		return true
	case object.FALSE:
		return false
	default:
		return true
	}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func (e *Evaluator) evalIdentifier(
	node *ast.Identifier,
	env *object.Environment,
) object.Object {
	if val, ok := env.Get(node.Value); ok {
		if lazyObj, ok := val.(*object.LazyObject); ok {
			unwrappedVal := e.Eval(lazyObj.Right, env)
			env.Assign(node.Value, unwrappedVal)
			return unwrappedVal
		}

		return val
	}

	if builtin, ok := e.getBuiltins(e.out)[node.Value]; ok {
		return builtin
	}

	return object.NewError("identifier not found: " + node.Value)
}

func (e *Evaluator) applyUserFunction(fn *object.Function, args []object.Object) object.Object {
	extendedEnv := e.extendFunctionEnv(fn, args)
	evaluated := e.Eval(fn.Body, extendedEnv)
	return unwrapReturnValue(evaluated)
}

func (e *Evaluator) applyFunction(fn object.Object, args []object.Object, env *object.Environment) object.Object {
	switch fn := fn.(type) {

	case *object.Function:
		return e.applyUserFunction(fn, args)

	case *object.Method:
		newEnv := e.createMethodEnv(fn, args, env)
		evaluated := e.Eval(fn.StructMethod.FunctionLiteral.Body, newEnv)

		if err := e.handleEvolve(fn.StructInstance, env); err != nil {
			return err
		}

		return unwrapReturnValue(evaluated)

	case *object.Builtin:
		return fn.Fn(args...)

	default:
		return object.NewError("not a function: %s", fn.Type())
	}
}

func (e *Evaluator) handleEvolve(instance *object.StructInstance, env *object.Environment) object.Object {
	if instance.IsMethod("evolve") {
		evolveMethod := instance.GetMethod("evolve").(*object.Method)
		newEnv := e.createMethodEnv(evolveMethod, []object.Object{}, env)
		other := e.Eval(evolveMethod.StructMethod.FunctionLiteral.Body, newEnv)
		other = unwrapReturnValue(other)
		if other.Type() != object.NULL_OBJ {
			new, ok := other.(*object.StructInstance)
			if !ok {
				return object.NewError("evolve method must return NO! or a nac instance, returned %s: %s", other.Type(), other.Inspect())
			}
			instance.EvolveInto(new)
		}
	}

	return nil
}

func (e *Evaluator) createMethodEnv(
	method *object.Method,
	args []object.Object,
	env *object.Environment,
) *object.Environment {
	newEnv := object.OnlyStructs(env)

	functionLiteral := method.StructMethod.FunctionLiteral
	// if the first arg is 'selfish' we need to pass in the struct instance for that
	if len(functionLiteral.Parameters) > 0 && functionLiteral.Parameters[0].Value == "selfish" {
		newEnv.Set("selfish", method.StructInstance)

		for paramIdx, param := range functionLiteral.Parameters[1:] {
			newEnv.Set(param.Value, args[paramIdx])
		}
	} else {
		for paramIdx, param := range functionLiteral.Parameters {
			newEnv.Set(param.Value, args[paramIdx])
		}
	}

	newEnv.SetCurrentStructInstance(method.StructInstance)

	return newEnv
}

func (e *Evaluator) extendFunctionEnv(
	fn *object.Function,
	args []object.Object,
) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}

	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}

	return obj
}

func (e *Evaluator) evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return e.evalArrayIndexExpression(left, index)
	case left.Type() == object.HASH_OBJ:
		return e.evalHashIndexExpression(left, index)
	default:
		return object.NewError("index operator not supported: %s", left.Type())
	}
}

func (e *Evaluator) evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	if idx < 0 || idx > max {
		return object.NULL
	}

	return arrayObject.Elements[idx]
}

func (e *Evaluator) evalHashLiteral(
	node *ast.HashLiteral,
	env *object.Environment,
) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for keyNode, valueNode := range node.Pairs {
		key := e.Eval(keyNode, env)
		if isError(key) {
			return key
		}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return object.NewError("unusable as hash key: %s", key.Type())
		}

		value := e.Eval(valueNode, env)
		if isError(value) {
			return value
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}

	return &object.Hash{Pairs: pairs}
}

func (e *Evaluator) evalHashIndexExpression(hash, index object.Object) object.Object {
	hashObject := hash.(*object.Hash)

	key, ok := index.(object.Hashable)
	if !ok {
		return object.NewError("unusable as hash key: %s", index.Type())
	}

	pair, ok := hashObject.Pairs[key.HashKey()]
	if !ok {
		return object.NULL
	}

	return pair.Value
}

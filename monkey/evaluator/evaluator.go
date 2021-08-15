package evaluator

import (
	"fmt"
	"monkey/ast"
	"monkey/object"
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {

	// Statements
	case *ast.Program:
		return evalProgram(node, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.PrefixExpression:
		return evalPrefixExpression(node.Operator, node.Right, env)

	case *ast.InfixExpression:
		switch node.Operator {
		case "=":
			return evalAssignmentExpression(node.Left, node.Right, env)
		case "&&":
			left := Eval(node.Left, env)
			if isError(left) {
				return left
			}

			if !isTruthy(left) {
				return nativeBoolToBooleanObject(false)
			}

			right := Eval(node.Right, env)
			if isError(right) {
				return right
			}

			return nativeBoolToBooleanObject(isTruthy(right))
		case "||":
			left := Eval(node.Left, env)
			if isError(left) {
				return left
			}

			if isTruthy(left) {
				return nativeBoolToBooleanObject(true)
			}

			right := Eval(node.Right, env)
			if isError(right) {
				return right
			}

			return nativeBoolToBooleanObject(isTruthy(right))
		}

		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		return evalInfixExpression(node.Operator, left, right)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.SwitchExpression:
		return evalSwitchExpression(node, env)

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Env: env, Body: body}

	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyFunction(function, args, env)

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}

	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)

	case *ast.HashLiteral:
		return evalHashLiteral(node, env)

	case *ast.Struct:
		return evalStructDefinition(node, env)

	case *ast.StructInstantiation:
		return evalStructInstantiation(node, env)

	case *ast.StructMemberAccessExpression:
		return evalStructMemberAccess(node, env)

	case *ast.NullLiteral:
		return object.NULL

	case *ast.LazyExpression:
		return evalLazyExpression(node)

	case nil:
		// TODO: I'm not actually sure why this would ever be nil. Might need to investigate
		return object.NULL

	default:
		return object.NewError("Unknown expression type: %T", node)
	}

	return nil
}

func evalLazyExpression(node *ast.LazyExpression) object.Object {
	return &object.LazyObject{
		Right: node.Right,
	}
}

// this is only for when used as an expression
func evalStructMemberAccess(node *ast.StructMemberAccessExpression, env *object.Environment) object.Object {
	left := Eval(node.Left, env)
	structInstance, ok := left.(*object.StructInstance)
	if !ok {
		return object.NewError(fmt.Sprintf("`%s` is not a nac", node.Left.String()))
	}

	if structInstance.IsField(node.MemberName) {
		if !structInstance.IsPublicField(node.MemberName) && !env.IsCurrentStructInstance(structInstance) {
			return object.NewError(fmt.Sprintf("`%s` is a private field on nac %s", node.MemberName, structInstance.Struct.Name))
		}
		return structInstance.GetFieldValue(node.MemberName)
	} else if structInstance.IsMethod(node.MemberName) {
		if !structInstance.IsPublicMethod(node.MemberName) && !env.IsCurrentStructInstance(structInstance) {
			return object.NewError(fmt.Sprintf("`%s` is a private method on nac %s", node.MemberName, structInstance.Struct.Name))
		}
		return structInstance.GetMethod(node.MemberName)
	} else {
		return object.NewError(fmt.Sprintf("undefined field for nac %s: %s", structInstance.Struct.Name, node.MemberName))
	}
}

// TODO: support referring to global variables from within a struct
func evalStructDefinition(structDef *ast.Struct, env *object.Environment) object.Object {
	env.SetStruct(structDef)
	return object.NULL
}

func evalStructInstantiation(node *ast.StructInstantiation, env *object.Environment) object.Object {
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

func evalExpressions(
	exps []ast.Expression,
	env *object.Environment,
) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

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

func evalPrefixExpression(operator string, rightNode ast.Node, env *object.Environment) object.Object {
	switch operator {
	case "!":
		right := Eval(rightNode, env)
		if isError(right) {
			return right
		}
		return evalBangOperatorExpression(right)
	case "-":
		right := Eval(rightNode, env)
		if isError(right) {
			return right
		}
		return evalMinusPrefixOperatorExpression(right)
	case "lazy":
		return &object.LazyObject{Right: rightNode}
	default:
		return object.NewError("unknown operator: %s for %s", operator, rightNode.String())
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
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

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return object.NewError("unknown operator: -%s", right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return object.NewError("type mismatch: %s %s %s",
			left.Type(), operator, right.Type())
	default:
		return object.NewError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalAssignmentExpression(left ast.Expression, right ast.Expression, env *object.Environment) object.Object {
	val := Eval(right, env)
	if isError(val) {
		return val
	}

	switch v := left.(type) {
	case *ast.Identifier:
		return env.Assign(v.Value, val)
	case *ast.IndexExpression:
		// I can just evaluate the left entirely and that will leave me with an object
		key := Eval(v.Index, env)
		if isError(key) {
			return key
		}
		leftVal := Eval(v.Left, env)
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
		leftVal := Eval(v.Left, env)
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
		if !structInstance.IsPublicField(v.MemberName) && !env.IsCurrentStructInstance(structInstance) {
			return object.NewError(fmt.Sprintf("`%s` is a private field on nac %s", v.MemberName, structInstance.Struct.Name))
		}

		structInstance.SetFieldValue(v.MemberName, val)

	default:
		return object.NewError("LHS must be an identifier or index expression")
	}

	return val
}

func evalStringInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	switch operator {
	case "+":
		leftVal := left.(*object.String).Value
		rightVal := right.(*object.String).Value
		return &object.String{Value: leftVal + rightVal}
	case "==":
		leftVal := left.(*object.String).Value
		rightVal := right.(*object.String).Value
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		leftVal := left.(*object.String).Value
		rightVal := right.(*object.String).Value
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return object.NewError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(
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
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return object.NewError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return object.NULL
	}
}

func evalSwitchExpression(se *ast.SwitchExpression, env *object.Environment) object.Object {
	subject := Eval(se.Subject, env)
	if isError(subject) {
		return subject
	}

	for _, c := range se.Cases {
		value := Eval(c.Value, env)
		if value.Type() != subject.Type() {
			return object.NewError("mismatched types in switch statement: %s %s",
				subject.Type(), value.Type())
		}
		test := evalInfixExpression("==", subject, value)
		if test == object.TRUE {
			return Eval(c.Block, env)
		}
	}

	if se.Default != nil {
		return Eval(se.Default, env)
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

func evalIdentifier(
	node *ast.Identifier,
	env *object.Environment,
) object.Object {
	if val, ok := env.Get(node.Value); ok {
		if lazyObj, ok := val.(*object.LazyObject); ok {
			unwrappedVal := Eval(lazyObj.Right, env)
			env.Assign(node.Value, unwrappedVal)
			return unwrappedVal
		}

		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return object.NewError("identifier not found: " + node.Value)
}

func applyFunction(fn object.Object, args []object.Object, env *object.Environment) object.Object {
	switch fn := fn.(type) {

	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)

	case *object.Method:
		newEnv := createMethodEnv(fn, args, env)
		evaluated := Eval(fn.StructMethod.FunctionLiteral.Body, newEnv)

		if err := handleEvolve(fn.StructInstance, env); err != nil {
			return err
		}

		return unwrapReturnValue(evaluated)

	case *object.Builtin:
		return fn.Fn(args...)

	default:
		return object.NewError("not a function: %s", fn.Type())
	}
}

func handleEvolve(instance *object.StructInstance, env *object.Environment) object.Object {
	if instance.IsMethod("evolve") {
		evolveMethod := instance.GetMethod("evolve").(*object.Method)
		newEnv := createMethodEnv(evolveMethod, []object.Object{}, env)
		other := Eval(evolveMethod.StructMethod.FunctionLiteral.Body, newEnv)
		other = unwrapReturnValue(other)
		if other.Type() != object.NULL_OBJ {
			new, ok := other.(*object.StructInstance)
			if !ok {
				return object.NewError("evolve method must return nil or a nac instance, returned %s: %s", other.Type(), other.Inspect())
			}
			instance.EvolveInto(new)
		}
	}

	return nil
}

func createMethodEnv(
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

func extendFunctionEnv(
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

func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.HASH_OBJ:
		return evalHashIndexExpression(left, index)
	default:
		return object.NewError("index operator not supported: %s", left.Type())
	}
}

func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	if idx < 0 || idx > max {
		return object.NULL
	}

	return arrayObject.Elements[idx]
}

func evalHashLiteral(
	node *ast.HashLiteral,
	env *object.Environment,
) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env)
		if isError(key) {
			return key
		}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return object.NewError("unusable as hash key: %s", key.Type())
		}

		value := Eval(valueNode, env)
		if isError(value) {
			return value
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}

	return &object.Hash{Pairs: pairs}
}

func evalHashIndexExpression(hash, index object.Object) object.Object {
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

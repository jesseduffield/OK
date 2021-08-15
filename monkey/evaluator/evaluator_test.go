package evaluator

import (
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func testEval(t *testing.T, input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d",
			result.Value, expected)
		return false
	}

	return true
}

func testErrorObject(t *testing.T, obj object.Object, expected string) bool {
	result, ok := obj.(*object.Error)
	if !ok {
		t.Errorf("object is not Error. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Message != expected {
		t.Errorf("object has wrong value. got=%s, want=%s",
			result.Message, expected)
		return false
	}

	return true
}

func testStringObject(t *testing.T, obj object.Object, expected string) bool {
	result, ok := obj.(*object.String)
	if !ok {
		t.Errorf("object is not String. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%s, want=%s",
			result.Value, expected)
		return false
	}

	return true
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
		{"\"a\" == \"a\"", true},
		{"\"a\" != \"b\"", true},
		{"true && false", false},
		{"true && true", true},
		{"false && true", false},
		{"false && false", false},
		{"true || false", true},
		{"true || true", true},
		{"false || true", true},
		{"false || false", false},
		{"false || false || true", true},
		{"true && true && true", true},
		{"3 > 2 && 3 < 4", true},
		{"3 > 4 || 3 < 4", true},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t",
			result.Value, expected)
		return false
	}
	return true
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != object.NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + true;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true",
			"unknown operator: -BOOLEAN",
		},
		{
			"true + false;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { true + false; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`
if (10 > 1) {
if (10 > 1) {
	return true + false;
}

return 1;
}
`,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar",
			"identifier not found: foobar",
		},
		{
			`"Hello" - "World"`,
			"unknown operator: STRING - STRING",
		},
		{
			`{"name": "Monkey"}[fn(x) { x }];`,
			"unusable as hash key: FUNCTION",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)

		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T(%+v)",
				evaluated, evaluated)
			continue
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q",
				tt.expectedMessage, errObj.Message)
		}
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(t, tt.input), tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 2; };"

	evaluated := testEval(t, input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v",
			fn.Parameters)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0])
	}

	expectedBody := "(x + 2)"

	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = fn(x) { x; }; identity(5);", 5},
		{"let identity = fn(x) { return x; }; identity(5);", 5},
		{"let double = fn(x) { x * 2; }; double(5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5, 5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"fn(x) { x; }(5)", 5},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(t, tt.input), tt.expected)
	}
}

func TestClosures(t *testing.T) {
	input := `
let newAdder = fn(x) {
fn(y) { x + y };
};

let addTwo = newAdder(2);
addTwo(2);`

	testIntegerObject(t, testEval(t, input), 4)
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello World!"`

	evaluated := testEval(t, input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + " " + "World!"`

	evaluated := testEval(t, input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, "argument to `len` not supported, got INTEGER"},
		{`len("one", "two")`, "wrong number of arguments. got=2, want=1"},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)",
					evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message. expected=%q, got=%q",
					expected, errObj.Message)
			}
		}
	}
}

func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	evaluated := testEval(t, input)
	result, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
	}

	if len(result.Elements) != 3 {
		t.Fatalf("array has wrong num of elements. got=%d",
			len(result.Elements))
	}

	testIntegerObject(t, result.Elements[0], 1)
	testIntegerObject(t, result.Elements[1], 4)
	testIntegerObject(t, result.Elements[2], 6)
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			"[1, 2, 3][0]",
			1,
		},
		{
			"[1, 2, 3][1]",
			2,
		},
		{
			"[1, 2, 3][2]",
			3,
		},
		{
			"let i = 0; [1][i];",
			1,
		},
		{
			"[1, 2, 3][1 + 1];",
			3,
		},
		{
			"let myArray = [1, 2, 3]; myArray[2];",
			3,
		},
		{
			"let myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];",
			6,
		},
		{
			"let myArray = [1, 2, 3]; let i = myArray[0]; myArray[i]",
			2,
		},
		{
			"[1, 2, 3][3]",
			nil,
		},
		{
			"[1, 2, 3][-1]",
			nil,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestHashLiterals(t *testing.T) {
	input := `let two = "two";
	{
			"one": 10 - 9,
			two: 1 + 1,
			"thr" + "ee": 6 / 2,
			4: 4,
			true: 5,
			false: 6
	}`

	evaluated := testEval(t, input)
	result, ok := evaluated.(*object.Hash)
	if !ok {
		t.Fatalf("Eval didn't return Hash. got=%T (%+v)", evaluated, evaluated)
	}

	expected := map[object.HashKey]int64{
		(&object.String{Value: "one"}).HashKey():   1,
		(&object.String{Value: "two"}).HashKey():   2,
		(&object.String{Value: "three"}).HashKey(): 3,
		(&object.Integer{Value: 4}).HashKey():      4,
		object.TRUE.HashKey():                      5,
		object.FALSE.HashKey():                     6,
	}

	if len(result.Pairs) != len(expected) {
		t.Fatalf("Hash has wrong num of pairs. got=%d", len(result.Pairs))
	}

	for expectedKey, expectedValue := range expected {
		pair, ok := result.Pairs[expectedKey]
		if !ok {
			t.Errorf("no pair for given key in Pairs")
		}

		testIntegerObject(t, pair.Value, expectedValue)
	}
}

func TestHashIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`{"foo": 5}["foo"]`,
			5,
		},
		{
			`{"foo": 5}["bar"]`,
			nil,
		},
		{
			`let key = "foo"; {"foo": 5}[key]`,
			5,
		},
		{
			`{}["foo"]`,
			nil,
		},
		{
			`{5: 5}[5]`,
			5,
		},
		{
			`{true: 5}[true]`,
			5,
		},
		{
			`{false: 5}[false]`,
			5,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestSwitchExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`switch 5 { case 5: 1 }`,
			1,
		},
		{
			`switch 5 { case 4: 1 }`,
			nil,
		},
		{
			`switch 5 { case 4: 1; default: 12 }`,
			12,
		},
		{
			`switch 6 { case 4: 1; case 6: 2; default: 12 }`,
			2,
		},
		{
			`switch 6 { case 4 + 2: 1; case 6: 2; default: 12 }`,
			1,
		},
		{
			`switch 1+5 { case 4 + 2: 1; case 6: 2; default: 12 }`,
			1,
		},
		{
			`switch "a" { case "a": "a"; default: "b" }`,
			"a",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		switch v := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(v))
		case string:
			testStringObject(t, evaluated, v)
		default:
			testNullObject(t, evaluated)
		}
	}
}

func TestAssignment(t *testing.T) {
	tests := []struct {
		input       string
		expected    interface{}
		expectedErr string
	}{
		{
			`let x = 1; x = 2; x`,
			2,
			"",
		},
		{
			`let x = 1; x = 2; x = x + 1`,
			3,
			"",
		},
		{
			`let x = [1,2]; x[0] = 2; x[0]`,
			2,
			"",
		},
		{
			`let x = {"one":1,"two":2}; x["one"] = 2; x["one"]`,
			2,
			"",
		},
		{
			`let x = [{"one":1}]; x[0]["one"] = 2; x[0]["one"]`,
			2,
			"",
		},
		{
			`let x = [[1],[2]]; x[1][0] = 3; x[1][0]`,
			3,
			"",
		},
		{
			`let x = [0]; x[1] = 1;`,
			nil,
			"Index 1 is out of bounds (array length 1)",
		},
		{
			`let x = [0]; x[-1] = 1;`,
			nil,
			"Index must be positive",
		},
		{
			`let x = [0]; x["1"] = 1;`,
			nil,
			"Index must be an integer",
		},
		{
			`let x = {}; x["a"]["b"] = 2`,
			nil,
			"Attempted index of NULL object",
		},
		{
			`let foo = fn() { return 1 }; foo()["a"] = 2`,
			nil,
			"`foo()` is neither a hash nor array so you cannot index into it",
		},
		{
			`x = 1`,
			nil,
			"x has not been declared",
		},
		{
			`notAClass Person { public foo fn() { return 5; } };let x = new Person(); x.foo()`,
			5,
			"",
		},
		{
			`notAClass Person { field email public getEmail fn(selfish) { return selfish.email } };let x = new Person(); x.getEmail()`,
			nil,
			"",
		},
		{
			`
			notAClass Person {
				field email
				public getEmail fn(selfish) { return selfish.email }
				public setEmail fn(selfish, value) { selfish.email = value }
			}

			let x = new Person();
			x.setEmail("test")
			x.getEmail()`,
			"test",
			"",
		},
		{
			`
			notAClass Person {
				public add fn(a, b) { return a + b }
			}

			let x = new Person();
			x.add(1, 2)`,
			3,
			"",
		},
		{
			`
			notAClass Person {
				add fn(a, b) { return a + b }
			}

			let x = new Person();
			x.add(1, 2)`,
			nil,
			"`add` is a private method on nac Person",
		},
		{
			`
			notAClass Person {
				field email

				public foo fn(selfish) {
					selfish.email = "haha"
					let cl = fn() { return selfish.email }
					return cl
				}
			}

			let x = new Person();
			let cl = x.foo()
			cl()`,
			"haha",
			"",
		},
		{
			`
			notAClass Person {
				field email
			}

			let x = new Person();
			x.add(1, 2)`,
			nil,
			"undefined field for nac Person: add",
		},
		{
			`
			notAClass Person {
				field email
			}

			let x = new Person();
			x.email = "test"`,
			nil,
			"`email` is a private field on nac Person",
		},
		{
			`
			notAClass Person {
				public foo fn() { return 5 }
			}

			let x = new Person();
			x.foo = "test"`,
			nil,
			"`foo` is a method, not a field, on nac Person. You cannot reassign it",
		},
		{
			`
			let divide = fn(a, b) {
				return switch b {
					case 0: [nil, "cannot divide by zero"];
					default: [a / b, ""];
				}
			};
			let result = divide(5, 0);
			let x = switch result[1] {
				case "": result[0]
				default: result[1]
			};
			x`,
			"cannot divide by zero",
			"",
		},
		{
			`
			switch true { case true: 1;2;3; case false: 2;3;4 }`,
			1,
			"",
		},
		{
			`
			switch true { case true: ; case false: 2; }`,
			nil,
			"",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		if tt.expectedErr != "" {
			testErrorObject(t, evaluated, tt.expectedErr)
			continue
		}
		switch v := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(v))
		case string:
			testStringObject(t, evaluated, v)
		default:
			testNullObject(t, evaluated)
		}
	}
}

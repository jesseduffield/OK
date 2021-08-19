package evaluator

import (
	"testing"

	"github.com/jesseduffield/OK/ok/lexer"
	"github.com/jesseduffield/OK/ok/object"
	"github.com/jesseduffield/OK/ok/parser"
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
		{"5 + 2 * 10", 70},
		{"20 + 2 * -10", -220},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 20},
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
		{"1 >= 2", false},
		{"2 >= 1", true},
		{"NO! >= NO!", true},
		{"true >= true", true},
		{"false >= false", true},
		{"true >= false", false},
		{"(1 >= 2) >= true", false},
		{"(1 >= 2) >= false", true},
		{"\"a\" >= \"a\"", true},
		{"\"a\" >= \"b\"", false},
		{"\"b\" >= \"a\"", true},
		{"let x = true; let y = false; x || y", true},
		{"let x = true; let y = false; x || x", true},
		{"let x = true; let y = false; y || x", true},
		{"let x = true; let y = false; y || y", false},
		{"let x = false; let y = false; let z = true; x || y || z", true},
		{"let x = true; let y = true; let z = true; x && y && z", true},
		{"let x = 3 >= 2; let y = 5 >= 4; x && y", true},
		{"switch 1 { case 2: true; default: false }", false},
		{"switch 1 { case 2: false; default: true }", true},
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
		{"if (1 >= 2) { 10 }", nil},
		{"if (1 >= 2) { 10 } else { 20 }", 20},
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
			"if (10 >= 1) { true + false; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`
if (10 >= 1) {
if (10 >= 1) {
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
			`{"name": "OK"}[fn(x) { x }];`,
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
let newadder = fn(x) {
fn(y) { x + y };
};

let addtwo = newadder(2);
addtwo(2);`

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
			"let myarray = [1, 2, 3]; myarray[2];",
			3,
		},
		{
			"let myarray = [1, 2, 3]; myarray[0] + myarray[1] + myarray[2];",
			6,
		},
		{
			"let myarray = [1, 2, 3]; let i = myarray[0]; myarray[i]",
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
			`notaclass person { public foo fn() { return 5; } };let x = new person(); x.foo()`,
			5,
			"",
		},
		{
			`notaclass person { field email public getemail fn(selfish) { return selfish.email } };let x = new person(); x.getemail()`,
			nil,
			"",
		},
		{
			`
			notaclass person {
				field email
				public getemail fn(selfish) { return selfish.email }
				public setemail fn(selfish, value) { selfish.email = value }
			}

			let x = new person();
			x.setemail("test")
			x.getemail()`,
			"test",
			"",
		},
		{
			`
			notaclass person {
				public add fn(a, b) { return a + b }
			}

			let x = new person();
			x.add(1, 2)`,
			3,
			"",
		},
		{
			`
			notaclass person {
				add fn(a, b) { return a + b }
			}

			let x = new person();
			x.add(1, 2)`,
			nil,
			"`add` is a private method on nac person",
		},
		{
			`
			notaclass person {
				field email

				public foo fn(selfish) {
					selfish.email = "haha"
					let cl = fn() { return selfish.email }
					return cl
				}
			}

			let x = new person();
			let cl = x.foo()
			cl()`,
			"haha",
			"",
		},
		{
			`
			notaclass person {
				field email
			}

			let x = new person();
			x.add(1, 2)`,
			nil,
			"undefined field for nac person: add",
		},
		{
			`
			notaclass person {
				field email
			}

			let x = new person();
			x.email = "test"`,
			nil,
			"`email` is a private field on nac person",
		},
		{
			`
			let x = NO!;
			ayok?(x);
			`,
			false,
			"",
		},
		{
			`
			let x = 10;
			ayok?(x);
			`,
			true,
			"",
		},
		{
			`
			NO! >= NO!
			`,
			true,
			"",
		},
		{
			`
			NO! >= 10
			`,
			false,
			"",
		},
		{
			`
			notaclass person {
				pack "this is bad"

				field email
			}

			let x = new person();

			// I acknowledge that this is bad
			x.email = "test";
			x.email;`,
			"test",
			"",
		},
		{
			`
			notaclass person {
				pack "this is bad"

				field email
			}

			let x = new person();

			// I do not acknowledge that this is bad
			x.email = "test";
			x.email;`,
			nil,
			"`email` is a private field on nac person",
		},
		{
			`
			notaclass person {
				public foo fn() { return 5 }
			}

			let x = new person();
			x.foo = "test"`,
			nil,
			"`foo` is a method, not a field, on nac person. You cannot reassign it",
		},
		{
			`
			let divide = fn(a, b) {
				return switch b {
					case 0: [NO!, "cannot divide by zero"];
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
			`switch true { case true: NO!; case false: 2; }`,
			nil,
			"",
		},
		{
			`let x = lazy 3; x`,
			3,
			"",
		},
		{
			`let x = lazy 3; x`,
			3,
			"",
		},
		{
			`
			let r = 0;
			let x = fn() { r = r + 1; return true };
			let y = fn() { r = r + 2; return true };
			let lx = lazy x();
			let ly = lazy y();
			lx || ly;
			r;
			`,
			1,
			"",
		},
		{
			`
			let r = 0;
			fn() { r = r + 1 }();
			r;
			`,
			1,
			"",
		},
		{
			`
			let r = 0;
			r = 5 + 5;
			r;
			`,
			10,
			"",
		},
		{
			`
			notaclass brgousie {
				public whoami fn(selfish) {
					return "a good-for-nothing aristocrat who likes classes"
				}
			};

			notaclass person {
				field name
				field email
				field likeclas

				public init fn(selfish, name, email) {
					selfish.name = name;
					selfish.email = email;
					selfish.likeclas = false;
				}

				public whoami fn(selfish) {
					return selfish.name;
				}

				public makeold fn(selfish) {
					selfish.likeclas = true;
				}

				evolve fn(selfish) {
					switch selfish.likeclas {
						case true:
							return new brgousie()
						default:
							return NO!
					}
				}
			};

			let p = new person();
			p.init("John", "")
			p.makeold();
			p.whoami();
			`,
			"a good-for-nothing aristocrat who likes classes",
			"",
		},
		{
			`
			let arr = [1,2,3];
			arr = map(arr, fn(e) { e * 2 }); // [2,4,6]
			arr[1]
			`,
			4,
			"",
		},
		{
			`
			let result = map([0,1], fn(e, i) {
				switch i {
				case 0:
					return 5 * 2;
				case 1:
					return 10 * 3;
				}
			})
			result[0]
			`,
			10,
			"",
		},
		{
			`
			let every = fn(arr, check) {
				let fail = false;
				map(arr, fn(e) {
					switch check(e) { case true: fail = true; } }
				)
				return !fail;
			};

			every([5,2,4,1,3], fn(e) { return e >= 2 })
			`,
			false,
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
		case bool:
			testBooleanObject(t, evaluated, v)
		default:
			testNullObject(t, evaluated)
		}
	}
}

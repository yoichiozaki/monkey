package compiler

import (
	"fmt"
	"monkey/ast"
	"monkey/code"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"
)

type compilerTestCase struct {
	input                string
	expectedConstants    []interface{}       // which we expect in constant pool
	expectedInstructions []code.Instructions // which we expect the compiler generate
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "1 + 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpAdd),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "1 - 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpSub),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "1 * 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpMul),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "2 / 1",
			expectedConstants: []interface{}{2, 1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpDiv),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "-1",
			expectedConstants: []interface{}{1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpMinus),
				code.Make(code.OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}

func runCompilerTests(t *testing.T, tests []compilerTestCase) {
	t.Helper()
	for _, tt := range tests {
		program := parse(tt.input)
		compiler := New()
		err := compiler.Compile(program)
		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}
		bytecode := compiler.Bytecode()
		err = testInstructions(tt.expectedInstructions, bytecode.Instructions)
		if err != nil {
			t.Fatalf("testInstructions failed: %s", err)
		}
		err = testConstants(tt.expectedConstants, bytecode.Constants)
		if err != nil {
			t.Fatalf("testConstants failed: %s", err)
		}
	}
}

func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}

func testInstructions(expected []code.Instructions, actual code.Instructions) error {
	concatted := concatInstructions(expected)
	if len(actual) != len(concatted) {
		return fmt.Errorf("wrong instructions length.\nwant=%q\ngot=%q", concatted, actual)
	}
	for i, ins := range concatted {
		if actual[i] != ins {
			return fmt.Errorf("wrong instruction at %d.\nwant=%q\ngot=%q", i, concatted, actual)
		}
	}
	return nil
}

func concatInstructions(s []code.Instructions) code.Instructions {
	out := code.Instructions{}
	for _, ins := range s {
		out = append(out, ins...)
	}
	return out
}

func testConstants(expected []interface{}, actual []object.Object) error {
	if len(expected) != len(actual) {
		return fmt.Errorf("wrong number of constants. got=%d, want=%d", len(actual), len(expected))
	}
	for i, constant := range expected {
		switch constant := constant.(type) {
		case int:
			err := testIntegerObject(int64(constant), actual[i])
			if err != nil {
				return fmt.Errorf("constant %d - testIntegerObject failed: %s", i, err)
			}
		case string:
			err := testStringObject(constant, actual[i])
			if err != nil {
				return fmt.Errorf("constant %d - testStringObject failed: %s", i, err)
			}
		}
	}
	return nil
}

func testIntegerObject(expected int64, actual object.Object) error {
	result, ok := actual.(*object.Integer)
	if !ok {
		return fmt.Errorf("object is not Integer. got=%T (%+v)", actual, actual)
	}
	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
	}
	return nil
}

func testStringObject(expected string, actual object.Object) error {
	result, ok := actual.(*object.String)
	if !ok {
		return fmt.Errorf("object is not String. got=%T (%+v)", actual, actual)
	}
	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%q, want=%q", result.Value, expected)
	}
	return nil
}

func TestBooleanExpressions(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "true",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "false",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpFalse),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "1 > 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpGreaterThan),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "1 < 2",
			expectedConstants: []interface{}{2, 1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpGreaterThan),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "1 == 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpEqual),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "1 != 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpNotEqual),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "true == false",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpFalse),
				code.Make(code.OpEqual),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "true != false",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpFalse),
				code.Make(code.OpNotEqual),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "!true",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpBang),
				code.Make(code.OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}

func TestConditionals(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
				if(true) { 
					10;
				}
				3333;
`,
			expectedConstants: []interface{}{10, 3333},
			expectedInstructions: []code.Instructions{

				// +======================+=================+
				// | ADDRESS in byte size |  INSTRUCTIONS   |
				// +======================+=================+
				// |                 0000 | OpTrue          |
				// +----------------------+-----------------+
				// |                 0001 | OpJumpNotTruthy |+
				// +----------------------+-----------------+|
				// |                 0002 | 00000000        |+-----+
				// +----------------------+-----------------+|     |
				// |                 0003 | 00001010        |+     |
				// +----------------------+-----------------+      |
				// |                 0004 | OpConstant      |      |
				// +----------------------+-----------------+      |
				// |                 0005 | 00000000        |      |
				// +----------------------+-----------------+      |
				// |                 0006 | 00000000        |      |
				// +----------------------+-----------------+      |
				// |                 0007 | OpJump          |+     |
				// +----------------------+-----------------+|     |
				// |                 0008 | 00000000        |+-----------+
				// +----------------------+-----------------+|     |     |
				// |                 0009 | 00001011        |+     |     |
				// +----------------------+-----------------+      |     |
				// |                 0010 | OpNull          |<-----+     |
				// +----------------------+-----------------+            |
				// |                 0011 | OpPop           |<-----------+
				// +----------------------+-----------------+
				// |                 0012 | OpConstant      |
				// +----------------------+-----------------+
				// |                 0013 | 00000000        |
				// +----------------------+-----------------+
				// |                 0014 | 00000001        |
				// +----------------------+-----------------+
				// |                 0015 | OpPop           |
				// +----------------------+-----------------+

				// 0000
				code.Make(code.OpTrue),
				// 0001
				code.Make(code.OpJumpNotTruthy, 10),
				// 0004
				code.Make(code.OpConstant, 0),
				// 0007
				code.Make(code.OpJump, 11),
				// 0010
				code.Make(code.OpNull), // means no Alternative
				// 0011
				code.Make(code.OpPop),
				// 0012
				code.Make(code.OpConstant, 1),
				// 0015
				code.Make(code.OpPop),
			},
		},
		{
			input: `
				if(true) { 
					10;
				} else {
					20;
				}
				3333;
`,
			expectedConstants: []interface{}{10, 20, 3333},
			expectedInstructions: []code.Instructions{

				// +======================+=================+
				// | ADDRESS in byte size |  INSTRUCTIONS   |
				// +======================+=================+
				// |                 0000 | OpTrue          |
				// +----------------------+-----------------+
				// |                 0001 | OpJumpNotTruthy |+
				// +----------------------+-----------------+|
				// |                 0002 | 00000000        |+-----+
				// +----------------------+-----------------+|     |
				// |                 0003 | 00001010        |+     |
				// +----------------------+-----------------+      |
				// |                 0004 | OpConstant      |      |
				// +----------------------+-----------------+      |
				// |                 0005 | 00000000        |      |
				// +----------------------+-----------------+      |
				// |                 0006 | 00000000        |      |
				// +----------------------+-----------------+      |
				// |                 0007 | OpJump          |+     |
				// +----------------------+-----------------+|     |
				// |                 0008 | 00000000        |+-----------+
				// +----------------------+-----------------+|     |     |
				// |                 0009 | 00001011        |+     |     |
				// +----------------------+-----------------+      |     |
				// |                 0010 | OpNull          |<-----+     |
				// +----------------------+-----------------+            |
				// |                 0011 | OpPop           |<-----------+
				// +----------------------+-----------------+
				// |                 0012 | OpConstant      |
				// +----------------------+-----------------+
				// |                 0013 | 00000000        |
				// +----------------------+-----------------+
				// |                 0014 | 00000001        |
				// +----------------------+-----------------+
				// |                 0015 | OpPop           |
				// +----------------------+-----------------+

				// 0000
				code.Make(code.OpTrue),
				// 0001
				code.Make(code.OpJumpNotTruthy, 10),
				// 0004
				code.Make(code.OpConstant, 0),
				// 0007
				code.Make(code.OpJump, 13),
				// 0010
				code.Make(code.OpConstant, 1),
				// 0011
				code.Make(code.OpPop),
				// 0012
				code.Make(code.OpConstant, 2),
				// 0015
				code.Make(code.OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}

func TestGlobalLetStatements(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
					let one = 1;
					let two = 2;
`,
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpSetGlobal, 1),
			},
		},
		{
			input: `
					let one = 1;
					one;
`,
			expectedConstants: []interface{}{1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input: `
					let one = 1;
					let two = one;
					two;
`,
			expectedConstants: []interface{}{1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpSetGlobal, 1),
				code.Make(code.OpGetGlobal, 1),
				code.Make(code.OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}

// assertions about the Define method.
func TestDefine(t *testing.T) {
	expected := map[string]Symbol{
		"a": {Name: "a", Scope: GlobalScope, Index: 0},
		"b": {Name: "b", Scope: GlobalScope, Index: 1},
	}
	global := NewSymbolTable()
	a := global.Define("a")
	if a != expected["a"] {
		t.Errorf("expected a=%+v, got=%+v", expected["a"], a)
	}
	b := global.Define("b")
	if b != expected["b"] {
		t.Errorf("expected b=%+v, got=%v", expected["b"], b)
	}
}

// assertions about the Resolve method.
func TestResolveGlobal(t *testing.T) {
	global := NewSymbolTable()
	global.Define("a")
	global.Define("b")
	expected := []Symbol{
		{Name: "a", Scope: GlobalScope, Index: 0},
		{Name: "b", Scope: GlobalScope, Index: 1},
	}
	for _, sym := range expected {
		result, ok := global.Resolve(sym.Name)
		if !ok {
			t.Errorf("name %s not resolved", sym.Name)
			continue
		}
		if result != sym {
			t.Errorf("expected %s to resolve %+v, got=%+v", sym.Name, sym, result)
		}
	}
}

func TestStringExpressions(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             `"monkey"`,
			expectedConstants: []interface{}{"monkey"},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input:             `"mon" + "key"`,
			expectedConstants: []interface{}{"mon", "key"},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpAdd),
				code.Make(code.OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}

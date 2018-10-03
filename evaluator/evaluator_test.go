package evaluator

import (
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"
)

// Integerを正しく評価できているかをテスト
func TestEvalIntegerExpression(t *testing.T) {

	// テストセット
	tests := []struct {
		input    string
		expected int64
	}{
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

	// 各テストセットに対して
	for _, tt := range tests {

		// inputを評価して
		evaluated := testEval(tt.input)

		// 結果を確認
		testIntegerObject(t, evaluated, tt.expected)
	}
}

// 入力をレキサ・パーサに通して得られたASTをObjectに変換して返すヘルパー関数
func testEval(input string) object.Object {

	// 入力で初期化したレキサを生成
	l := lexer.New(input)

	// レキサをセットしたパーサを生成
	p := parser.New(l)

	// プログラムをパース
	program := p.ParseProgram()

	// パースした結果得られるASTを評価
	return Eval(program)
}

// 引数objがIntegerObject型で、かつ格納されている値が期待したものになっていることを確認するヘルパー関数
func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {

	// 引数がInteger型であることを確認
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T(%+v)", obj, obj)
		return false
	}

	// 格納してる値が期待したものになっていることを確認
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
		return false
	}

	return true
}

// Booleanを正しく評価できているかをテスト
func TestEvalBooleanExpression(t *testing.T) {

	// テストケース
	tests := []struct {
		input    string
		expected bool
	}{
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
	}

	// 各テストケースに対して
	for _, tt := range tests {

		// inputを評価して
		evaluated := testEval(tt.input)

		// 結果を確認
		testBooleanObject(t, evaluated, tt.expected)
	}
}

// 引数objが期待するBooleanObjectであることを確認するヘルパー関数
func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {

	// Boolean型であることを確認
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T(%+v)", obj, obj)
		return false
	}

	// 格納している値が期待したものであることを確認
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t", result.Value, expected)
		return false
	}

	return true
}

// !演算子の評価をテスト
func TestBangOperator(t *testing.T) {

	// テストケース
	// 5はtruthyに扱う
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

	// 各テストケースに対して
	for _, tt := range tests {

		// inputを評価して
		evaluated := testEval(tt.input)

		// 結果を確認
		testBooleanObject(t, evaluated, tt.expected)
	}
}

// If-Else式の評価をテスト
func TestIfElseExpressions(t *testing.T) {

	// テストケース
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

	// 各テストケースについて
	for _, tt := range tests {

		// inputを評価
		evaluated := testEval(tt.input)

		// 型アサーション
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}

	}
}

// 引数objがNullObjectであるかを確認するヘルパー関数
func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not Null. got=%T(%+v)", obj, obj)
		return false
	}
	return true
}

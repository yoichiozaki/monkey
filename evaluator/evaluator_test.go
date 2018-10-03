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
	}

	// 各テストセットに対して
	for _, tt := range tests {

		// inputを評価して
		evaluated := testEval(tt.input)

		// 結果を確認
		testIntegerObject(t, evaluated, tt.expected)
	}
}

// 入力をレキサ・パーサに通して得られたASTをObjectに変換して返す
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

// Object型の引数がInteger型で、かつ格納されている値が期待したものになっていることを確認するヘルパー関数
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
	}

	// 各テストケースに対して
	for _, tt := range tests {

		// inputを評価して
		evaluated := testEval(tt.input)

		// 結果を確認
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {

	// Boolean型であることを確認
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T(%+v)", obj, obj)
		return false
	}

	// 格納している値が期待したものであることを確認
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%", result.Value, expected)
		return false
	}

	return true
}

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

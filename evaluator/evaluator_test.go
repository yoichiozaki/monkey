package evaluator

import (
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"
)

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

		// 得られたObjectが期待したものであることを確認
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

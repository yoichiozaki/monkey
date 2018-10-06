package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"testing"
)

// LET文のパースをテストする
func TestLetStatements(t *testing.T) {
	input := `
let x = 5;
let y = 10;
let foobar = 838383;
`
	// inputで初期化されたレキサを生成
	l := lexer.New(input)

	// レキサをセットしたパーサを生成
	p := New(l)

	// プログラムをパースする
	program := p.ParseProgram()

	// パース中のエラーを出力
	checkParserErrors(t, p)

	// パースした結果得られるASTを調べていく
	// rootノードが存在するか
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	// 正しい個数だけStatementノードが生成されているか
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d",
			len(program.Statements))
	}

	// 期待する結果との照合
	tests := []struct {
		expectedIdentifier string // 期待する結果
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}
	for i, tt := range tests {
		stmt := program.Statements[i]

		// Statement型のASTノードたるstmtが期待したノードになっているのかを確認する
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

// Statement型のASTノードと期待する識別子名を引数に受けて、正しくLetStatement型のASTノードになっているかを確認するヘルパー関数
func testLetStatement(t *testing.T, s ast.Statement, name string) bool {

	// 引数であるStatementノードがLET文のものでなかったらダメなのでまずそのノードのToken.Literalを確認
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", s.TokenLiteral())
		return false
	}

	// LetStatement型にキャスト
	letStmt, ok := s.(*ast.LetStatement)

	// LetStatement型にキャストできないとダメ
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%t", s)
		return false
	}

	// そのASTノードに格納されている識別子名が期待したものでないとダメ
	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}

	// そのASTノードに格納されているトークンのリテラルが期待したものでないとダメ
	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral() not '%s. got=%s", name, letStmt.Name.TokenLiteral())
		return false
	}
	return true
}

// パース中にエラーを検知するとnilのASTノードが返ると同時にパーサ内のerrorsにエラーメッセージが追加される
// これを体裁を整えて出力する関数
func checkParserErrors(t *testing.T, p *Parser) {

	// パーサ内に記録されたエラーを取り出す
	errors := p.errors

	// エラーがなければそれで終了
	if len(errors) == 0 {
		return
	}

	// 各情報を出力
	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

// RETURN文のパースをテストする
func TestReturnStatements(t *testing.T) {
	input := `
return 5;
return 10;
return 993322;
`
	// inputで初期化されたレキサを生成
	l := lexer.New(input)

	// レキサをセットしたパーサを生成
	p := New(l)

	// プログラムをパースする
	program := p.ParseProgram()

	// パース中のエラーを出力
	checkParserErrors(t, p)

	// パースした結果得られるASTを調べる
	// 正しい個数だけStatementノードが生成されているか
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d",
			len(program.Statements))
	}

	// それぞれのStatementノードについて、
	// - その型がReturnStatement型であること
	// - そこに格納されているリテラルが正しいものになっていること
	// を確認する
	for _, stmt := range program.Statements {

		// ReturnStatement型のチェック
		returnStmt, ok := stmt.(*ast.ReturnStatement)

		// 型が正しくないのでダメ
		if !ok {
			t.Errorf("stmt not *ast.returnStatement. got=%T", stmt)
			continue
		}

		// リテラルの確認
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return', got %q",
				returnStmt.TokenLiteral())
		}
	}
}

// 式文としての識別子のパースをテスト
func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	// inputで初期化されたレキサを生成
	l := lexer.New(input)

	// レキサをセットしたパーサを生成
	p := New(l)

	// プログラムをパースする
	program := p.ParseProgram()

	// パース中のエラーを出力
	checkParserErrors(t, p)

	// パースして得られたASTに正しい個数のStatementノードが含まれるかをチェック
	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}

	// ASTに含まれるStatement型のASTノードがExpressionStatement型であることを確認
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	// そのASTノードにぶら下がっているExpression型のASTノードがIdentifier型であることを確認
	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T",
			stmt.Expression)
	}

	// ぶら下がっているIdentifier型のノードの格納している値が期待したものであるかを確認
	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %s. got=%s",
			"foobar", ident.Value)
	}

	// ぶら下がっているIdentifier型のノードの格納しているリテラルが期待したものであるかを確認
	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral is not %s. got=%s",
			"foobar", ident.TokenLiteral())
	}
}

// 式文としての整数リテラルのパースをテスト
func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	// inputによって初期化されたレキサを生成
	l := lexer.New(input)

	// レキサをセットしたパーサを生成
	p := New(l)

	// プログラムをパース
	program := p.ParseProgram()

	// パース中のエラーを出力
	checkParserErrors(t, p)

	// パースした結果得られるASTに正しい個数のStatementノードが含まれることを確認
	if len(program.Statements) != 1 {
		t.Fatalf("program has not ecnough statements. got=%d",
			len(program.Statements))
	}

	// StatementノードがExpressionStatement型のASTノードであることを確認
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	// ExpressionStatement型のASTノードにぶら下がっているノードがIntegerLiteral型であることを確認
	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp not *ast.IntegerLiteral. got=%T",
			stmt.Expression)
	}

	// ぶら下がっていたIntegerLiteral型のノードに格納されていた値が期待するものであるかを確認
	if literal.Value != 5 {
		t.Errorf("literal.Valur not %d. got=%d",
			5, literal.Value)
	}

	// ぶら下がっていたIntegerLiteral型のノードに格納されていたリテラルが期待するものであるかを確認
	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral not %s. got=%s",
			"5", literal.TokenLiteral())
	}
}

// 前置演算子式のASTノードのパースをテスト
func TestParsingPrefixExpressions(t *testing.T) {

	// テストセットを定義
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
		{"!true", "!", true},
		{"!false", "!", false},
	}

	// それぞれのテストセットに対して
	for _, tt := range prefixTests {

		// inputで初期化されたレキサを生成
		l := lexer.New(tt.input)

		// レキサをセットしたパーサを生成
		p := New(l)

		// プログラムをパースする
		program := p.ParseProgram()

		// パース中のエラーを出力
		checkParserErrors(t, p)

		// パースした結果得られるASTに正しい個数のStatementノードが含まれていることを確認
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d",
				1, len(program.Statements))
		}

		// Program型のノードにぶら下がっているStatement型のノードがExpressionStatement型であることを確認
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("progmram.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		// ExpressionStatement型のノードstmtがPrefixExpression型であることを確認
		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T",
				stmt.Expression)
		}

		// PrefixExpression型のstmtに格納されているOperatorが期待されたものであることを確認
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s",
				tt.operator, exp.Operator)
		}

		// 現在検討している前置演算子の右側にあるリテラルのトークンが正しいものになっているかを確認
		if !testLiteralExpression(t, exp.Right, tt.value) {
			return
		}
	}
}

// Expression型のASTノードである引数expが期待したリテラルであることを確認するヘルパー関数
func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {

	// 期待した型に応じて適切なヘルパー関数を呼び出す
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

// Expression型のASTノードである引数ilが期待するIntegerLiteral型のASTノードであることを確認するヘルパー関数
func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {

	// IntegerLiteral型であることを確認
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}

	// 格納している値について期待しているものであるかを確認
	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false

	}

	// リテラルの表現が期待しているものであることを確認
	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLitearl not %d. got=%s",
			value, integ.TokenLiteral())
		return false
	}
	return true
}

// 中置演算子式のASTノードのパースをテスト
func TestParsingInfixExpression(t *testing.T) {

	// テストセット定義
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5", 5, "+", 5},
		{"5 - 5", 5, "-", 5},
		{"5 * 5", 5, "*", 5},
		{"5 / 5", 5, "/", 5},
		{"5 > 5", 5, ">", 5},
		{"5 < 5", 5, "<", 5},
		{"5 == 5", 5, "==", 5},
		{"5 != 5", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	// それぞれのテストセットに対して
	for _, tt := range infixTests {

		// inputで初期化されたレキサを生成
		l := lexer.New(tt.input)

		// レキサをセットしたパーサを生成
		p := New(l)

		// プログラムをパース
		program := p.ParseProgram()

		// パース中のエラーを出力
		checkParserErrors(t, p)

		// パースの結果得られたASTのprogramノードに正しい個数のStatement型のノードがぶら下がっているkと尾を確認
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d",
				1, len(program.Statements))
		}

		// ぶら下がっているStatement型のノードがExpressionStatement型であることを確認
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		// ExpressionStatement型のstmtにぶら下がっているASTノードが期待するInfixExpression型のASTノードになっているかを確認
		if !testInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"(5 + 5) * 2 * (5 + 5)",
			"(((5 + 5) * 2) * (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b, c) + d",
			"((a + add(b, c)) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
		{
			"a * [1, 2, 3, 4][b * c] * d",
			"((a * ([1, 2, 3, 4][(b * c)])) * d)",
		},
		{
			"add(a * b[2], b[1], 2 * [1, 2][1])",
			"add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

// Expression型のASTノードを引数にとって、そのノードがIdentifier型のASTノードであるかを確認するヘルパー関数
func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {

	// Identifier型であることを確認
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Expression. got=%T", exp)
		return false
	}

	// 格納している値が期待したものとなっているかを確認
	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}

	// リテラルが期待したものとなっているかを確認
	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got=%s", value, ident.TokenLiteral())
		return false
	}
	return true
}

// Expression型のASTノードを引数にとって、そのノードがInfixExpression型のASTノードであるかを確認するヘルパー関数
func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {

	// InfixExpression型であることを確認
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.InfixExpression. got=%T(%s)", exp, exp)
		return false
	}

	// InfixExpression型のASTノードのLeftフィールドが期待したASTノードであることを確認
	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	// InfixExpression型のASTノードのOperatorフィールドが期待したものであることを確認
	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
		return false
	}

	// InfixExpression型のASTノードのRightフィールドが期待したASTノードであることを確認
	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}
	return true
}

// Boolean式のパースのテスト
func TestBooleanExpression(t *testing.T) {

	// テストセット
	tests := []struct {
		input           string
		expectedBoolean bool
	}{
		{"true;", true},
		{"false;", false},
	}

	// それぞれのテストセットに対して
	for _, tt := range tests {

		// inputで初期化されたレキサを生成
		l := lexer.New(tt.input)

		// レキサをセットしたパーサを生成
		p := New(l)

		// プログラムをパース
		program := p.ParseProgram()

		// パース中のエラーを出力
		checkParserErrors(t, p)

		// ノードprogramにぶら下がっているStatement型のASTノードが正しい個数かを確認
		if len(program.Statements) != 1 {
			t.Fatalf("program has not enough statements. got=%d",
				len(program.Statements))
		}

		// Statement型のASTノードにぶら下がっているノードがExpressionStatement型であることを確認
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		// ExpressionStatement型のASTノードにぶら下がっているノードがBoolean型であることを確認
		boolean, ok := stmt.Expression.(*ast.Boolean)
		if !ok {
			t.Fatalf("exp not *ast.Boolean. got=%T", stmt.Expression)
		}

		// ExpressionStatement型のASTノードにぶら下がっているノードに格納されている値が期待したものであるかを確認
		if boolean.Value != tt.expectedBoolean {
			t.Errorf("boolean.Value not %T. got=%T", tt.expectedBoolean, boolean.Value)
		}
	}
}

// Expression型のASTノードを引数にとって、そのノードがBoolean型のASTノードであるかを確認するヘルパー関数
func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {

	// Expression型の引数expがBoolean型であることを確認
	bo, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp not *ast.Boolean. got=%T", exp)
		return false
	}

	// 格納されている値が期待したものになっていることを確認
	if bo.Value != value {
		t.Errorf("bo.Value not %t. got=%t", value, bo.Value)
		return false
	}

	// リテラルが期待したものになっていることを確認
	if bo.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("bo.TokenLiteral not %t. got=%s", value, bo.TokenLiteral())
		return false
	}
	return true
}

// If式のパースをテスト
func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

	// 入力によって初期化されたレキサを生成
	l := lexer.New(input)

	// レキサをセットしたパーサを生成
	p := New(l)

	// プログラムをパース
	program := p.ParseProgram()

	// パース中のエラーを出力
	checkParserErrors(t, p)

	// ノードprogramにぶら下がっているStatement型のノードの個数が期待したものになっているかを確認
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	// ぶら下がっているノードがExpressionStatement型のASTノードであることを確認
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	// ExpressionStatement型のASTノードにぶら下がっているのがIfExpression型のASTノードであることを確認
	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T",
			stmt.Expression)
	}

	// ぶら下がっているIfExpression型のノードが期待したものになっているかを確認
	// まずCondition部に対して
	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	// 次にConsequence部に対して、ぶら下がっているノードの個数の確認
	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statement. got=%d\n",
			len(exp.Consequence.Statements))
	}

	// ぶら下がっているノードの型の確認
	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.Consequence.Statements[0])
	}

	// Consequence部のExpressionにぶら下がっているノードの型を確認
	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	// Alternativeにぶら下がっているノードの型を確認
	if exp.Alternative != nil {
		t.Errorf("exp.Alternative.Statements was not nil. got=%+v",
			exp.Alternative)
	}
}

// If-Else式のパースをテスト
func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	// inputによって初期化されたレキサを生成
	l := lexer.New(input)

	// レキサをセットしたパーサを生成
	p := New(l)

	// プログラムをパース
	program := p.ParseProgram()

	// パース中のエラーを出力
	checkParserErrors(t, p)

	// パースした結果得られたASTのルートであるprogramのStatementフィールドに正しい個数のノードがぶら下がっていることを確認する
	if len(program.Statements) != 1 {
		t.Fatalf("program.Body does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	// ぶら下がっているノードがExpressionStatement型であることを確認
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	// ぶら下がっているノードがIfExpression型であることを確認
	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", stmt.Expression)
	}

	// ぶら下がっているIfExpression型のノードのConditionフィールドに
	// 正しくパースされた結果のExpression型のノードがぶら下がっていることの確認
	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	// ぶら下がっているIfExpression型のノードのStatementフィールドに
	// 正しい個数のBlockStatement型のノードがぶら下がっていることを確認
	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n",
			len(exp.Consequence.Statements))
	}

	// BlockStatement型のノードにぶら下がっているノードがExpressionStatement型であることを確認
	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.Consequence.Statements[0])
	}

	// BlockStatement型のノードにぶら下がっているノードが正しくパースされて得られたものかを確認
	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	// IfExpression型のノードのAlternativeフィールドに正しい個数のBlockStatementがぶら下がっていることを確認
	if len(exp.Alternative.Statements) != 1 {
		t.Errorf("exp.Alternative.Statements does not contain 1 statements. got=%d\n",
			len(exp.Alternative.Statements))
	}

	// ぶら下がっているのがExpressionStatement型であることを確認
	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.Alternative.Statements[0])
	}

	// ぶら下がっているExpressionStatement型のノードが正しくパースされた結果であることを確認
	if !testIdentifier(t, alternative.Expression, "y") {
		return
	}
}

// 関数リテラルのパースをテスト
func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x, y) { x + y; }`

	// inputで初期化されたレキサを生成
	l := lexer.New(input)

	// レキサをセットしたパーサを生成
	p := New(l)

	// プログラムをパース
	program := p.ParseProgram()

	// パース中のエラーを出力
	checkParserErrors(t, p)

	// programにぶら下がっているStatement型のASTノードの個数を確認
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statement. got=%d\n",
			1, len(program.Statements))
	}

	// ぶら下がっているノードがExpressionStatement型であることを確認
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			stmt.Expression)
	}

	// ぶら下がっているノードがFunctionLiteral型であることを確認
	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T",
			stmt.Expression)
	}

	// FunctionLiteral型のノードにぶら下がっている引数を表すASTノードの個数が正しいことを確認
	if len(function.Parameters) != 2 {
		t.Fatalf("function literal parameters wrong. want 2, got=%d\n",
			len(function.Parameters))
	}

	// 引数を表すASTノードが正しくパースされた結果であるかを確認
	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")

	// FunctionLiteral型のノードにぶら下がっている関数本体を表すASTノードに
	// ぶら下がっているStatementノードが正しい個数であるかを確認
	if len(function.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements has not 1 statement. got=%d\n",
			len(function.Body.Statements))
	}

	// ノードfunctionのBodyフィールドにぶら下がっているノードがExpressionStatement型であることを確認
	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function body stmt is not ast.ExpressionStatement. got=%T",
			function.Body.Statements[0])
	}

	// ノードfunctionのBodyフィールドにぶら下がっているノードが正しくInfixExpressionとしてパースされたかを確認
	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

// 関数の引数リストを正しくパースできているかをテスト
func TestFunctionParameterParsing(t *testing.T) {

	// テストケース
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "fn() {};", expectedParams: []string{}},
		{input: "fn(x) {};", expectedParams: []string{"x"}},
		{input: "fn(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
	}

	// それぞれのテストに対して
	for _, tt := range tests {

		// inputで初期化したレキサを生成
		l := lexer.New(tt.input)

		// レキサをセットしたパーサを生成
		p := New(l)

		// プログラムをパース
		program := p.ParseProgram()

		// パース中のエラーを出力
		checkParserErrors(t, p)

		// 正しくパースされたかを確認していく
		stmt := program.Statements[0].(*ast.ExpressionStatement)
		function := stmt.Expression.(*ast.FunctionLiteral)

		// パラメータの個数が等しいか
		if len(function.Parameters) != len(tt.expectedParams) {
			t.Errorf("length parameters wrong. want %d, fot=%d\n",
				len(tt.expectedParams), len(function.Parameters))
		}

		// パラメータのリテラルが正しいか
		for i, ident := range tt.expectedParams {
			testLiteralExpression(t, function.Parameters[i], ident)
		}
	}
}

// 関数呼び出しをパースするテスト
func TestCallExpressionParsing(t *testing.T) {
	input := `add(1, 2 * 3, 4 + 5)`

	// inputで初期化されたレキサを生成
	l := lexer.New(input)

	// レキサをセットしたパーサを生成
	p := New(l)

	// プログラムをパース
	program := p.ParseProgram()

	// パース中のエラーを出力
	checkParserErrors(t, p)

	// programノードに格納されている文の個数を確認
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statement. got=%d\n",
			1, len(program.Statements))
	}

	// 格納されている文がExpressionStatementであることを確認
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	// 格納されていた文がCallExpressionであることを確認
	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T",
			stmt.Expression)
	}

	// 格納されていたCallExpression型のASTノードのFunctionフィールドに
	// 正しいIdentifier型のASTノードが格納されているかを確認
	if !testIdentifier(t, exp.Function, "add") {
		return
	}

	// 格納されていたCallExpression型のASTノードのArgumentsフィールドに
	// 正しい個数の実引数を表すASTノードが格納されていることを確認
	if len(exp.Arguments) != 3 {
		t.Fatalf("wrong length of arguments. got=%d",
			len(exp.Arguments))
	}

	// 実引数それぞれのリテラルが期待したものであるかを確認
	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}

// Let文を正しくパースできるかをテスト
func TestLetStatement(t *testing.T) {

	// テストケース
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}

	// 各テストに対して
	for _, tt := range tests {

		// inputで初期化されたレキサを生成
		l := lexer.New(tt.input)

		// レキサをセットしたパーサを生成
		p := New(l)

		// プログラムをパース
		program := p.ParseProgram()

		// パース中のエラーを出力
		checkParserErrors(t, p)

		// パースした結果が正しいかを検証
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement. got=%d",
				len(program.Statements))
		}

		stmt := program.Statements[0]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}

		val := stmt.(*ast.LetStatement).Value
		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

// StringLiteralを正しくパースできるかテスト
func TestStringLiteralExpression(t *testing.T) {
	input := `"hello world"`

	// inputで初期化されたレキサを生成
	l := lexer.New(input)

	// レキサをセットしたパーサを生成
	p := New(l)

	// プログラムをパース
	program := p.ParseProgram()

	// パース中のエラーを出力
	checkParserErrors(t, p)

	// 文字列は式文なので
	// パースして得られたASTにぶら下がっているStatementをExpressionStatementにキャスト
	stmt := program.Statements[0].(*ast.ExpressionStatement)

	// stmtにぶら下がっているExpressionがStringLiteral型であるかを確認
	literal, ok := stmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("exp not *ast.StringLiteral. got=%T", stmt.Expression)
	}

	// stmtにぶら下がっているExpressionに格納されているリテラルが正しいかを確認
	if literal.Value != "hello world" {
		t.Errorf("literal.Value not %q. got=%q",
			"hello world", literal.Value)
	}
}

// 配列リテラルを正しくパースできるかをテスト
func TestParsingArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	// inputで初期化されたレキサを生成
	l := lexer.New(input)

	// レキサをセットしたパーサを生成
	p := New(l)

	// プログラムをパース
	program := p.ParseProgram()

	// パース中のエラーを出力
	checkParserErrors(t, p)

	// 正しい型のASTノードが得られたかを確認
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	array, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("exp not astArrayLiteral. gpt=%T", stmt.Expression)
	}
	if len(array.Elements) != 3 {
		t.Fatalf("len(array.Elements) not 3. got=%d", len(array.Elements))
	}

	// ぶら下がっているASTノードが正しいものであるかを確認
	testIntegerLiteral(t, array.Elements[0], 1)
	testInfixExpression(t, array.Elements[1], 2, "*", 2)
	testInfixExpression(t, array.Elements[2], 3, "+", 3)
}

// IndexExpressionを正しくパースできるかをテスト
func TestParsingIndexExpression(t *testing.T) {
	input := "myArray[1 + 1]"

	// inputで初期化されたレキサを生成
	l := lexer.New(input)

	// レキサをセットしたパーサを生成
	p := New(l)

	// プログラムをパース
	program := p.ParseProgram()

	// パース中のエラーを出力
	checkParserErrors(t, p)

	// 正しい型のASTノードが得られたかを確認
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	indexExp, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("exp not ast.IndexExpression. got=%T", stmt.Expression)
	}
	if !testIdentifier(t, indexExp.Left, "myArray") {
		return
	}
	if !testInfixExpression(t, indexExp.Index, 1, "+", 1) {
		return
	}
}

// 文字列キーのハッシュリテラルを正しくパースできるかのテスト
func TestParsingHashLiteralsStringKeys(t *testing.T) {
	input := `{"one": 1, "two": 2, "three": 3}`

	// inputで初期化されたレキサを生成
	l := lexer.New(input)

	// レキサをセットしたパーサを生成
	p := New(l)

	// プログラムをパース
	program := p.ParseProgram()

	// パース中のエラーを出力
	checkParserErrors(t, p)

	// 正しい型のASTノードが得られたかを確認
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp not ast.HashLiteral. got=%T",
			stmt.Expression)
	}
	if len(hash.Pairs) != 3 {
		t.Errorf("hash.Paris has wrong length. got=%d", len(hash.Pairs))
	}

	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral. got=%T", key)
		}
		expectedValue := expected[literal.String()]
		testIntegerLiteral(t, value, expectedValue)
	}
}

// 空のハッシュリテラルを正しくパースできるかをテスト
func TestParsingEmptyHashLiteral(t *testing.T) {
	input := "{}"

	// inputで初期化されたレキサを生成
	l := lexer.New(input)

	// レキサをセットしたパーサを生成
	p := New(l)

	// プログラムをパース
	program := p.ParseProgram()

	// パース中のエラーを出力
	checkParserErrors(t, p)

	// 正しい型のASTノードが得られたかを確認
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp not ast.HashLiteral. got=%T",
			stmt.Expression)
	}
	if len(hash.Pairs) != 0 {
		t.Errorf("hash.Paris has wrong length. got=%d", len(hash.Pairs))
	}
}

// ハッシュリテラルの値が任意の指揮をとりうることを確認するテスト
func TestParsingHashLiteralsWithExpressions(t *testing.T) {
	input := `{"one": 0 + 1, "two": 10 - 8, "three": 15 / 5}`

	// inputで初期化されたレキサを生成
	l := lexer.New(input)

	// レキサをセットしたパーサを生成
	p := New(l)

	// プログラムをパース
	program := p.ParseProgram()

	// パース中のエラーを出力
	checkParserErrors(t, p)

	// 正しい型のASTノードが得られたかを確認
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp not ast.HashLiteral. got=%T",
			stmt.Expression)
	}
	if len(hash.Pairs) != 3 {
		t.Errorf("hash.Paris has wrong length. got=%d", len(hash.Pairs))
	}

	tests := map[string]func(ast.Expression){
		"one": func(e ast.Expression) {
			testInfixExpression(t, e, 0, "+", 1)
		},
		"two": func(e ast.Expression) {
			testInfixExpression(t, e, 10, "-", 8)
		},
		"three": func(e ast.Expression) {
			testInfixExpression(t, e, 15, "/", 5)
		},
	}

	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral. got=%T", key)
			continue
		}
		testFunc, ok := tests[literal.String()]
		if !ok {
			t.Errorf("No test function for key %q found", literal.String())
			continue
		}
		testFunc(value)
	}
}

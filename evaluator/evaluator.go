package evaluator

import (
	"monkey/ast"
	"monkey/object"
)

// ast.Node型を受け取り評価して、適切なobject.Objectを返す
func Eval(node ast.Node) object.Object {

	// 引数nodeの型によって処理を振り分ける
	switch node := node.(type) {

	// 文だった
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	// 式だった
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	}
	return nil
}

// プログラム内のすべての式を評価するヘルパー関数
func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object
	for _, statement := range stmts {
		result = Eval(statement)
	}
	return result
}

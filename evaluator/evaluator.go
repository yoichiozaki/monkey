package evaluator

import (
	"monkey/ast"
	"monkey/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

// ast.Node型を受け取り評価して、適切なobject.Objectを返す
func Eval(node ast.Node) object.Object {

	// 引数nodeの型によって処理を振り分ける
	switch node := node.(type) {

	// 文だった
	case *ast.Program:
		return evalProgram(node)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.BlockStatement:
		return evalBlockStatement(node)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue)
		return &object.ReturnValue{Value: val}

	// 式だった
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalInfixExpression(node.Operator, left, right)
	case *ast.IfExpression:
		return evalIfExpression(node)
	}

	return nil
}

// // プログラムやブロック内のすべての式を評価するヘルパー関数
// func evalStatements(stmts []ast.Statement) object.Object {
// 	var result object.Object
// 	for _, statement := range stmts {
// 		result = Eval(statement)
//
// 		// returnに出くわしたら今評価した値で処理を中断する
// 		if returnValue, ok := result.(*object.ReturnValue); ok {
// 			return returnValue.Value
// 		}
// 	}
// 	return result
// }

// bool値に対して適切なBooleanオブジェクトを返す
func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

// operatorがサポート対象の演算子であることを確認するヘルパー関数
func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return NULL
	}
}

// 演算子!を評価して適切なObjectを返すヘルパー関数
// この関数が!の挙動を決定している
func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

// 演算子-を評価して適切なObjectを返すヘルパー関数
func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return NULL
	}
	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

// 中置式を構成するオペランドに応じて適切な評価関数へ処理を振り分けるヘルパー関数
func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	default:
		return NULL
	}
}

// 整数による中置式を評価してObjectを返すヘルパー関数
func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
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
		return NULL
	}
}

// IfExpression型のASTノードを引数にとって評価して適切なObjectを返すヘルパー関数
func evalIfExpression(ie *ast.IfExpression) object.Object {
	condition := Eval(ie.Condition)
	if isTruthy(condition) {
		return Eval(ie.Consequence)
	} else if Eval(ie.Alternative) != nil {
		return Eval(ie.Alternative)
	} else {
		return NULL
	}
}

// 引数objがTruthyであるかを確認するヘルパー関数
func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

// プログラムを評価してObjectを返すヘルパー関数
func evalProgram(program *ast.Program) object.Object {
	var result object.Object
	for _, statement := range program.Statements {

		// プログラムを構成する一文一文を一つずつ評価していく
		result = Eval(statement)

		// 評価した結果得られたObjectがReturnValue型であったならばそれを返す
		if returnValue, ok := result.(*object.ReturnValue); ok {
			return returnValue.Value
		}
	}
	return result
}

// ブロック文を評価してObjectを返すヘルパー関数
func evalBlockStatement(block *ast.BlockStatement) object.Object {
	var result object.Object

	// ブロックに含まれている各文を評価していく
	for _, statement := range block.Statements {
		result = Eval(statement)

		// 評価した結果がReturnValue型であったならばあんラップせずにresultを返す
		if result != nil && result.Type() == object.RETURN_VALUE_OBJ {
			return result
		}
	}
	return result
}

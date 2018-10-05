package evaluator

import (
	"fmt"
	"monkey/ast"
	"monkey/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

// ast.Node型を受け取り評価して、適切なobject.Objectを返す
func Eval(node ast.Node, env *object.Environment) object.Object {

	// 引数nodeの型によって処理を振り分ける
	switch node := node.(type) {

	// 文だった
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	// 式だった
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.Identifier:
		return evalIdentifier(node, env)
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
	case "!": // 演算子!を評価するヘルパー関数に処理を譲渡
		return evalBangOperatorExpression(right)
	case "-": // 演算子-を評価するヘルパー関数に処理を譲渡
		return evalMinusPrefixOperatorExpression(right)
	default: // サポートしていない演算子に遭遇したらErrorObjectを返す
		return newError("unknown operator: %s%s", operator, right.Type())
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

	// 演算子-のサポートしていない型に対して作用させようとしているときにはErrorObjectを返す
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
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
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s",
			left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
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
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

// IfExpression型のASTノードを引数にとって評価して適切なObjectを返すヘルパー関数
func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}
	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if Eval(ie.Alternative, env) != nil {
		return Eval(ie.Alternative, env)
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
func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object
	for _, statement := range program.Statements {

		// プログラムを構成する一文一文を一つずつ評価していく
		result = Eval(statement, env)

		// 評価した結果得られたObjectがReturnValue型であったならばそれを返す
		switch result := result.(type) {
		case *object.ReturnValue: // 評価した結果得られたObjectがError型であったならばそれを返す
			return result.Value
		case *object.Error: // 評価した結果得られたObjectがReturnValue型であったならばそれを返す
			return result
		}
	}
	return result
}

// ブロック文を評価してObjectを返すヘルパー関数
func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	// ブロックに含まれている各文を評価していく
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

// フォーマットと内容を引数にエラーメッセージを格納したErrorObjectを返すヘルパー関数
func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

// 引数objがError型であるかを確認するヘルパー関数
func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

// Identifier型のASTノードを引数に環境内に登録されている対応するObjectを返すヘルパー関数
func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	val, ok := env.Get(node.Value)
	if !ok {
		return newError("identifier not found: " + node.Value)
	}
	return val
}

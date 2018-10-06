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
	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Body: body, Env: env}
	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args)
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}
	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)
	case *ast.HashLiteral:
		return evalHashLiteral(node, env)
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
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
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
	} else if ie.Alternative != nil {
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
		case *object.ReturnValue: // 評価した結果得られたObjectがReturnValue型であったならばそれを返す
			return result.Value
		case *object.Error: // 評価した結果得られたObjectがError型であったならばそれを返す
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
	if val, ok := env.Get(node.Value); ok {
		return val
	}
	if builtin, ok := builtin[node.Value]; ok {
		return builtin
	}
	return newError("identifier not found: " + node.Value)
}

// 一連の式を評価し適切なオブジェクトのスライスを返すヘルパー関数
func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {

	// 返すObjectのスライス
	var result []object.Object

	// 各式に対してい
	for _, e := range exps {

		// 評価しObjectを得る
		evaluated := Eval(e, env)

		// エラーが起きたらそこで一連の評価を中断しエラーのみを一つ含むスライスを返す
		if isError(evaluated) {
			return []object.Object{evaluated}
		}

		// 追加
		result = append(result, evaluated)
	}
	return result
}

// 関数を引数に対して適応させ得られたObjectを返すヘルパー関数
func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		// 関数の持っている環境で環境を拡張する
		extendedEnv := extendFunctionEnv(fn, args)

		// 関数を引数に対して適応
		evaluated := Eval(fn.Body, extendedEnv)

		// ReturnValueObjectでったらならば皮を剥いでObject.Objectにする必要がある
		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		return fn.Fn(args...)
	default:
		return newError("not a function: %s", fn.Type())
	}
}

// 関数ごとに拡張された環境を返すヘルパー関数
func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {

	// まず関数自体の属している環境から拡張する環境を用意
	env := object.NewEnclosedEnvironment(fn.Env)

	// 拡張した環境に関数独自の変数を登録していく
	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}
	return env
}

// 関数呼び出しから戻ってくるReturnValueObjectをObjectに脱がせてやるヘルパー関数
// これがいないと関数からのReturnがプログラム全体のReturnとして扱われてしまう
func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

// 文字列による中置式を評価して適切なObjectを返すヘルパーヘルパー関数
func evalStringInfixExpression(operator string, left, right object.Object) object.Object {

	// 文字列に対して+しかサポートしていない
	// TODO: 文字列に対する!=演算子をサポートするならここに書く
	if operator != "+" {
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}

	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value
	return &object.String{Value: leftVal + rightVal}
}

// 組み込み関数を表すオブジェクトを登録するmap
var builtin = map[string]*object.Builtin{

	// USAGE:
	// len("string") -> 6
	// len([1, 23, 4]) -> 3
	"len": {
		Fn: func(args ...object.Object) object.Object {

			// ERROR: len("123", "234")
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1",
					len(args))
			}
			switch arg := args[0].(type) {

			// len("string")
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}

			// len([1, 2, 3])
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}

			// ERROR: len(123) etc.
			default:
				return newError("argument to `len` not supported, got=%s",
					args[0].Type())
			}
		},
	},

	// USAGE:
	// first(["A", 123, "54"]) -> "A"
	"first": {
		Fn: func(args ...object.Object) object.Object {

			// ERROR: first(["A", 123, "54"], [45, "45"])
			if len(args) != 1 {
				return newError("wrong number if arguments. got=%d, want=1",
					len(args))
			}

			// ERROR: first("array")
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `first` must be ARRAY, got %s",
					args[0].Type())
			}

			arr := args[0].(*object.Array)
			if len(arr.Elements) > 0 {
				return arr.Elements[0]
			}

			return NULL
		},
	},

	// USAGE:
	// last(["A", 123, "54"]) -> "54"
	"last": {
		Fn: func(args ...object.Object) object.Object {

			// ERROR: last(["A", 123, "54"], [45, "45"])
			if len(args) != 1 {
				return newError("wrong number if arguments. got=%d, want=1",
					len(args))
			}

			// ERROR: last("array")
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `last` must be ARRAY, got %s",
					args[0].Type())
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length > 0 {
				return arr.Elements[length-1]
			}

			return NULL
		},
	},

	// USAGE:
	// rest(["A", 123, "54"]) -> [123, "54"]
	"rest": {
		Fn: func(args ...object.Object) object.Object {

			// ERROR: rest(["A", 123, "54"], [45, "45"])
			if len(args) != 1 {
				return newError("wrong number if arguments. got=%d, want=1",
					len(args))
			}

			// ERROR: rest("array")
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `rest` must be ARRAY, got %s",
					args[0].Type())
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length > 0 {

				// 組み込み関数restは非破壊的な関数で、新たに割り当てられたArrayを返す
				newElements := make([]object.Object, length-1, length-1)
				copy(newElements, arr.Elements[1:length])
				return &object.Array{Elements: newElements}
			}

			return NULL
		},
	},

	// USAGE:
	// push(["A", 123, "54"], 45) -> ["A", 123, "54", 45]
	"push": {
		Fn: func(args ...object.Object) object.Object {

			// ERROR: push(["A", 123, "54"], 45, 45)
			if len(args) != 2 {
				return newError("wrong number if arguments. got=%d, want=2",
					len(args))
			}

			// ERROR: push("array")
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `push` must be ARRAY, got %s",
					args[0].Type())
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			newElements := make([]object.Object, length+1, length+1)
			copy(newElements, arr.Elements)
			newElements[length] = args[1]
			return &object.Array{Elements: newElements}
		},
	},

	// USAGE:
	// puts("Hello World") -> "Hello World"
	"puts": {
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}
			return NULL
		},
	},
}

// 添字演算子式が適切なオペランドに対して用いられているかを確認しつつ、適切なObjectに評価するヘルパー関数
func evalIndexExpression(left object.Object, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpressions(left, index)
	case left.Type() == object.HASH_OBJ:
		return evalHashIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}

// 配列に対する添字演算子式を適切なObjectに評価するヘルパーヘルパー関数
func evalArrayIndexExpressions(array, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	// 配列に格納している要素数を超えたインデックスに対してはNULLObjectを返す
	if idx < 0 || max < idx {
		return NULL
	}
	return arrayObject.Elements[idx]
}

// ハッシュリテラルを評価してObjectを返す関数
// リテラルのペアに対するHashKeyを生成して、リテラルのペアとそのHashKeyの組をObjectとして保存しておく
// {"one": 1, "two": 2}というリテラルのハッシュに対してこれを評価した結果得られるのは
// {「"one"-1」というペアとこれに対するHashKey、「"two"-2」というペアとこれに対するHashKey}というObject
func evalHashLiteral(node *ast.HashLiteral, env *object.Environment) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)
	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env)
		if isError(key) {
			return key
		}
		hashKey, ok := key.(object.Hashable)
		if !ok {
			return newError("unusable as hash key: %s", key.Type())
		}
		value := Eval(valueNode, env)
		if isError(value) {
			return value
		}
		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}
	return &object.Hash{Pairs: pairs}
}

func evalHashIndexExpression(hash, index object.Object) object.Object {
	hashObject := hash.(*object.Hash)
	key, ok := index.(object.Hashable)
	if !ok {
		return newError("unusable as hash key: %s", index.Type())
	}
	pair, ok := hashObject.Pairs[key.HashKey()]
	if !ok {
		return NULL
	}
	return pair.Value
}

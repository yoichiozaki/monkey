package evaluator

import (
	"monkey/object"
)

// 組み込み関数を表すオブジェクトを登録するmap
var builtin = map[string]*object.Builtin{

	// USAGE:
	// len("string") -> 6
	// len([1, 23, 4]) -> 3
	// "len": {
	// 	Fn: func(args ...object.Object) object.Object {
	//
	// 		// ERROR: len("123", "234")
	// 		if len(args) != 1 {
	// 			return newError("wrong number of arguments. got=%d, want=1",
	// 				len(args))
	// 		}
	// 		switch arg := args[0].(type) {
	//
	// 		// len("string")
	// 		case *object.String:
	// 			return &object.Integer{Value: int64(len(arg.Value))}
	//
	// 			// len([1, 2, 3])
	// 		case *object.Array:
	// 			return &object.Integer{Value: int64(len(arg.Elements))}
	//
	// 			// ERROR: len(123) etc.
	// 		default:
	// 			return newError("argument to `len` not supported, got=%s",
	// 				args[0].Type())
	// 		}
	// 	},
	// },
	"len": object.GetBuiltinByName("len"),

	// USAGE:
	// first(["A", 123, "54"]) -> "A"
	// "first": {
	// 	Fn: func(args ...object.Object) object.Object {
	//
	// 		// ERROR: first(["A", 123, "54"], [45, "45"])
	// 		if len(args) != 1 {
	// 			return newError("wrong number if arguments. got=%d, want=1",
	// 				len(args))
	// 		}
	//
	// 		// ERROR: first("array")
	// 		if args[0].Type() != object.ARRAY_OBJ {
	// 			return newError("argument to `first` must be ARRAY, got %s",
	// 				args[0].Type())
	// 		}
	//
	// 		arr := args[0].(*object.Array)
	// 		if len(arr.Elements) > 0 {
	// 			return arr.Elements[0]
	// 		}
	//
	// 		return NULL
	// 	},
	// },
	"first": object.GetBuiltinByName("first"),

	// USAGE:
	// last(["A", 123, "54"]) -> "54"
	// "last": {
	// 	Fn: func(args ...object.Object) object.Object {
	//
	// 		// ERROR: last(["A", 123, "54"], [45, "45"])
	// 		if len(args) != 1 {
	// 			return newError("wrong number if arguments. got=%d, want=1",
	// 				len(args))
	// 		}
	//
	// 		// ERROR: last("array")
	// 		if args[0].Type() != object.ARRAY_OBJ {
	// 			return newError("argument to `last` must be ARRAY, got %s",
	// 				args[0].Type())
	// 		}
	//
	// 		arr := args[0].(*object.Array)
	// 		length := len(arr.Elements)
	// 		if length > 0 {
	// 			return arr.Elements[length-1]
	// 		}
	//
	// 		return NULL
	// 	},
	// },
	"last": object.GetBuiltinByName("last"),

	// USAGE:
	// rest(["A", 123, "54"]) -> [123, "54"]
	// "rest": {
	// 	Fn: func(args ...object.Object) object.Object {
	//
	// 		// ERROR: rest(["A", 123, "54"], [45, "45"])
	// 		if len(args) != 1 {
	// 			return newError("wrong number if arguments. got=%d, want=1",
	// 				len(args))
	// 		}
	//
	// 		// ERROR: rest("array")
	// 		if args[0].Type() != object.ARRAY_OBJ {
	// 			return newError("argument to `rest` must be ARRAY, got %s",
	// 				args[0].Type())
	// 		}
	//
	// 		arr := args[0].(*object.Array)
	// 		length := len(arr.Elements)
	// 		if length > 0 {
	//
	// 			// 組み込み関数restは非破壊的な関数で、新たに割り当てられたArrayを返す
	// 			newElements := make([]object.Object, length-1, length-1)
	// 			copy(newElements, arr.Elements[1:length])
	// 			return &object.Array{Elements: newElements}
	// 		}
	//
	// 		return NULL
	// 	},
	// },
	"rest": object.GetBuiltinByName("rest"),

	// USAGE:
	// push(["A", 123, "54"], 45) -> ["A", 123, "54", 45]
	// "push": {
	// 	Fn: func(args ...object.Object) object.Object {
	//
	// 		// ERROR: push(["A", 123, "54"], 45, 45)
	// 		if len(args) != 2 {
	// 			return newError("wrong number if arguments. got=%d, want=2",
	// 				len(args))
	// 		}
	//
	// 		// ERROR: push("array")
	// 		if args[0].Type() != object.ARRAY_OBJ {
	// 			return newError("argument to `push` must be ARRAY, got %s",
	// 				args[0].Type())
	// 		}
	//
	// 		arr := args[0].(*object.Array)
	// 		length := len(arr.Elements)
	// 		newElements := make([]object.Object, length+1, length+1)
	// 		copy(newElements, arr.Elements)
	// 		newElements[length] = args[1]
	// 		return &object.Array{Elements: newElements}
	// 	},
	// },
	"push": object.GetBuiltinByName("push"),

	// USAGE:
	// puts("Hello World") -> "Hello World"
	// "puts": {
	// 	Fn: func(args ...object.Object) object.Object {
	// 		for _, arg := range args {
	// 			fmt.Println(arg.Inspect())
	// 		}
	// 		return NULL
	// 	},
	// },
	"puts": object.GetBuiltinByName("puts"),
}

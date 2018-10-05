package object

import (
	"fmt"
)

// -----------------------------------------------------
// Objectの定義
type ObjectType string

// Monkeyに登場する値はすべてObjectとする
type Object interface {
	Type() ObjectType
	Inspect() string
}

const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VAL"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
)

// -----------------------------------------------------

// -----------------------------------------------------
// Integerの定義
type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

// -----------------------------------------------------

// -----------------------------------------------------
// Booleanの定義
type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

// -----------------------------------------------------

// -----------------------------------------------------
// Nullの定義
type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "Null" }

// -----------------------------------------------------

// -----------------------------------------------------
// Returnの定義
type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

// -----------------------------------------------------

// -----------------------------------------------------
// Errorの定義
type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }

// -----------------------------------------------------

// -----------------------------------------------------
// Environmentの定義
type Environment struct {

	// 識別子に対応するObjectを保存する
	store map[string]Object
}

// 新しい環境を生成する
// 常に一つの環境を使いまわしたいのでポインタで渡す
func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s}
}

// 環境内にnameという名前で登録されているObjectを持ってくる
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	return obj, ok
}

// 環境内にnameという名前でObjectを登録する
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

// -----------------------------------------------------

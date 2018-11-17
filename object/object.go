package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"monkey/ast"
	"monkey/code"
	"strings"
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
	INTEGER_OBJ              = "INTEGER"
	BOOLEAN_OBJ              = "BOOLEAN"
	NULL_OBJ                 = "NULL"
	RETURN_VALUE_OBJ         = "RETURN_VAL"
	ERROR_OBJ                = "ERROR"
	FUNCTION_OBJ             = "FUNCTION"
	STRING_OBJ               = "STRING"
	BUILTIN_OBJ              = "BUILTIN"
	ARRAY_OBJ                = "ARRAY"
	HASH_OBJ                 = "HASH"
	COMPILED_FUNCTION_OBJECT = "COMPILED_FUNCTION_OBJECT"
	CLOSURE_OBJ              = "CLOSURE"
)

// ハッシュテーブルにおける管理用オブジェクトとしてのHashKey
type HashKey struct {
	Type  ObjectType
	Value uint64
}

// Monkeyのハッシュテーブルに格納できるものはHashableインタフェースを満たさなくてはならない
type Hashable interface {
	HashKey() HashKey
}

// -----------------------------------------------------

// -----------------------------------------------------
// Integerの定義
type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

// -----------------------------------------------------

// -----------------------------------------------------
// Booleanの定義
type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) HashKey() HashKey {
	var value uint64
	if b.Value {
		value = 1
	} else {
		value = 0
	}
	return HashKey{Type: b.Type(), Value: value}
}

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
// Functionの定義
type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}
	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")
	return out.String()
}

// -----------------------------------------------------

// -----------------------------------------------------
// コンパイルされた関数を表現するオブジェクトの定義
type CompiledFunction struct {
	Instructions  code.Instructions // この関数をコンパイルして得られる命令列
	NumLocals     int               // 関数内で使われるローカル変数の個数
	NumParameters int               // 関数リテラルが実行しようとしているときに保持している引数の個数
}

func (cf *CompiledFunction) Type() ObjectType { return COMPILED_FUNCTION_OBJECT }
func (cf *CompiledFunction) Inspect() string {
	return fmt.Sprintf("CompiledFunction[%p]", cf)
}

// -----------------------------------------------------

// -----------------------------------------------------
// Stringの定義
type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }
func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

// -----------------------------------------------------

// -----------------------------------------------------
// Builtinの定義
type BuiltinFunction func(args ...Object) Object
type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "builtin function" }

// -----------------------------------------------------

// -----------------------------------------------------
// Arrayオブジェクトの定義
type Array struct {
	Elements []Object
}

func (ao *Array) Type() ObjectType { return ARRAY_OBJ }
func (ao *Array) Inspect() string {
	var out bytes.Buffer
	elements := []string{}
	for _, e := range ao.Elements {
		elements = append(elements, e.Inspect())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

// -----------------------------------------------------

// -----------------------------------------------------
// Hashオブジェクトの定義
type HashPair struct {
	Key   Object // RELPにおいてハッシュを表示するときにキーと値のペアを表示するために必要
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType { return HASH_OBJ }
func (h *Hash) Inspect() string {
	var out bytes.Buffer
	pairs := []string{}
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}

// -----------------------------------------------------

// -----------------------------------------------------
// Closureオブジェクトの定義
type Closure struct {
	Fn   *CompiledFunction
	Free []Object // store for free variables.
}

func (c *Closure) Type() ObjectType { return CLOSURE_OBJ }
func (c *Closure) Inspect() string {
	return fmt.Sprintf("Closure[%p]", c)
}

// -----------------------------------------------------

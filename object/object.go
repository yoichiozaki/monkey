package object

import "fmt"

// -----------------------------------------------------
// Objectの定義
type ObjectType string

// Monkeyに登場する値はすべてObjectとする
type Object interface {
	Type() ObjectType
	Inspect() string
}

const (
	INTEGER_OBJ = "INTEGER"
	BOOLEAN_OBJ = "BOOLEAN"
	NULL_OBJ    = "NULL"
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

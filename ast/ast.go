package ast

import (
	"bytes"
	"monkey/token"
	"strings"
)

// Nodeは以下の関数を実装している
type Node interface {
	TokenLiteral() string
	String() string
}

// 文ノード: 値を返さない
type Statement interface {
	Node
	statementNode()
}

// 式ノード: 値を返す
type Expression interface {
	Node
	expressionNode()
}

// -----------------------------------------------------
// プログラムを表すASTノード: 文の集合
type Program struct {
	Statements []Statement
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

// -----------------------------------------------------

// -----------------------------------------------------
// LET文を表すASTノード
// let <identifier> = <expression>;
// let x = 5;
type LetStatement struct {
	Token token.Token // token.LET = "let"
	Name  *Identifier // x
	Value Expression  // 5
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")
	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	out.WriteString(";")
	return out.String() // "let x = 5;"
}

// -----------------------------------------------------

// -----------------------------------------------------
// 識別子を表すASTノード
// 「let x = 5;」における「x」
type Identifier struct {
	Token token.Token // token.IDENT
	Value string      // x
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

// -----------------------------------------------------

// -----------------------------------------------------
// RETURN文を表すASTノード
// return <expression>;
// return 5;
type ReturnStatement struct {
	Token       token.Token // token.RETURN = "return"
	ReturnValue Expression  // 5
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(rs.TokenLiteral() + " ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String() // "return 5;"
}

// -----------------------------------------------------

// -----------------------------------------------------
// 式文を表すASTノード
// 式単体で文扱い。要するに「式のwrapper」としての型
// let x = 5;
// x + 10; <- これ
type ExpressionStatement struct {
	Token      token.Token // 式の最初のトークン
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// -----------------------------------------------------

// -----------------------------------------------------
// 整数リテラルを表すASTノード(整数値も式)
// 5
type IntegerLiteral struct {
	Token token.Token // token.INT = int
	Value int64       // 5
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

// -----------------------------------------------------

// -----------------------------------------------------
// 前置演算子を表すASTノード
// <prefix operator> <expression>;
// !true
// -5
type PrefixExpression struct {
	Token    token.Token // 前置トークン
	Operator string      // 「!」「-」のどちらか
	Right    Expression  // 前置演算子の右隣に来る式を表現するASTノード
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String() // (!true), (-5)
}

// -----------------------------------------------------

// -----------------------------------------------------
// 中置演算子を表すASTノード
// <expression> <infix operator> <expression>;
// 5 + 5
type InfixExpression struct {
	Token    token.Token // 演算子トークン
	Left     Expression  // 5
	Operator string      // *
	Right    Expression  // 5
}

func (oe *InfixExpression) expressionNode()      {}
func (oe *InfixExpression) TokenLiteral() string { return oe.Token.Literal }
func (oe *InfixExpression) String() string {
	var out bytes.Buffer
	// out.WriteString("(")
	out.WriteString(oe.Left.String())
	out.WriteString(" " + oe.Operator + " ")
	out.WriteString(oe.Right.String())
	// out.WriteString(")")
	return out.String() // "(5 * 5)"
}

// -----------------------------------------------------

// -----------------------------------------------------
// BOOLEAN型のトークンを表すASTノード
// false
type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }

// -----------------------------------------------------

// -----------------------------------------------------
// IF文を表すASTノード
// if ( <condition> ) <consequence> else <alternative>
// if(x < y) { return x; } else { return y; }
type IfExpression struct {
	Token       token.Token     // 'if' トークン
	Condition   Expression      // x < y
	Consequence *BlockStatement // return x;
	Alternative *BlockStatement // return y;
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())
	if ie.Alternative != nil { // <alternative>があればelse節
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}
	return out.String() // "if(x < y) { return x;} else { return y;}
}

// -----------------------------------------------------

// -----------------------------------------------------
// ブロック文を表すASTノード
// ブロックは複数の文で成る
type BlockStatement struct {
	Token      token.Token // '{' トークン
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	for _, s := range bs.Statements {
		out.WriteString("\n")
		out.WriteString("\t" + s.String())
		out.WriteString("\n")
	}
	return out.String()
}

// -----------------------------------------------------

// -----------------------------------------------------
// 関数リテラルを表すASTノード
// fn <parameters> <block statement>
// fn(x, y) { x + y; }
type FunctionLiteral struct {
	Token      token.Token     // 'fn' トークン
	Parameters []*Identifier   // x, y
	Body       *BlockStatement // x + y;
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}
	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())
	return out.String()
}

// -----------------------------------------------------

// -----------------------------------------------------
// 関数呼び出し式を表すASTノード
// <expression> ( <comma separated expressions> )
// add(2, 3)
// fn(x, y) { x + y }(2, 3)
type CallExpression struct {
	Token     token.Token // '(' トークン
	Function  Expression  // Identifier または FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out bytes.Buffer
	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}
	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	return out.String()
}

// -----------------------------------------------------

// -----------------------------------------------------
// 文字列を表すASTノード
// 文字列は式であって文ではない
// <sequence of characters>
// "hello world"
type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return sl.Token.Literal }

// -----------------------------------------------------

// -----------------------------------------------------
// 配列リテラルを表すASTノード
// [ <sequence of Expressions> ]
// ["hello", 123, fn(name) = { return "Hi there, " + name; }];
type ArrayLiteral struct {
	Token    token.Token // '[' トークン
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer
	elements := []string{}
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

// -----------------------------------------------------

// -----------------------------------------------------
// 添字演算子式を表すASTノード
// <expression> [ <expression> ]
// myArray[3]
// [1, 2, 3, 4][3] => 4
type IndexExpression struct {
	Token token.Token // '[' トークン
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("]")
	out.WriteString(")")
	return out.String()
}

// -----------------------------------------------------

// -----------------------------------------------------
// ハッシュリテラルを表すASTノード
// { <expression> : <expression>, <expression> : <expression>, ... }
type HashLiteral struct {
	Token token.Token // '{' トークン
	Pairs map[Expression]Expression
}

func (hl *HashLiteral) expressionNode()      {}
func (hl *HashLiteral) TokenLiteral() string { return hl.Token.Literal }
func (hl *HashLiteral) String() string {
	var out bytes.Buffer
	pairs := []string{}
	for key, value := range hl.Pairs {
		pairs = append(pairs, key.String()+": "+value.String())
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}

// -----------------------------------------------------

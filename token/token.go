package token

type TokenType string
type Token struct {
	Type TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF = "EOF"

	// 識別子 + リテラル
	IDENT = "IDENT" // add, result, x, y, etc.
	INT = "INT" // 12, 34, ...

	// 演算子
	ASSIGN = "="
	PLUS = "+"
	MINUS = "-"
	BANG = "!"
	ASTERISK = "*"
	SLASH = "/"

	LT = "<" // Less Than
	GT = ">" // Greater Than

	EQ = "=="
	NOT_EQ = "!="

	// デリミタ
	COMMA = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// キーワード
	FUNCTION = "FUNCTION"
	LET = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
)

// ユーザー定義の識別子と言語のキーワードを区別する機能
var keywords = map[string]TokenType {
	"fn": FUNCTION,
	"let": LET,
	"true": TRUE,
	"false": FALSE,
	"if": IF,
	"else": ELSE,
	"return": RETURN,
}

// 渡された識別子とされるものがキーワードではないかを確認する
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok // それはキーワードだった
	}
	return IDENT // それは識別子だった
}
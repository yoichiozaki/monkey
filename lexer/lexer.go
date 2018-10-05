package lexer

import "monkey/token"

type Lexer struct {
	input        string
	position     int  // 入力における現在の位置
	readPosition int  // これから読み込む文字の位置（すなわち現在の文字の次の文字）
	ch           byte // 現在検査中の文字
}

// 入力によって初期化済みの字句解析器を与える
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// 文字を一つ読み込む
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

// 読み込んだ文字を判別して対応するトークンを返す
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhiteSpace()

	switch l.ch {
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case 0: // 終端
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			// ここで識別子であろうとされているtok.Literalがキーワードでないことを確認する
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

// NextToken()内で、読み込んだ文字に対応するトークンを返すヘルパー関数
func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func (l *Lexer) readIdentifier() string {
	// 非英字まで読み進めていく
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
	// 「|| ch == '_'」の部分は、意外と重要。ここではアンダースコアを英文字として認めるということ
	// isLetter()の認める範囲内で識別子に用いることができる文字種に制限がかかる
}

// 空白とか改行文字を読み飛ばしていく
func (l *Lexer) skipWhiteSpace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// readPositionの文字を処理する前に覗き見peekする関数
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

// 文字列として扱われるべき部分まで読み進めていき、得られた文字列を返す関数
func (l *Lexer) readString() string {

	// TODO: 以下の実装では「"hello \"world\""」や「hello\n world」「hello\t\t\tworld」などには対応できていない

	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}

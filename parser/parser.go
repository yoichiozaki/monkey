package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
	"strconv"
)

const (
	// 優先順位の定義
	_ int = iota
	LOWEST
	EQUALS     // ==
	LESSGRATER // > or <
	SUM        // +
	PRODUCT    // *
	PREFIX     // -x or !x
	CALL       // myFunction(x)
)

// 優先順位テーブル
var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGRATER,
	token.GT:       LESSGRATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
}

// パーサの定義
type Parser struct {
	l         *lexer.Lexer // 字句解析器を内部に含む
	errors    []string     // エラー
	curToken  token.Token  // 今見ているトークン
	peekToken token.Token  // 次見るべきトークン

	// Pratt構文解析器のアイディアの核心
	prefixParseFns map[token.TokenType]prefixParseFn // 特定の前置演算子トークンとそれを解析する関数のマップ
	infixParseFns  map[token.TokenType]infixParseFn  // 特定の中置演算子トークンとそれを解析する関数のマップ
}

// パーサーを生成する
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}
	p.nextToken()
	p.nextToken()

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	return p
}

// エラーを返す
func (p *Parser) Errors() []string {
	return p.errors
}

// 次に来るべきトークンが来ていないならばエラーメッセージを追加
func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

// 見るトークンを一つ進める
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// プログラムをパースしてProgram型のASTノードを返す
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF { // EOF型のトークンに遭遇するまで

		// 現在見ている文をパースした結果得られるStatement型のASTノードstmtを
		// ノードprogramのStatementsフィールドに追加する
		stmt := p.parseStatement()

		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken() // 調べるトークンを一つ進める
	}
	return program // パースして得られたProgram型のASTノードを返す
}

// 文をパースしてStatement型のASTノードを返す
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type { // 現在見ているトークンのタイプによって処理が分かれる
	case token.LET: // LET文: let <identifier> = <expression>;
		return p.parseLetStatement()
	case token.RETURN: // RETURN文: return <expression>;
		return p.parseReturnStatement()
	default: // その他は式文
		return p.parseExpressionStatement()
	}
}

// LET文をパースしてLetStatement型のASTノードを返す
func (p *Parser) parseLetStatement() *ast.LetStatement {
	// let <identifier> = <expression>;
	// let x = 5;

	// LetStatement型のASTノードを生成
	stmt := &ast.LetStatement{Token: p.curToken}

	// 後続するトークンにアサーションを設けつつパースを進めていく
	if !p.expectPeek(token.IDENT) { // let = 5;みたいなやつはだめ
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.expectPeek(token.ASSIGN) { // let x 5;みたいなやつはだめ
		return nil
	}

	// ここに到達しているのでLET文としての体裁は整っているはず
	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// 今見ているトークンのタイプをチェックするヘルパー関数
func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

// 次来るトークンのタイプをチェックするヘルパー関数
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// 次来るトークンのタイプが期待したものであるかどうかを返すヘルパー関数
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t) // パーサがエラーを検知ししたのでエラーメッセージを追加
		return false
	}
}

// RETURN文をパースしてReturnStatement型のASTノードを返す
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	// return <expression>;
	// return 5;

	// ReturnStatement型のASTノードを生成
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	// returnに続くトークンをパースした結果得られるExpression型のASTノードをstmtのReturnValueとして追加
	stmt.ReturnValue = p.parseExpression(LOWEST)

	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

// Prattの考え方の革新的なところの一つ
// 各トークンにそのトークンを解析する2関数を関連付けさせる
// それぞれの使い分けはトークンの出現位置で判別する
type (
	prefixParseFn func() ast.Expression               // 関連付けられたトークンタイプが前置で出現した場合に呼ばれる
	infixParseFn  func(ast.Expression) ast.Expression // 関連付けられたトークンタイプが中置で出現した場合に呼ばれる
)

// 特定の前置演算子に対してそのトークンを解析する関数を紐づけるヘルパー関数
func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

// 特定の中置演算子に対してそのトークンを解析する関数を紐づけるヘルパー関数
func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

// 式文をパースしてExpressionStatement型のASTノードを返す
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	// defer untrace(trace("parseExpressionStatement"))

	// ExpressionStatement型のASTノードを生成
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// 式をパースしてExpression型のASTノードを返す
func (p *Parser) parseExpression(precedence int) ast.Expression {
	// defer untrace(trace("parseExpression"))

	// 現在見ているトークンに関連付けられた構文解析関数が存在するかを確認する
	prefix := p.prefixParseFns[p.curToken.Type]

	// なければnoPrefixParserErrorを吐いてパーサ内にエラーメッセージを記録してnilのASTノードを返す
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}

	// 現在見ているトークンを解析する
	leftExp := prefix()

	// 現在見ているトークンの右結合力（precedence）と左結合力（peekPrecedence()）を確認
	// 左結合力が高いということは1つネストするということになる
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}
	return leftExp
}

// 識別子をパースしてExpression型のASTノードを返す
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

// 整数リテラルをパースしてExpression型のASTノードを返す
func (p *Parser) parseIntegerLiteral() ast.Expression {
	// defer untrace(trace("parseIntegerLiteral"))

	// IntegerLiteral型のASTノードを生成
	lit := &ast.IntegerLiteral{Token: p.curToken}

	// 今見ているトークンのリテラルが整数リテラルであることを確認する
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)

	// 整数リテラルでなければエラーメッセージをパーサ内に記録したのちnilのExpression型ASTノードを返す
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer",
			p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	// IntegerLiteral型のASTノードに値を格納
	lit.Value = value
	return lit
}

// 該当する前置演算子トークンに対してそれをパースする関数が紐づけられていなかった時にエラーメッセージを出力するヘルパー関数
func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

// 前置演算子トークンをパースしてExpression型のASTノードを返す
func (p *Parser) parsePrefixExpression() ast.Expression {
	// defer untrace(trace("parsePrefixExpression"))

	// PrefixExpression型のASTノードを生成
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	// 前置演算子の作用する式のパースに移る
	p.nextToken()

	// 前置演算子の右隣をパースした結果をPrefixExpression型のASTノードexpressionのRightフィールドに格納
	expression.Right = p.parseExpression(PREFIX)

	return expression
}

// 次に見るべきトークンの優先順位を返すヘルパー関数
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

// 現在見ているトークンの優先順位を返すヘルパー関数
func (p *Parser) currPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

// 現在見ている中置演算子の左にある式を表現するExpression型のASTノードを引数に、
// その中置演算子トークンをパースしてExpression型のASTノードを返す
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	// defer untrace(trace("parseInfixExpression"))

	// InfixExpression型のASTノードを生成する
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	// 現在見ている中置演算子の優先度をprecedenceに格納する
	precedence := p.currPrecedence()

	// トークンを一つ進める
	// すなわちここを通過した暁に見ているトークンは中置演算子の右側にあるトークン
	p.nextToken()

	// 中置演算子の右側にあるトークンをパースしてその結果得られるExpression型のASTノードを
	// InfixExpression型のASTノードexpressionのRightフィールドに格納する
	expression.Right = p.parseExpression(precedence)

	return expression
}

// Boolean型のトークンをパースしてExpression型のASTノードを返す
func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

// 丸括弧でまとめられたトークンをパースしてExpression型のASTノードを返す
func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return exp
}

// IF式をパースしてExpression型のASTノードを返す
func (p *Parser) parseIfExpression() ast.Expression {
	// if (<condition>) <consequence> else <alternative>;
	// if (x > y) { return x; } else { return y; }
	// let foobar = if (x > y) { x; } else { y; }

	// IfExpression型のASTノードを生成
	expression := &ast.IfExpression{Token: p.curToken}

	// 「(」が来るはず
	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()

	// Condition部をパースしてIfExpression型のASTノードexpressionのConditionフィールドに追加
	expression.Condition = p.parseExpression(LOWEST)

	// 「)」が来るはず
	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	// 「{」が来るはず
	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	//Consequence部をパースIfExpression型のASTノードexpressionのConsequenceフィールドに追加
	expression.Consequence = p.parseBlockStatement()

	// 続けてelseが来るかどうかを確認する
	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		// 「{」が来るはず
		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		//Alternative部をパースIfExpression型のASTノードexpressionのAlternativeフィールドに追加
		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

// ブロック文をパースしてBlockStatement型のASTノードを返す
func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	// { statement1; statement2; ... }

	// BlockStatement型のASTノードを生成
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	// 「}」かEOFに到達するまでに遭遇する文をパースしながらblockのStatementフィールドにその結果のASTを追加していく
	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}
	return block
}

// 関数リテラルをパースしてExpression型のASTノードを返す
func (p *Parser) parseFunctionLiteral() ast.Expression {
	// fn (<parameter1>, <parameter2>, ...) <block statement>;
	// fn () <blocks tatement>;

	// FunctionLiteral型のASTノードを生成
	lit := &ast.FunctionLiteral{Token: p.curToken}

	// 「(」が来るはず
	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	// 関数の引数リストをパースして得られるASTをFunctionLiteral型のASTノードlitのParametersフィールドに登録
	lit.Parameters = p.parseFunctionParameters()

	// 「{」が来るはず
	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	// block statementである関数の本体をパースして得られるASTをFunctionLiteral型のASTノードlitのBodyフィールドに登録
	lit.Body = p.parseBlockStatement()

	return lit
}

// 関数リテラルの引数リストを解析してIdentifier型のASTノードのスライスを返すヘルパー関数
func (p *Parser) parseFunctionParameters() []*ast.Identifier {

	// 関数の引数リストは識別子の集まり
	identifiers := []*ast.Identifier{}

	// fn()の時
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	// 一つ目の識別子に遭遇
	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// Identifier型のASTノードを生成したので追加
	identifiers = append(identifiers, ident)

	// コンマごとにident見つけてidentifiersに追加していく
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
	}

	// 「)」が来るはず
	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return identifiers
}

// 関数呼び出し式をパースしてExpression型のASTノードを返す
func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {

	// CallExpression型のASTノードを生成
	exp := &ast.CallExpression{Token: p.curToken, Function: function}

	// expのArgumentsフィールドに実引数を格納する
	exp.Arguments = p.parseCallArguments()
	return exp
}

// 関数呼び出し式における実引数リストを解析してExpression型のASTノードを返すヘルパー関数
func (p *Parser) parseCallArguments() []ast.Expression {

	// 返すべき実引数リストを表現するExpression型のASTノードのスライスを用意
	args := []ast.Expression{}

	// hello()みたいな関数の時は空の引数リスト
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return args
	}

	p.nextToken()

	// 実引数に遭遇したのでパースしてargsに追加
	args = append(args, p.parseExpression(LOWEST))

	// コンマに遭遇するごとに同じことを繰り返す
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	// 「)」が来るはず
	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return args
}

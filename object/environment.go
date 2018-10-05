package object

// -----------------------------------------------------
// Environmentの定義
type Environment struct {

	// 識別子に対応するObjectを保存する
	store map[string]Object

	// 拡張環境
	outer *Environment
}

// 新しい環境を生成する
// 常に一つの環境を使いまわしたいのでポインタで渡す
func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s}
}

// 環境内にnameという名前で登録されているObjectを持ってくる
func (e *Environment) Get(name string) (Object, bool) {

	// 内側の環境(スコープ)でまず探す
	obj, ok := e.store[name]

	// 内側で見つからなくてかつ拡張されているならば外側を探す
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

// 環境内にnameという名前でObjectを登録する
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

// 拡張環境をセットする
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

// -----------------------------------------------------

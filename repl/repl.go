package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey/compiler"
	"monkey/lexer"
	// "monkey/object"
	"monkey/parser"
	"monkey/vm"
)

const PROMPT = ">> "

const MONKEY = `    ___
　 彡_＿ ＼_　 n
　 (・・) ○) ((
　 /‥ ( 　｜　))
　( Θ　)　 ＼//_
　 ￣/￣　 (⌒⌒)ヽ_
　　｜ ∧　 ＼／ (ミ)
　　｜｜ ＼　　　ノ
　　 ＼_) ( (￣) )
　　　　　(＿)(＿)

`

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	// env := object.NewEnvironment()

	for {

		// プロンプト「>>」の出力
		fmt.Printf(PROMPT)

		// 入力
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()

		// inputで初期化されたレキサを生成
		l := lexer.New(line)

		// for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
		// 	fmt.Printf("%+v\n", tok)
		// }

		// レキサをセットしたパーサを生成
		p := parser.New(l)

		// プログラムをパース
		program := p.ParseProgram()

		// パース中のエラーを出力
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		// io.WriteString(out, program.String())
		// io.WriteString(out, "\n")

		// // パースした結果得られたASTを評価器に通してObjectを得る
		// evaluated := evaluator.Eval(program, env)
		//
		// if evaluated != nil {
		// 	io.WriteString(out, evaluated.Inspect())
		// 	io.WriteString(out, "\n")
		// }

		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			io.WriteString(out, MONKEY)
			fmt.Fprintf(out, "Woops! Complation failed:\n\t%s\n", err)
			continue
		}
		machine := vm.New(comp.Bytecode())
		err = machine.Run()
		if err != nil {
			io.WriteString(out, MONKEY)
			fmt.Fprintf(out, "Woops! Executing bytecode failed:\n\t%s\n", err)
			continue
		}
		stackTop := machine.LastPoppedStackElem()
		io.WriteString(out, stackTop.Inspect())
		io.WriteString(out, "\n")
	}
}

// パース中のエラーを出力するヘルパー関数
func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, MONKEY)
	io.WriteString(out, "Woops! We ran into some monkey business here!\n")
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}

package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/akyrey/monkey-programming-language/lexer"
	"github.com/akyrey/monkey-programming-language/token"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprintf(out, PROMPT)
		// Read from the input source until encountering a newline
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		// Take the just read line and pass it to an instance of our lexer
		line := scanner.Text()
		l := lexer.New(line)

		// Print all the tokens the lexer gives us until we encounter EOF
		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Fprintf(out, "%+v\n", tok)
		}
	}
}

package main

import (
	"bufio"
	"fmt"
	"os"

	gojs "github.com/ruiconti/gojs/internal"
	"github.com/ruiconti/gojs/lex"
)

func prettyPrintMap(tokens []lex.Token) ([]string, error) {
	var result []string
	for _, token := range tokens {
		name, err := lex.ResolveName(token)
		if err != nil {
			return []string{}, err
		}
		pretty := fmt.Sprintf("%s(%s)", name, token.Lexeme)
		if err != nil {
			return []string{}, err
		}
		result = append(result, pretty)
	}
	return result, nil
}

func main() {
	s := bufio.NewScanner(os.Stdin)
	// fmt.Printf("> ")
	logger := gojs.NewSimpleLogger(gojs.ModeInfo | gojs.ModeWarn | gojs.ModeError)
	for {
		fmt.Print("\n> ")
		n := s.Scan()
		if !n {
			break
		}
		src := s.Text()
		scanner := lex.NewScanner(src, logger)
		tokens, scanErr := scanner.Scan()
		if scanErr != nil {
			fmt.Printf("SyntaxError: %s is not valid: %s.", src, scanErr)
			continue
		}
		ptokens, err := prettyPrintMap(tokens)
		if err != nil {
			fmt.Printf("ReplError: failed to pretty print tokens %v: %s.", tokens, err)
			panic(1)
		}
		if len(tokens) == 0 {
			fmt.Printf("SyntaxError: %s is not valid.", src)
		} else if len(tokens) == 1 {
			fmt.Printf("%v", ptokens)
		} else {
			fmt.Print("[\n")
			for _, token := range ptokens {
				fmt.Printf("  %v\n", token)
			}
			fmt.Print("]")
		}
	}

	if err := s.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

}

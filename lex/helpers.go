package lex

import (
	"fmt"
	"strings"
	"testing"

	gojs "github.com/ruiconti/gojs/internal"
)

var defaultLogger = gojs.NewSimpleLogger(gojs.ModeError)

// ResolveName returns the name of a token type.
func ResolveName(t Token) (string, error) {
	dicts := []map[TokenType]string{
		LiteralNames, ReservedWordNames,
		PunctuationNames}
	for _, dict := range dicts {
		if name, ok := dict[t.T]; ok {
			return name, nil
		}
	}
	return "", fmt.Errorf("token name not found: %d", t.T)
}

func assertLiterals(t *testing.T, got, expected []Token) {
	assertInternal(
		t,
		got,
		expected, func(a, b Token) bool { return a.Lexeme == b.Lexeme && a.Literal == b.Literal },
		func(a Token) string {
			if a.Literal == nil {
				return ""
			}
			return fmt.Sprintf("%v", a.Literal)
		},
		// inline print
		func(a Token) string {
			return fmt.Sprintf("%v ", a.Lexeme)
		},
	)
}

func assertTokens(t *testing.T, got, expected []Token) {
	assertInternal(
		t,
		got,
		expected, func(a, b Token) bool { return a.Lexeme == b.Lexeme && a.Literal == b.Literal && a.T == b.T },
		func(a Token) string {
			if a.Literal == nil {
				return ""
			}
			name, err := ResolveName(a)
			if err != nil {
				name = fmt.Sprintf("T%d", a.T)
			}
			return fmt.Sprintf("nam:%v lit:%v lex:%v", name, a.Literal, a.Lexeme)
		},
		// inline print
		func(a Token) string {
			return fmt.Sprintf("%v ", a.Lexeme)
		},
	)
}

func assertLexemes(t *testing.T, got, expected []Token) {
	t.Helper()
	assertInternal(
		t,
		got,
		expected, func(a, b Token) bool { return a.Lexeme == b.Lexeme },
		// readable print
		func(a Token) string {
			if a.Lexeme == "" {
				return ""
			}
			return fmt.Sprintf("%v ", a.Lexeme)
		},
		// inline print
		func(a Token) string {
			return fmt.Sprintf("%v ", a.Lexeme)
		},
	)
}

func assertInternal(
	t *testing.T,
	got,
	expected []Token,
	callbackEqualCmp func(a,
		b Token) bool,
	callbackPrint func(a Token) string,
	callbackPrintInline func(a Token) string) {
	t.Helper()
	// Need to be in the same order
	failure := false
	if len(got) != len(expected) {
		failure = true
	}
	diffStrs := []string{}
	for i, expectedToken := range expected {
		var gotToken Token
		if i >= len(got) {
			gotToken = Token{T: TEOF, Lexeme: `<EOF>`, Literal: nil, Line: 0, Column: 0}
		} else {
			gotToken = got[i]
		}
		if !callbackEqualCmp(expectedToken, gotToken) {
			failure = true
			diffStrs = append(diffStrs, fmt.Sprintf("(index:%d)\tgot:\t\t%v len(%d)", i, callbackPrint(gotToken), len(gotToken.Lexeme)))
			diffStrs = append(diffStrs, fmt.Sprintf("(index:%d)\texpected:\t%v len(%d)", i, callbackPrint(expected[i]), len(expected[i].Lexeme)))
		}
	}

	if failure {
		strGot, strExpect := strings.Builder{}, strings.Builder{}
		for i, expectedToken := range expected {
			var gotToken Token
			if i >= len(got) {
				gotToken = Token{T: TEOF, Lexeme: ``, Literal: nil, Line: 0, Column: 0}
			} else {
				gotToken = got[i]
			}
			strGot.Write([]byte(callbackPrintInline(gotToken)))
			strExpect.Write([]byte(callbackPrintInline(expectedToken)))
		}
		t.Errorf("Diff on lines:")
		for _, diffStr := range diffStrs {
			t.Errorf(diffStr)
		}

		t.Errorf("Whole string diff:")
		t.Errorf("got:\t\t%v", strGot.String())
		t.Errorf("expected:\t%v", strExpect.String())
		t.Fail()
	}
}

func contains(arr []int, val int) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}

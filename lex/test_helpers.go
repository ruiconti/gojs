package lex

import (
	"fmt"
	"strings"
	"testing"

	gojs "github.com/ruiconti/gojs/internal"
)

var defaultLogger = gojs.NewSimpleLogger(gojs.ModeDebug)

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
	)
}

func assertLexemes(t *testing.T, got, expected []Token) {
	assertInternal(
		t,
		got,
		expected, func(a, b Token) bool { return a.Lexeme == b.Lexeme },
		func(a Token) string {
			if a.Lexeme == "" {
				return ""
			}
			return fmt.Sprintf("%v ", a.Lexeme)
		},
	)
}

func assertInternal(t *testing.T, got, expected []Token, callbackEqualCmp func(a, b Token) bool, callbackPrint func(a Token) string) {
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
			strGot.Write([]byte(callbackPrint(gotToken)))
			strExpect.Write([]byte(callbackPrint(expectedToken)))
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

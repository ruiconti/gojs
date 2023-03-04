package lex

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/ruiconti/gojs/internal"
)

// ResolveName returns the name of a token type.
func resolveName(t TokenType) (string, error) {
	dicts := []map[TokenType]string{
		LiteralNames,
		ReservedWordNames,
		PunctuationNames,
	}
	for _, dict := range dicts {
		if name, ok := dict[t]; ok {
			return name, nil
		}
	}
	return "", fmt.Errorf("token name not found: %d", t)
}

func ResolveName(t TokenType) string {
	tokName, err := resolveName(t)
	if err != nil {
		return "Unknown"
	}
	return tokName
}

func assertLiterals(t *testing.T, logger internal.Logger, got, expected []Token) {
	assertInternal(
		t,
		logger,
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

func assertTokens(t *testing.T, logger internal.Logger, got, expected []Token) {
	assertInternal(
		t,
		logger,
		got,
		expected, func(a, b Token) bool { return a.Lexeme == b.Lexeme && a.Literal == b.Literal && a.T == b.T },
		func(a Token) string {
			if a.Literal == nil {
				return ""
			}
			return fmt.Sprintf("%v (%v)", a.Literal, ResolveName(a.T))
		},
		// inline print
		func(a Token) string {
			return fmt.Sprintf("%v ", a.Literal)
		},
	)
}

func assertLexemes(t *testing.T, logger internal.Logger, got, expected []Token) {
	t.Helper()
	assertInternal(
		t,
		logger,
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
	logger internal.Logger,
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
			diffStrs = append(diffStrs, fmt.Sprintf("[%d]\tgot:\t%v", i, callbackPrint(gotToken)))
			diffStrs = append(diffStrs, fmt.Sprintf("[%d]\texp:\t%v", i, callbackPrint(expected[i])))
		}
	}

	if failure {
		strGot, strExpect := strings.Builder{}, strings.Builder{}
		logger.EmitStdout()
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
		t.Errorf("got:\t%v", strGot.String())
		t.Errorf("exp:\t%v", strExpect.String())
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

func assertErrors(t *testing.T, logger internal.Logger, eexp, egot error, in string, out []Token) {
	if !errors.Is(egot, eexp) {
		logger.EmitStdout()
		t.Errorf("expected error %q, got %v", eexp, egot)
		if len(out) > 0 {
			t.Errorf("values for reference: in=%q, out=%v", in, out)
		}
	}
}

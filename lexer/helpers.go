package lexer

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/ruiconti/gojs/internal"
)

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
		expected, func(a, b Token) bool { return a.Lexeme == b.Lexeme && a.Literal == b.Literal && a.Type == b.Type },
		func(a Token) string {
			return fmt.Sprintf("%v (%v)", a.Lexeme, a.Type.S())
		},
		// inline print
		func(a Token) string {
			return fmt.Sprintf("%v ", a.Lexeme)
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
			gotToken = Token{Type: TUnknown, Lexeme: `<EOF>`, Literal: nil, Line: 0, Column: 0}
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
		logger.DumpLogs()
		for i, expectedToken := range expected {
			var gotToken Token
			if i >= len(got) {
				gotToken = Token{Type: TUnknown, Lexeme: ``, Literal: nil, Line: 0, Column: 0}
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

func assertErrors(t *testing.T, logger internal.Logger, exp error, got []error, in string, out []Token) {
	found := false
	for _, e := range got {
		if errors.Is(exp, e) {
			found = true
		}
	}
	if found {
		logger.DumpLogs()
		t.Errorf("expected error %q, got %+v", exp, got)
		if len(out) > 0 {
			t.Errorf("values for reference:\nin=%q,\nout=%v", in, out)
		}
	}
}

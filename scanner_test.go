package main

import (
	"fmt"
	"testing"
)

func TestScanSimplePunctuators(t *testing.T) {
	src := `;():{}[];`

	expected := []Token{
		{T: TSemicolon, Lexeme: ";", Literal: nil, Line: 0, Column: 0},
		{T: TLeftParen, Lexeme: "(", Literal: nil, Line: 0, Column: 0},
		{T: TRightParen, Lexeme: ")", Literal: nil, Line: 0, Column: 0},
		{T: TColon, Lexeme: ":", Literal: nil, Line: 0, Column: 0},
		{T: TLeftBrace, Lexeme: "{", Literal: nil, Line: 0, Column: 0},
		{T: TRightBrace, Lexeme: "}", Literal: nil, Line: 0, Column: 0},
		{T: TLeftBracket, Lexeme: "[", Literal: nil, Line: 0, Column: 0},
		{T: TRightBracket, Lexeme: "]", Literal: nil, Line: 0, Column: 0},
		{T: TSemicolon, Lexeme: ";", Literal: nil, Line: 0, Column: 0},
	}

	scanner := NewScanner(src)
	got := scanner.Scan()
	assertLexemes(t, got, expected)
}

func TestScanDoublePunctuators(t *testing.T) {
	t.Run("tokens > and <", func(t *testing.T) {
		src := `> >= >> >>= >>> >>>= < << <<= <=`
		expected := []Token{
			{T: TGreaterThan, Lexeme: ">", Literal: nil, Line: 0, Column: 0},
			{T: TGreaterThanEqual, Lexeme: ">=", Literal: nil, Line: 0, Column: 0},
			{T: TRightShift, Lexeme: ">>", Literal: nil, Line: 0, Column: 0},
			{T: TRightShiftAssign, Lexeme: ">>=", Literal: nil, Line: 0, Column: 0},
			{T: TUnsignedRightShift, Lexeme: ">>>", Literal: nil, Line: 0, Column: 0},
			{T: TUnsignedRightShiftAssign, Lexeme: ">>>=", Literal: nil, Line: 0, Column: 0},
			{T: TLessThan, Lexeme: "<", Literal: nil, Line: 0, Column: 0},
			{T: TLeftShift, Lexeme: "<<", Literal: nil, Line: 0, Column: 0},
			{T: TLeftShiftAssign, Lexeme: "<<=", Literal: nil, Line: 0, Column: 0},
			{T: TLessThanEqual, Lexeme: "<=", Literal: nil, Line: 0, Column: 0},
		}

		scanner := NewScanner(src)
		got := scanner.Scan()
		assertLexemes(t, got, expected)
	})
	t.Run("operators: ! and =", func(t *testing.T) {
		src := `. ... ? ?? ! != !== = == === =>`
		expected := []Token{
			{T: TPeriod, Lexeme: ".", Literal: nil, Line: 0, Column: 0},
			{T: TEllipsis, Lexeme: "...", Literal: nil, Line: 0, Column: 0},
			{T: TQuestionMark, Lexeme: "?", Literal: nil, Line: 0, Column: 0},
			{T: TDoubleQuestionMark, Lexeme: "??", Literal: nil, Line: 0, Column: 0},
			{T: TBang, Lexeme: "!", Literal: nil, Line: 0, Column: 0},
			{T: TNotEqual, Lexeme: "!=", Literal: nil, Line: 0, Column: 0},
			{T: TStrictNotEqual, Lexeme: "!==", Literal: nil, Line: 0, Column: 0},
			{T: TAssign, Lexeme: "=", Literal: nil, Line: 0, Column: 0},
			{T: TEqual, Lexeme: "==", Literal: nil, Line: 0, Column: 0},
			{T: TStrictEqual, Lexeme: "===", Literal: nil, Line: 0, Column: 0},
			{T: TArrow, Lexeme: "=>", Literal: nil, Line: 0, Column: 0},
		}

		scanner := NewScanner(src)
		got := scanner.Scan()
		assertLexemes(t, got, expected)
	})
	t.Run("operators & and I", func(t *testing.T) {
		src := `& && &= &&= | || |= ||=`
		expected := []Token{
			{T: TAnd, Lexeme: "&", Literal: nil, Line: 0, Column: 0},
			{T: TLogicalAnd, Lexeme: "&&", Literal: nil, Line: 0, Column: 0},
			{T: TAndAssign, Lexeme: "&=", Literal: nil, Line: 0, Column: 0},
			{T: TLogicalAndAssign, Lexeme: "&&=", Literal: nil, Line: 0, Column: 0},
			{T: TOr, Lexeme: "|", Literal: nil, Line: 0, Column: 0},
			{T: TLogicalOr, Lexeme: "||", Literal: nil, Line: 0, Column: 0},
			{T: TOrAssign, Lexeme: "|=", Literal: nil, Line: 0, Column: 0},
			{T: TLogicalOrAssign, Lexeme: "||=", Literal: nil, Line: 0, Column: 0},
		}

		scanner := NewScanner(src)
		got := scanner.Scan()
		assertLexemes(t, got, expected)
	})

	t.Run("operators + and -", func(t *testing.T) {
		src := `+ ++ += - -- -=`
		expected := []Token{
			{T: TPlus, Lexeme: "+", Literal: nil, Line: 0, Column: 0},
			{T: TPlusPlus, Lexeme: "++", Literal: nil, Line: 0, Column: 0},
			{T: TPlusAssign, Lexeme: "+=", Literal: nil, Line: 0, Column: 0},
			{T: TMinus, Lexeme: "-", Literal: nil, Line: 0, Column: 0},
			{T: TMinusMinus, Lexeme: "--", Literal: nil, Line: 0, Column: 0},
			{T: TMinusAssign, Lexeme: "-=", Literal: nil, Line: 0, Column: 0},
		}

		scanner := NewScanner(src)
		got := scanner.Scan()
		assertLexemes(t, got, expected)
	})

	t.Run("operators * and /", func(t *testing.T) {
		src := `* *= / /=`
		expected := []Token{
			{T: TStar, Lexeme: "*", Literal: nil, Line: 0, Column: 0},
			{T: TStarAssign, Lexeme: "*=", Literal: nil, Line: 0, Column: 0},
			{T: TSlash, Lexeme: "/", Literal: nil, Line: 0, Column: 0},
			{T: TSlashAssign, Lexeme: "/=", Literal: nil, Line: 0, Column: 0},
		}

		scanner := NewScanner(src)
		got := scanner.Scan()
		assertLexemes(t, got, expected)
	})

}

func TestScanDigits(t *testing.T) {
	t.Run("decimal literal: single digits", func(t *testing.T) {
		src := `0 1 2 3 4 5 6 7 8 9`
		expected := []Token{
			{T: TNumericLiteral, Lexeme: "0", Literal: "0", Line: 0, Column: 0},
			{T: TNumericLiteral, Lexeme: "1", Literal: "1", Line: 0, Column: 0},
			{T: TNumericLiteral, Lexeme: "2", Literal: "2", Line: 0, Column: 0},
			{T: TNumericLiteral, Lexeme: "3", Literal: "3", Line: 0, Column: 0},
			{T: TNumericLiteral, Lexeme: "4", Literal: "4", Line: 0, Column: 0},
			{T: TNumericLiteral, Lexeme: "5", Literal: "5", Line: 0, Column: 0},
			{T: TNumericLiteral, Lexeme: "6", Literal: "6", Line: 0, Column: 0},
			{T: TNumericLiteral, Lexeme: "7", Literal: "7", Line: 0, Column: 0},
			{T: TNumericLiteral, Lexeme: "8", Literal: "8", Line: 0, Column: 0},
			{T: TNumericLiteral, Lexeme: "9", Literal: "9", Line: 0, Column: 0},
		}

		scanner := NewScanner(src)
		got := scanner.Scan()
		assertLexemes(t, got, expected)
	})

	t.Run("decimal literal: N digits", func(t *testing.T) {
		src := `01345 10 12345678901 100000000000000000000000000000000001`
		expected := []Token{
			{T: TNumericLiteral, Lexeme: "01345", Literal: "01345", Line: 0, Column: 0},
			{T: TNumericLiteral, Lexeme: "10", Literal: "10", Line: 0, Column: 0},
			{T: TNumericLiteral, Lexeme: "12345678901", Literal: "12345678901", Line: 0, Column: 0},
			{T: TNumericLiteral, Lexeme: "100000000000000000000000000000000001", Literal: "100000000000000000000000000000000001", Line: 0, Column: 0},
		}

		scanner := NewScanner(src)
		got := scanner.Scan()
		assertLexemes(t, got, expected)
	})

	// TODO:
	// TODO: Error handling: make sure we don't parse 000111 numbers as valid
	// Err: errInvalidNumber

	t.Run("decimal digit simple: without leading period", func(t *testing.T) {
		src := `10.33340 0.000000001 3939.333393 9999.11100`
		expected := []Token{
			{T: TNumericLiteral, Lexeme: "10.33340", Literal: "10.33340", Line: 0, Column: 0},
			{T: TNumericLiteral, Lexeme: "0.000000001", Literal: "0.000000001", Line: 0, Column: 0},
			{T: TNumericLiteral, Lexeme: "3939.333393", Literal: "3939.333393", Line: 0, Column: 0},
			{T: TNumericLiteral, Lexeme: "9999.11100", Literal: "9999.11100", Line: 0, Column: 0},
		}

		scanner := NewScanner(src)
		got := scanner.Scan()
		assertLexemes(t, got, expected)
	})

	// TODO: Error handling: make sure that we only accept one period while parsing the number
	// e.g 10.00.000.00 as valid numbers
	// Err: errUnexpectedToken

	t.Run("decimal digit complex: with leading period", func(t *testing.T) {
		src := `.0000001 .00000000000000`
		expected := []Token{
			{T: TNumericLiteral, Lexeme: ".0000001", Literal: ".0000001", Line: 0, Column: 0},
			{T: TNumericLiteral, Lexeme: ".00000000000000", Literal: ".00000000000000", Line: 0, Column: 0},
		}

		scanner := NewScanner(src)
		got := scanner.Scan()
		assertLexemes(t, got, expected)
	})

	// TODO: Error handling: make sure that we only accept one period while parsing the number
	// e.g:
	// ....1
	// .00.
	// .0.0.
	// Err: errUnexpectedToken

	t.Run("decimal digit complex: with exp and digits and signs", func(t *testing.T) {
		src := `0e20 .3e25 0.e25 1.e+50 0.3e+50 .1e-50 0e00001 .0e00001 .5E+50 0.E-50`
		expected := []Token{
			{T: TNumericLiteral, Lexeme: "0e20", Literal: "0e20", Line: 0, Column: 0},
			{T: TNumericLiteral, Lexeme: ".3e25", Literal: ".3e25", Line: 0, Column: 0},
			{T: TNumericLiteral, Lexeme: "0.e25", Literal: "0.e25", Line: 0, Column: 0},
			{T: TNumericLiteral, Lexeme: "1.e+50", Literal: "1.e+50", Line: 0, Column: 0},
			{T: TNumericLiteral, Lexeme: "0.3e+50", Literal: "0.3e+50", Line: 0, Column: 0},
			{T: TNumericLiteral, Lexeme: ".1e-50", Literal: ".1e-50", Line: 0, Column: 0},
			{T: TNumericLiteral, Lexeme: "0e00001", Literal: "0e00001", Line: 0, Column: 0},
			{T: TNumericLiteral, Lexeme: ".0e00001", Literal: ".0e00001", Line: 0, Column: 0},
			{T: TNumericLiteral, Lexeme: ".5E+50", Literal: ".5E+50", Line: 0, Column: 0},
			{T: TNumericLiteral, Lexeme: "0.E-50", Literal: "0.E-50", Line: 0, Column: 0},
		}

		scanner := NewScanner(src)
		got := scanner.Scan()
		assertLexemes(t, got, expected)
	})
	// TODO: Error handling: make sure that we only accept one period while parsing the number
	// e.g:
	// e20
	// 01e20
	// 0e20e
	// 0e.20
	// 0e20.
	// .e03
	// .0e20.
	// 0e20++3
	// 0e20--3
	// 0e20+-4
	// 0e20+4.
	// 0e20+.4
	// Err: errUnexpectedToken

	t.Run("hexadecimal digit", func(t *testing.T) {
		src := `0x1 0xA 0x1234567890abcdef 0X1234567890ABCDEF`
		expected := []Token{
			{T: TNumericLiteral, Lexeme: "0x1", Literal: "0x1", Line: 0, Column: 0},
			{T: TNumericLiteral, Lexeme: "0xA", Literal: "0xA", Line: 0, Column: 0},
			{T: TNumericLiteral, Lexeme: "0x1234567890abcdef", Literal: "0x1234567890abcdef", Line: 0, Column: 0},
			{T: TNumericLiteral, Lexeme: "0X1234567890ABCDEF", Literal: "0X1234567890ABCDEF", Line: 0, Column: 0},
		}

		scanner := NewScanner(src)
		got := scanner.Scan()
		assertLexemes(t, got, expected)
	})

	// TODO: Error handling: make sure that we only accept one period while parsing the number
	// e.g:
	// 0x
	// 1x000AB
	// 0x0AABx
	// 0xx303
	// 0Xx23
	// 0X+3
	// 0X-3
	// 0Xe3
}

// /////////////
//
//	UTILS
//
// /////////////
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
	// Need to be in the same order
	failure := false
	if len(got) != len(expected) {
		failure = true
	}
	for i, expectedToken := range expected {
		var gotToken Token
		if i >= len(got) {
			gotToken = Token{T: TEOF, Lexeme: `<EOF>`, Literal: nil, Line: 0, Column: 0}
		} else {
			gotToken = got[i]
		}
		if !callbackEqualCmp(expectedToken, gotToken) {
			failure = true
			t.Errorf("(%d)\tgot:`%v`  expected:`%v`", i, callbackPrint(gotToken), callbackPrint(expected[i]))
		}
	}

	if failure {
		strGot, strExpect := ``, ``
		for i, expectedToken := range expected {
			var gotToken Token
			if i >= len(got) {
				gotToken = Token{T: TEOF, Lexeme: ``, Literal: nil, Line: 0, Column: 0}
			} else {
				gotToken = got[i]
			}
			strGot += callbackPrint(gotToken) + ``
			strExpect += callbackPrint(expectedToken) + ``
		}

		t.Errorf("got:`%v`  expected:`%v`", strGot, strExpect)
		t.Fail()
	}
}

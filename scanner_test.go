package main

import (
	"testing"
)

func TestScanSimplePunctuators(t *testing.T) {
	src := `;():{}[];`

	expected := []Token{
		TSemicolon,
		TLeftParen,
		TRightParen,
		TColon,
		TLeftBrace,
		TRightBrace,
		TLeftBracket,
		TRightBracket,
		TSemicolon,
	}

	got := Scan(src)
	assertTokens(t, got, expected)
}

func TestScanDoublePunctuators(t *testing.T) {
	t.Run("tokens > and <", func(t *testing.T) {
		src := `> >= >> >>= >>> >>>= < << <<= <=`
		expected := []Token{
			TGreaterThan,
			TGreaterThanEqual,
			TRightShift,
			TRightShiftAssign,
			TUnsignedRightShift,
			TUnsignedRightShiftAssign,
			TLessThan,
			TLeftShift,
			TLeftShiftAssign,
			TLessThanEqual,
		}

		got := Scan(src)
		assertTokens(t, got, expected)
	})
	t.Run("operators: ! and =", func(t *testing.T) {
		src := `. ... ? ?? ! != !== = == === =>`
		expected := []Token{
			TPeriod,
			TEllipsis,
			TQuestionMark,
			TDoubleQuestionMark,
			TBang,
			TNotEqual,
			TStrictNotEqual,
			TAssign,
			TEqual,
			TStrictEqual,
			TArrow,
		}

		got := Scan(src)
		assertTokens(t, got, expected)
	})
	t.Run("operators & and |", func(t *testing.T) {
		src := `& && &= &&= | || |= ||=`
		expected := []Token{
			TAnd,
			TLogicalAnd,
			TAndAssign,
			TLogicalAndAssign,
			TOr,
			TLogicalOr,
			TOrAssign,
			TLogicalOrAssign,
		}

		got := Scan(src)
		assertTokens(t, got, expected)
	})

	t.Run("operators + and -", func(t *testing.T) {
		src := `+ ++ += - -- -=`
		expected := []Token{
			TPlus,
			TPlusPlus,
			TPlusAssign,
			TMinus,
			TMinusMinus,
			TMinusAssign,
		}

		got := Scan(src)
		assertTokens(t, got, expected)
	})

	t.Run("operators * and /", func(t *testing.T) {
		src := `* *= / /=`
		expected := []Token{
			TStar,
			TStarAssign,
			TSlash,
			TSlashAssign,
		}
		got := Scan(src)
		assertTokens(t, got, expected)
	})

}

func assertTokens(t *testing.T, got, expected []Token) {
	// Need to be in the same order
	failure := false
	if len(got) != len(expected) {
		failure = true
	}
	for i, gotToken := range got {
		if i >= len(expected) {
			failure = true
			t.Errorf("%d: got %v, expected <nothing>", i, TokenMap[gotToken])
		}
		if gotToken != expected[i] {
			failure = true
			t.Errorf("%d: got %v, expected %v", i, TokenMap[gotToken], TokenMap[expected[i]])
		}
	}

	if failure {
		strGot, strExpect := ``, ``
		for i, gotToken := range got {
			strGot += TokenMap[gotToken] + ` `
			strExpect += TokenMap[expected[i]] + ` `
		}

		t.Errorf("got %v, expected %v", strGot, strExpect)
		t.Fail()
	}
}

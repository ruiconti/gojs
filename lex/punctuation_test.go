package lex

import (
	"testing"
)

func TestScanSimplePunctuators(t *testing.T) {
	src := `;():{}[];,~`

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
		{T: TComma, Lexeme: ",", Literal: nil, Line: 0, Column: 0},
		{T: TTilde, Lexeme: "~", Literal: nil, Line: 0, Column: 0},
	}

	scanner := NewScanner(src, defaultLogger)
	got, _ := scanner.Scan()
	assertLexemes(t, got, expected)
}

func TestScanDoublePunctuators_GThanLThan(t *testing.T) {
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

	scanner := NewScanner(src, defaultLogger)
	got, _ := scanner.Scan()
	assertLexemes(t, got, expected)
}

func TestScanDoublePunctuators_BangEq(t *testing.T) {
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

	scanner := NewScanner(src, defaultLogger)
	got, err := scanner.Scan()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertLexemes(t, got, expected)
}

func TestScanDoublePunctuators_AndOr(t *testing.T) {
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

	scanner := NewScanner(src, defaultLogger)
	got, _ := scanner.Scan()
	assertLexemes(t, got, expected)
}

func TestScanDoublePunctuators_PlusMinus(t *testing.T) {
	src := `+ ++ += - -- -=`
	expected := []Token{
		{T: TPlus, Lexeme: "+", Literal: nil, Line: 0, Column: 0},
		{T: TPlusPlus, Lexeme: "++", Literal: nil, Line: 0, Column: 0},
		{T: TPlusAssign, Lexeme: "+=", Literal: nil, Line: 0, Column: 0},
		{T: TMinus, Lexeme: "-", Literal: nil, Line: 0, Column: 0},
		{T: TMinusMinus, Lexeme: "--", Literal: nil, Line: 0, Column: 0},
		{T: TMinusAssign, Lexeme: "-=", Literal: nil, Line: 0, Column: 0},
	}

	scanner := NewScanner(src, defaultLogger)
	got, _ := scanner.Scan()
	assertLexemes(t, got, expected)
}
func TestScanDoublePunctuators_StarSlash(t *testing.T) {
	src := `* *= / /=`
	expected := []Token{
		{T: TStar, Lexeme: "*", Literal: nil, Line: 0, Column: 0},
		{T: TStarAssign, Lexeme: "*=", Literal: nil, Line: 0, Column: 0},
		{T: TSlash, Lexeme: "/", Literal: nil, Line: 0, Column: 0},
		{T: TSlashAssign, Lexeme: "/=", Literal: nil, Line: 0, Column: 0},
	}

	scanner := NewScanner(src, defaultLogger)
	got, _ := scanner.Scan()
	assertLexemes(t, got, expected)
}

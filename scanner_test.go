package main

import (
	"fmt"
	"strings"
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

// Numeric Literals
//
// | NumericLiteral ::
// | DecimalLiteral
// | DecimalBigIntegerLiteral
// | NonDecimalIntegerLiteral[+Sep]
// | NonDecimalIntegerLiteral[+Sep] BigIntLiteralSuffix
// | LegacyOctalIntegerLiteral
func TestScanDigits_DecimalLiterals(t *testing.T) {
	// TODO: Add failing tests:
	// Leading 0s, >1 punctuator, mixing punctuators, etc.
	// Err: errUnexpectedToken
	t.Run("1P: DecimalIntegerLiteral . DecimalDigits? ExponentialPart?", func(t *testing.T) {
		// DecimalIntegerLiteral ::
		// | 0
		// | NonZeroDigit DecimalDigits?
		// | NonZeroDigit NumericLiteralSeparator? DecimalDigits
		// | NonOctalDecimalIntegerLiteral ((tested below))
		//
		// ExponentialPart ::
		// | e SignedInteger
		// | E SignedInteger
		//
		// SignedInteger ::
		// | + DecimalDigits
		// | - DecimalDigits
		//
		// DecimalDigits[Sep] ::
		// | DecimalDigit
		// | DecimalDigits DecimalDigit
		// | DecimalDigits NumericLiteralSeparator DecimalDigit
		//
		// DecimalDigit :: one of 0 1 2 3 4 5 6 7 8 9
		// NonZeroDigit :: one of 1 2 3 4 5 6 7 8 9
		src := `0 10.33340 0.000000001 3939.333393 9999.11100 10_000_000 10_0.30_0 9.30_0E30_034 0e20 0.e25 1.e+50 0.3e+50 0e00001 0.E-50`
		expected := []Token{
			{T: TNumericLiteral, Lexeme: "0", Literal: "0", Line: 0, Column: 0},
			{T: TNumericLiteral, Lexeme: "10.33340", Literal: "10.33340", Line: 0, Column: 0},
			{T: TNumericLiteral, Lexeme: "0.000000001", Literal: "0.000000001", Line: 0, Column: 0},
			{T: TNumericLiteral, Lexeme: "3939.333393", Literal: "3939.333393", Line: 0, Column: 0},
			{T: TNumericLiteral, Lexeme: "9999.11100", Literal: "9999.11100", Line: 0, Column: 0},
			{T: TNumericLiteral, Lexeme: "10_000_000", Literal: "10_000_000", Line: 0, Column: 0},
			{T: TNumericLiteral, Lexeme: "10_0.30_0", Literal: "10_0.30_0", Line: 0, Column: 0},
			{T: TNumericLiteral, Lexeme: "9.30_0E30_034", Literal: "9.30_0E30_034", Line: 0, Column: 0},
			{T: TNumericLiteral, Lexeme: "0e20", Literal: "0e20", Line: 0, Column: 0},
			{T: TNumericLiteral, Lexeme: "0.e25", Literal: "0.e25", Line: 0, Column: 0},
			{T: TNumericLiteral, Lexeme: "1.e+50", Literal: "1.e+50", Line: 0, Column: 0},
			{T: TNumericLiteral, Lexeme: "0.3e+50", Literal: "0.3e+50", Line: 0, Column: 0},
			{T: TNumericLiteral, Lexeme: "0e00001", Literal: "0e00001", Line: 0, Column: 0},
			{T: TNumericLiteral, Lexeme: "0.E-50", Literal: "0.E-50", Line: 0, Column: 0},
		}

		scanner := NewScanner(src)
		got := scanner.Scan()
		assertLexemes(t, got, expected)

		// NonOctalDecimalIntegerLiteral ::
		// | 0 NonOctalDigit
		// | LegacyOctalLikeDecimalIntegerLiteral  NonOctalDigit
		// | NonOctalDecimalIntegerLiteral DecimalDigit
		//
		// LegacyOctalLikeDecimalIntegerLiteral ::
		// | 0 OctalDigit
		// | LegacyOctalLikeDecimalIntegerLiteral OctalDigit
		//
		// NonOctalDigit :: one of 8 9
		// OctalDigit :: one of 0 1 2 3 4 5 6 7
		//
		// (We could simplify the production in one form)
		// NonOctalDecimalIntegerLiteral ::
		// | 0 DecimalDigit
		// | NonOctalDecimalIntegerLiteral DecimalDigit
		//
		// FYI: These are not valid in strict mode, though we need to be overly permissive
		src = `0000008989 01234567 0777`
		expected = []Token{
			{T: TNumericLiteral, Lexeme: "0000008989", Literal: "0000008989", Line: 0, Column: 0},
			{T: TNumericLiteral, Lexeme: "01234567", Literal: "01234567", Line: 0, Column: 0},
			{T: TNumericLiteral, Lexeme: "0777", Literal: "0777", Line: 0, Column: 0},
		}

	})

	t.Run("2P: . DecimalDigits ExponentialPart?", func(t *testing.T) {
		src := `.33340 0.0000_0000_1 .3_0E0_034 .1e+2_0 .3e-2_5 .0000e25 .1e-50 .5E+50 .9E-50`
		expected := []Token{
			{T: TNumericLiteral, Lexeme: ".33340", Literal: ".33340", Line: 0, Column: 0},
			{T: TNumericLiteral, Lexeme: "0.0000_0000_1", Literal: "0.0000_0000_1", Line: 0, Column: 0},
			{T: TNumericLiteral, Lexeme: ".3_0E0_034", Literal: ".3_0E0_034", Line: 0, Column: 0},
			{T: TNumericLiteral, Lexeme: ".1e+2_0", Literal: ".1e+2_0", Line: 0, Column: 0},
			{T: TNumericLiteral, Lexeme: ".3e-2_5", Literal: ".3e-2_5", Line: 0, Column: 0},
			{T: TNumericLiteral, Lexeme: ".0000e25", Literal: ".0000e25", Line: 0, Column: 0},
			{T: TNumericLiteral, Lexeme: ".1e-50", Literal: ".1e-50", Line: 0, Column: 0},
			{T: TNumericLiteral, Lexeme: ".5E+50", Literal: ".5E+50", Line: 0, Column: 0},
			{T: TNumericLiteral, Lexeme: ".9E-50", Literal: ".9E-50", Line: 0, Column: 0},
		}
		scanner := NewScanner(src)
		got := scanner.Scan()
		assertLexemes(t, got, expected)
	})

	t.Run("3P: DecimalIntegerLiteral ExponentialPart", func(t *testing.T) {
		// DecimalIntegerLiteral ::
		// | 0
		// | NonZeroDigit
		// | NonZeroDigit NumericLiteralSeparator? DecimalDigits
		// | NonOctalDecimalIntegerLiteral
		//
		src := `0_E25 1_35E-50_0 00000E-50_00000 000000e000000 007654321e+1 000e+1`
		expected := []Token{
			{T: TNumericLiteral, Lexeme: "0_E25", Literal: "0_E25", Line: 0, Column: 0},
			{T: TNumericLiteral, Lexeme: "1_35E-50_0", Literal: "1_35E-50_0", Line: 0, Column: 0},
			{T: TNumericLiteral, Lexeme: "00000E-50_00000", Literal: "00000E-50_00000", Line: 0, Column: 0},
			{T: TNumericLiteral, Lexeme: "000000e000000", Literal: "000000e000000", Line: 0, Column: 0},
			{T: TNumericLiteral, Lexeme: "007654321e+1", Literal: "007654321e+1", Line: 0, Column: 0},
			{T: TNumericLiteral, Lexeme: "000e+1", Literal: "000e+1", Line: 0, Column: 0},
		}
		scanner := NewScanner(src)
		got := scanner.Scan()
		assertLexemes(t, got, expected)
	})
}

func TestScanDigits_HexIntegerLiteral(t *testing.T) {
	// HexIntegerLiteral ::
	// | 0x HexDigits
	// | 0X HexDigits
	src := `0x1 0xA 0x1234567890abcdef 0X1234567890ABCDEF 0xB_AAB_445`
	expected := []Token{
		{T: TNumericLiteral, Lexeme: "0x1", Literal: "0x1", Line: 0, Column: 0},
		{T: TNumericLiteral, Lexeme: "0xA", Literal: "0xA", Line: 0, Column: 0},
		{T: TNumericLiteral, Lexeme: "0x1234567890abcdef", Literal: "0x1234567890abcdef", Line: 0, Column: 0},
		{T: TNumericLiteral, Lexeme: "0X1234567890ABCDEF", Literal: "0X1234567890ABCDEF", Line: 0, Column: 0},
		{T: TNumericLiteral, Lexeme: "0xB_AAB_445", Literal: "0xB_AAB_445", Line: 0, Column: 0},
	}

	scanner := NewScanner(src)
	got := scanner.Scan()
	assertLexemes(t, got, expected)

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

func TestScanDigits_DecimalBigIntegerLiteral(t *testing.T) {
	// DecimalBigIntegerLiteral ::
	// | 0 BigIntLiteralSuffix
	// | NonZeroDigit DecimalDigits? BigIntLiteralSuffix
	// | NonZeroDigit NumericLiteralSeparator DecimalDigits BigIntLiteralSuffix
	//
	// BigIntLiteralSuffix :: n
	src := `0n 8n 84981283n 1_923_921_839_1273n`
	expected := []Token{
		{T: TNumericLiteral, Lexeme: "0n", Literal: "0n", Line: 0, Column: 0},
		{T: TNumericLiteral, Lexeme: "8n", Literal: "8n", Line: 0, Column: 0},
		{T: TNumericLiteral, Lexeme: "84981283n", Literal: "84981283n", Line: 0, Column: 0},
		{T: TNumericLiteral, Lexeme: "1_923_921_839_1273n", Literal: "1_923_921_839_1273n", Line: 0, Column: 0},
	}

	scanner := NewScanner(src)
	got := scanner.Scan()
	assertLexemes(t, got, expected)
}

func TestScanDigits_BinaryIntegerLiteral(t *testing.T) {
	// BinaryIntegerLiteral ::
	// | 0b BinaryDigits
	// | 0B BinaryDigits
	//
	// BinaryDigits ::
	// | BinaryDigit BinaryDigits
	// | BinaryDigits NumericLiteralSeparator BinaryDigit
	//
	// BinaryDigit :: one of 0 1
	//
	src := `0b0 0b1 0B0 0B1 0B0101010 0b101010 0b1010_0101_0110 0b0100_0101_0110`
	expected := []Token{
		{T: TNumericLiteral, Lexeme: "0b0", Literal: "0b0", Line: 0, Column: 0},
		{T: TNumericLiteral, Lexeme: "0b1", Literal: "0b1", Line: 0, Column: 0},
		{T: TNumericLiteral, Lexeme: "0B0", Literal: "0B0", Line: 0, Column: 0},
		{T: TNumericLiteral, Lexeme: "0B1", Literal: "0B1", Line: 0, Column: 0},
		{T: TNumericLiteral, Lexeme: "0B0101010", Literal: "0B0101010", Line: 0, Column: 0},
		{T: TNumericLiteral, Lexeme: "0b101010", Literal: "0b101010", Line: 0, Column: 0},
		{T: TNumericLiteral, Lexeme: "0b1010_0101_0110", Literal: "0b1010_0101_0110", Line: 0, Column: 0},
		{T: TNumericLiteral, Lexeme: "0b0100_0101_0110", Literal: "0b0100_0101_0110", Line: 0, Column: 0},
	}

	scanner := NewScanner(src)
	got := scanner.Scan()
	assertLexemes(t, got, expected)
}

func TestScanDigits_OctalIntegerLiteral(t *testing.T) {
	// OctalIntegerLiteral ::
	// | 0o OctalDigits
	// | 0O OctalDigits
	//
	// OctalDigits ::
	// | OctalDigit
	// | OctalDigits OctalDigit
	// | OctalDigits NumericLiteralSeparator OctalDigit
	//
	// OctalDigit :: one of 0 1 2 3 4 5 6 7
	src := `0o0 0o1 0O0 0O7 0O6 0O2112_2234_6670 0o1234_5672_5012`
	expected := []Token{
		{T: TNumericLiteral, Lexeme: "0o0", Literal: "0o0", Line: 0, Column: 0},
		{T: TNumericLiteral, Lexeme: "0o1", Literal: "0o1", Line: 0, Column: 0},
		{T: TNumericLiteral, Lexeme: "0O0", Literal: "0O0", Line: 0, Column: 0},
		{T: TNumericLiteral, Lexeme: "0O7", Literal: "0O7", Line: 0, Column: 0},
		{T: TNumericLiteral, Lexeme: "0O6", Literal: "0O6", Line: 0, Column: 0},
		{T: TNumericLiteral, Lexeme: "0O2112_2234_6670", Literal: "0O2112_2234_6670", Line: 0, Column: 0},
		{T: TNumericLiteral, Lexeme: "0o1234_5672_5012", Literal: "0o1234_5672_5012", Line: 0, Column: 0},
	}

	scanner := NewScanner(src)
	got := scanner.Scan()
	assertLexemes(t, got, expected)
}

// String Literals
//
// StringLiteral ::
// | " DoubleString? "
// | ' SingleString? '
func TestScanString_DoubleString_1P(t *testing.T) {
	// DoubleStringCharacters ::
	// | DoubleStringCharacter DoubleStringCharactersopt
	//
	// DoubleStringCharacter ::
	// | SourceCharacter but not one of " or \ or LineTerminator
	// | <LS>
	// | <PS>
	// | \ EscapeSequence
	// | LineContinuation
	//
	// SourceCharacter :: any Unicode code point
	// All Unicode code point values from U+0000 to U+10FFFF
	//
	// LineContinuation ::
	// | \ LineTerminatorSequence
	//
	// LineTerminatorSequence ::
	// | <LF>
	// | <CR> [lookahead ≠ <LF>]
	// | <LS>
	// | <PS>
	// | <CR> <LF>
	//
	// LineTerminator ::
	// | <LF>
	// | <CR>
	// | <LS>
	// | <PS>
	// t.Skip()
	src := `"   " "abcdefghijklmnopqrstuvwxyz" "ABCDEFGHIJKLMNOPQRSTUVWXYZ" "0123456789" "!" "#" "\na\n\n$\n\r\n\t\n\v\n\f"`
	expected := []Token{
		{T: TStringLiteral, Lexeme: `"   "`, Literal: `""`, Line: 0, Column: 0},
		{T: TStringLiteral, Lexeme: `"abcdefghijklmnopqrstuvwxyz"`, Literal: `"abcdefghijklmnopqrstuvwxyz"`, Line: 0, Column: 0},
		{T: TStringLiteral, Lexeme: `"ABCDEFGHIJKLMNOPQRSTUVWXYZ"`, Literal: `"ABCDEFGHIJKLMNOPQRSTUVWXYZ"`, Line: 0, Column: 0},
		{T: TStringLiteral, Lexeme: `"0123456789"`, Literal: `"0123456789"`, Line: 0, Column: 0},
		{T: TStringLiteral, Lexeme: `"!"`, Literal: `"!"`, Line: 0, Column: 0},
		{T: TStringLiteral, Lexeme: `"#"`, Literal: `"#"`, Line: 0, Column: 0},
		{T: TStringLiteral, Lexeme: `"\na\n\n$\n\r\n\t\n\v\n\f"`, Literal: `"\na\n\n$\n\r\n\t\n\v\n\f"`, Line: 0, Column: 0},
	}
	scanner := NewScanner(src)
	got := scanner.Scan()
	assertLexemes(t, got, expected)
}

func TestScanString_DoubleString_Escapes(t *testing.T) {
	// DoubleStringCharacter ::
	// | \ EscapeSequence
	//
	// EscapeSequence ::
	// | CharacterEscapeSequence
	// | 0 [lookahead ∉ DecimalDigit]
	// | LegacyOctalEscapeSequence
	// | NonOctalDecimalEscapeSequence
	// | HexEscapeSequence
	// | UnicodeEscapeSequence
	//
	// CharacterEscapeSequence ::
	// | SingleEscapeCharacter
	// | NonEscapeCharacter
	//
	// SingleEscapeCharacter :: one of ' " \ b f n r t v
	// NonEscapeCharacter :: any Unicode point but not one of EscapeCharacter or LineTerminator
	//
	// LegacyOctalEscapeSequence ::
	// | 0 [lookahead ∈ { 8, 9 }]
	// | NonZeroOctalDigit [lookahead ∉ OctalDigit]
	// | ZeroToThree OctalDigit [lookahead ∉ OctalDigit]
	// | FourToSeven OctalDigit
	// | ZeroToThree OctalDigit OctalDigit
	//
	// NonOctalDecimalEscapeSequence :: one of 8 9
	//
	// HexEscapeSequence ::
	// x HexDigit HexDigit
	//
	// UnicodeEscapeSequence ::
	// | u Hex4Digits
	// | u{ CodePoint }
	src := `"\\\\\\" "\"\'\\a\b\c\d\e\f\g\h\i\j\k\l\m\n\o\p\q\r\s\t\v\w\y\z" "\00\10\20\30\40\50\60\70\80\90" "\x10\x20\x30\x40\x50\x60\x70\x80\x90\xA0\xB0\xC0\xD0\xE0\xF0" "\u0000\u0001\u00005\u99999"`
	expected := []Token{
		{T: TStringLiteral, Lexeme: `"\\\\\\"`, Literal: `"\\\\\\"`, Line: 0, Column: 0},
		{T: TStringLiteral, Lexeme: `"\"\'\\a\b\c\d\e\f\g\h\i\j\k\l\m\n\o\p\q\r\s\t\v\w\y\z"`, Literal: `"\a\b\c\d\e\f\g\h\i\j\k\l\m\n\o\p\q\r\s\t\v\w\y\z"`, Line: 0, Column: 0},
		{T: TStringLiteral, Lexeme: `"\00\10\20\30\40\50\60\70\80\90"`, Literal: `"\00\10\20\30\40\50\60\70\80\90"`, Line: 0, Column: 0},
		{T: TStringLiteral, Lexeme: `"\x10\x20\x30\x40\x50\x60\x70\x80\x90\xA0\xB0\xC0\xD0\xE0\xF0"`, Literal: `"\x10\x20\x30\x40\x50\x60\x70\x80\x90\xA0\xB0\xC0\xD0\xE0\xF0"`, Line: 0, Column: 0},
		{T: TStringLiteral, Lexeme: `"\u0000\u0001\u00005\u99999"`, Literal: `"\u0000\u0001\u00005\u99999"`, Line: 0, Column: 0},
	}
	scanner := NewScanner(src)
	got := scanner.Scan()
	assertLexemes(t, got, expected)
}

// SingleStringCharacters ::
// | SingleStringCharacter SingleStringCharactersopt
//
// SingleStringCharacter ::
// | SourceCharacter but not one of ' or \ or LineTerminator
// | <LS>
// | <PS>
// | \ EscapeSequence
// | LineContinuation
func TestStringLiteral_Single(t *testing.T) {
}

// Template Literals
//
// Template ::
// | NoSubstitutionTemplate
// | TemplateHead
func TestScanTemplateLiteral_NoSubstitutionTemplate(t *testing.T) {
	t.Skip("NotImplemented")
	src := "`$end$` `$$$` `        ` `scan\n\nthis\n\ntoo!` `bla\n\n\nbla`"
	expected := []Token{
		{T: TTemplateLiteral, Lexeme: "`$end$`", Literal: "$end$", Line: 0, Column: 0},
		{T: TTemplateLiteral, Lexeme: "`$$$`", Literal: "$$$", Line: 0, Column: 0},
		{T: TTemplateLiteral, Lexeme: "`        `", Literal: "        ", Line: 0, Column: 0},
		{T: TTemplateLiteral, Lexeme: "`scan\n\nthis\n\ntoo!`", Literal: "scan\n\nthis\n\ntoo!", Line: 0, Column: 0},
		{T: TTemplateLiteral, Lexeme: "`bla\n\n\nbla`", Literal: "bla\n\n\nbla", Line: 0, Column: 0},
	}
	scanner := NewScanner(src)
	got := scanner.Scan()
	assertLexemes(t, got, expected)
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
	t.Helper()
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
			t.Errorf("(index:%d)\tgot:\t\t%v", i, callbackPrint(gotToken))
			t.Errorf("(index:%d)\texpected:\t%v", i, callbackPrint(expected[i]))
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

		t.Errorf("got:\t\t%v", strGot.String())
		t.Errorf("expected:\t%v", strExpect.String())
		t.Fail()
	}
}

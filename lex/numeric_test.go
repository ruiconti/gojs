package lex

import (
	"errors"
	"testing"
)

// Numeric Literals
//
// | NumericLiteral ::
// | DecimalLiteral
// | DecimalBigIntegerLiteral
// | NonDecimalIntegerLiteral[+Sep]
// | NonDecimalIntegerLiteral[+Sep] BigIntLiteralSuffix
// | LegacyOctalIntegerLiteral
//
// TODO: Add failing tests:
// Leading 0s, >1 punctuator, mixing punctuators, etc.
// Err: errUnexpectedToken
func TestScanDigits_DecimalLiterals_Prod1(t *testing.T) {
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

	scanner := NewScanner(src, defaultLogger)
	got, err := scanner.Scan()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
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

}

func TestScanDigits_DecimalLiterals_FirstProd_Err(t *testing.T) {
	src := []string{
		"0a01",
		"1b01",
		"0c01",
		"0d01",
		"0f01",
		"0g01",
		"0h01",
		"0i01",
		"0j01",
		"0k01",
		"0l01",
		"0m01",
		// "1n01n", // TODO: Fix this case
		"1o01",
		"0p01",
		"0q01",
		"0r01",
		"0s01",
		"0t01",
		"0u01",
		"0v01",
		"0w01",
		"1x01",
		"0y01",
		"0z01",
		"0A01",
		"00a01",
		"00b01",
		"00c01",
		"00d01",
		"00f01",
		"00g01",
		"00h01",
		"00i01",
		"00j01",
		"00k01",
		"00l01",
		"00m01",
		// "00n01", // TODO: Fix this case
		"00o01",
		"00p01",
		"00q01",
		"00r01",
		"00s01",
		"00t01",
		"00u01",
		"00v01",
		"00w01",
		"00x01",
		"00y01",
		"00z01",
		"00A01",
	}

	for i, s := range src {
		scanner := NewScanner(s, defaultLogger)
		got, err := scanner.Scan()
		if !errors.Is(err, errNoLiteralAfterNumber) {
			t.Errorf("src:%s got:%v expected error %q, got %v", src[i], got, errNoLiteralAfterNumber, err)
		}
	}

}

func TestScanDigits_DecimalLiterals_Prod2_Err(t *testing.T) {
	src := []string{
		"10e",
		"10E",
		".0e",
		"0.e",
		".e",
		"10_e",
		"10_",
		"10_.",
		"10e.",
		"10+",
		"10-",
		"._3",
		".+3",
		".-3",
		"10-0",
		"10+0",
		"10e++3",
		"10e--3", // TODO: That's a bad case -- 10-- is legal and should leave the digit parser
	}

	for i, s := range src {
		scanner := NewScanner(s, defaultLogger)
		got, err := scanner.Scan()
		if !errors.Is(err, errDigitExpected) {
			t.Errorf("src:%s got:%v expected error %q, got %v", src[i], got, errNoLiteralAfterNumber, err)
		}
	}

}

func TestScanDigits_DecimalLiterals_Prod2(t *testing.T) {
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
	scanner := NewScanner(src, defaultLogger)
	got, _ := scanner.Scan()
	assertLexemes(t, got, expected)
}
func TestScanDigits_DecimalLiterals_Prod3(t *testing.T) {
	// DecimalIntegerLiteral ::
	// | 0
	// | NonZeroDigit
	// | NonZeroDigit NumericLiteralSeparator? DecimalDigits
	// | NonOctalDecimalIntegerLiteral
	//
	src := `1_35E-50_0 00000E-50_00000 000000e000000 007654321e+1 000e+1`
	expected := []Token{
		// {T: TNumericLiteral, Lexeme: "0_E25", Literal: "0_E25", Line: 0, Column: 0},
		{T: TNumericLiteral, Lexeme: "1_35E-50_0", Literal: "1_35E-50_0", Line: 0, Column: 0},
		{T: TNumericLiteral, Lexeme: "00000E-50_00000", Literal: "00000E-50_00000", Line: 0, Column: 0},
		{T: TNumericLiteral, Lexeme: "000000e000000", Literal: "000000e000000", Line: 0, Column: 0},
		{T: TNumericLiteral, Lexeme: "007654321e+1", Literal: "007654321e+1", Line: 0, Column: 0},
		{T: TNumericLiteral, Lexeme: "000e+1", Literal: "000e+1", Line: 0, Column: 0},
	}
	scanner := NewScanner(src, defaultLogger)
	got, err := scanner.Scan()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertLexemes(t, got, expected)
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

	scanner := NewScanner(src, defaultLogger)
	got, _ := scanner.Scan()
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

	scanner := NewScanner(src, defaultLogger)
	got, _ := scanner.Scan()
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

	scanner := NewScanner(src, defaultLogger)
	got, _ := scanner.Scan()
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

	scanner := NewScanner(src, defaultLogger)
	got, _ := scanner.Scan()
	assertLexemes(t, got, expected)
}

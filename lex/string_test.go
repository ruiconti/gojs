package lex

import "testing"

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
	scanner := NewScanner(src, defaultLogger)
	got, _ := scanner.Scan()
	assertLexemes(t, got, expected)
}

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
func TestScanString_DoubleString_Escapes(t *testing.T) {
	src := `"\\\\\\" "\"\'\\a\b\c\d\e\f\g\h\i\j\k\l\m\n\o\p\q\r\s\t\v\w\y\z" "\00\10\20\30\40\50\60\70\80\90" "\x10\x20\x30\x40\x50\x60\x70\x80\x90\xA0\xB0\xC0\xD0\xE0\xF0" "\u0000\u0001\u00005\u99999"`
	expected := []Token{
		{T: TStringLiteral, Lexeme: `"\\\\\\"`, Literal: `"\\\\\\"`, Line: 0, Column: 0},
		{T: TStringLiteral, Lexeme: `"\"\'\\a\b\c\d\e\f\g\h\i\j\k\l\m\n\o\p\q\r\s\t\v\w\y\z"`, Literal: `"\a\b\c\d\e\f\g\h\i\j\k\l\m\n\o\p\q\r\s\t\v\w\y\z"`, Line: 0, Column: 0},
		{T: TStringLiteral, Lexeme: `"\00\10\20\30\40\50\60\70\80\90"`, Literal: `"\00\10\20\30\40\50\60\70\80\90"`, Line: 0, Column: 0},
		{T: TStringLiteral, Lexeme: `"\x10\x20\x30\x40\x50\x60\x70\x80\x90\xA0\xB0\xC0\xD0\xE0\xF0"`, Literal: `"\x10\x20\x30\x40\x50\x60\x70\x80\x90\xA0\xB0\xC0\xD0\xE0\xF0"`, Line: 0, Column: 0},
		{T: TStringLiteral, Lexeme: `"\u0000\u0001\u00005\u99999"`, Literal: `"\u0000\u0001\u00005\u99999"`, Line: 0, Column: 0},
	}
	scanner := NewScanner(src, defaultLogger)
	got, _ := scanner.Scan()
	assertLexemes(t, got, expected)
}

// SingleStringCharacters ::
// | SingleStringCharacter SingleStringCharacters?
//
// SingleStringCharacter ::
// | SourceCharacter but not one of ' or \ or LineTerminator
// | <LS>
// | <PS>
// | \ EscapeSequence
// | LineContinuation
func TestStringLiteral_Single(t *testing.T) {
	src := `'   ' '\\\\\\\x01' '\\\\\\\u01' '"\'\\a\b\c\d\e\f\g\h\i\j\k\l\m\n\o\p\q\r\s\t\v\w\y\z' '' '\u0000\u0001\u00005\u99999'`

	expected := []Token{
		{T: TStringLiteral, Lexeme: `'   '`, Literal: `'   '`, Line: 0, Column: 0},
		{T: TStringLiteral, Lexeme: `'\\\\\\\x01'`, Literal: `\\\\\\\x01'`, Line: 0, Column: 0},
		{T: TStringLiteral, Lexeme: `'\\\\\\\u01'`, Literal: `\\\\\\\u01'`, Line: 0, Column: 0},
		{T: TStringLiteral, Lexeme: `'"\'\\a\b\c\d\e\f\g\h\i\j\k\l\m\n\o\p\q\r\s\t\v\w\y\z'`, Literal: `'"\'\\a\b\c\d\e\f\g\h\i\j\k\l\m\n\o\p\q\r\s\t\v\w\y\z'`, Line: 0, Column: 0},
		{T: TStringLiteral, Lexeme: `''`, Literal: `''`, Line: 0, Column: 0},
		{T: TStringLiteral, Lexeme: `'\u0000\u0001\u00005\u99999'`, Literal: `'\u0000\u0001\u00005\u99999'`, Line: 0, Column: 0},
	}
	scanner := NewScanner(src, defaultLogger)
	got, _ := scanner.Scan()
	assertLexemes(t, got, expected)
}

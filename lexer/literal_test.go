package lexer

import (
	"testing"

	gojs "github.com/ruiconti/gojs/internal"
)

// String Literals
//
// StringLiteral ::
// | " DoubleString? "
// | ' SingleString? '
func TestString_DoubleQuote(t *testing.T) {
	src := `"" "   " "abcdefghijklmnopqrstuvwxyz" "ABCDEFGHIJKLMNOPQRSTUVWXYZ" "0123456789" "!" "#" "\na\n\n$\n\r\n\t\n\v\n\f"`
	expected := []Token{
		{Type: TStringLiteral_DoubleQuote, Lexeme: `""`, Literal: `""`, Line: 0, Column: 0},
		{Type: TStringLiteral_DoubleQuote, Lexeme: `"   "`, Literal: `"   "`, Line: 0, Column: 0},
		{Type: TStringLiteral_DoubleQuote, Lexeme: `"abcdefghijklmnopqrstuvwxyz"`, Literal: `"abcdefghijklmnopqrstuvwxyz"`, Line: 0, Column: 0},
		{Type: TStringLiteral_DoubleQuote, Lexeme: `"ABCDEFGHIJKLMNOPQRSTUVWXYZ"`, Literal: `"ABCDEFGHIJKLMNOPQRSTUVWXYZ"`, Line: 0, Column: 0},
		{Type: TStringLiteral_DoubleQuote, Lexeme: `"0123456789"`, Literal: `"0123456789"`, Line: 0, Column: 0},
		{Type: TStringLiteral_DoubleQuote, Lexeme: `"!"`, Literal: `"!"`, Line: 0, Column: 0},
		{Type: TStringLiteral_DoubleQuote, Lexeme: `"#"`, Literal: `"#"`, Line: 0, Column: 0},
		{Type: TStringLiteral_DoubleQuote, Lexeme: `"\na\n\n$\n\r\n\t\n\v\n\f"`, Literal: `"\na\n\n$\n\r\n\t\n\v\n\f"`, Line: 0, Column: 0},
	}

	logger := gojs.NewSimpleLogger(gojs.ModeDebug)
	lexer := NewLexer(src, logger)
	got, err := lexer.ScanAll()
	if err != nil {
		logger.DumpLogs()
		t.Fatalf("unexpected error: %v", err)
	}
	assertLexemes(t, logger, got, expected)
}

// StringLiteral ::
// | " (SourceCharacter | \ EscapeSequence)? "
func TestString_DoubleQuote_Escaped(t *testing.T) {
	src := `"\\\\" "\"\'\\a\b\c\d\e\f\g\h\i\j\k\l\m\n\o\p\q\r\s\t\v\w\y\z" "\00\01\02\03\04\05\06\07\08\09" "\x10\x20\x30\x40\x50\x60\x70\x80\x90\xA0\xB0\xC0\xD0\xE0\xF0" "\u0000\u0001\u0005\u9999"`
	expected := []Token{
		{Type: TStringLiteral_DoubleQuote, Lexeme: `"\\\\"`, Literal: `"\\\\"`, Line: 0, Column: 0},
		{Type: TStringLiteral_DoubleQuote, Lexeme: `"\"\'\\a\b\c\d\e\f\g\h\i\j\k\l\m\n\o\p\q\r\s\t\v\w\y\z"`, Literal: `"\a\b\c\d\e\f\g\h\i\j\k\l\m\n\o\p\q\r\s\t\v\w\y\z"`, Line: 0, Column: 0},
		{Type: TStringLiteral_DoubleQuote, Lexeme: `"\00\01\02\03\04\05\06\07\08\09"`, Literal: `"\00\01\02\03\04\05\06\07\08\09"`, Line: 0, Column: 0},
		{Type: TStringLiteral_DoubleQuote, Lexeme: `"\x10\x20\x30\x40\x50\x60\x70\x80\x90\xA0\xB0\xC0\xD0\xE0\xF0"`, Literal: `"\x10\x20\x30\x40\x50\x60\x70\x80\x90\xA0\xB0\xC0\xD0\xE0\xF0"`, Line: 0, Column: 0},
		{Type: TStringLiteral_DoubleQuote, Lexeme: `"\u0000\u0001\u0005\u9999"`, Literal: `"\u0000\u0001\u0005\u9999"`, Line: 0, Column: 0},
	}
	logger := gojs.NewSimpleLogger(gojs.ModeDebug)
	lexer := NewLexer(src, logger)
	got, err := lexer.ScanAll()
	if err != nil {
		logger.DumpLogs()
		t.Fatalf("unexpected error: %v", err)
	}
	assertLexemes(t, logger, got, expected)
}

// SingleStringCharacters ::
// | (SourceCharacter | \ EscapeSequence | LineContinuation)*
func TestString_SingleQuote_Escaped(t *testing.T) {
	src := `'   ' '\\\\\\\x01' '\\\\\\\u011a' '"\'\\a\b\c\d\e\f\g\h\i\j\k\l\m\n\o\p\q\r\s\t\v\w\y\z' '' '\u0000\u0001\u00005\u99999'`

	expected := []Token{
		{Type: TStringLiteral_SingleQuote, Lexeme: `'   '`, Literal: `'   '`, Line: 0, Column: 0},
		{Type: TStringLiteral_SingleQuote, Lexeme: `'\\\\\\\x01'`, Literal: `'\\\\\\\x01'`, Line: 0, Column: 0},
		{Type: TStringLiteral_SingleQuote, Lexeme: `'\\\\\\\u011a'`, Literal: `'\\\\\\\u011a'`, Line: 0, Column: 0},
		{Type: TStringLiteral_SingleQuote, Lexeme: `'"\'\\a\b\c\d\e\f\g\h\i\j\k\l\m\n\o\p\q\r\s\t\v\w\y\z'`, Literal: `"\'\\a\b\c\d\e\f\g\h\i\j\k\l\m\n\o\p\q\r\s\t\v\w\y\z'`, Line: 0, Column: 0},
		{Type: TStringLiteral_SingleQuote, Lexeme: `''`, Literal: `''`, Line: 0, Column: 0},
		{Type: TStringLiteral_SingleQuote, Lexeme: `'\u0000\u0001\u00005\u99999'`, Literal: `'\u0000\u0001\u00005\u99999'`, Line: 0, Column: 0},
	}
	logger := gojs.NewSimpleLogger(gojs.ModeDebug)
	lexer := NewLexer(src, logger)
	got, err := lexer.ScanAll()
	if err != nil {
		logger.DumpLogs()
		t.Fatalf("unexpected error: %v", err)
	}
	assertLexemes(t, logger, got, expected)
}

// NumericLiteral ::
// | DecimalLiteral
// TODO: Add failing tests:
func TestLiteral_Digit_Decimal_Prod1(t *testing.T) {
	src := `0 10.33340 0.000000001 3939.333393 9999.11100 10_000_000 10_0.30_0 9.30_0E30_034 0e20 0.e25 1.e+50 0.3e+50 0e00001 0.E-50`
	expected := []Token{
		{Type: TNumericLiteral, Lexeme: "0", Literal: "0", Line: 0, Column: 0},
		{Type: TNumericLiteral, Lexeme: "10.33340", Literal: "10.33340", Line: 0, Column: 0},
		{Type: TNumericLiteral, Lexeme: "0.000000001", Literal: "0.000000001", Line: 0, Column: 0},
		{Type: TNumericLiteral, Lexeme: "3939.333393", Literal: "3939.333393", Line: 0, Column: 0},
		{Type: TNumericLiteral, Lexeme: "9999.11100", Literal: "9999.11100", Line: 0, Column: 0},
		{Type: TNumericLiteral, Lexeme: "10_000_000", Literal: "10_000_000", Line: 0, Column: 0},
		{Type: TNumericLiteral, Lexeme: "10_0.30_0", Literal: "10_0.30_0", Line: 0, Column: 0},
		{Type: TNumericLiteral, Lexeme: "9.30_0E30_034", Literal: "9.30_0E30_034", Line: 0, Column: 0},
		{Type: TNumericLiteral, Lexeme: "0e20", Literal: "0e20", Line: 0, Column: 0},
		{Type: TNumericLiteral, Lexeme: "0.e25", Literal: "0.e25", Line: 0, Column: 0},
		{Type: TNumericLiteral, Lexeme: "1.e+50", Literal: "1.e+50", Line: 0, Column: 0},
		{Type: TNumericLiteral, Lexeme: "0.3e+50", Literal: "0.3e+50", Line: 0, Column: 0},
		{Type: TNumericLiteral, Lexeme: "0e00001", Literal: "0e00001", Line: 0, Column: 0},
		{Type: TNumericLiteral, Lexeme: "0.E-50", Literal: "0.E-50", Line: 0, Column: 0},
	}

	logger := gojs.NewSimpleLogger(gojs.ModeDebug)
	lexer := NewLexer(src, logger)
	got, err := lexer.ScanAll()
	if err != nil {
		logger.DumpLogs()
		t.Fatalf("unexpected error: %v", err)
	}
	assertLexemes(t, logger, got, expected)

	// FYI: These are not valid in strict mode, though we need to be overly permissive
	src = `0000008989 01234567 0777`
	expected = []Token{
		{Type: TNumericLiteral, Lexeme: "0000008989", Literal: "0000008989", Line: 0, Column: 0},
		{Type: TNumericLiteral, Lexeme: "01234567", Literal: "01234567", Line: 0, Column: 0},
		{Type: TNumericLiteral, Lexeme: "0777", Literal: "0777", Line: 0, Column: 0},
	}

}

func TestLiteral_Digit_Decimal_Prod1_Err(t *testing.T) {
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

	logger := gojs.NewSimpleLogger(gojs.ModeDebug)
	for _, s := range src {
		lexer := NewLexer(s, logger)
		got, errs := lexer.ScanAll()
		assertErrors(t, logger, errNoLiteralAfterNumber, errs, s, got)
	}

}

func TestLiteral_Digit_Decimal_Prod2_Err(t *testing.T) {
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
		"._3",
		".+3",
		".-3",
		"10e++3",
		"10e--3", // TODO: That's a bad case -- 10-- is legal and should leave the digit parser
	}

	logger := gojs.NewSimpleLogger(gojs.ModeDebug)
	for _, s := range src {
		lexer := NewLexer(s, logger)
		got, egot := lexer.ScanAll()
		eexp := errUnexpectedToken
		assertErrors(t, logger, eexp, egot, s, got)
	}

}

func TestLiteral_Digit_Decimal_Prod2(t *testing.T) {
	src := `.33340 0.0000_0000_1 .3_0E0_034 .1e+2_0 .3e-2_5 .0000e25 .1e-50 .5E+50 .9E-50`
	expected := []Token{
		{Type: TNumericLiteral, Lexeme: ".33340", Literal: ".33340", Line: 0, Column: 0},
		{Type: TNumericLiteral, Lexeme: "0.0000_0000_1", Literal: "0.0000_0000_1", Line: 0, Column: 0},
		{Type: TNumericLiteral, Lexeme: ".3_0E0_034", Literal: ".3_0E0_034", Line: 0, Column: 0},
		{Type: TNumericLiteral, Lexeme: ".1e+2_0", Literal: ".1e+2_0", Line: 0, Column: 0},
		{Type: TNumericLiteral, Lexeme: ".3e-2_5", Literal: ".3e-2_5", Line: 0, Column: 0},
		{Type: TNumericLiteral, Lexeme: ".0000e25", Literal: ".0000e25", Line: 0, Column: 0},
		{Type: TNumericLiteral, Lexeme: ".1e-50", Literal: ".1e-50", Line: 0, Column: 0},
		{Type: TNumericLiteral, Lexeme: ".5E+50", Literal: ".5E+50", Line: 0, Column: 0},
		{Type: TNumericLiteral, Lexeme: ".9E-50", Literal: ".9E-50", Line: 0, Column: 0},
	}
	logger := gojs.NewSimpleLogger(gojs.ModeDebug)
	lexer := NewLexer(src, logger)
	got, err := lexer.ScanAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertLexemes(t, logger, got, expected)
}
func TestLiteral_Digit_Decimal_Prod3(t *testing.T) {
	src := `1_35E-50_0 00000E-50_00000 000000e000000 007654321e+1 000e+1`
	expected := []Token{
		// {Type: TNumericLiteral, Lexeme: "0_E25", Literal: "0_E25", Line: 0, Column: 0},
		{Type: TNumericLiteral, Lexeme: "1_35E-50_0", Literal: "1_35E-50_0", Line: 0, Column: 0},
		{Type: TNumericLiteral, Lexeme: "00000E-50_00000", Literal: "00000E-50_00000", Line: 0, Column: 0},
		{Type: TNumericLiteral, Lexeme: "000000e000000", Literal: "000000e000000", Line: 0, Column: 0},
		{Type: TNumericLiteral, Lexeme: "007654321e+1", Literal: "007654321e+1", Line: 0, Column: 0},
		{Type: TNumericLiteral, Lexeme: "000e+1", Literal: "000e+1", Line: 0, Column: 0},
	}
	logger := gojs.NewSimpleLogger(gojs.ModeDebug)
	lexer := NewLexer(src, logger)
	got, err := lexer.ScanAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertLexemes(t, logger, got, expected)
}

func TestLiteral_Digit_Hex(t *testing.T) {
	src := `0x1 0xA 0x1234567890abcdef 0X1234567890ABCDEF 0xB_AAB_445`
	expected := []Token{
		{Type: TNumericLiteral, Lexeme: "0x1", Literal: "0x1", Line: 0, Column: 0},
		{Type: TNumericLiteral, Lexeme: "0xA", Literal: "0xA", Line: 0, Column: 0},
		{Type: TNumericLiteral, Lexeme: "0x1234567890abcdef", Literal: "0x1234567890abcdef", Line: 0, Column: 0},
		{Type: TNumericLiteral, Lexeme: "0X1234567890ABCDEF", Literal: "0X1234567890ABCDEF", Line: 0, Column: 0},
		{Type: TNumericLiteral, Lexeme: "0xB_AAB_445", Literal: "0xB_AAB_445", Line: 0, Column: 0},
	}

	logger := gojs.NewSimpleLogger(gojs.ModeDebug)
	lexer := NewLexer(src, logger)
	got, _ := lexer.ScanAll()
	assertLexemes(t, logger, got, expected)

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

func TestLiteral_Digit_BigInt(t *testing.T) {
	src := `0n 8n 84981283n 1_923_921_839_1273n`
	expected := []Token{
		{Type: TNumericLiteral, Lexeme: "0n", Literal: "0n", Line: 0, Column: 0},
		{Type: TNumericLiteral, Lexeme: "8n", Literal: "8n", Line: 0, Column: 0},
		{Type: TNumericLiteral, Lexeme: "84981283n", Literal: "84981283n", Line: 0, Column: 0},
		{Type: TNumericLiteral, Lexeme: "1_923_921_839_1273n", Literal: "1_923_921_839_1273n", Line: 0, Column: 0},
	}

	logger := gojs.NewSimpleLogger(gojs.ModeDebug)
	lexer := NewLexer(src, logger)
	got, err := lexer.ScanAll()
	if err != nil {
		logger.DumpLogs()
		t.Fatalf("unexpected error: %v", err)
	}
	assertLexemes(t, logger, got, expected)
}

func TestLiteral_Digit_Binary(t *testing.T) {
	src := `0b0 0b1 0B0 0B1 0B0101010 0b101010 0b1010_0101_0110 0b0100_0101_0110`
	expected := []Token{
		{Type: TNumericLiteral, Lexeme: "0b0", Literal: "0b0", Line: 0, Column: 0},
		{Type: TNumericLiteral, Lexeme: "0b1", Literal: "0b1", Line: 0, Column: 0},
		{Type: TNumericLiteral, Lexeme: "0B0", Literal: "0B0", Line: 0, Column: 0},
		{Type: TNumericLiteral, Lexeme: "0B1", Literal: "0B1", Line: 0, Column: 0},
		{Type: TNumericLiteral, Lexeme: "0B0101010", Literal: "0B0101010", Line: 0, Column: 0},
		{Type: TNumericLiteral, Lexeme: "0b101010", Literal: "0b101010", Line: 0, Column: 0},
		{Type: TNumericLiteral, Lexeme: "0b1010_0101_0110", Literal: "0b1010_0101_0110", Line: 0, Column: 0},
		{Type: TNumericLiteral, Lexeme: "0b0100_0101_0110", Literal: "0b0100_0101_0110", Line: 0, Column: 0},
	}

	logger := gojs.NewSimpleLogger(gojs.ModeDebug)
	lexer := NewLexer(src, logger)
	got, err := lexer.ScanAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertLexemes(t, logger, got, expected)
}

func TestLiteral_Digit_Octal(t *testing.T) {
	src := `0o0 0o1 0O0 0O7 0O6 0O2112_2234_6670 0o1234_5672_5012`
	expected := []Token{
		{Type: TNumericLiteral, Lexeme: "0o0", Literal: "0o0", Line: 0, Column: 0},
		{Type: TNumericLiteral, Lexeme: "0o1", Literal: "0o1", Line: 0, Column: 0},
		{Type: TNumericLiteral, Lexeme: "0O0", Literal: "0O0", Line: 0, Column: 0},
		{Type: TNumericLiteral, Lexeme: "0O7", Literal: "0O7", Line: 0, Column: 0},
		{Type: TNumericLiteral, Lexeme: "0O6", Literal: "0O6", Line: 0, Column: 0},
		{Type: TNumericLiteral, Lexeme: "0O2112_2234_6670", Literal: "0O2112_2234_6670", Line: 0, Column: 0},
		{Type: TNumericLiteral, Lexeme: "0o1234_5672_5012", Literal: "0o1234_5672_5012", Line: 0, Column: 0},
	}

	logger := gojs.NewSimpleLogger(gojs.ModeDebug)
	lexer := NewLexer(src, logger)
	got, err := lexer.ScanAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertLexemes(t, logger, got, expected)
}

// IdentifierName ::
// | IdentifierStart
// | IdentifierName IdentifierPart
func TestIdentifier(t *testing.T) {
	src := `abcdefghjijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTWUVWXYZ _012bx $02213 $$$$$ $\u0000\u0001\u0003 _____ \u3417\u93f0x$$\u0122a_`
	expected := []Token{
		{Type: TIdentifier, Lexeme: `abcdefghjijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTWUVWXYZ`, Literal: `abcdefghjijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTWUVWXYZ`, Line: 0, Column: 0},
		{Type: TIdentifier, Lexeme: `_012bx`, Literal: `_012bx`, Line: 0, Column: 0},
		{Type: TIdentifier, Lexeme: `$02213`, Literal: `$02213`, Line: 0, Column: 0},
		{Type: TIdentifier, Lexeme: `$$$$$`, Literal: `$$$$$`, Line: 0, Column: 0},
		{Type: TIdentifier, Lexeme: `$\u0000\u0001\u0003`, Literal: `$\u0000\u0001\u0003`, Line: 0, Column: 0},
		{Type: TIdentifier, Lexeme: `_____`, Literal: `_____`, Line: 0, Column: 0},
		{Type: TIdentifier, Lexeme: `\u3417\u93f0x$$\u0122a_`, Literal: `\u3417\u93f0x$$\u0122a_`, Line: 0, Column: 0},
	}

	logger := gojs.NewSimpleLogger(gojs.ModeDebug)
	lexer := NewLexer(src, logger)
	got, _ := lexer.ScanAll()
	assertTokens(t, logger, got, expected)
}

func TestIdentifier_Keywords(t *testing.T) {
	src := `await async break case catch class const continue debugger default delete do else enum export extends false finally for function if import in let new null return super switch this throw true try typeof undefined var void while with yield`
	expected := []Token{
		{Type: TAwait, Lexeme: `await`, Line: 0, Column: 0},
		{Type: TAsync, Lexeme: `async`, Line: 0, Column: 0},
		{Type: TBreak, Lexeme: `break`, Line: 0, Column: 0},
		{Type: TCase, Lexeme: `case`, Line: 0, Column: 0},
		{Type: TCatch, Lexeme: `catch`, Line: 0, Column: 0},
		{Type: TClass, Lexeme: `class`, Line: 0, Column: 0},
		{Type: TConst, Lexeme: `const`, Line: 0, Column: 0},
		{Type: TContinue, Lexeme: `continue`, Line: 0, Column: 0},
		{Type: TDebugger, Lexeme: `debugger`, Line: 0, Column: 0},
		{Type: TDefault, Lexeme: `default`, Line: 0, Column: 0},
		{Type: TDelete, Lexeme: `delete`, Line: 0, Column: 0},
		{Type: TDo, Lexeme: `do`, Line: 0, Column: 0},
		{Type: TElse, Lexeme: `else`, Line: 0, Column: 0},
		{Type: TEnum, Lexeme: `enum`, Line: 0, Column: 0},
		{Type: TExport, Lexeme: `export`, Line: 0, Column: 0},
		{Type: TExtends, Lexeme: `extends`, Line: 0, Column: 0},
		{Type: TFalse, Lexeme: `false`, Line: 0, Column: 0},
		{Type: TFinally, Lexeme: `finally`, Line: 0, Column: 0},
		{Type: TFor, Lexeme: `for`, Line: 0, Column: 0},
		{Type: TFunction, Lexeme: `function`, Line: 0, Column: 0},
		{Type: TIf, Lexeme: `if`, Line: 0, Column: 0},
		// {Type: TImplements, Lexeme: `implements` 0, Column: 0},
		{Type: TImport, Lexeme: `import`, Line: 0, Column: 0},
		{Type: TIn, Lexeme: `in`, Line: 0, Column: 0},
		// {Type: TInstanceOf, Lexeme: `instanceof` 0, Column: 0},
		// {Type: TInterface, Lexeme: `interface` 0, Column: 0},
		{Type: TLet, Lexeme: `let`, Line: 0, Column: 0},
		{Type: TNew, Lexeme: `new`, Line: 0, Column: 0},
		{Type: TNull, Lexeme: `null`, Line: 0, Column: 0},
		// {Type: TPackage, Lexeme: `package` 0, Column: 0},
		// {Type: TPrivate, Lexeme: `private` 0, Column: 0},
		// {Type: TProtected, Lexeme: `protected` 0, Column: 0},
		// {Type: TPublic, Lexeme: `public` 0, Column: 0},
		{Type: TReturn, Lexeme: `return`, Line: 0, Column: 0},
		// {Type: TStatic, Lexeme: `static` 0, Column: 0},
		{Type: TSuper, Lexeme: `super`, Line: 0, Column: 0},
		{Type: TSwitch, Lexeme: `switch`, Line: 0, Column: 0},
		{Type: TThis, Lexeme: `this`, Line: 0, Column: 0},
		{Type: TThrow, Lexeme: `throw`, Line: 0, Column: 0},
		{Type: TTrue, Lexeme: `true`, Line: 0, Column: 0},
		{Type: TTry, Lexeme: `try`, Line: 0, Column: 0},
		{Type: TTypeof, Lexeme: `typeof`, Line: 0, Column: 0},
		{Type: TUndefined, Lexeme: `undefined`, Line: 0, Column: 0},
		{Type: TVar, Lexeme: `var`, Line: 0, Column: 0},
		{Type: TVoid, Lexeme: `void`, Line: 0, Column: 0},
		{Type: TWhile, Lexeme: `while`, Line: 0, Column: 0},
		{Type: TWith, Lexeme: `with`, Line: 0, Column: 0},
		{Type: TYield, Lexeme: `yield`, Line: 0, Column: 0},
	}

	logger := gojs.NewSimpleLogger(gojs.ModeDebug)
	lexer := NewLexer(src, logger)
	got, err := lexer.ScanAll()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		t.Fail()
	}
	assertTokens(t, logger, got, expected)
}

// Template Literals
//
// Template ::
// | NoSubstitutionTemplate
// | TemplateHead
func TestLiteral_Template(t *testing.T) {
	t.Skip("NotImplemented")
	src := "`$end$` `$$$` `        ` `scan\n\nthis\n\ntoo!` `bla\n\n\nbla`"
	expected := []Token{
		{Type: TTemplateLiteral, Lexeme: "`$end$`", Literal: "$end$", Line: 0, Column: 0},
		{Type: TTemplateLiteral, Lexeme: "`$$$`", Literal: "$$$", Line: 0, Column: 0},
		{Type: TTemplateLiteral, Lexeme: "`        `", Literal: "        ", Line: 0, Column: 0},
		{Type: TTemplateLiteral, Lexeme: "`scan\n\nthis\n\ntoo!`", Literal: "scan\n\nthis\n\ntoo!", Line: 0, Column: 0},
		{Type: TTemplateLiteral, Lexeme: "`bla\n\n\nbla`", Literal: "bla\n\n\nbla", Line: 0, Column: 0},
	}
	logger := gojs.NewSimpleLogger(gojs.ModeDebug)
	lexer := NewLexer(src, logger)
	got, _ := lexer.ScanAll()
	assertTokens(t, logger, got, expected)
}

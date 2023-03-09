package lex

import (
	"testing"

	gojs "github.com/ruiconti/gojs/internal"
)

// Identifiers
// IdentifierName ::
// | IdentifierStart
// | IdentifierName IdentifierPart
//
// IdentifierStart ::
// | IdentifierStartChar
// | \ UnicodeEscapeSequence
//
// IdentifierStartChar ::
// | UnicodeIDStart
// | $
// | _
//
// IdentifierPart ::
// | IdentifierPartChar
// | \ UnicodeEscapeSequence
//
// IdentifierPartChar ::
// | UnicodeIDContinue
// | $
// | <ZWNJ>
// | <ZWJ>
//
// UnicodeIDStart ::
// | any Unicode code point with the Unicode property “ID_Start”
// UnicodeIDContinue ::
// | any Unicode code point with the Unicode property “ID_Continue”

func TestScan_IdentifierName(t *testing.T) {
	src := `abcdefghjijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTWUVWXYZ _012bx $02213 $$$$$ $\u0000\u0001\u0003 _____ \u3417\u93f0x$$\u0122a_`
	expected := []Token{
		{T: TIdentifier, Lexeme: `abcdefghjijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTWUVWXYZ`, Literal: `abcdefghjijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTWUVWXYZ`, Line: 0, Column: 0},
		{T: TIdentifier, Lexeme: `_012bx`, Literal: `_012bx`, Line: 0, Column: 0},
		{T: TIdentifier, Lexeme: `$02213`, Literal: `$02213`, Line: 0, Column: 0},
		{T: TIdentifier, Lexeme: `$$$$$`, Literal: `$$$$$`, Line: 0, Column: 0},
		{T: TIdentifier, Lexeme: `$\u0000\u0001\u0003`, Literal: `$\u0000\u0001\u0003`, Line: 0, Column: 0},
		{T: TIdentifier, Lexeme: `_____`, Literal: `_____`, Line: 0, Column: 0},
		{T: TIdentifier, Lexeme: `\u3417\u93f0x$$\u0122a_`, Literal: `\u3417\u93f0x$$\u0122a_`, Line: 0, Column: 0},
	}

	logger := gojs.NewSimpleLogger(gojs.ModeDebug)
	scanner := NewScanner(src, logger)
	got, _ := scanner.Scan()
	assertTokens(t, logger, got, expected)
}

func TestScanReservedKeywords(t *testing.T) {
	src := `break case catch class const continue debugger default delete do else enum export extends false finally for function if import in let new null return super switch this throw true try typeof undefined var void while with yield`
	expected := []Token{
		{T: TBreak, Lexeme: `break`, Literal: `break`, Line: 0, Column: 0},
		{T: TCase, Lexeme: `case`, Literal: `case`, Line: 0, Column: 0},
		{T: TCatch, Lexeme: `catch`, Literal: `catch`, Line: 0, Column: 0},
		{T: TClass, Lexeme: `class`, Literal: `class`, Line: 0, Column: 0},
		{T: TConst, Lexeme: `const`, Literal: `const`, Line: 0, Column: 0},
		{T: TContinue, Lexeme: `continue`, Literal: `continue`, Line: 0, Column: 0},
		{T: TDebugger, Lexeme: `debugger`, Literal: `debugger`, Line: 0, Column: 0},
		{T: TDefault, Lexeme: `default`, Literal: `default`, Line: 0, Column: 0},
		{T: TDelete, Lexeme: `delete`, Literal: `delete`, Line: 0, Column: 0},
		{T: TDo, Lexeme: `do`, Literal: `do`, Line: 0, Column: 0},
		{T: TElse, Lexeme: `else`, Literal: `else`, Line: 0, Column: 0},
		{T: TEnum, Lexeme: `enum`, Literal: `enum`, Line: 0, Column: 0},
		{T: TExport, Lexeme: `export`, Literal: `export`, Line: 0, Column: 0},
		{T: TExtends, Lexeme: `extends`, Literal: `extends`, Line: 0, Column: 0},
		{T: TFalse, Lexeme: `false`, Literal: `false`, Line: 0, Column: 0},
		{T: TFinally, Lexeme: `finally`, Literal: `finally`, Line: 0, Column: 0},
		{T: TFor, Lexeme: `for`, Literal: `for`, Line: 0, Column: 0},
		{T: TFunction, Lexeme: `function`, Literal: `function`, Line: 0, Column: 0},
		{T: TIf, Lexeme: `if`, Literal: `if`, Line: 0, Column: 0},
		// {T: TImplements, Lexeme: `implements`, Literal: ``, Line: 0, Column: 0},
		{T: TImport, Lexeme: `import`, Literal: `import`, Line: 0, Column: 0},
		{T: TIn, Lexeme: `in`, Literal: `in`, Line: 0, Column: 0},
		// {T: TInstanceOf, Lexeme: `instanceof`, Literal: ``, Line: 0, Column: 0},
		// {T: TInterface, Lexeme: `interface`, Literal: ``, Line: 0, Column: 0},
		{T: TLet, Lexeme: `let`, Literal: `let`, Line: 0, Column: 0},
		{T: TNew, Lexeme: `new`, Literal: `new`, Line: 0, Column: 0},
		{T: TNull, Lexeme: `null`, Literal: `null`, Line: 0, Column: 0},
		// {T: TPackage, Lexeme: `package`, Literal: ``, Line: 0, Column: 0},
		// {T: TPrivate, Lexeme: `private`, Literal: ``, Line: 0, Column: 0},
		// {T: TProtected, Lexeme: `protected`, Literal: ``, Line: 0, Column: 0},
		// {T: TPublic, Lexeme: `public`, Literal: ``, Line: 0, Column: 0},
		{T: TReturn, Lexeme: `return`, Literal: `return`, Line: 0, Column: 0},
		// {T: TStatic, Lexeme: `static`, Literal: ``, Line: 0, Column: 0},
		{T: TSuper, Lexeme: `super`, Literal: `super`, Line: 0, Column: 0},
		{T: TSwitch, Lexeme: `switch`, Literal: `switch`, Line: 0, Column: 0},
		{T: TThis, Lexeme: `this`, Literal: `this`, Line: 0, Column: 0},
		{T: TThrow, Lexeme: `throw`, Literal: `throw`, Line: 0, Column: 0},
		{T: TTrue, Lexeme: `true`, Literal: `true`, Line: 0, Column: 0},
		{T: TTry, Lexeme: `try`, Literal: `try`, Line: 0, Column: 0},
		{T: TTypeof, Lexeme: `typeof`, Literal: `typeof`, Line: 0, Column: 0},
		{T: TUndefined, Lexeme: `undefined`, Literal: `undefined`, Line: 0, Column: 0},
		{T: TVar, Lexeme: `var`, Literal: `var`, Line: 0, Column: 0},
		{T: TVoid, Lexeme: `void`, Literal: `void`, Line: 0, Column: 0},
		{T: TWhile, Lexeme: `while`, Literal: `while`, Line: 0, Column: 0},
		{T: TWith, Lexeme: `with`, Literal: `with`, Line: 0, Column: 0},
		{T: TYield, Lexeme: `yield`, Literal: `yield`, Line: 0, Column: 0},
	}

	logger := gojs.NewSimpleLogger(gojs.ModeDebug)
	scanner := NewScanner(src, logger)
	got, err := scanner.Scan()
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
	logger := gojs.NewSimpleLogger(gojs.ModeDebug)
	scanner := NewScanner(src, logger)
	got, _ := scanner.Scan()
	assertTokens(t, logger, got, expected)
}

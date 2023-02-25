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
	src := `abcdefghjijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTWUVWXYZ _012bx $02213 $$$$$ $\u0000\u0001\u0003 _____`
	expected := []Token{
		{T: TIdentifier, Lexeme: `abcdefghjijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTWUVWXYZ`, Literal: `abcdefghjijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTWUVWXYZ`, Line: 0, Column: 0},
		{T: TIdentifier, Lexeme: `_012bx`, Literal: `_012bx`, Line: 0, Column: 0},
		{T: TIdentifier, Lexeme: `$02213`, Literal: `$02213`, Line: 0, Column: 0},
		{T: TIdentifier, Lexeme: `$$$$$`, Literal: `$$$$$`, Line: 0, Column: 0},
		{T: TIdentifier, Lexeme: `$\u0000\u0001\u0003`, Literal: `$\u0000\u0001\u0003`, Line: 0, Column: 0},
		{T: TIdentifier, Lexeme: `_____`, Literal: `_____`, Line: 0, Column: 0},
	}

	debugLogger := gojs.NewSimpleLogger(gojs.ModeDebug)
	scanner := NewScanner(src, debugLogger)
	got, _ := scanner.Scan()
	assertTokens(t, got, expected)
}

func TestScanReservedKeywords(t *testing.T) {
	src := `break case catch class const continue debugger default delete do else enum export extends false finally for function if import in let new null return super switch this throw true try typeof undefined var void while with yield`
	expected := []Token{
		{T: TBreak, Lexeme: `break`, Literal: nil, Line: 0, Column: 0},
		{T: TCase, Lexeme: `case`, Literal: nil, Line: 0, Column: 0},
		{T: TCatch, Lexeme: `catch`, Literal: nil, Line: 0, Column: 0},
		{T: TClass, Lexeme: `class`, Literal: nil, Line: 0, Column: 0},
		{T: TConst, Lexeme: `const`, Literal: nil, Line: 0, Column: 0},
		{T: TContinue, Lexeme: `continue`, Literal: nil, Line: 0, Column: 0},
		{T: TDebugger, Lexeme: `debugger`, Literal: nil, Line: 0, Column: 0},
		{T: TDefault, Lexeme: `default`, Literal: nil, Line: 0, Column: 0},
		{T: TDelete, Lexeme: `delete`, Literal: nil, Line: 0, Column: 0},
		{T: TDo, Lexeme: `do`, Literal: nil, Line: 0, Column: 0},
		{T: TElse, Lexeme: `else`, Literal: nil, Line: 0, Column: 0},
		{T: TEnum, Lexeme: `enum`, Literal: nil, Line: 0, Column: 0},
		{T: TExport, Lexeme: `export`, Literal: nil, Line: 0, Column: 0},
		{T: TExtends, Lexeme: `extends`, Literal: nil, Line: 0, Column: 0},
		{T: TFalse, Lexeme: `false`, Literal: nil, Line: 0, Column: 0},
		{T: TFinally, Lexeme: `finally`, Literal: nil, Line: 0, Column: 0},
		{T: TFor, Lexeme: `for`, Literal: nil, Line: 0, Column: 0},
		{T: TFunction, Lexeme: `function`, Literal: nil, Line: 0, Column: 0},
		{T: TIf, Lexeme: `if`, Literal: nil, Line: 0, Column: 0},
		// {T: TImplements, Lexeme: `implements`, Literal: nil, Line: 0, Column: 0},
		{T: TImport, Lexeme: `import`, Literal: nil, Line: 0, Column: 0},
		{T: TIn, Lexeme: `in`, Literal: nil, Line: 0, Column: 0},
		// {T: TInstanceOf, Lexeme: `instanceof`, Literal: nil, Line: 0, Column: 0},
		// {T: TInterface, Lexeme: `interface`, Literal: nil, Line: 0, Column: 0},
		{T: TLet, Lexeme: `let`, Literal: nil, Line: 0, Column: 0},
		{T: TNew, Lexeme: `new`, Literal: nil, Line: 0, Column: 0},
		{T: TNull, Lexeme: `null`, Literal: nil, Line: 0, Column: 0},
		// {T: TPackage, Lexeme: `package`, Literal: nil, Line: 0, Column: 0},
		// {T: TPrivate, Lexeme: `private`, Literal: nil, Line: 0, Column: 0},
		// {T: TProtected, Lexeme: `protected`, Literal: nil, Line: 0, Column: 0},
		// {T: TPublic, Lexeme: `public`, Literal: nil, Line: 0, Column: 0},
		{T: TReturn, Lexeme: `return`, Literal: nil, Line: 0, Column: 0},
		// {T: TStatic, Lexeme: `static`, Literal: nil, Line: 0, Column: 0},
		{T: TSuper, Lexeme: `super`, Literal: nil, Line: 0, Column: 0},
		{T: TSwitch, Lexeme: `switch`, Literal: nil, Line: 0, Column: 0},
		{T: TThis, Lexeme: `this`, Literal: nil, Line: 0, Column: 0},
		{T: TThrow, Lexeme: `throw`, Literal: nil, Line: 0, Column: 0},
		{T: TTrue, Lexeme: `true`, Literal: nil, Line: 0, Column: 0},
		{T: TTry, Lexeme: `try`, Literal: nil, Line: 0, Column: 0},
		{T: TTypeof, Lexeme: `typeof`, Literal: nil, Line: 0, Column: 0},
		{T: TUndefined, Lexeme: `undefined`, Literal: nil, Line: 0, Column: 0},
		{T: TVar, Lexeme: `var`, Literal: nil, Line: 0, Column: 0},
		{T: TVoid, Lexeme: `void`, Literal: nil, Line: 0, Column: 0},
		{T: TWhile, Lexeme: `while`, Literal: nil, Line: 0, Column: 0},
		{T: TWith, Lexeme: `with`, Literal: nil, Line: 0, Column: 0},
		{T: TYield, Lexeme: `yield`, Literal: nil, Line: 0, Column: 0},
	}

	scanner := NewScanner(src, defaultLogger)
	got, err := scanner.Scan()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		t.Fail()
	}
	assertTokens(t, got, expected)

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
	scanner := NewScanner(src, defaultLogger)
	got, _ := scanner.Scan()
	assertTokens(t, got, expected)
}

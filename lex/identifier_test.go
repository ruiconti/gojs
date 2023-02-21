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
	assertLexemes(t, got, expected)
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
	assertLexemes(t, got, expected)
}

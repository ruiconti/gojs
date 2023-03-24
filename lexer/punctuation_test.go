package lexer

import (
	"testing"

	gojs "github.com/ruiconti/gojs/internal"
)

func TestPunctuation_Single(t *testing.T) {
	src := `;()~:{}[];,~^%##`

	expected := []Token{
		{Type: TSemicolon, Lexeme: ";", Literal: nil, Line: 0, Column: 0},
		{Type: TLeftParen, Lexeme: "(", Literal: nil, Line: 0, Column: 0},
		{Type: TRightParen, Lexeme: ")", Literal: nil, Line: 0, Column: 0},
		{Type: TTilde, Lexeme: "~", Literal: nil, Line: 0, Column: 0},
		{Type: TColon, Lexeme: ":", Literal: nil, Line: 0, Column: 0},
		{Type: TLeftBrace, Lexeme: "{", Literal: nil, Line: 0, Column: 0},
		{Type: TRightBrace, Lexeme: "}", Literal: nil, Line: 0, Column: 0},
		{Type: TLeftBracket, Lexeme: "[", Literal: nil, Line: 0, Column: 0},
		{Type: TRightBracket, Lexeme: "]", Literal: nil, Line: 0, Column: 0},
		{Type: TSemicolon, Lexeme: ";", Literal: nil, Line: 0, Column: 0},
		{Type: TComma, Lexeme: ",", Literal: nil, Line: 0, Column: 0},
		{Type: TTilde, Lexeme: "~", Literal: nil, Line: 0, Column: 0},
		{Type: TXor, Lexeme: "^", Literal: nil, Line: 0, Column: 0},
		{Type: TPercent, Lexeme: "%", Literal: nil, Line: 0, Column: 0},
		{Type: TNumberSign, Lexeme: "#", Literal: nil, Line: 0, Column: 0},
		{Type: TNumberSign, Lexeme: "#", Literal: nil, Line: 0, Column: 0},
	}

	logger := gojs.NewSimpleLogger(gojs.ModeDebug)
	lexer := NewLexer(src, logger)
	got, err := lexer.ScanAll()
	if err != nil {
		lexer.logger.DumpLogs()
		t.Fatalf("unexpected error: %v", err)
	}
	assertLexemes(t, logger, got, expected)
}

func TestPunctuation_GreaterLessShift(t *testing.T) {
	src := `> >= >> >>= >>> >>>= < << <<= <=`
	expected := []Token{
		{Type: TGreaterThan, Lexeme: ">", Literal: nil, Line: 0, Column: 0},
		{Type: TGreaterThanEqual, Lexeme: ">=", Literal: nil, Line: 0, Column: 0},
		{Type: TRightShift, Lexeme: ">>", Literal: nil, Line: 0, Column: 0},
		{Type: TRightShiftAssign, Lexeme: ">>=", Literal: nil, Line: 0, Column: 0},
		{Type: TUnsignedRightShift, Lexeme: ">>>", Literal: nil, Line: 0, Column: 0},
		{Type: TUnsignedRightShiftAssign, Lexeme: ">>>=", Literal: nil, Line: 0, Column: 0},
		{Type: TLessThan, Lexeme: "<", Literal: nil, Line: 0, Column: 0},
		{Type: TLeftShift, Lexeme: "<<", Literal: nil, Line: 0, Column: 0},
		{Type: TLeftShiftAssign, Lexeme: "<<=", Literal: nil, Line: 0, Column: 0},
		{Type: TLessThanEqual, Lexeme: "<=", Literal: nil, Line: 0, Column: 0},
	}

	logger := gojs.NewSimpleLogger(gojs.ModeDebug)
	lexer := NewLexer(src, logger)
	got, err := lexer.ScanAll()
	if err != nil {
		lexer.logger.DumpLogs()
		t.Fatalf("unexpected error: %v", err)
	}
	assertLexemes(t, logger, got, expected)
}

func TestPunctuation_BangEq(t *testing.T) {
	src := `. ... ? ?? ! != !== = == === =>`
	expected := []Token{
		{Type: TPeriod, Lexeme: ".", Literal: nil, Line: 0, Column: 0},
		{Type: TEllipsis, Lexeme: "...", Literal: nil, Line: 0, Column: 0},
		{Type: TQuestionMark, Lexeme: "?", Literal: nil, Line: 0, Column: 0},
		{Type: TDoubleQuestionMark, Lexeme: "??", Literal: nil, Line: 0, Column: 0},
		{Type: TBang, Lexeme: "!", Literal: nil, Line: 0, Column: 0},
		{Type: TNotEqual, Lexeme: "!=", Literal: nil, Line: 0, Column: 0},
		{Type: TStrictNotEqual, Lexeme: "!==", Literal: nil, Line: 0, Column: 0},
		{Type: TAssign, Lexeme: "=", Literal: nil, Line: 0, Column: 0},
		{Type: TEqual, Lexeme: "==", Literal: nil, Line: 0, Column: 0},
		{Type: TStrictEqual, Lexeme: "===", Literal: nil, Line: 0, Column: 0},
		{Type: TArrow, Lexeme: "=>", Literal: nil, Line: 0, Column: 0},
	}

	logger := gojs.NewSimpleLogger(gojs.ModeDebug)
	lexer := NewLexer(src, logger)
	got, err := lexer.ScanAll()
	if err != nil {
		lexer.logger.DumpLogs()
		t.Fatalf("unexpected error: %v", err)
	}
	assertLexemes(t, logger, got, expected)
}

func TestPunctuation_AndOrBit(t *testing.T) {
	src := `& && &= &&= | || |= ||=`
	expected := []Token{
		{Type: TAnd, Lexeme: "&", Literal: nil, Line: 0, Column: 0},
		{Type: TLogicalAnd, Lexeme: "&&", Literal: nil, Line: 0, Column: 0},
		{Type: TAndAssign, Lexeme: "&=", Literal: nil, Line: 0, Column: 0},
		{Type: TLogicalAndAssign, Lexeme: "&&=", Literal: nil, Line: 0, Column: 0},
		{Type: TOr, Lexeme: "|", Literal: nil, Line: 0, Column: 0},
		{Type: TLogicalOr, Lexeme: "||", Literal: nil, Line: 0, Column: 0},
		{Type: TOrAssign, Lexeme: "|=", Literal: nil, Line: 0, Column: 0},
		{Type: TLogicalOrAssign, Lexeme: "||=", Literal: nil, Line: 0, Column: 0},
	}

	logger := gojs.NewSimpleLogger(gojs.ModeDebug)
	lexer := NewLexer(src, logger)
	got, _ := lexer.ScanAll()
	assertLexemes(t, logger, got, expected)
}

func TestPunctuation_PlusMinusAssign(t *testing.T) {
	src := `+ ++ += - -- -=`
	expected := []Token{
		{Type: TPlus, Lexeme: "+", Literal: nil, Line: 0, Column: 0},
		{Type: TPlusPlus, Lexeme: "++", Literal: nil, Line: 0, Column: 0},
		{Type: TPlusAssign, Lexeme: "+=", Literal: nil, Line: 0, Column: 0},
		{Type: TMinus, Lexeme: "-", Literal: nil, Line: 0, Column: 0},
		{Type: TMinusMinus, Lexeme: "--", Literal: nil, Line: 0, Column: 0},
		{Type: TMinusAssign, Lexeme: "-=", Literal: nil, Line: 0, Column: 0},
	}

	logger := gojs.NewSimpleLogger(gojs.ModeDebug)
	lexer := NewLexer(src, logger)
	got, _ := lexer.ScanAll()
	assertLexemes(t, logger, got, expected)
}
func TestPunctuation_StarSlashAssign(t *testing.T) {
	src := `* *= / /= **`
	expected := []Token{
		{Type: TStar, Lexeme: "*", Literal: nil, Line: 0, Column: 0},
		{Type: TStarAssign, Lexeme: "*=", Literal: nil, Line: 0, Column: 0},
		{Type: TSlash, Lexeme: "/", Literal: nil, Line: 0, Column: 0},
		{Type: TSlashAssign, Lexeme: "/=", Literal: nil, Line: 0, Column: 0},
		{Type: TStarStar, Lexeme: "**", Literal: nil, Line: 0, Column: 0},
	}

	logger := gojs.NewSimpleLogger(gojs.ModeDebug)
	lexer := NewLexer(src, logger)
	got, _ := lexer.ScanAll()
	assertLexemes(t, logger, got, expected)
}

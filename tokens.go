package main

// Map all tokens from the specification
// https://262.ecma-international.org/#sec-tokens
// Tokens

type Token int

const (
	TIdentifier Token = iota

	// Literals
	TNumericLiteral
	TStringLiteral
	TRegularExpressionLiteral

	// Reserved words
	TAwait
	TBreak
	TCase
	TCatch
	TClass
	TConst
	TContinue
	TDebugger
	TDefault
	TDelete
	TDo
	TElse
	TEnum
	TExport
	TExtends
	TFalse
	TFinally
	TFor
	TFunction
	TIf
	TImport
	TIn
	TInstanceof
	TLet
	TNew
	TNull
	TReturn
	TSuper
	TSwitch
	TThis
	TThrow
	TTrue
	TTry
	TTypeof
	TVar
	TVoid
	TWhile
	TWith
	TYield

	// Punctuators
	TLeftBrace
	TRightBrace
	TLeftParen
	TRightParen
	TLeftBracket
	TRightBracket
	TPeriod
	TEllipsis
	TSemicolon
	TComma
	TLessThan
	TGreaterThan
	TLessThanEqual
	TGreaterThanEqual
	TEqual
	TNotEqual
	TStrictEqual
	TStrictNotEqual
	TPlus
	TMinus
	TStar
	TPercent
	TPlusPlus
	TMinusMinus
	TLeftShift
	TRightShift
	TUnsignedRightShift
	TAnd
	TOr
	TXor
	TBang
	TBitwiseNot
	TLogicalAnd
	TLogicalOr
	TQuestionMark
	TDoubleQuestionMark
	TColon
	TAssign
	TPlusAssign
	TMinusAssign
	TStarAssign
	TPercentAssign
	TLeftShiftAssign
	TRightShiftAssign
	TUnsignedRightShiftAssign
	TAndAssign
	TOrAssign
	TXorAssign
	TLogicalAndAssign
	TLogicalOrAssign
	TArrow
	TSlash
	TSlashAssign
)

var TokenMap = map[Token]string{
	TTrue:                     "true",
	TFalse:                    "false",
	TLeftBrace:                "{",
	TRightBrace:               "}",
	TLeftParen:                "(",
	TRightParen:               ")",
	TLeftBracket:              "[",
	TRightBracket:             "]",
	TPeriod:                   ".",
	TEllipsis:                 "...",
	TSemicolon:                ";",
	TComma:                    ",",
	TLessThan:                 "<",
	TGreaterThan:              ">",
	TLessThanEqual:            "<=",
	TGreaterThanEqual:         ">=",
	TEqual:                    "==",
	TNotEqual:                 "!=",
	TStrictEqual:              "===",
	TStrictNotEqual:           "!==",
	TPlus:                     "+",
	TMinus:                    "-",
	TStar:                     "*",
	TPercent:                  "%",
	TPlusPlus:                 "++",
	TMinusMinus:               "--",
	TLeftShift:                "<<",
	TRightShift:               ">>",
	TUnsignedRightShift:       ">>>",
	TAnd:                      "&",
	TOr:                       "|",
	TXor:                      "^",
	TBang:                     "!",
	TBitwiseNot:               "~",
	TLogicalAnd:               "&&",
	TLogicalOr:                "||",
	TQuestionMark:             "?",
	TDoubleQuestionMark:       "??",
	TColon:                    ":",
	TAssign:                   "=",
	TPlusAssign:               "+=",
	TMinusAssign:              "-=",
	TStarAssign:               "*=",
	TPercentAssign:            "%=",
	TLeftShiftAssign:          "<<=",
	TRightShiftAssign:         ">>=",
	TUnsignedRightShiftAssign: ">>>=",
	TAndAssign:                "&=",
	TOrAssign:                 "|=",
	TXorAssign:                "^=",
	TLogicalAndAssign:         "&&=",
	TLogicalOrAssign:          "||=",
	TArrow:                    "=>",
	TSlash:                    "/",
	TSlashAssign:              "/=",
	TAwait:                    "await",
	TBreak:                    "break",
	TCase:                     "case",
	TCatch:                    "catch",
	TClass:                    "class",
	TConst:                    "const",
	TContinue:                 "continue",
	TDebugger:                 "debugger",
	TDefault:                  "default",
	TDelete:                   "delete",
	TDo:                       "do",
	TElse:                     "else",
	TEnum:                     "enum",
	TExport:                   "export",
	TExtends:                  "extends",
	TFinally:                  "finally",
	TFor:                      "for",
	TFunction:                 "function",
	TIf:                       "if",
	TImport:                   "import",
	TIn:                       "in",
	TInstanceof:               "instanceof",
	TLet:                      "let",
	TNew:                      "new",
	TNull:                     "null",
	TReturn:                   "return",
	TSuper:                    "super",
	TSwitch:                   "switch",
	TThis:                     "this",
	TThrow:                    "throw",
	TTry:                      "try",
	TTypeof:                   "typeof",
	TVar:                      "var",
	TVoid:                     "void",
	TWhile:                    "while",
	TWith:                     "with",
	TYield:                    "yield",
}

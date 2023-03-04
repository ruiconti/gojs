package lex

type Token struct {
	T       TokenType
	Lexeme  string
	Literal interface{}
	Line    int
	Column  int
}

// Map all tokens from the specification
// https://262.ecma-international.org/#sec-tokens
// Tokens
type TokenType int

const (
	TIdentifier TokenType = iota

	// Literals
	TNumericLiteral
	TStringLiteral_SingleQuote
	TStringLiteral_DoubleQuote
	TRegularExpressionLiteral
	TTemplateLiteral
	TUnknown

	// Other
	TAnd
	TAndAssign
	TArrow
	TAssign
	TAwait
	TBang
	TBitwiseNot
	TBreak
	TCase
	TCatch
	TClass
	TColon
	TComma
	TConst
	TContinue
	TDebugger
	TDefault
	TDelete
	TDo
	TDoubleQuestionMark
	TEllipsis
	TElse
	TEnum
	TEOF
	TEqual
	TExport
	TExtends
	TFalse
	TFinally
	TFor
	TFunction
	TGreaterThan
	TGreaterThanEqual
	TIf
	TImport
	TIn
	TInstanceof
	TLeftBrace
	TLeftBracket
	TLeftParen
	TLeftShift
	TLeftShiftAssign
	TLessThan
	TLessThanEqual
	TLet
	TLogicalAnd
	TLogicalAndAssign
	TLogicalOr
	TLogicalOrAssign
	TMinus
	TMinusAssign
	TMinusMinus
	TNew
	TNotEqual
	TNull
	TOr
	TOrAssign
	TPercent
	TPercentAssign
	TPeriod
	TPlus
	TPlusAssign
	TPlusPlus
	TQuestionMark
	TReturn
	TRightBrace
	TRightBracket
	TRightParen
	TRightShift
	TRightShiftAssign
	TSemicolon
	TSlash
	TSlashAssign
	TStar
	TStarAssign
	TStrictEqual
	TStrictNotEqual
	TSuper
	TSwitch
	TThis
	TThrow
	TTilde
	TTrue
	TTry
	TTypeof
	TUndefined
	TUnsignedRightShift
	TUnsignedRightShiftAssign
	TVar
	TVoid
	TWhile
	TWith
	TXor
	TXorAssign
	TYield
)

var LiteralNames = map[TokenType]string{
	TIdentifier:                "Identifier",
	TNumericLiteral:            "NumericLiteral",
	TStringLiteral_SingleQuote: "StringLiteral_SimpleQuote",
	TStringLiteral_DoubleQuote: "StringLiteral_DoubleQuote",
	TRegularExpressionLiteral:  "RegularExpressionLiteral",
	TTemplateLiteral:           "TemplateLiteral",
}

var ReservedWordNames = map[TokenType]string{
	TTrue:       "true",
	TFalse:      "false",
	TAwait:      "await",
	TBreak:      "break",
	TCase:       "case",
	TCatch:      "catch",
	TClass:      "class",
	TConst:      "const",
	TContinue:   "continue",
	TDebugger:   "debugger",
	TDefault:    "default",
	TDelete:     "delete",
	TDo:         "do",
	TElse:       "else",
	TEnum:       "enum",
	TExport:     "export",
	TExtends:    "extends",
	TFinally:    "finally",
	TFor:        "for",
	TFunction:   "function",
	TIf:         "if",
	TImport:     "import",
	TIn:         "in",
	TInstanceof: "instanceof",
	TLet:        "let",
	TNew:        "new",
	TNull:       "null",
	TReturn:     "return",
	TSuper:      "super",
	TSwitch:     "switch",
	TThis:       "this",
	TThrow:      "throw",
	TTry:        "try",
	TTypeof:     "typeof",
	TUndefined:  "undefined",
	TVar:        "var",
	TVoid:       "void",
	TWhile:      "while",
	TWith:       "with",
	TYield:      "yield",
	TEOF:        "TEOF",
}

var PunctuationNames = map[TokenType]string{
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
	TTilde:                    "~",
}

package parser

import (
	"fmt"
	"testing"

	"github.com/ruiconti/gojs/internal"
	"github.com/ruiconti/gojs/lex"
)

func binExpr(left, right AstNode, op lex.TokenType) *ExprBinaryOp {
	return &ExprBinaryOp{
		left:     left,
		right:    right,
		operator: op,
	}
}

func idExpr(name string) *ExprIdentifierReference {
	return &ExprIdentifierReference{
		reference: name,
	}
}

func TestCondExpr_MemberExpr(t *testing.T) {
	t.Run("properly parses binary expr", func(t *testing.T) {
		logger := internal.NewSimpleLogger(internal.ModeDebug)
		binOperators := []lex.TokenType{lex.TLogicalOr,
			lex.TLogicalAnd,
			lex.TOr,
			lex.TAnd,
			lex.TXor,
			lex.TEqual,
			lex.TStrictEqual,
			lex.TNotEqual,
			lex.TStrictNotEqual,
			lex.TGreaterThan,
			lex.TGreaterThanEqual,
			lex.TLessThan,
			lex.TLessThanEqual,
			lex.TLeftShift,
			lex.TRightShift,
			lex.TPlus,
			lex.TMinus,
			lex.TStar,
			lex.TPercent,
			lex.TSlash,
		}

		for _, binOperator := range binOperators {
			lexeme := lex.ResolveName(binOperator)
			src := fmt.Sprintf(`a %s b %s c %s d %s e %s f`, lexeme, lexeme, lexeme, lexeme, lexeme)
			binExpr := func(left, right AstNode) *ExprBinaryOp {
				return binExpr(left, right, binOperator)
			}

			expected := &ExprRootNode{
				children: []AstNode{
					binExpr(
						binExpr(
							binExpr(
								binExpr(
									binExpr(
										idExpr("a"),
										idExpr("b"),
									),
									idExpr("c"),
								),
								idExpr("d"),
							),
							idExpr("e"),
						),
						idExpr("f"),
					),
				},
			}
			got := Parse(logger, src)
			AssertExprEqual(t, logger, got, expected)

		}
	})

	t.Run("logical AND has precedence over OR expression", func(t *testing.T) {
		opsLogicalOr := []lex.TokenType{lex.TLogicalOr}
		opsLogicalAnd := []lex.TokenType{lex.TLogicalAnd}

		assertExprPrecedence(
			t,
			opsLogicalAnd, /* higher */
			opsLogicalOr,  /* lower */
		)
	})

	t.Run("bitwise OR has precedence over logical AND", func(t *testing.T) {
		opsOr := []lex.TokenType{lex.TOr}
		opsLogicalAnd := []lex.TokenType{lex.TLogicalAnd}

		assertExprPrecedence(
			t,
			opsOr,         /* higher */
			opsLogicalAnd, /* lower */
		)
	})

	t.Run("bitwise XOR has precedence over bitwise OR", func(t *testing.T) {
		opsXor := []lex.TokenType{lex.TXor}
		opsOr := []lex.TokenType{lex.TOr}

		assertExprPrecedence(
			t,
			opsXor, /* higher */
			opsOr,  /* lower */
		)
	})

	t.Run("bitwise AND has precedence over bitwise XOR", func(t *testing.T) {
		opsAnd := []lex.TokenType{lex.TAnd}
		opsXor := []lex.TokenType{lex.TXor}

		assertExprPrecedence(
			t,
			opsAnd, /* higher */
			opsXor, /* lower */
		)
	})

	t.Run("equality comparison has precedence over bitwise AND", func(t *testing.T) {
		opsEq := []lex.TokenType{lex.TEqual, lex.TNotEqual, lex.TStrictEqual, lex.TStrictNotEqual}
		opsBitwise := []lex.TokenType{lex.TAnd, lex.TOr, lex.TXor}

		assertExprPrecedence(
			t,
			opsEq,      /* higher */
			opsBitwise, /* lower */
		)
	})
	t.Run("relational comparison has precedence over equality comparison", func(t *testing.T) {
		opsEq := []lex.TokenType{lex.TEqual, lex.TNotEqual, lex.TStrictEqual, lex.TStrictNotEqual}
		opsRelational := []lex.TokenType{lex.TLessThan, lex.TLessThanEqual, lex.TGreaterThan, lex.TGreaterThanEqual, lex.TIn, lex.TInstanceof}

		assertExprPrecedence(
			t,
			opsRelational, /* higher */
			opsEq,         /* lower */
		)
	})
	t.Run("shift operation has precedence over relational comparison", func(t *testing.T) {
		opsRelational := []lex.TokenType{lex.TLessThan, lex.TLessThanEqual, lex.TGreaterThan, lex.TGreaterThanEqual, lex.TIn, lex.TInstanceof}
		opsShift := []lex.TokenType{lex.TLeftShift, lex.TRightShift, lex.TUnsignedRightShift}

		assertExprPrecedence(
			t,
			opsShift,      /* higher */
			opsRelational, /* lower */
		)
	})

	t.Run("additive operation has precedence over shift operation", func(t *testing.T) {
		opsShift := []lex.TokenType{lex.TLeftShift, lex.TRightShift, lex.TUnsignedRightShift}
		opsAdd := []lex.TokenType{lex.TPlus, lex.TMinus}

		assertExprPrecedence(
			t,
			opsAdd,   /* higher */
			opsShift, /* lower */
		)
	})

	t.Run("multiplicative operation has precedence over additive operation", func(t *testing.T) {
		opsMult := []lex.TokenType{lex.TStar, lex.TSlash, lex.TPercent}
		opsAdd := []lex.TokenType{lex.TPlus, lex.TMinus}

		assertExprPrecedence(
			t,
			opsMult, /* higher */
			opsAdd,  /* lower */
		)
	})

	t.Run("exponential operation has precedence over multiplicative operation", func(t *testing.T) {
		opsMult := []lex.TokenType{lex.TStar, lex.TSlash, lex.TPercent}
		opsExp := []lex.TokenType{lex.TStarStar}

		assertExprPrecedence(
			t,
			opsExp,  /* higher */
			opsMult, /* lower */
		)
	})
}

// assertExprPrecedence asserts that the given operators have the correct precedence.
func assertExprPrecedence(
	t *testing.T,
	opsHigherPrecedence []lex.TokenType,
	opsLowerPrecedence []lex.TokenType,
) {
	logger := internal.NewSimpleLogger(internal.ModeDebug)

	for _, opHigher := range opsHigherPrecedence {
		lexemeHigherPrecedence := lex.ResolveName(opHigher)
		binHigherExpr := func(left, right AstNode) *ExprBinaryOp {
			return binExpr(left, right, opHigher)
		}
		for _, opLower := range opsLowerPrecedence {
			lexemeLowerPrecedence := lex.ResolveName(opLower)
			binLowerExpr := func(left, right AstNode) *ExprBinaryOp {
				return binExpr(left, right, opLower)
			}

			// act
			src := fmt.Sprintf(
				`a %s b %s c %s d %s e %s f`,
				lexemeLowerPrecedence,
				lexemeHigherPrecedence,
				lexemeLowerPrecedence,
				lexemeHigherPrecedence,
				lexemeLowerPrecedence,
			)
			// for example, if opLower is TPlus and opHigher is TRightShift
			// equals to: ((a + (b / c)) + (d / e)) + f
			expected := &ExprRootNode{
				children: []AstNode{
					binLowerExpr(
						binLowerExpr(
							binLowerExpr(
								idExpr("a"),
								binHigherExpr(
									idExpr("b"),
									idExpr("c"),
								),
							),
							binHigherExpr(
								idExpr("d"),
								idExpr("e"),
							),
						),
						idExpr("f"),
					),
				},
			}
			got := Parse(logger, src)
			AssertExprEqual(t, logger, got, expected)
		}
	}
}

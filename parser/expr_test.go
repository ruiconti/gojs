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

func TestBinaryOperators(t *testing.T) {
	t.Run("properly parses simple, same-precedence, binary expr", func(t *testing.T) {
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

		assertBinaryExprPrecedence(
			t,
			opsLogicalAnd, /* higher */
			opsLogicalOr,  /* lower */
		)
	})

	t.Run("bitwise OR has precedence over logical AND", func(t *testing.T) {
		opsOr := []lex.TokenType{lex.TOr}
		opsLogicalAnd := []lex.TokenType{lex.TLogicalAnd}

		assertBinaryExprPrecedence(
			t,
			opsOr,         /* higher */
			opsLogicalAnd, /* lower */
		)
	})

	t.Run("bitwise XOR has precedence over bitwise OR", func(t *testing.T) {
		opsXor := []lex.TokenType{lex.TXor}
		opsOr := []lex.TokenType{lex.TOr}

		assertBinaryExprPrecedence(
			t,
			opsXor, /* higher */
			opsOr,  /* lower */
		)
	})

	t.Run("bitwise AND has precedence over bitwise XOR", func(t *testing.T) {
		opsAnd := []lex.TokenType{lex.TAnd}
		opsXor := []lex.TokenType{lex.TXor}

		assertBinaryExprPrecedence(
			t,
			opsAnd, /* higher */
			opsXor, /* lower */
		)
	})

	t.Run("equality comparison has precedence over bitwise AND", func(t *testing.T) {
		opsEq := []lex.TokenType{lex.TEqual, lex.TNotEqual, lex.TStrictEqual, lex.TStrictNotEqual}
		opsBitwise := []lex.TokenType{lex.TAnd, lex.TOr, lex.TXor}

		assertBinaryExprPrecedence(
			t,
			opsEq,      /* higher */
			opsBitwise, /* lower */
		)
	})
	t.Run("relational comparison has precedence over equality comparison", func(t *testing.T) {
		opsEq := []lex.TokenType{lex.TEqual, lex.TNotEqual, lex.TStrictEqual, lex.TStrictNotEqual}
		opsRelational := []lex.TokenType{lex.TLessThan, lex.TLessThanEqual, lex.TGreaterThan, lex.TGreaterThanEqual, lex.TIn, lex.TInstanceof}

		assertBinaryExprPrecedence(
			t,
			opsRelational, /* higher */
			opsEq,         /* lower */
		)
	})
	t.Run("shift operation has precedence over relational comparison", func(t *testing.T) {
		opsRelational := []lex.TokenType{lex.TLessThan, lex.TLessThanEqual, lex.TGreaterThan, lex.TGreaterThanEqual, lex.TIn, lex.TInstanceof}
		opsShift := []lex.TokenType{lex.TLeftShift, lex.TRightShift, lex.TUnsignedRightShift}

		assertBinaryExprPrecedence(
			t,
			opsShift,      /* higher */
			opsRelational, /* lower */
		)
	})

	t.Run("additive operation has precedence over shift operation", func(t *testing.T) {
		opsShift := []lex.TokenType{lex.TLeftShift, lex.TRightShift, lex.TUnsignedRightShift}
		opsAdd := []lex.TokenType{lex.TPlus, lex.TMinus}

		assertBinaryExprPrecedence(
			t,
			opsAdd,   /* higher */
			opsShift, /* lower */
		)
	})

	t.Run("multiplicative operation has precedence over additive operation", func(t *testing.T) {
		opsMult := []lex.TokenType{lex.TStar, lex.TSlash, lex.TPercent}
		opsAdd := []lex.TokenType{lex.TPlus, lex.TMinus}

		assertBinaryExprPrecedence(
			t,
			opsMult, /* higher */
			opsAdd,  /* lower */
		)
	})

	t.Run("exponential operation has precedence over multiplicative operation", func(t *testing.T) {
		opsMult := []lex.TokenType{lex.TStar, lex.TSlash, lex.TPercent}
		opsExp := []lex.TokenType{lex.TStarStar}

		assertBinaryExprPrecedence(
			t,
			opsExp,  /* higher */
			opsMult, /* lower */
		)
	})

	// TODO: this should err because the grammar doesn't support this operation
	// TODO: without parentheses
	// t.Run("unary operation has precedence over exponential operation", func(t *testing.T) {
	// 	multOpExpr := func(left AstNode, right AstNode) *ExprBinaryOp {
	// 		return &ExprBinaryOp{
	// 			operator: lex.TStarStar,
	// 			left:     left,
	// 			right:    right,
	// 		}
	// 	}

	// 	for _, operator := range UnaryOperators {
	// 		unaryOpExpr := func(binding string) *ExprUnaryOp {
	// 			return &ExprUnaryOp{
	// 				operator: operator,
	// 				operand: &ExprIdentifierReference{
	// 					reference: binding,
	// 				},
	// 			}
	// 		}

	// 		lexeme := lex.ResolveName(operator)
	// 		src := fmt.Sprintf("%s a ** b ** %s c ** d ** %s e ** f", lexeme, lexeme, lexeme)
	// 		// equals: delete a ** (b ** (delete c ** (d ** (delete e ** f) ) ) )
	// 		exp := &ExprRootNode{
	// 			children: []AstNode{
	// 				multOpExpr(
	// 					unaryOpExpr("a"),
	// 					multOpExpr(
	// 						idExpr("b"),
	// 						multOpExpr(
	// 							unaryOpExpr("c"),
	// 							multOpExpr(
	// 								idExpr("d"),
	// 								multOpExpr(
	// 									unaryOpExpr("e"),
	// 									idExpr("f"),
	// 								),
	// 							),
	// 						),
	// 					),
	// 				),
	// 			},
	// 		}

	// 		logger := internal.NewSimpleLogger(internal.ModeDebug)
	// 		got := Parse(logger, src)
	// 		AssertExprEqual(t, logger, got, exp)
	// 	}
	// })

	t.Run("unary expr called recursively", func(t *testing.T) {
		logger := internal.NewSimpleLogger(internal.ModeDebug)
		for _, operator := range UnaryOperators {
			operatorName := lex.ResolveName(operator)
			src := fmt.Sprintf("%s %s %s %s bar", operatorName, operatorName, operatorName, operatorName)
			got := Parse(logger, src)
			exp := &ExprRootNode{
				children: []AstNode{
					&ExprUnaryOp{
						operand: &ExprUnaryOp{
							operand: &ExprUnaryOp{
								operand: &ExprUnaryOp{
									operand: &ExprIdentifierReference{
										reference: "bar",
									},
									operator: operator,
								},
								operator: operator,
							},
							operator: operator,
						},
						operator: operator,
					},
				},
			}
			AssertExprEqual(t, logger, got, exp)
		}
	})

	t.Run("unary expression called with update expression", func(t *testing.T) {
		logger := internal.NewSimpleLogger(internal.ModeDebug)
		for _, unaryOp := range UnaryOperators {
			for _, updateOp := range UpdateOperators {
				src := fmt.Sprintf("%s %s foo", lex.ResolveName(unaryOp), lex.ResolveName(updateOp))
				got := Parse(logger, src)
				exp := &ExprRootNode{
					children: []AstNode{
						&ExprUnaryOp{
							operator: unaryOp,
							operand: &ExprUnaryOp{
								operator: updateOp,
								operand:  idExpr("foo"),
							},
						},
					},
				}

				AssertExprEqual(t, logger, got, exp)
			}
		}

	})

	t.Run("new expression called with member expression", func(t *testing.T) {
		logger := internal.NewSimpleLogger(internal.ModeDebug)
		src := fmt.Sprintf("new foo.bar[baz][foo2].bar2")
		got := Parse(logger, src)
		exp := &ExprRootNode{
			children: []AstNode{
				&ExprNew{
					callee: &ExprMemberAccess{
						object: &ExprMemberAccess{
							object: &ExprMemberAccess{
								object: &ExprMemberAccess{
									object:   idExpr("foo"),
									property: idExpr("bar"),
								},
								property: idExpr("baz"),
							},
							property: idExpr("foo2"),
						},
						property: idExpr("bar2"),
					},
				},
			},
		}

		AssertExprEqual(t, logger, got, exp)
	})
}

// assertBinaryExprPrecedence asserts that the given operators have the correct precedence.
func assertBinaryExprPrecedence(
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

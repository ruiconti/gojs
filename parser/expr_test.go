package parser

import (
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
	binOrExpr := func(left, right AstNode) *ExprBinaryOp {
		return binExpr(left, right, lex.TLogicalOr)
	}
	t.Run("logical OR expression", func(t *testing.T) {
		logger := internal.NewSimpleLogger(internal.ModeDebug)
		src := `a || b || c || d || e || f`
		expected := &ExprRootNode{
			children: []AstNode{
				binOrExpr(
					binOrExpr(
						binOrExpr(
							binOrExpr(
								binOrExpr(idExpr("a"), idExpr("b")),
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
	})

	binAndExpr := func(left, right AstNode) *ExprBinaryOp {
		return binExpr(left, right, lex.TLogicalAnd)
	}
	t.Run("logical AND expression", func(t *testing.T) {
		logger := internal.NewSimpleLogger(internal.ModeDebug)
		src := `a && b && c && d && e && f`
		expected := &ExprRootNode{
			children: []AstNode{
				binAndExpr(
					binAndExpr(
						binAndExpr(
							binAndExpr(
								binAndExpr(idExpr("a"), idExpr("b")),
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
	})
	t.Run("logical AND has precedence over OR expression 1", func(t *testing.T) {
		logger := internal.NewSimpleLogger(internal.ModeDebug)
		src := `a || b && c || d && e || f`
		// equals to: ((a || (b && c)) || (d && e)) || f
		expected := &ExprRootNode{
			children: []AstNode{
				binOrExpr(
					binOrExpr(
						binOrExpr(
							idExpr("a"),
							binAndExpr(idExpr("b"), idExpr("c")),
						),
						binAndExpr(idExpr("d"), idExpr("e")),
					),
					idExpr("f"),
				),
			},
		}
		got := Parse(logger, src)
		AssertExprEqual(t, logger, got, expected)
	})

	binBitOrExpr := func(left, right AstNode) *ExprBinaryOp {
		return binExpr(left, right, lex.TOr)
	}
	t.Run("bitwise OR has precedence over logical AND", func(t *testing.T) {
		logger := internal.NewSimpleLogger(internal.ModeDebug)
		src := `a | b && c | d && e | f`
		// equals to: ((a | b) && (c | d)) && (e | f))
		expected := &ExprRootNode{
			children: []AstNode{
				binAndExpr(
					binAndExpr(
						binBitOrExpr(idExpr("a"), idExpr("b")),
						binBitOrExpr(idExpr("c"), idExpr("d")),
					),
					binBitOrExpr(idExpr("e"), idExpr("f")),
				),
			},
		}
		got := Parse(logger, src)
		AssertExprEqual(t, logger, got, expected)
	})

	binBitXorExpr := func(left, right AstNode) *ExprBinaryOp {
		return binExpr(left, right, lex.TXor)
	}
	t.Run("bitwise XOR has precedence over bitwise OR", func(t *testing.T) {
		logger := internal.NewSimpleLogger(internal.ModeDebug)
		src := `a | b ^ c | d ^ e | f`
		// equals to: ((a | (b ^ c)) | (d ^ e)) | f
		expected := &ExprRootNode{
			children: []AstNode{
				binBitOrExpr(
					binBitOrExpr(
						binBitOrExpr(
							idExpr("a"),
							binBitXorExpr(idExpr("b"), idExpr("c")),
						),
						binBitXorExpr(idExpr("d"), idExpr("e"))),
					idExpr("f")),
			},
		}
		got := Parse(logger, src)
		AssertExprEqual(t, logger, got, expected)
	})
}

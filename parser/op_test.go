package parser

import (
	"fmt"
	"testing"

	"github.com/ruiconti/gojs/internal"
	"github.com/ruiconti/gojs/lex"
)

func TestUnaryOp(t *testing.T) {
	t.Run("simple identifier reference", func(t *testing.T) {
		logger := internal.NewSimpleLogger(internal.ModeDebug)
		for _, operator := range UnaryOperators {
			operatorName := lex.ResolveName(operator)
			src := fmt.Sprintf("%s foo", operatorName)
			got := Parse(logger, src)
			exp := &ExprRootNode{
				children: []AstNode{
					&ExprUnaryOp{
						operand: &ExprIdentifierReference{
							reference: "foo",
						},
						operator: operator,
					},
				},
			}
			AssertExprEqual(t, logger, got, exp)
		}
	})

	t.Run("simple literal", func(t *testing.T) {
	})
}

func TestUpdateExpr(t *testing.T) {
	t.Run("simple identifier reference", func(t *testing.T) {
		logger := internal.NewSimpleLogger(internal.ModeDebug)
		for _, operator := range UpdateOperators {
			src := fmt.Sprintf("%s foo", lex.ResolveName(operator))
			exp := &ExprRootNode{
				children: []AstNode{
					&ExprUnaryOp{
						operand: &ExprIdentifierReference{
							reference: "foo",
						},
						operator: operator,
					},
				},
			}
			got := Parse(logger, src)
			AssertExprEqual(t, logger, got, exp)
		}
	})
}

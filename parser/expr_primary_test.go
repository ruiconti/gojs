package parser

import (
	"testing"

	"github.com/ruiconti/gojs/internal"
)

func TestPrimaryLiterals(t *testing.T) {
	t.Run("literals basic", func(t *testing.T) {
		logger := internal.NewSimpleLogger(internal.ModeDebug)
		src := "123 true false null undefined \"foo\" 'bar'"
		got := Parse(logger, src)
		exp := &ExprRootNode{
			children: []AstNode{
				&ExprNumeric{
					value: 123,
				},
				&ExprBoolean{
					value: true,
				},
				&ExprBoolean{
					value: false,
				},
				&ExprNullLiteral{},
				&ExprUndefinedLiteral{},
				&ExprStringLiteral{
					value: "foo",
				},
				&ExprStringLiteral{
					value: "bar",
				},
			},
		}

		AssertExprEqual(t, logger, got, exp)
	})
	t.Run("literals unicode", func(t *testing.T) {
		logger := internal.NewSimpleLogger(internal.ModeDebug)
		src := `\u3034baz; \u9023\u4930\u1102x; b\u400e\u99a0`
		got := Parse(logger, src)
		exp := &ExprRootNode{
			children: []AstNode{
				&ExprIdentifierReference{
					reference: `\u3034baz`,
				},
				&ExprIdentifierReference{
					reference: `\u9023\u4930\u1102x`,
				},
				&ExprIdentifierReference{
					reference: `b\u400e\u99a0`,
				},
			},
		}
		AssertExprEqual(t, logger, got, exp)
	})

}

func TestExprIdentifierReference(t *testing.T) {
	logger := internal.NewSimpleLogger(internal.ModeDebug)
	src := "foo; bar; baz"

	got := Parse(logger, src)
	exp := &ExprRootNode{
		children: []AstNode{
			&ExprIdentifierReference{
				reference: "foo",
			},
			&ExprIdentifierReference{
				reference: "bar",
			},
			&ExprIdentifierReference{
				reference: "baz",
			},
		},
	}

	AssertExprEqual(t, logger, got, exp)
}

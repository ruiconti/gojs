package parser

import (
	"testing"

	"github.com/ruiconti/gojs/internal"
	l "github.com/ruiconti/gojs/lexer"
)

func TestParseVariableStatement(t *testing.T) {
	t.Run("variable statement", func(t *testing.T) {
		logger := internal.NewSimpleLogger(internal.ModeDebug)
		src := `var x = 10, y = 20;`
		kind := l.TVar

		exp := &NodeRoot{
			children: []Node{
				&VariableStatement{
					kind: kind.Token(),
					declarations: []*VariableDeclaration{
						{
							identifier: &ExprIdentifier{name: "x"},
							init:       intExpr(10),
						},
						{
							pattern: &ExprIdentifier{name: "y"},
							init:    intExpr(20),
						},
					},
				},
			},
		}

		got := Parse(logger, src)
		AssertStmtEqual(t, logger, got, exp)
	})

	t.Run("lexical declaration", func(t *testing.T) {
		logger := internal.NewSimpleLogger(internal.ModeDebug)
		src := `const a = 5, b = 10;`
		kind := l.TConst

		exp := &NodeRoot{
			children: []Node{
				&VariableStatement{
					kind: kind.Token(),
					declarations: []*VariableDeclaration{
						{
							identifier: &ExprIdentifier{name: "a"},
							init:       intExpr(5),
						},
						{
							identifier: &ExprIdentifier{name: "b"},
							init:       intExpr(10),
						},
					},
				},
			},
		}

		got := Parse(logger, src)
		AssertStmtEqual(t, logger, got, exp)
	})
}

func TestParseBindingPattern(t *testing.T) {
	t.Run("binding pattern", func(t *testing.T) {
		logger := internal.NewSimpleLogger(internal.ModeDebug)
		src := `let {a, b: x} = obj;`
		kind := l.TLet

		exp := &NodeRoot{
			children: []Node{
				&VariableStatement{
					kind: kind.Token(),
					declarations: []*VariableDeclaration{
						{
							pattern: &ExprObject{
								properties: []*PropertyDefinition{
									{
										key:       idExpr("a"),
										value:     idExpr("a"),
										shorthand: true,
									},
									{
										key:   idExpr("b"),
										value: idExpr("x"),
									},
								},
							},
							init: idExpr("obj"),
						},
					},
				},
			},
		}

		got := Parse(logger, src)
		AssertStmtEqual(t, logger, got, exp)
	})
}

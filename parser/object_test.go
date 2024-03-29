package parser

import (
	"testing"

	"github.com/ruiconti/gojs/internal"
	l "github.com/ruiconti/gojs/lexer"
)

var assignt = l.TAssign

func TestObjectInitialization(t *testing.T) {
	// Helper function to create an identifier expression
	t.Run("empty object", func(t *testing.T) {
		logger := internal.NewSimpleLogger(internal.ModeDebug)
		src := `a = {}`
		exp := &NodeRoot{
			children: []Node{
				&ExprAssign{
					operator: assignt.Token(),
					left:     idExpr("a"),
					right: &ExprObject{
						properties: []*PropertyDefinition{},
					},
				},
			},
		}
		got := Parse(logger, src)
		AssertExprEqual(t, logger, got, exp)
	})

	t.Run("single property", func(t *testing.T) {
		logger := internal.NewSimpleLogger(internal.ModeDebug)
		src := `a = {
    foo: 42
}`
		exp := &NodeRoot{
			children: []Node{
				&ExprAssign{
					operator: assignt.Token(),
					left:     idExpr("a"),
					right: &ExprObject{
						properties: []*PropertyDefinition{
							{
								key:   idExpr("foo"),
								value: intExpr(42),
							},
						},
					},
				},
			},
		}
		got := Parse(logger, src)
		AssertExprEqual(t, logger, got, exp)
	})

	t.Run("multiple properties", func(t *testing.T) {
		logger := internal.NewSimpleLogger(internal.ModeDebug)
		src := `a = {
			foo: "bar",
			num: 42,
			[2 + 2]: true
		}`
		op := l.TPlus
		exp := &NodeRoot{
			children: []Node{
				&ExprAssign{
					operator: assignt.Token(),
					left:     idExpr("a"),
					right: &ExprObject{
						properties: []*PropertyDefinition{
							{
								key:   idExpr("foo"),
								value: stringExpr(`"bar"`),
							},
							{
								key:   idExpr("num"),
								value: intExpr(42),
							},
							{
								computed: true,
								key: &ExprBinaryOp{
									operator: op.Token(),
									left:     intExpr(2),
									right:    intExpr(2),
								},
								value: MakeLiteralExpr(l.TTrue),
							},
						},
					},
				},
			},
		}
		got := Parse(logger, src)
		AssertExprEqual(t, logger, got, exp)
	})

	t.Run("single shorthand property", func(t *testing.T) {
		logger := internal.NewSimpleLogger(internal.ModeDebug)
		src := `a = {foo}`
		exp := &NodeRoot{
			children: []Node{
				&ExprAssign{
					operator: assignt.Token(),
					left:     idExpr("a"),
					right: &ExprObject{
						properties: []*PropertyDefinition{
							{
								key:       &ExprIdentifier{name: "foo"},
								value:     &ExprIdentifier{name: "foo"},
								computed:  false,
								method:    false,
								shorthand: true,
							},
						},
					},
				},
			},
		}
		got := Parse(logger, src)
		AssertExprEqual(t, logger, got, exp)
	})

	t.Run("spread operator", func(t *testing.T) {
		logger := internal.NewSimpleLogger(internal.ModeDebug)
		src := `a = {...foo, ...bar, baz, [foo > 'bar']: {...bar}}`
		exp := &NodeRoot{
			children: []Node{
				// write the expected AST here all at once
				&ExprAssign{
					operator: assignt.Token(),
					left:     idExpr("a"),
					right: &ExprObject{
						properties: []*PropertyDefinition{
							{
								key:   idExpr("foo"),
								value: &SpreadElement{argument: idExpr("foo")},
							},
							{
								key:   idExpr("bar"),
								value: &SpreadElement{argument: idExpr("bar")},
							},
							{
								key:       idExpr("baz"),
								value:     idExpr("baz"),
								shorthand: true,
							},
							{
								key: binExpr(idExpr("foo"), stringExpr(`'bar'`), l.TGreaterThan),
								value: &ExprObject{
									properties: []*PropertyDefinition{
										{
											key:   idExpr("bar"),
											value: &SpreadElement{argument: idExpr("bar")},
										},
									},
								},
								computed: true,
							},
						},
					},
				},
			},
		}
		got := Parse(logger, src)
		AssertExprEqual(t, logger, got, exp)
	})
}

package parser

import (
	"testing"

	"github.com/ruiconti/gojs/internal"
	l "github.com/ruiconti/gojs/lexer"
)

func TestParseArray_Simple(t *testing.T) {
	t.Run("empty array", func(t *testing.T) {
		logger := internal.NewSimpleLogger(internal.ModeDebug)
		src := `[]`
		expected := &ExprRootNode{
			children: []Node{
				&ExprArray{},
			},
		}
		got := Parse(logger, src)
		AssertExprEqual(t, logger, got, expected)
	})
	t.Run("full of elisions", func(t *testing.T) {
		logger := internal.NewSimpleLogger(internal.ModeDebug)
		src := `[,,, ,,   , ]`
		// src := `[null,null,null,null,null,null,]`
		expected := &ExprRootNode{
			children: []Node{
				&ExprArray{
					elements: []Node{
						ExprLitNull,
						ExprLitNull,
						ExprLitNull,
						ExprLitNull,
						ExprLitNull,
						ExprLitNull,
					},
				},
			},
		}
		got := Parse(logger, src)
		AssertExprEqual(t, logger, got, expected)
	})
	t.Run("full of primary expressions", func(t *testing.T) {
		logger := internal.NewSimpleLogger(internal.ModeDebug)
		src := `[1,2,true,\u3340xa,undefined, null,'foo', "bar",]`
		expected := &ExprRootNode{
			children: []Node{
				&ExprArray{
					elements: []Node{
						&ExprLiteral[float64]{l.Token{Type: l.TNumericLiteral, Literal: "1"}},
						&ExprLiteral[float64]{l.Token{Type: l.TNumericLiteral, Literal: "2"}},
						ExprLitTrue,
						&ExprIdentifier{
							name: `\u3340xa`,
						},
						ExprLitUndefined,
						ExprLitNull,
						&ExprLiteral[string]{l.Token{Type: l.TStringLiteral_SingleQuote, Literal: "'foo'"}},
						&ExprLiteral[string]{l.Token{Type: l.TStringLiteral_DoubleQuote, Literal: `"bar"`}},
					},
				},
			},
		}
		got := Parse(logger, src)
		AssertExprEqual(t, logger, got, expected)
	})
}

// ConditionalExpression is way too big
// An ArrayLiteral is made up of almost all JS expressions, so it does make sense to do them later

// Weirdly enough, to get to the raw value of an array element, you have to expand
// | ConditionalExpression
// | ShortCircuitExpression
// | LogicalORExpression
// | LogicalANDExpression
// | BitwiseORExpression
// | BitwiseXORExpression
// | BitwiseANDExpression
// | EqualityExpression
// | RelationalExpression
// | ShiftExpression
// | AdditiveExpression
// | MultiplicativeExpression
// | ExponentiationExpression
// | UnaryExpression
// | UpdateExpression
// | LeftHandSideExpression
// | NewExpression
// | MemberExpression
// | PrimaryExpression
// | IdentifierReference | ArrayLiteral | ObjectLiteral | FunctionExpression | ClassExpression | GeneratorExpression | RegularExpressionLiteral | TemplateLiteral

func TestParseArrayElementList_Assignment_Cond(t *testing.T) {
	t.Skip()
	// src := `[, a ? b : c, a ?? b, a?.b ?? c, d !== a ? b : c, a === b ? c : d]`
	// expected := ExprRootNode{}
	// got := Parse(src)
	// CompareRootChildren(
	// 	t,
	// 	src,
	// 	(got.children[0]).(*ExprArray).elements,
	// 	(expected.children[0]).(*ExprArray).elements,
	// )
}

func TestParseArrayElementList_Assignment_Yield(t *testing.T) {
	t.Skip()
	// src := `[, yield a]`
	// expected := ExprRootNode{}
	// got := Parse(src)
	// CompareRootChildren(
	// 	t,
	// 	src,
	// 	(got.children[0]).(*ExprArray).elements,
	// 	(expected.children[0]).(*ExprArray).elements,
	// )
}

func TestParseArrayElementList_Assignment_ArrowFunc(t *testing.T) {
	t.Skip()
	// src := `[, (a) => ({}), a => {}, ([a,b,{c}]) => c]`
	// expected := ExprRootNode{}
	// got := Parse(src)
	// CompareRootChildren(
	// 	t,
	// 	src,
	// 	(got.children[0]).(*ExprArray).elements,
	// 	(expected.children[0]).(*ExprArray).elements,
	// )
}

func TestParseArrayElementList_Assignment_AsyncArrowFunc(t *testing.T) {
	t.Skip()
	// src := `[, async (a) => ({}), async a => {}, async b => await b]`
	// expected := ExprRootNode{}
	// got := Parse(src)
	// CompareRootChildren(
	// 	t,
	// 	src,
	// 	(got.children[0]).(*ExprArray).elements,
	// 	(expected.children[0]).(*ExprArray).elements,
	// )
}

func TestParseArrayElementList_Assignment_LeftHS_NewExp1(t *testing.T) {
	t.Run("new class", func(t *testing.T) {
		logger := internal.NewSimpleLogger(internal.ModeDebug)
		src := `[, new Map([1, 2]), ]`
		exp := &ExprRootNode{
			children: []Node{
				&ExprArray{
					elements: []Node{
						ExprLitNull,
						&ExprNew{
							callee: &ExprIdentifier{
								name: "Map",
							},
							arguments: []Node{
								&ExprArray{
									elements: []Node{idExpr("1"), idExpr("2")},
								},
							},
						},
					},
				},
			},
		}
		got := Parse(logger, src)
		AssertExprEqual(t, logger, got, exp)
	})
	t.Run("new this", func(t *testing.T) {
		// logger := internal.NewSimpleLogger(internal.ModeDebug)
		// src := `[new this]`
	})
	t.Run("new call expr and member access", func(t *testing.T) {
		logger := internal.NewSimpleLogger(internal.ModeDebug)
		src := `[new t.p, new t.p(...x), a[b[c[d[e]]]]`
		exp := &ExprRootNode{
			children: []Node{
				&ExprArray{
					elements: []Node{
						&ExprNew{
							callee: &ExprMemberAccess{
								object:   idExpr("t"),
								property: idExpr("p"),
							},
						},
						&ExprNew{
							callee: &ExprMemberAccess{
								object:   idExpr("t"),
								property: idExpr("p"),
							},
							arguments: []Node{
								&SpreadElement{argument: idExpr("x")},
							},
						},
						&ExprMemberAccess{
							object: idExpr("a"),
							property: &ExprMemberAccess{
								object: idExpr("b"),
								property: &ExprMemberAccess{
									object: idExpr("c"),
									property: &ExprMemberAccess{
										object:   idExpr("d"),
										property: idExpr("e"),
									},
								},
							},
						},
					},
				},
			},
		}

		got := Parse(logger, src)
		AssertExprEqual(t, logger, got, exp)
	})
	t.Run("import and super expressions", func(t *testing.T) {
		logger := internal.NewSimpleLogger(internal.ModeDebug)
		src := `[import(a), super(a,...b,)]`
		exp := &ExprRootNode{
			children: []Node{
				&ExprArray{
					elements: []Node{
						&ExprImportCall{
							source: idExpr("a"),
						},
						&ExprCall{
							callee: MakeLiteralExpr(l.TSuper),
							arguments: []Node{
								idExpr("a"),
								&SpreadElement{argument: idExpr("b")},
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

func TestParseArrayElementList_Assignment_LeftHS_NewExp2(t *testing.T) {
	t.Skip()
	// src := `[, new Map() = 1, new this = 3, new t.p = 1, new t.p() = a, a[b] = x, a[b[c[d[e]]]] = \u8888(0,1,), () => import(a).x = x, super(a,b,...c) = \u4444]`
	// expected := ExprRootNode{}
	// got := Parse(src)
	// CompareRootChildren(
	// 	t,
	// 	src,
	// 	(got.children[0]).(*ExprArray).elements,
	// 	(expected.children[0]).(*ExprArray).elements,
	// )
}

func TestParseArrayElementList_Assignment_LeftHS_NewExp3(t *testing.T) {
	t.Skip()
	// src := `[, new Map() += 1, new this *= 3, new t.p &&= 1, new t.p() ||= a, a[b] /= x, a[b[c[d[e]]]] *= \u8888(0,1,), () => import(a).x &&&= x, super(a,b,...c) -= \u4444]`
	// expected := ExprRootNode{}
	// got := Parse(src)
	// CompareRootChildren(
	// 	t,
	// 	src,
	// 	(got.children[0]).(*ExprArray).elements,
	// 	(expected.children[0]).(*ExprArray).elements,
	// )
}

func TestParseArrayElementList_SpreadElement(t *testing.T) {
	t.Skip()
	// src := `[, ...new Map(), ...new this, ...new t.p, ...new t.p(), ...a[b], ...a[b[c[d[e]]]], ...() => import(a).x, ...super(a,b,...c)]`
	// expected := ExprRootNode{}
	// got := Parse(src)
	// CompareRootChildren(
	// 	t,
	// 	src,
	// 	(got.children[0]).(*ExprArray).elements,
	// 	(expected.children[0]).(*ExprArray).elements,
	// )
}

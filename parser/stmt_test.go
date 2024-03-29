package parser

import (
	"testing"

	"github.com/ruiconti/gojs/internal"
	l "github.com/ruiconti/gojs/lexer"
)

func TestParseVariableStatement(t *testing.T) {
	t.Run("empty block", func(t *testing.T) {
		logger := internal.NewSimpleLogger(internal.ModeDebug)
		src := `{}`
		exp := &NodeRoot{
			children: []Node{
				&BlockStatement{},
			},
		}
		got := Parse(logger, src)
		AssertStmtEqual(t, logger, got, exp)

	})

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
	t.Run("object binding pattern", func(t *testing.T) {
		logger := internal.NewSimpleLogger(internal.ModeDebug)
		src := `let {u, a: y, b: x, ...a} = obj;`
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
										key:       idExpr("u"),
										value:     idExpr("u"),
										shorthand: true,
									},
									{
										key:   idExpr("a"),
										value: idExpr("y"),
									},
									{
										key:   idExpr("b"),
										value: idExpr("x"),
									},
									{
										key:   idExpr("a"),
										value: &SpreadElement{argument: idExpr("a")},
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

	t.Run("array binding pattern", func(t *testing.T) {
		logger := internal.NewSimpleLogger(internal.ModeDebug)
		src := `let [f,b,...q] = arr;`
		kind := l.TLet

		exp := &NodeRoot{
			children: []Node{
				&VariableStatement{
					kind: kind.Token(),
					declarations: []*VariableDeclaration{
						{
							pattern: &ExprArray{
								elements: []Expr{
									idExpr("f"),
									idExpr("b"),
									&SpreadElement{argument: idExpr("q")},
								},
							},
							init: idExpr("arr"),
						},
					},
				},
			},
		}

		got := Parse(logger, src)
		AssertStmtEqual(t, logger, got, exp)
	})

}

func TestParseIfStatement(t *testing.T) {
	t.Run("if statement with else", func(t *testing.T) {
		logger := internal.NewSimpleLogger(internal.ModeDebug)
		src := `if (x > 10) { a = 1; } else { let b = 2; }`

		tassign := l.TAssign
		tlet := l.TLet

		exp := &NodeRoot{
			children: []Node{
				&IfStatement{
					Condition: binExpr(idExpr("x"), intExpr(10), l.TGreaterThan),
					ThenStmt: &BlockStatement{
						Stmts: []Stmt{
							&ExpressionStatement{
								expression: &ExprAssign{
									operator: tassign.Token(),
									left:     idExpr("a"),
									right:    intExpr(1),
								},
							},
						},
					},
					ElseStmt: &BlockStatement{
						Stmts: []Stmt{
							&VariableStatement{
								kind: tlet.Token(),
								declarations: []*VariableDeclaration{
									{
										identifier: idExpr("b"),
										init:       intExpr(2),
									},
								},
							},
						},
					},
				},
			},
		}

		got := Parse(logger, src)
		AssertStmtEqual(t, logger, got, exp)
	})

	t.Run("if statement without else", func(t *testing.T) {
		logger := internal.NewSimpleLogger(internal.ModeDebug)
		src := `if (x > 10) { a = b; }`

		tassign := l.TAssign
		exp := &NodeRoot{
			children: []Node{
				&IfStatement{
					Condition: binExpr(idExpr("x"), intExpr(10), l.TGreaterThan),
					ThenStmt: &BlockStatement{
						Stmts: []Stmt{
							&ExpressionStatement{
								expression: &ExprAssign{
									operator: tassign.Token(),
									left:     idExpr("a"),
									right:    idExpr("b"),
								},
							},
						},
					},
					ElseStmt: nil,
				},
			},
		}

		got := Parse(logger, src)
		AssertStmtEqual(t, logger, got, exp)
	})
}

func TestFunctionDeclaration(t *testing.T) {
	t.Run("regular function declaration", func(t *testing.T) {
		logger := internal.NewSimpleLogger(internal.ModeDebug)
		src := `function testFunction(a, b) {
			return a + b;
		}`
		exp := &NodeRoot{
			children: []Node{
				&FunctionDeclarationStmt{
					BindingIdentifier: idExpr("testFunction"),
					Params: []Node{
						idExpr("a"),
						idExpr("b"),
					},
					Body: []Stmt{
						&ReturnStatement{
							expr: binExpr(idExpr("a"), idExpr("b"), l.TPlus),
						},
					},
				},
			},
		}
		got := Parse(logger, src)
		AssertStmtEqual(t, logger, got, exp)
	})

	t.Run("function expression with patterns", func(t *testing.T) {
		logger := internal.NewSimpleLogger(internal.ModeDebug)
		src := `let fn = function({a, b:c}, [d], ...{e}) {
			return e + d;
		}`
		tlet := l.TLet
		exp := &NodeRoot{
			children: []Node{
				&VariableStatement{
					kind: tlet.Token(),
					declarations: []*VariableDeclaration{
						{
							identifier: idExpr("fn"),
							pattern:    nil,
							init: &ExprFunction{
								BindingIdentifier: nil,
								Params: []Node{
									&ExprObject{
										properties: []*PropertyDefinition{
											{key: idExpr("a"), value: idExpr("a"), shorthand: true},
											{key: idExpr("b"), value: idExpr("c")},
										},
									},
									&ExprArray{[]Expr{idExpr("d")}},
									&SpreadElement{&ExprObject{
										properties: []*PropertyDefinition{
											{key: idExpr("e"), value: idExpr("e"), shorthand: true},
										},
									}},
								},
								Body: []Stmt{
									&ReturnStatement{binExpr(idExpr("e"), idExpr("d"), l.TPlus)},
								},
							},
						},
					},
				},
			},
		}
		got := Parse(logger, src)
		AssertStmtEqual(t, logger, got, exp)
	})

	t.Run("simple function expression", func(t *testing.T) {
		logger := internal.NewSimpleLogger(internal.ModeDebug)
		src := `const fn = function() {
			return a;
		}`
		tconst := l.TConst
		exp := &NodeRoot{
			children: []Node{
				&VariableStatement{
					kind: tconst.Token(),
					declarations: []*VariableDeclaration{
						{
							identifier: idExpr("fn"),
							pattern:    nil,
							init: &ExprFunction{
								BindingIdentifier: nil,
								Params:            []Node{},
								Body: []Stmt{
									&ReturnStatement{idExpr("a")},
								},
							},
						},
					},
				},
			},
		}
		got := Parse(logger, src)
		AssertStmtEqual(t, logger, got, exp)
	})
}

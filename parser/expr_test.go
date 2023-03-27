package parser

import (
	"fmt"
	"testing"

	"github.com/ruiconti/gojs/internal"
	l "github.com/ruiconti/gojs/lexer"
)

// /////////////////////
// PrimaryExpression //
// /////////////////////
func TestPrimaryLiterals(t *testing.T) {
	t.Run("literals basic", func(t *testing.T) {
		logger := internal.NewSimpleLogger(internal.ModeDebug)
		src := "123 true false null undefined \"foo\" 'bar'"
		got := Parse(logger, src)
		exp := &ExprRootNode{
			children: []Node{
				&ExprLiteral[float64]{l.Token{Type: l.TNumericLiteral, Literal: "123"}},
				ExprLitTrue,
				ExprLitFalse,
				ExprLitNull,
				ExprLitUndefined,
				&ExprLiteral[string]{
					l.Token{Type: l.TStringLiteral_DoubleQuote, Literal: `"foo"`},
				},
				&ExprLiteral[string]{
					l.Token{Type: l.TStringLiteral_SingleQuote, Literal: `'bar'`},
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
			children: []Node{
				&ExprIdentifier{
					name: `\u3034baz`,
				},
				&ExprIdentifier{
					name: `\u9023\u4930\u1102x`,
				},
				&ExprIdentifier{
					name: `b\u400e\u99a0`,
				},
			},
		}
		AssertExprEqual(t, logger, got, exp)
	})

}

//////////////////////////
// IdentifierExpression //
//////////////////////////

func TestExprIdentifierReference(t *testing.T) {
	logger := internal.NewSimpleLogger(internal.ModeDebug)
	src := "foo; bar; baz"

	got := Parse(logger, src)
	exp := &ExprRootNode{
		children: []Node{
			&ExprIdentifier{
				name: "foo",
			},
			&ExprIdentifier{
				name: "bar",
			},
			&ExprIdentifier{
				name: "baz",
			},
		},
	}

	AssertExprEqual(t, logger, got, exp)
}

func binExpr(left, right Node, op l.TokenType) *ExprBinaryOp {
	return &ExprBinaryOp{
		left:     left,
		right:    right,
		operator: op.Token(),
	}
}

func stringExpr(s string) *ExprLiteral[string] {
	var st l.TokenType
	if s[0] == '"' {
		st = l.TStringLiteral_DoubleQuote
	} else {
		st = l.TStringLiteral_SingleQuote
	}
	return &ExprLiteral[string]{
		tok: l.Token{
			Literal: s,
			Lexeme:  s,
			Type:    st,
		},
	}
}

func idExpr(name string) *ExprIdentifier {
	return &ExprIdentifier{
		name: name,
	}
}
func spreadExpr(expr Node) *SpreadElement {
	return &SpreadElement{
		argument: expr,
	}
}

func idPrivateExpr(name string) *ExprPrivateIdentifier {
	return &ExprPrivateIdentifier{
		name: name,
	}
}

func intExpr(n int32) *ExprLiteral[float64] {
	return &ExprLiteral[float64]{
		tok: l.Token{
			Literal: float64(n),
			Lexeme:  fmt.Sprintf("%d", n),
			Type:    l.TNumericLiteral,
		},
	}
}

// //////////////
// Operations //
// /////////////

// Binary operators
func TestBinaryOperators(t *testing.T) {
	t.Run("properly parses simple, same-precedence, binary expr", func(t *testing.T) {
		logger := internal.NewSimpleLogger(internal.ModeDebug)
		binOperators := []l.TokenType{l.TLogicalOr,
			l.TLogicalAnd,
			l.TOr,
			l.TAnd,
			l.TXor,
			l.TEqual,
			l.TStrictEqual,
			l.TNotEqual,
			l.TStrictNotEqual,
			l.TGreaterThan,
			l.TGreaterThanEqual,
			l.TLessThan,
			l.TLessThanEqual,
			l.TLeftShift,
			l.TRightShift,
			l.TPlus,
			l.TMinus,
			l.TStar,
			l.TPercent,
			l.TSlash,
		}

		for _, binOperator := range binOperators {
			lexeme := binOperator.S()
			src := fmt.Sprintf(`a %s b %s c %s d %s e %s f`, lexeme, lexeme, lexeme, lexeme, lexeme)
			binExpr := func(left, right Node) *ExprBinaryOp {
				return binExpr(left, right, binOperator)
			}

			expected := &ExprRootNode{
				children: []Node{
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
		opsLogicalOr := []l.TokenType{l.TLogicalOr}
		opsLogicalAnd := []l.TokenType{l.TLogicalAnd}

		assertBinaryExprPrecedence(
			t,
			opsLogicalAnd, /* higher */
			opsLogicalOr,  /* lower */
		)
	})

	t.Run("bitwise OR has precedence over logical AND", func(t *testing.T) {
		opsOr := []l.TokenType{l.TOr}
		opsLogicalAnd := []l.TokenType{l.TLogicalAnd}

		assertBinaryExprPrecedence(
			t,
			opsOr,         /* higher */
			opsLogicalAnd, /* lower */
		)
	})

	t.Run("bitwise XOR has precedence over bitwise OR", func(t *testing.T) {
		opsXor := []l.TokenType{l.TXor}
		opsOr := []l.TokenType{l.TOr}

		assertBinaryExprPrecedence(
			t,
			opsXor, /* higher */
			opsOr,  /* lower */
		)
	})

	t.Run("bitwise AND has precedence over bitwise XOR", func(t *testing.T) {
		opsAnd := []l.TokenType{l.TAnd}
		opsXor := []l.TokenType{l.TXor}

		assertBinaryExprPrecedence(
			t,
			opsAnd, /* higher */
			opsXor, /* lower */
		)
	})

	t.Run("equality comparison has precedence over bitwise AND", func(t *testing.T) {
		opsEq := []l.TokenType{l.TEqual, l.TNotEqual, l.TStrictEqual, l.TStrictNotEqual}
		opsBitwise := []l.TokenType{l.TAnd, l.TOr, l.TXor}

		assertBinaryExprPrecedence(
			t,
			opsEq,      /* higher */
			opsBitwise, /* lower */
		)
	})
	t.Run("relational comparison has precedence over equality comparison", func(t *testing.T) {
		opsEq := []l.TokenType{l.TEqual, l.TNotEqual, l.TStrictEqual, l.TStrictNotEqual}
		opsRelational := []l.TokenType{l.TLessThan, l.TLessThanEqual, l.TGreaterThan, l.TGreaterThanEqual, l.TIn, l.TInstanceof}

		assertBinaryExprPrecedence(
			t,
			opsRelational, /* higher */
			opsEq,         /* lower */
		)
	})
	t.Run("shift operation has precedence over relational comparison", func(t *testing.T) {
		opsRelational := []l.TokenType{l.TLessThan, l.TLessThanEqual, l.TGreaterThan, l.TGreaterThanEqual, l.TIn, l.TInstanceof}
		opsShift := []l.TokenType{l.TLeftShift, l.TRightShift, l.TUnsignedRightShift}

		assertBinaryExprPrecedence(
			t,
			opsShift,      /* higher */
			opsRelational, /* lower */
		)
	})

	t.Run("additive operation has precedence over shift operation", func(t *testing.T) {
		opsShift := []l.TokenType{l.TLeftShift, l.TRightShift, l.TUnsignedRightShift}
		opsAdd := []l.TokenType{l.TPlus, l.TMinus}

		assertBinaryExprPrecedence(
			t,
			opsAdd,   /* higher */
			opsShift, /* lower */
		)
	})

	t.Run("multiplicative operation has precedence over additive operation", func(t *testing.T) {
		opsMult := []l.TokenType{l.TStar, l.TSlash, l.TPercent}
		opsAdd := []l.TokenType{l.TPlus, l.TMinus}

		assertBinaryExprPrecedence(
			t,
			opsMult, /* higher */
			opsAdd,  /* lower */
		)
	})

	t.Run("exponential operation has precedence over multiplicative operation", func(t *testing.T) {
		opsMult := []l.TokenType{l.TStar, l.TSlash, l.TPercent}
		opsExp := []l.TokenType{l.TStarStar}

		assertBinaryExprPrecedence(
			t,
			opsExp,  /* higher */
			opsMult, /* lower */
		)
	})
}

// Unary operators
func TestUnaryOperators(t *testing.T) {
	t.Run("unary operator with single reference binding", func(t *testing.T) {
		logger := internal.NewSimpleLogger(internal.ModeDebug)
		for _, operator := range UnaryOperators {
			src := fmt.Sprintf("%s foo", operator.S())
			got := Parse(logger, src)
			exp := &ExprRootNode{
				children: []Node{
					&ExprUnaryOp{
						operand: &ExprIdentifier{
							name: "foo",
						},
						operator: l.Token{
							Type:    operator,
							Literal: operator.S(),
						},
					},
				},
			}
			AssertExprEqual(t, logger, got, exp)
		}
	})

	t.Run("update expr with single reference binding", func(t *testing.T) {
		logger := internal.NewSimpleLogger(internal.ModeDebug)
		for _, operator := range UpdateOperators {
			src := fmt.Sprintf("%s foo", operator.S())
			exp := &ExprRootNode{
				children: []Node{
					&ExprUnaryOp{
						operand: &ExprIdentifier{
							name: "foo",
						},
						operator: operator.Token(),
					},
				},
			}
			got := Parse(logger, src)
			AssertExprEqual(t, logger, got, exp)
		}
	})
	// TODO: this should err because the grammar doesn't support this operation
	// TODO: without parentheses
	// t.Run("unary operation has precedence over exponential operation", func(t *testing.T) {
	// 	multOpExpr := func(left Node, right Node) *ExprBinaryOp {
	// 		return &ExprBinaryOp{
	// 			operator: l.TStarStar,
	// 			left:     left,
	// 			right:    right,
	// 		}
	// 	}

	// 	for _, operator := range UnaryOperators {
	// 		unaryOpExpr := func(binding string) *ExprUnaryOp {
	// 			return &ExprUnaryOp{
	// 				operator: operator,
	// 				operand: &ExprIdentifier{
	// 					name: binding,
	// 				},
	// 			}
	// 		}

	// 		lexeme := l.ResolveName(operator)
	// 		src := fmt.Sprintf("%s a ** b ** %s c ** d ** %s e ** f", lexeme, lexeme, lexeme)
	// 		// equals: delete a ** (b ** (delete c ** (d ** (delete e ** f) ) ) )
	// 		exp := &ExprRootNode{
	// 			children: []Node{
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
			operatorName := operator.S()
			operatorToken := operator.Token()
			src := fmt.Sprintf("%s %s %s %s bar", operatorName, operatorName, operatorName, operatorName)
			got := Parse(logger, src)
			exp := &ExprRootNode{
				children: []Node{
					&ExprUnaryOp{
						operand: &ExprUnaryOp{
							operand: &ExprUnaryOp{
								operand: &ExprUnaryOp{
									operand: &ExprIdentifier{
										name: "bar",
									},
									operator: operatorToken,
								},
								operator: operatorToken,
							},
							operator: operatorToken,
						},
						operator: operatorToken,
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
				src := fmt.Sprintf("%s %s foo", unaryOp.S(), updateOp.S())
				got := Parse(logger, src)
				exp := &ExprRootNode{
					children: []Node{
						&ExprUnaryOp{
							operator: unaryOp.Token(),
							operand: &ExprUnaryOp{
								operator: updateOp.Token(),
								operand:  idExpr("foo"),
							},
						},
					},
				}

				AssertExprEqual(t, logger, got, exp)
			}
		}

	})
}

func TestMemberAndNewExpressions(t *testing.T) {
	t.Run("primary expression", func(t *testing.T) {
		logger := internal.NewSimpleLogger(internal.ModeDebug)
		src := `foo`
		got := Parse(internal.NewSimpleLogger(internal.ModeDebug), src)
		exp := &ExprRootNode{
			children: []Node{
				idExpr("foo"),
			},
		}
		AssertExprEqual(t, logger, got, exp)
	})

	t.Run("computed property access", func(t *testing.T) {
		logger := internal.NewSimpleLogger(internal.ModeDebug)
		src := `foo[bar]`
		got := Parse(internal.NewSimpleLogger(internal.ModeDebug), src)
		exp := &ExprRootNode{
			children: []Node{
				&ExprMemberAccess{
					object:   idExpr("foo"),
					property: idExpr("bar"),
				},
			},
		}
		AssertExprEqual(t, logger, got, exp)
	})

	t.Run("static property access", func(t *testing.T) {
		logger := internal.NewSimpleLogger(internal.ModeDebug)
		src := `foo.bar`
		got := Parse(internal.NewSimpleLogger(internal.ModeDebug), src)
		exp := &ExprRootNode{
			children: []Node{
				&ExprMemberAccess{
					object:   idExpr("foo"),
					property: idExpr("bar"),
				},
			},
		}
		AssertExprEqual(t, logger, got, exp)
	})

	t.Run("template literal as tagged", func(t *testing.T) {
		t.Skip()
		logger := internal.NewSimpleLogger(internal.ModeDebug)
		src := "foo`bar`"
		got := Parse(internal.NewSimpleLogger(internal.ModeDebug), src)
		exp := &ExprRootNode{
			children: []Node{
				// &ExprTaggedTemplate{
				// 	tag:      idExpr("foo"),
				// 	template: templateExpr("bar"),
				// },
			},
		}
		AssertExprEqual(t, logger, got, exp)
	})

	t.Run("super property", func(t *testing.T) {
		logger := internal.NewSimpleLogger(internal.ModeDebug)
		src := `super.foo`
		got := Parse(internal.NewSimpleLogger(internal.ModeDebug), src)
		exp := &ExprRootNode{
			children: []Node{
				&ExprMemberAccess{
					object:   MakeLiteralExpr(l.TSuper),
					property: idExpr("foo"),
				},
			},
		}
		AssertExprEqual(t, logger, got, exp)
	})

	t.Run("meta property: new.target", func(t *testing.T) {
		logger := internal.NewSimpleLogger(internal.ModeDebug)
		src := `new.target`
		got := Parse(internal.NewSimpleLogger(internal.ModeDebug), src)
		exp := &ExprRootNode{
			children: []Node{
				&ExprMetaProperty{
					meta:     idExpr("new"),
					property: idExpr("target"),
				},
			},
		}
		AssertExprEqual(t, logger, got, exp)
	})

	t.Run("meta property: import.meta", func(t *testing.T) {
		logger := internal.NewSimpleLogger(internal.ModeDebug)
		src := `import.meta`
		got := Parse(internal.NewSimpleLogger(internal.ModeDebug), src)
		exp := &ExprRootNode{
			children: []Node{
				&ExprMetaProperty{
					meta:     idExpr("import"),
					property: idExpr("meta"),
				},
			},
		}
		AssertExprEqual(t, logger, got, exp)
	})
	t.Run("new expression with arguments", func(t *testing.T) {
		logger := internal.NewSimpleLogger(internal.ModeDebug)
		src := `new foo(bar)`
		got := Parse(internal.NewSimpleLogger(internal.ModeDebug), src)
		exp := &ExprRootNode{
			children: []Node{
				&ExprNew{
					callee: idExpr("foo"),
					arguments: []Node{
						idExpr("bar"),
					},
				},
			},
		}
		AssertExprEqual(t, logger, got, exp)
	})

	t.Run("private identifier", func(t *testing.T) {
		logger := internal.NewSimpleLogger(internal.ModeDebug)
		src := `foo.#bar`
		got := Parse(internal.NewSimpleLogger(internal.ModeDebug), src)
		exp := &ExprRootNode{
			children: []Node{
				&ExprMemberAccess{
					object: idExpr("foo"),
					property: &ExprPrivateIdentifier{
						name: "bar",
					},
				},
			},
		}
		AssertExprEqual(t, logger, got, exp)
	})
	t.Run("member expressions combined", func(t *testing.T) {
		logger := internal.NewSimpleLogger(internal.ModeDebug)
		src := `foo.bar[baz][foo2].bar2`
		got := Parse(logger, src)
		exp := &ExprRootNode{
			children: []Node{
				&ExprMemberAccess{
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
		}
		AssertExprEqual(t, logger, got, exp)
	})

	t.Run("new expression with member expression", func(t *testing.T) {
		logger := internal.NewSimpleLogger(internal.ModeDebug)
		src := `new foo.bar[baz][foo2].bar2`
		got := Parse(logger, src)
		exp := &ExprRootNode{
			children: []Node{
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

	t.Run("new expression recursive", func(t *testing.T) {
		logger := internal.NewSimpleLogger(internal.ModeDebug)
		src := `new new new new foo[0 >> 2]`
		got := Parse(logger, src)
		exp := &ExprRootNode{
			children: []Node{
				&ExprNew{
					callee: &ExprNew{
						callee: &ExprNew{
							callee: &ExprNew{
								callee: &ExprMemberAccess{
									object:   idExpr("foo"),
									property: binExpr(intExpr(0), intExpr(2), l.TRightShift),
								},
							},
						},
					},
				},
			},
		}

		AssertExprEqual(t, logger, got, exp)
	})

	t.Run("member expression with optional chaining", func(t *testing.T) {
		logger := internal.NewSimpleLogger(internal.ModeDebug)
		srcs := []string{
			`foo?.bar`,
			`foo?.#bar`,
			`foo?.[bar]`,
			`foo?.['bar']`,
		}
		expectedProps := []Node{
			idExpr("bar"),
			idPrivateExpr("bar"),
			idExpr("bar"),
			stringExpr("'bar'"),
		}

		for i := 0; i < len(srcs); i++ {
			src, expectedProp := srcs[i], expectedProps[i]
			exp := &ExprRootNode{
				children: []Node{
					&ExprMemberAccess{
						object:   idExpr("foo"),
						property: expectedProp,
						optional: true,
					},
				},
			}
			got := Parse(logger, src)
			AssertExprEqual(t, logger, got, exp)
		}
	})

}

func TestCallExpression(t *testing.T) {
	t.Run("simple call expression", func(t *testing.T) {
		logger := internal.NewSimpleLogger(internal.ModeDebug)
		src := `foo(bar)`
		got := Parse(logger, src)
		exp := &ExprRootNode{
			children: []Node{
				&ExprCall{
					callee: idExpr("foo"),
					arguments: []Node{
						idExpr("bar"),
					},
				},
			},
		}
		AssertExprEqual(t, logger, got, exp)
	})

	t.Run("call expression with spread", func(t *testing.T) {
		logger := internal.NewSimpleLogger(internal.ModeDebug)
		src := `foo(bar, baz, ...qux)`
		got := Parse(logger, src)
		exp := &ExprRootNode{
			children: []Node{
				&ExprCall{
					callee: idExpr("foo"),
					arguments: []Node{
						idExpr("bar"), idExpr("baz"), spreadExpr(idExpr("qux")),
					},
				},
			},
		}
		AssertExprEqual(t, logger, got, exp)
	})

	t.Run("call expression with template literal", func(t *testing.T) {
		t.Skip()
		logger := internal.NewSimpleLogger(internal.ModeDebug)
		src := "foo`bar`"
		got := Parse(logger, src)
		exp := &ExprRootNode{
			children: []Node{
				&ExprCall{
					callee:    idExpr("foo"),
					arguments: []Node{
						// templateExpr("bar"),
					},
				},
			},
		}
		AssertExprEqual(t, logger, got, exp)
	})

	t.Run("call expression with computed property", func(t *testing.T) {
		logger := internal.NewSimpleLogger(internal.ModeDebug)
		src := `foo[bar()]`
		got := Parse(logger, src)
		exp := &ExprRootNode{
			children: []Node{
				&ExprMemberAccess{
					object: idExpr("foo"),
					property: &ExprCall{
						callee:    idExpr("bar"),
						arguments: []Node{},
					},
				},
			},
		}
		AssertExprEqual(t, logger, got, exp)
	})

	t.Run("call expression with super", func(t *testing.T) {
		logger := internal.NewSimpleLogger(internal.ModeDebug)
		src := `super.foo()`
		got := Parse(logger, src)
		exp := &ExprRootNode{
			children: []Node{
				&ExprCall{
					callee: &ExprMemberAccess{
						object:   idExpr("super"),
						property: idExpr("foo"),
					},
					arguments: []Node{},
				},
			},
		}
		AssertExprEqual(t, logger, got, exp)
	})

	t.Run("call expression with import", func(t *testing.T) {
		logger := internal.NewSimpleLogger(internal.ModeDebug)
		src := `import("foo.js")`
		got := Parse(logger, src)
		exp := &ExprRootNode{
			children: []Node{
				&ExprImportCall{
					source: stringExpr(`"foo.js"`),
				},
			},
		}
		AssertExprEqual(t, logger, got, exp)
	})

	t.Run("nested call expression with import", func(t *testing.T) {
		logger := internal.NewSimpleLogger(internal.ModeDebug)
		src := `import("foo.js")(bar,'baz')`
		got := Parse(logger, src)
		exp := &ExprRootNode{
			children: []Node{
				&ExprCall{
					callee: &ExprImportCall{
						source: stringExpr(`"foo.js"`),
					},
					arguments: []Node{
						idExpr("bar"),
						stringExpr(`'baz'`),
					},
				},
			},
		}
		AssertExprEqual(t, logger, got, exp)
	})

	t.Run("call expression with private identifier", func(t *testing.T) {
		logger := internal.NewSimpleLogger(internal.ModeDebug)
		src := `foo.#bar(a, b)`
		got := Parse(logger, src)
		exp := &ExprRootNode{
			children: []Node{
				&ExprCall{
					callee: &ExprMemberAccess{
						object: idExpr("foo"),
						property: &ExprPrivateIdentifier{
							name: "bar",
						},
					},
					arguments: []Node{idExpr("a"), idExpr("b")},
				},
			},
		}
		AssertExprEqual(t, logger, got, exp)
	})

	t.Run("nested call expression", func(t *testing.T) {
		logger := internal.NewSimpleLogger(internal.ModeDebug)
		src := `foo(bar(baz(qux)))`
		got := Parse(logger, src)
		exp := &ExprRootNode{
			children: []Node{
				&ExprCall{
					callee: idExpr("foo"),
					arguments: []Node{
						&ExprCall{
							callee: idExpr("bar"),
							arguments: []Node{
								&ExprCall{
									callee: idExpr("baz"),
									arguments: []Node{
										idExpr("qux"),
									},
								},
							},
						},
					},
				},
			},
		}
		AssertExprEqual(t, logger, got, exp)
	})

	t.Run("call expression with optional chaining", func(t *testing.T) {
		logger := internal.NewSimpleLogger(internal.ModeDebug)
		srcs := []string{
			`foo?.(bar)`,
			`foo?.(a)?.(b)`,
			`foo?.()(bar)?.(baz)`,
		}
		expectedNodes := []Node{
			&ExprCall{
				callee:    idExpr("foo"),
				arguments: []Node{idExpr("bar")},
				optional:  true,
			},
			&ExprCall{
				callee: &ExprCall{
					callee:    idExpr("foo"),
					arguments: []Node{idExpr("a")},
					optional:  true,
				},
				arguments: []Node{idExpr("b")},
				optional:  true,
			},
			&ExprCall{
				callee: &ExprCall{
					callee: &ExprCall{
						callee:    idExpr("foo"),
						arguments: []Node{},
						optional:  true,
					},
					arguments: []Node{idExpr("bar")},
					optional:  false,
				},
				arguments: []Node{idExpr("baz")},
				optional:  true,
			},
		}

		for i := 0; i < len(srcs); i++ {
			src, expectedNode := srcs[i], expectedNodes[i]
			exp := &ExprRootNode{
				children: []Node{expectedNode},
			}
			got := Parse(logger, src)
			AssertExprEqual(t, logger, got, exp)
		}
	})

	t.Run("call expression with optional chaining and private identifier", func(t *testing.T) {
		logger := internal.NewSimpleLogger(internal.ModeDebug)
		src := `foo?.#bar(a, b)?.c`
		exp := &ExprRootNode{
			children: []Node{
				&ExprMemberAccess{
					property: idExpr("c"),
					optional: true,
					object: &ExprCall{
						callee: &ExprMemberAccess{
							object:   idExpr("foo"),
							optional: true,
							property: &ExprPrivateIdentifier{
								name: "bar",
							},
						},
						arguments: []Node{idExpr("a"), idExpr("b")},
					},
				},
			},
		}
		got := Parse(logger, src)
		AssertExprEqual(t, logger, got, exp)
	})

}

// //////////
// Helpers //
// //////////
//
// assertBinaryExprPrecedence asserts that the given operators have the correct precedence.
func assertBinaryExprPrecedence(
	t *testing.T,
	opsHigherPrecedence []l.TokenType,
	opsLowerPrecedence []l.TokenType,
) {
	logger := internal.NewSimpleLogger(internal.ModeDebug)

	for _, opHigher := range opsHigherPrecedence {
		lexemeHigherPrecedence := opHigher.S()
		binHigherExpr := func(left, right Node) *ExprBinaryOp {
			return binExpr(left, right, opHigher)
		}
		for _, opLower := range opsLowerPrecedence {
			lexemeLowerPrecedence := opLower.S()
			binLowerExpr := func(left, right Node) *ExprBinaryOp {
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
				children: []Node{
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

package parser

import "testing"

func TestParseArrayElision(t *testing.T) {
	src := `[,,, ,,   , ]`
	expected := ExprRootNode{
		children: []AstNode{
			&ExprArray{
				elements: []AstNode{
					&ExprElision{},
					&ExprElision{},
					&ExprElision{},
					&ExprElision{},
					&ExprElision{},
					&ExprElision{},
				},
			},
		},
	}
	got := Parse(src)
	CompareRootChildren(
		t,
		src,
		(got.children[0]).(*ExprArray).elements,
		(expected.children[0]).(*ExprArray).elements,
	)
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
	src := `[, a ? b : c, a ?? b, a?.b ?? c, d !== a ? b : c, a === b ? c : d]`
	expected := ExprRootNode{}
	got := Parse(src)
	CompareRootChildren(
		t,
		src,
		(got.children[0]).(*ExprArray).elements,
		(expected.children[0]).(*ExprArray).elements,
	)
}

func TestParseArrayElementList_Assignment_Yield(t *testing.T) {
	t.Skip()
	src := `[, yield a]`
	expected := ExprRootNode{}
	got := Parse(src)
	CompareRootChildren(
		t,
		src,
		(got.children[0]).(*ExprArray).elements,
		(expected.children[0]).(*ExprArray).elements,
	)
}

func TestParseArrayElementList_Assignment_ArrowFunc(t *testing.T) {
	t.Skip()
	src := `[, (a) => ({}), a => {}, ([a,b,{c}]) => c]`
	expected := ExprRootNode{}
	got := Parse(src)
	CompareRootChildren(
		t,
		src,
		(got.children[0]).(*ExprArray).elements,
		(expected.children[0]).(*ExprArray).elements,
	)
}

func TestParseArrayElementList_Assignment_AsyncArrowFunc(t *testing.T) {
	t.Skip()
	src := `[, async (a) => ({}), async a => {}, async b => await b]`
	expected := ExprRootNode{}
	got := Parse(src)
	CompareRootChildren(
		t,
		src,
		(got.children[0]).(*ExprArray).elements,
		(expected.children[0]).(*ExprArray).elements,
	)
}

func TestParseArrayElementList_Assignment_LeftHS_NewExp1(t *testing.T) {
	t.Skip()
	src := `[, new Map(), new this, new t.p, new t.p(), a[b], a[b[c[d[e]]]], () => import(a, b[c[d]]), super(a,...b,)]`
	expected := ExprRootNode{}
	got := Parse(src)
	CompareRootChildren(
		t,
		src,
		(got.children[0]).(*ExprArray).elements,
		(expected.children[0]).(*ExprArray).elements,
	)
}

func TestParseArrayElementList_Assignment_LeftHS_NewExp2(t *testing.T) {
	t.Skip()
	src := `[, new Map() = 1, new this = 3, new t.p = 1, new t.p() = a, a[b] = x, a[b[c[d[e]]]] = \u8888(0,1,), () => import(a).x = x, super(a,b,...c) = \u4444]`
	expected := ExprRootNode{}
	got := Parse(src)
	CompareRootChildren(
		t,
		src,
		(got.children[0]).(*ExprArray).elements,
		(expected.children[0]).(*ExprArray).elements,
	)
}

func TestParseArrayElementList_Assignment_LeftHS_NewExp3(t *testing.T) {
	t.Skip()
	src := `[, new Map() += 1, new this *= 3, new t.p &&= 1, new t.p() ||= a, a[b] /= x, a[b[c[d[e]]]] *= \u8888(0,1,), () => import(a).x &&&= x, super(a,b,...c) -= \u4444]`
	expected := ExprRootNode{}
	got := Parse(src)
	CompareRootChildren(
		t,
		src,
		(got.children[0]).(*ExprArray).elements,
		(expected.children[0]).(*ExprArray).elements,
	)
}

func TestParseArrayElementList_SpreadElement(t *testing.T) {
	t.Skip()
	src := `[, ...new Map(), ...new this, ...new t.p, ...new t.p(), ...a[b], ...a[b[c[d[e]]]], ...() => import(a).x, ...super(a,b,...c)]`
	expected := ExprRootNode{}
	got := Parse(src)
	CompareRootChildren(
		t,
		src,
		(got.children[0]).(*ExprArray).elements,
		(expected.children[0]).(*ExprArray).elements,
	)
}

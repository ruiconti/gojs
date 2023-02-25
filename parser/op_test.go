package parser

import (
	"testing"

	"github.com/ruiconti/gojs/lex"
)

func TestUnaryOp_Delete(t *testing.T) {
	src := "delete foo"
	got := Parse(src)
	exp := &ExprRootNode{
		children: []AstNode{
			&ExprUnaryOp{
				operand: &ExprIdentifierReference{
					reference: "foo",
				},
				operator: lex.TDelete,
			},
		},
	}
	expt := exp.children[0].(*ExprUnaryOp)
	gott := got.children[0].(*ExprUnaryOp)
	if expt.operator != gott.operator {
		t.Errorf("Expected %s, got %s", lex.ReservedWordNames[expt.operator], lex.ReservedWordNames[gott.operator])
		t.Fail()
	}

	expop := expt.operand.(*ExprIdentifierReference)
	gotop := gott.operand.(*ExprIdentifierReference)
	if expop.reference != gotop.reference {
		t.Errorf("Expected %s, got %s", expop.reference, gotop.reference)
		t.Fail()
	}
}

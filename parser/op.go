package parser

import (
	"errors"
	"fmt"

	"github.com/ruiconti/gojs/lex"
)

// ////////////////////////////
// Unary Operator Expression
// ////////////////////////////
const EUnaryOp ExprType = "EUnaryOp"

var errNotUnaryOperator = errors.New("current token is not an unary operator")
var errNotUpdateOperator = errors.New("current token is not an update operator")

type ExprUnaryOp struct {
	operand  Node
	operator lex.TokenType
}

func (e *ExprUnaryOp) Name() string {
	return "UnaryOperatorExpression"
}

func (e *ExprUnaryOp) Type() ExprType {
	return EUnaryOp
}

func (e *ExprUnaryOp) Source() string {
	return e.operand.Source()
}

func (e *ExprUnaryOp) PrettyPrint() string {
	return fmt.Sprintf("(%s %s)", lex.ResolveName(e.operator), e.operand.PrettyPrint())
}

var UnaryOperators = []lex.TokenType{
	lex.TDelete,
	lex.TTypeof,
	lex.TVoid,
	lex.TPlus,
	lex.TMinus,
	lex.TBang,
	lex.TTilde,
}

// Parser

var UpdateOperators = []lex.TokenType{
	lex.TMinusMinus,
	lex.TPlusPlus,
}

// ////////////////////////////
// Binary Operator Expression
// ////////////////////////////
const EBinaryOp ExprType = "ExprBinaryOp"

type ExprBinaryOp struct {
	left     Node
	right    Node
	operator lex.TokenType
}

func (e *ExprBinaryOp) Type() ExprType {
	return EBinaryOp
}

func (e *ExprBinaryOp) Source() string {
	return fmt.Sprintf("%s %s %s", e.left.Source(), lex.ResolveName(e.operator), e.right.Source())
}

func (e *ExprBinaryOp) PrettyPrint() string {
	return fmt.Sprintf("(%s %s %s)", lex.ResolveName(e.operator), e.left.PrettyPrint(), e.right.PrettyPrint())
}

// /////////////////////////////
// New Expression
// /////////////////////////////
const ENew ExprType = "ExprBinaryOp"

type ExprNew struct {
	callee Node
}

func (e *ExprNew) Type() ExprType {
	return ENew
}

func (e *ExprNew) Source() string {
	return fmt.Sprintf("new %s", e.callee.Source())
}

func (e *ExprNew) PrettyPrint() string {
	return fmt.Sprintf("(new %s)", e.callee.PrettyPrint())
}

// /////////////////////////////
// Member Expression
// /////////////////////////////
const EMemberAccess ExprType = "ExprBinaryOp"

type ExprMemberAccess struct {
	object   Node
	property Node
}

func (e *ExprMemberAccess) Type() ExprType {
	return EMemberAccess
}

func (e *ExprMemberAccess) Source() string {
	return fmt.Sprintf("%s[%s]", e.property.PrettyPrint(), e.object.PrettyPrint())
}

func (e *ExprMemberAccess) PrettyPrint() string {
	return fmt.Sprintf("(<- %s %s)", e.property.PrettyPrint(), e.object.PrettyPrint())
}

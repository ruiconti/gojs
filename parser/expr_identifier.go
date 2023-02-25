package parser

import (
	"fmt"

	"github.com/ruiconti/gojs/lex"
)

// RootNode: Artificial node
type ExprRootNode struct {
	children []AstNode
}

func (e *ExprRootNode) Source() string {
	return ""
}

func (e *ExprRootNode) Type() ExprType {
	return ENodeRoot
}

// ExprIdentifier
const EIdentifierReference ExprType = "EIdentifierReference"

type ExprIdentifierReference struct {
	reference string
}

func (e *ExprIdentifierReference) Source() string {
	return e.reference
}

func (e *ExprIdentifierReference) Type() ExprType {
	return EIdentifierReference
}

// ExprNumeric
const ENumeric ExprType = "ENumeric"

type ExprNumeric struct {
	value string
}

func (e *ExprNumeric) Source() string {
	return e.value
}

func (e *ExprNumeric) Type() ExprType {
	return ENumeric
}

// ExprBoolean
const EBool ExprType = "EBool"

type ExprBoolean struct {
	value bool
}

func (e *ExprBoolean) Source() string {
	return fmt.Sprintf("%v", e.value)
}

func (e *ExprBoolean) Type() ExprType {
	return EBool
}

// ExprStringLiteral
const EStringLiteral ExprType = "EStringLiteral"

type ExprStringLiteral struct {
	value string
}

func (e *ExprStringLiteral) Source() string {
	return e.value
}

func (e *ExprStringLiteral) Type() ExprType {
	return EStringLiteral
}

// ExprNull
const ENull ExprType = "ENull"

type ExprNullLiteral struct{}

func (e *ExprNullLiteral) Source() string {
	return lex.ReservedWordNames[lex.TNull]
}

func (e *ExprNullLiteral) Type() ExprType {
	return ENull
}

// ExprUndefined
const EUndefined ExprType = " EUndefined"

type ExprUndefinedLiteral struct{}

func (e *ExprUndefinedLiteral) Source() string {
	return lex.ReservedWordNames[lex.TUndefined]
}

func (e *ExprUndefinedLiteral) Type() ExprType {
	return EUndefined
}

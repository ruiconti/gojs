package parser

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ruiconti/gojs/lex"
)

// --------------------------
// RootNode: Artificial node
// --------------------------
type ExprRootNode struct {
	children []Node
}

func (e *ExprRootNode) Source() string {
	src := strings.Builder{}
	for _, child := range e.children {
		src.Write([]byte(child.Source()))
	}
	return src.String()
}

func (e *ExprRootNode) Type() ExprType {
	return ENodeRoot
}

func (e *ExprRootNode) PrettyPrint() string {
	pp := strings.Builder{}
	pp.Write([]byte("("))
	for i, child := range e.children {
		pp.Write([]byte(child.PrettyPrint()))
		if i < len(e.children)-1 {
			pp.Write([]byte(" "))
		}

	}
	pp.Write([]byte(")"))
	return pp.String()
}

// --------------
// ExprIdentifier
// --------------
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

func (e *ExprIdentifierReference) PrettyPrint() string {
	return e.reference
}

// ------------------
// ---- LITERALS ----
// ------------------
var errNotPrimaryExpr = errors.New("current token is not a primary expression")

var LiteralsTokens = []lex.TokenType{
	lex.TNumericLiteral,
	lex.TRegularExpressionLiteral,
	lex.TStringLiteral_SingleQuote,
	lex.TStringLiteral_DoubleQuote,
	lex.TTrue,
	lex.TFalse,
	lex.TNull,
}

func isLiteralToken(token lex.Token) bool {
	for _, t := range LiteralsTokens {
		if token.Type == t {
			return true
		}
	}
	return false
}

// -----------
// ExprNumeric
// -----------
const ENumeric ExprType = "ENumeric"

type ExprNumeric struct {
	value float64
}

func (e *ExprNumeric) Source() string {
	return fmt.Sprintf("%f", e.value)
}

func (e *ExprNumeric) Type() ExprType {
	return ENumeric
}

func (e *ExprNumeric) PrettyPrint() string {
	return fmt.Sprintf("%f", e.value)
}

// -----------
// ExprBoolean
// -----------
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

func (e *ExprBoolean) PrettyPrint() string {
	return fmt.Sprintf("%v", e.value)
}

// -----------------
// ExprStringLiteral
// -----------------
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

func (e *ExprStringLiteral) PrettyPrint() string {
	return e.value
}

// --------
// ExprNull
// --------
const ENull ExprType = "ENull"

type ExprNullLiteral struct{}

func (e *ExprNullLiteral) Source() string {
	return lex.ReservedWordNames[lex.TNull]
}

func (e *ExprNullLiteral) Type() ExprType {
	return ENull
}

func (e *ExprNullLiteral) PrettyPrint() string {
	return "null"
}

// -------------
// ExprUndefined
// -------------
const EUndefined ExprType = " EUndefined"

type ExprUndefinedLiteral struct{}

func (e *ExprUndefinedLiteral) Source() string {
	return lex.ReservedWordNames[lex.TUndefined]
}

func (e *ExprUndefinedLiteral) Type() ExprType {
	return EUndefined
}

func (e *ExprUndefinedLiteral) PrettyPrint() string {
	return "undefined"
}

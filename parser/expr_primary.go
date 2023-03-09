package parser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/ruiconti/gojs/lex"
)

// --------------------------
// RootNode: Artificial node
// --------------------------
type ExprRootNode struct {
	children []AstNode
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
		if token.T == t {
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

// -------------
// ExprPrimary
// -------------
func (p *Parser) parsePrimaryExpr(cursor *int) (AstNode, error) {
	reject := false

	// current token position
	// some statement here
	// ˆ
	// p.cursor: 0
	token, err := p.peekN(*cursor)
	if err != nil {
		return nil, err
	}

	p.logger.Debug("[%d:%d] parser:primaryExpr %v", p.cursor, *cursor, token)
	// in primary expressions, we first process the operator
	var primaryExpr AstNode
	switch token.T {
	case lex.TIdentifier:
		primaryExpr = &ExprIdentifierReference{
			reference: token.Lexeme,
		}
	case lex.TNumericLiteral:
		num, err := strconv.ParseFloat(token.Lexeme, 64)
		if err != nil {
			return nil, err
		}
		primaryExpr = &ExprNumeric{
			value: num,
		}
	case lex.TStringLiteral_SingleQuote:
		primaryExpr = &ExprStringLiteral{
			value: token.Lexeme,
		}
	case lex.TStringLiteral_DoubleQuote:
		primaryExpr = &ExprStringLiteral{
			value: token.Lexeme,
		}
	case lex.TTrue:
		primaryExpr = &ExprBoolean{
			value: true,
		}
	case lex.TFalse:
		primaryExpr = &ExprBoolean{
			value: false,
		}
	case lex.TNull:
		primaryExpr = &ExprNullLiteral{}
	case lex.TUndefined:
		primaryExpr = &ExprUndefinedLiteral{}
	default:
		reject = true
		// not implemented yet
	}

	if reject {
		p.logger.Debug("[%d:%d] parser:primaryExpr:rej %v", p.cursor, *cursor, token)
		return nil, errNotPrimaryExpr
	}

	*cursor = *cursor + 1
	p.logger.Debug("[%d:%d] parser:primaryExpr:acc %v", p.cursor, *cursor, primaryExpr.PrettyPrint())
	return primaryExpr, nil
}

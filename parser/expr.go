package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ruiconti/gojs/lex"
)

type Expr interface{}

// /////////////////////////////
// RootNode: Artificial node //
// /////////////////////////////
type ExprRootNode struct {
	children []Node
}

func (e *ExprRootNode) Type() ExprType {
	return ENodeRoot
}

func (e *ExprRootNode) S() string {
	pp := strings.Builder{}
	pp.Write([]byte("(root"))
	for i, child := range e.children {
		if i == 0 {
			// first whitespace
			pp.Write([]byte(" "))
		}

		pp.Write([]byte(child.S()))
		if i < len(e.children)-1 {
			// subsequent whitespace
			pp.Write([]byte(" "))
		}

	}
	pp.Write([]byte(")"))
	return pp.String()
}

// /////////////////
// ExprIdentifier //
// /////////////////
const EIdentifierReference ExprType = "EIdentifierReference"

type ExprIdentifierReference struct {
	reference string
}

func (e *ExprIdentifierReference) Type() ExprType {
	return EIdentifierReference
}

func (e *ExprIdentifierReference) S() string {
	return e.reference
}

// ///////////////
// ExprLiterals //
// ///////////////
const ELiteral ExprType = "ExprLiteral"

var (
	ExprLitNull      = MakeLiteralExpr(lex.TNull)
	ExprLitUndefined = MakeLiteralExpr(lex.TUndefined)
	ExprLitTrue      = MakeLiteralExpr(lex.TTrue)
	ExprLitFalse     = MakeLiteralExpr(lex.TFalse)
)

type Literal interface {
	int | int64 | float64 | string | bool
}

type ExprLiteral[Value Literal] struct {
	typ   lex.TokenType
	value Value
}

func (e *ExprLiteral[Value]) Source() string {
	return fmt.Sprintf("%v", e.value)
}

func (e *ExprLiteral[Value]) Type() ExprType {
	return ELiteral
}

func (e *ExprLiteral[Value]) S() string {
	return fmt.Sprintf("%v", e.value)
}

func MakeLiteralExpr(typ lex.TokenType) *ExprLiteral[string] {
	literal := typ.S()
	if literal == lex.UnknownLiteral {
		return nil
	}
	return &ExprLiteral[string]{typ, literal}
}

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

// //////////////
// ExprUnaryOp //
// //////////////
const EUnaryOp ExprType = "ExprUnaryOp"

var UnaryOperators = []lex.TokenType{
	lex.TDelete,
	lex.TTypeof,
	lex.TVoid,
	lex.TPlus,
	lex.TMinus,
	lex.TBang,
	lex.TTilde,
}

type ExprUnaryOp struct {
	operand  Node
	operator lex.TokenType
}

func (e *ExprUnaryOp) Name() string {
	return "UnaryOpExpression"
}

func (e *ExprUnaryOp) Type() ExprType {
	return EUnaryOp
}

func (e *ExprUnaryOp) S() string {
	return fmt.Sprintf("(%s %s)", e.operator.S(), e.operand.S())
}

// Parser

var UpdateOperators = []lex.TokenType{
	lex.TMinusMinus,
	lex.TPlusPlus,
}

// ///////////////
// ExprBinaryOp //
// ///////////////
const EBinaryOp ExprType = "ExprBinaryOp"

type ExprBinaryOp struct {
	left     Node
	right    Node
	operator lex.TokenType
}

func (e *ExprBinaryOp) Type() ExprType {
	return EBinaryOp
}

func (e *ExprBinaryOp) S() string {
	return fmt.Sprintf("(%s %s %s)", e.operator.S(), e.left.S(), e.right.S())
}

// //////////
// ExprNew //
// //////////
const ENew ExprType = "ExprNew"

type ExprNew struct {
	callee Node
}

func (e *ExprNew) Type() ExprType {
	return ENew
}

func (e *ExprNew) S() string {
	return fmt.Sprintf("(new %s)", e.callee.S())
}

// ///////////////////
// ExprMemberAccess //
// ///////////////////
const EMemberAccess ExprType = "ExprMemberAccess"

type ExprMemberAccess struct {
	object   Node
	property Node
}

func (e *ExprMemberAccess) Type() ExprType {
	return EMemberAccess
}

func (e *ExprMemberAccess) S() string {
	return fmt.Sprintf("(. %s %s)", e.property.S(), e.object.S())
}

//
// Expressions productions
//

func (p *Parser) parseExpr(c *int) (Node, error) {
	return p.parseAssignExpr(c)
}

func (p *Parser) parseAssignExpr(c *int) (Node, error) {
	return p.parseCondExpr(c)
}

func (p *Parser) parseCondExpr(c *int) (Node, error) {
	return p.parseLogOrExpr(c)
}

func (p *Parser) parseBinaryExprGeneric(
	cursor *int,
	operators []lex.TokenType,
	left func(*int) (Node, error),
	right func(*int) (Node, error),
) (Node, error) {
	var opstr strings.Builder
	for _, op := range operators {
		opstr.Write([]byte(op.S()))
	}

	p.log(cursor, "parseBinExprGeneric2 (%s) ENTER", opstr.String())
	c := 0
	// try to parse left side
	expr, err := left(&c)
	if err != nil {
		p.log(cursor, "parseBinExprGeneric2 (%s) left-expr REJ", opstr.String())
		return &ExprBinaryOp{}, err
	}
	p.log(cursor, "parseBinExprGeneric2 (%s) left-expr ACC: %v", opstr.String(), expr.S())

	// advance cursor to parse operator
	p.advanceBy(c)
	if p.isEOF() {
		// we may have reached the right-hand-side of a binary expression
		return expr, nil
	}

	for p.matchAny(operators...) {
		// if it did match, then to parse the operator we need to look back
		// we know it won't err because we just matched it :)
		operator, _ := p.peekN(-1)

		str := operator.Type.S()
		c = 0
		rightArg, err := right(&c)
		if err != nil {
			p.log(cursor, "parseBinExprGeneric2 (loop:%s) right-expr (REJ): %v", str, err)
			// eof or didn't match
			break
		}
		p.log(cursor, "parseBinExprGeneric2 (loop:%s) right-expr (ACC): %v", str, rightArg.S())

		p.advanceBy(c)
		p.log(cursor, "parseBinExprGeneric2 (loop:%s) before: %v", str, expr.S())
		expr = &ExprBinaryOp{
			operator: operator.Type,
			left:     expr,
			right:    rightArg,
		}
		p.log(cursor, "parseBinExprGeneric2 (loop:%s) after: %v", str, expr.S())
	}

	p.log(cursor, "parseBinExprGeneric2 returning (ACC): %v", expr.S())
	return expr, nil
}

func (p *Parser) parseLogOrExpr(c *int) (Node, error) {
	return p.parseBinaryExprGeneric(
		c,
		[]lex.TokenType{lex.TLogicalOr},
		p.parseAndExpr,
		p.parseAndExpr,
	)
}

func (p *Parser) parseAndExpr(c *int) (Node, error) {
	return p.parseBinaryExprGeneric(
		c,
		[]lex.TokenType{lex.TLogicalAnd},
		p.parseBitOrExpr,
		p.parseBitOrExpr,
	)
}

func (p *Parser) parseBitOrExpr(c *int) (Node, error) {
	return p.parseBinaryExprGeneric(
		c,
		[]lex.TokenType{lex.TOr},
		p.parseBitXorExpr,
		p.parseBitXorExpr,
	)
}

func (p *Parser) parseBitXorExpr(c *int) (Node, error) {
	return p.parseBinaryExprGeneric(
		c,
		[]lex.TokenType{lex.TXor},
		p.parseBitAndExpr,
		p.parseBitAndExpr,
	)
}

func (p *Parser) parseBitAndExpr(c *int) (Node, error) {
	return p.parseBinaryExprGeneric(
		c,
		[]lex.TokenType{lex.TAnd},
		p.parseEqualityExpr,
		p.parseEqualityExpr,
	)
}

func (p *Parser) parseEqualityExpr(c *int) (Node, error) {
	return p.parseBinaryExprGeneric(
		c,
		[]lex.TokenType{lex.TEqual, lex.TNotEqual, lex.TStrictEqual, lex.TStrictNotEqual},
		p.parseRelationalExpr,
		p.parseRelationalExpr,
	)
}

func (p *Parser) parseRelationalExpr(c *int) (Node, error) {
	return p.parseBinaryExprGeneric(
		c,
		[]lex.TokenType{lex.TGreaterThan, lex.TGreaterThanEqual, lex.TLessThan, lex.TLessThanEqual, lex.TInstanceof, lex.TIn},
		p.parseShiftExpr,
		p.parseShiftExpr,
	)
}

func (p *Parser) parseShiftExpr(c *int) (Node, error) {
	return p.parseBinaryExprGeneric(
		c,
		[]lex.TokenType{lex.TLeftShift, lex.TRightShift, lex.TUnsignedRightShift},
		p.parseAdditiveExpr,
		p.parseAdditiveExpr,
	)
}

func (p *Parser) parseAdditiveExpr(c *int) (Node, error) {
	return p.parseBinaryExprGeneric(
		c,
		[]lex.TokenType{lex.TPlus, lex.TMinus},
		p.parseMultiplicativeExpr,
		p.parseMultiplicativeExpr,
	)
}

func (p *Parser) parseMultiplicativeExpr(c *int) (Node, error) {
	return p.parseBinaryExprGeneric(
		c,
		[]lex.TokenType{lex.TStar, lex.TSlash, lex.TPercent},
		p.parseExponentialExpr,
		p.parseExponentialExpr,
	)
}

func (p *Parser) parseExponentialExpr(c *int) (Node, error) {
	return p.parseBinaryExprGeneric(
		c,
		[]lex.TokenType{lex.TStarStar},
		p.parseUnaryOperator,
		p.parseUnaryOperator,
	)
}

func (p *Parser) parseUnaryOperator(cursor *int) (Node, error) {
	// UnaryExpr = UpdateExpr | ("delete" | "void" | "typeof" | "+" | "-" | "~" | "!") UnaryExpr
	// parse (UpdateExpression)
	p.log(cursor, "parseUnaryExpr ENTER")
	p.log(cursor, "parseUnaryExpr left-expr (updateExpr)")
	updateExpr, err := p.parseUpdateExpr(cursor)
	if err == nil {
		p.log(cursor, "parseUnaryExpr left-expr (updateExpr) returning (ACC): %v", updateExpr.S())
		return updateExpr, nil
	}

	var expr *ExprUnaryOp
	for p.matchAny(UnaryOperators...) {
		// parse (UnaryOperator)
		operator, _ := p.peekN(-1)
		str := operator.Type.S()

		p.log(cursor, "parseUnaryExpr (loop:%v) evaluating operand...", str)
		*cursor = 0
		operand, err := p.parseUnaryOperator(cursor)
		if err != nil {
			p.log(cursor, "parseUnaryExpr (loop:%v) operand eval rejected: %v", str, err)
			return nil, err
			// break
		}
		p.log(cursor, "parseUnaryExpr (loop:%v) operand eval accepted: %v", str, operand)

		// found a valid operand, walk it
		p.advanceBy(*cursor)
		expr = &ExprUnaryOp{
			operator: operator.Type,
			operand:  operand,
		}
		p.log(cursor, "parseUnaryExpr (loop:%v) expr: %v", str, expr.S())
	}

	if expr == nil {
		p.log(cursor, "parseUnaryExpr returning (REJ): expr nil")
		return expr, fmt.Errorf("parseUnaryExpr rejected")
	}

	p.log(cursor, "parseUnaryExpr returning (ACC): %v", expr.S())
	return expr, nil
}

func (p *Parser) parseUpdateExpr(cursor *int) (Node, error) {
	p.log(cursor, "parseUpdateExpr ENTER")
	var expr Node

	// parse (LeftHandSideExpression)
	leftExpr, err := p.parseLeftHandSideExpr(cursor)
	if err == nil {
		// TODO: this early bail may not parse valid update expressions
		return leftExpr, nil
	}

	for p.matchAny(UpdateOperators...) {
		// parse (UnaryOperator)
		operator, _ := p.peekN(-1)
		str := operator.Type.S()

		p.log(cursor, "parseUpdateExpr (loop:%v) evaluating operand...", str)
		*cursor = 0
		operand, err := p.parseUnaryOperator(cursor)
		if err != nil {
			p.log(cursor, "parseUpdateExpr (loop:%v) operand eval rejected: %v", str, err)
			return nil, err
			// break
		}
		p.log(cursor, "parseUpdateExpr (loop:%v) operand eval accepted: %v", str, operand)

		// found a valid operand, walk it
		p.advanceBy(*cursor)
		expr = &ExprUnaryOp{
			operator: operator.Type,
			operand:  operand,
		}
		p.log(cursor, "parseUpdateExpr (loop:%v) expr: %v", str, expr.S())
	}

	if expr == nil {
		p.log(cursor, "parseUpdateExpr returning (REJ): expr nil")
		return nil, fmt.Errorf("parseUpdateExpr rejected")
	}

	p.log(cursor, "parseUpdateExpr returning (ACC): %v", expr.S())
	return expr, nil
}

func (p *Parser) parseLeftHandSideExpr(cursor *int) (Node, error) {
	// parse (NewExpression)
	// NewExpression = MemberExpression | "new" NewExpression
	memberExpr, err := p.parseMemberExpr(cursor)
	if err == nil {
		return memberExpr, nil
	}

	var expr *ExprNew
	for p.matchAny(lex.TNew) {
		p.log(cursor, "parseNewExpr (loop) evaluating callee...")
		*cursor = 0
		callee, err := p.parseMemberExpr(cursor)
		if err != nil {
			p.log(cursor, "parseNewExpr (loop) callee eval rejected: %v", err)
			return nil, err
			// break
		}
		p.log(cursor, "parseNewExpr (loop) callee eval accepted: %v", callee)

		// found a valid operand, walk it
		p.advanceBy(*cursor)
		expr = &ExprNew{
			callee: callee,
		}
		p.log(cursor, "parseNewExpr (loop) expr: %v", expr.S())
	}

	if expr == nil {
		p.log(cursor, "parseNewExpr returning (REJ): expr nil")
		return nil, fmt.Errorf("parseNewExpr rejected")
	}

	p.log(cursor, "parseNewExpr returning (ACC): %v", expr.S())
	return expr, nil
}

func (p *Parser) parseMemberExpr(cursor *int) (Node, error) {
	p.log(cursor, "parseMemberExpr ENTER")
	primaryExpr, err := p.parsePrimaryExpr(cursor)
	if err != nil {
		p.log(cursor, "parseMemberExpr REJ: %v", err)
		return nil, err
		// try to parse SuperCall or SuperProperty
	}

	exprMember := primaryExpr
	parsed := false
	p.advanceBy(1)
	for p.matchAny(lex.TPeriod, lex.TLeftBracket) {
		current, err := p.peekN(-1)
		p.log(cursor, "parseMemberExpr (loop) matched %v", current.Type)
		if err != nil {
			p.log(cursor, "parseMemberExpr (loop) leaving early:%v", err)
			break
			// todo: error
		}

		switch current.Type {
		case lex.TPeriod:
			parsed = true
			p.log(cursor, "parseMemberExpr (loop) matched dot")
			if p.matchAny(lex.TIdentifier) {
				identifier, _ := p.peekN(-1)
				exprMember = &ExprMemberAccess{
					object: exprMember,
					property: &ExprIdentifierReference{
						reference: identifier.Lexeme,
					},
				}
				p.log(cursor, "parseMemberExpr (loop:ACC) %v", exprMember.S())
			} else {
				p.log(cursor, "parseMemberExpr (loop:REJ) parsed dot but failed to find identifier")
				break
			}
		case lex.TLeftBracket:
			parsed = true
			p.log(cursor, "parseMemberExpr (loop) matched left brace")
			expr, err := p.parseExpr(cursor)
			if err == nil {
				if p.matchAny(lex.TRightBracket) {
					exprMember = &ExprMemberAccess{
						object:   exprMember,
						property: expr,
					}
					p.log(cursor, "parseMemberExpr (loop:ACC) %v", exprMember.S())
				}
			} else {
				p.log(cursor, "parseMemberExpr (loop) leaving..")
				break
			}

		}
	}

	p.log(cursor, "parseMemberExpr ACC: %v", primaryExpr.S())
	if !parsed {
		p.advanceBy(-1)
	}
	// it's just a primary expression
	return exprMember, nil
}

// -------------
// ExprPrimary
// -------------
func (p *Parser) parsePrimaryExpr(cursor *int) (Node, error) {
	reject := false

	// current token position
	// some statement here
	// ˆ
	// p.cursor: 0
	token, err := p.peekN(*cursor)
	if err != nil {
		return nil, err
	}

	p.log(cursor, "primaryExpr ENTER: %v", token)
	// in primary expressions, we first process the operator
	var primaryExpr Node
	typ := token.Type
	switch typ {
	case lex.TIdentifier:
		primaryExpr = &ExprIdentifierReference{
			reference: token.Lexeme,
		}
	case lex.TNumericLiteral:
		num, err := strconv.ParseFloat(token.Lexeme, 64)
		if err != nil {
			return nil, err
		}
		primaryExpr = &ExprLiteral[float64]{
			value: num,
			typ:   typ,
		}
	case lex.TStringLiteral_SingleQuote:
		primaryExpr = &ExprLiteral[string]{
			value: token.Lexeme,
			typ:   typ,
		}
	case lex.TStringLiteral_DoubleQuote:
		primaryExpr = &ExprLiteral[string]{
			value: token.Lexeme,
			typ:   typ,
		}
	case lex.TTrue:
		primaryExpr = ExprLitTrue
	case lex.TFalse:
		primaryExpr = ExprLitFalse
	case lex.TNull:
		primaryExpr = ExprLitNull
	case lex.TUndefined:
		primaryExpr = ExprLitUndefined
	default:
		reject = true
		// not implemented yet
	}

	if reject {
		p.log(cursor, "primaryExpr returning (REJ): %v", token)
		return nil, fmt.Errorf("primaryExpr rejected")
	}

	*cursor = *cursor + 1
	p.log(cursor, "primaryExpr returning (ACC): %v", primaryExpr.S())
	return primaryExpr, nil
}

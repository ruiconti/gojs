package parser

import (
	"fmt"
	"strconv"
	"strings"

	l "github.com/ruiconti/gojs/lexer"
)

type ExprType string

const (
	ENodeRoot ExprType = "ENodeRoot"
	EElision  ExprType = "EElision"
)

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
	var s strings.Builder
	s.Write([]byte("(root"))
	for i, child := range e.children {
		if i == 0 {
			// first whitespace
			s.Write([]byte(" "))
		}

		if child == nil {
			panic(fmt.Sprintf("panic: nil child at index %d/%d. children: %+v", i, len(e.children)-1, s.String()))
		}

		s.Write([]byte(child.S()))
		if i < len(e.children)-1 {
			// subsequent whitespace
			s.Write([]byte(" "))
		}

	}
	s.Write([]byte(")"))
	return s.String()
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
	ExprLitNull      = MakeLiteralExpr(l.TNull)
	ExprLitUndefined = MakeLiteralExpr(l.TUndefined)
	ExprLitTrue      = MakeLiteralExpr(l.TTrue)
	ExprLitFalse     = MakeLiteralExpr(l.TFalse)
)

type Literal interface {
	int | int64 | float64 | string | bool
}

type ExprLiteral[Value Literal] struct {
	tok l.Token
}

func (e *ExprLiteral[Value]) Source() string {
	return fmt.Sprintf("%v", e.tok.Literal)
}

func (e *ExprLiteral[Value]) Type() ExprType {
	return ELiteral
}

func (e *ExprLiteral[Value]) S() string {
	return fmt.Sprintf("%v", e.tok.Literal)
}

func MakeLiteralExpr(typ l.TokenType) *ExprLiteral[string] {
	tok := l.Token{Type: typ, Literal: typ.S()}
	literal := tok.Type.S()
	if literal == l.UnknownLiteral {
		return nil
	}
	return &ExprLiteral[string]{tok}
}

var LiteralsTokens = []l.TokenType{
	l.TNumericLiteral,
	l.TRegularExpressionLiteral,
	l.TStringLiteral_SingleQuote,
	l.TStringLiteral_DoubleQuote,
	l.TTrue,
	l.TFalse,
	l.TNull,
}

func isLiteralToken(token l.Token) bool {
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

var UnaryOperators = []l.TokenType{
	l.TDelete,
	l.TTypeof,
	l.TVoid,
	l.TPlus,
	l.TMinus,
	l.TBang,
	l.TTilde,
}

type ExprUnaryOp struct {
	operand  Node
	operator l.Token
}

func (e *ExprUnaryOp) Name() string {
	return "UnaryOpExpression"
}

func (e *ExprUnaryOp) Type() ExprType {
	return EUnaryOp
}

func (e *ExprUnaryOp) S() string {
	return fmt.Sprintf("(%s %s)", e.operator.Type.S(), e.operand.S())
}

// Parser

var UpdateOperators = []l.TokenType{
	l.TMinusMinus,
	l.TPlusPlus,
}

// ///////////////
// ExprBinaryOp //
// ///////////////
const EBinaryOp ExprType = "ExprBinaryOp"

type ExprBinaryOp struct {
	left     Node
	right    Node
	operator l.Token
}

func (e *ExprBinaryOp) Type() ExprType {
	return EBinaryOp
}

func (e *ExprBinaryOp) S() string {
	return fmt.Sprintf("(%s %s %s)", e.operator.Type.S(), e.left.S(), e.right.S())
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
	if e == nil {
		panic("invalid object: nil")
	}
	if e.property == nil {
		panic("invalid nil access: property")
	}
	if e.object == nil {
		panic("invalid nil access: object")
	}
	return fmt.Sprintf("(. %s %s)", e.property.S(), e.object.S())
}

// //////////////////////////
// Expressions productions //
// //////////////////////////
func (p *Parser) parseExpr() (Node, error) {
	return p.parseAssignExpr()
}

func (p *Parser) parseAssignExpr() (Node, error) {
	return p.parseCondExpr()
}

func (p *Parser) parseCondExpr() (Node, error) {
	return p.parseLogOrExpr()
}

func newSet[C comparable](items ...C) map[C]struct{} {
	set := make(map[C]struct{}, len(items))
	for _, item := range items {
		set[item] = struct{}{}
	}
	return set
}

// parseBinaryOperators parses a binary expression with the following EBNF:
//
//	Expr ::= Expr operator RightHigherExpr | LeftHigherExpr
//
// which can be transformed, removing left-recursion, to:
//
//	Expr ::= LeftHigherExpr (operator RightHigherExpr)*
//
// for the more generic case where LeftHigherExpr == RightHigherExpr:
//
//	Expr ::= HigherExpr (operator HigherExpr)*
func (p *Parser) parseBinaryOperators(
	operators []l.TokenType,
	higherExprLeft func() (Node, error),
	higherExprRight func() (Node, error),
) (Node, error) {
	var (
		left  Node
		err   error
		opSet = newSet(operators...)
	)

	// Expr ::= HigherExpr
	if left, err = higherExprLeft(); err != nil {
		return nil, err
	}

	// (operator HigherExpr)*
	lastCursor := p.cursor
	for {
		token := p.Peek()
		p.Log("parseBinaryOperators: %s", token.Type.S())
		if _, ok := opSet[token.Type]; !ok {
			break
		}
		p.Next() // Consumes operator

		if right, err := higherExprRight(); err != nil {
			break
		} else {
			left = &ExprBinaryOp{
				operator: token,
				left:     left,
				right:    right,
			}
		}
		p.guardInfiniteLoop(&lastCursor)
	}

	return left, nil
}

func (p *Parser) parseLogOrExpr() (Node, error) {
	p.Log("parseLogOrExpr")
	return p.parseBinaryOperators(
		[]l.TokenType{l.TLogicalOr},
		p.parseAndExpr,
		p.parseAndExpr,
	)
}

func (p *Parser) parseAndExpr() (Node, error) {
	p.Log("parseAndExpr")
	return p.parseBinaryOperators(
		[]l.TokenType{l.TLogicalAnd},
		p.parseBitOrExpr,
		p.parseBitOrExpr,
	)
}

func (p *Parser) parseBitOrExpr() (Node, error) {
	p.Log("parseBitOrExpr")
	return p.parseBinaryOperators(
		[]l.TokenType{l.TOr},
		p.parseBitXorExpr,
		p.parseBitXorExpr,
	)
}

func (p *Parser) parseBitXorExpr() (Node, error) {
	p.Log("parseBitXorExpr")
	return p.parseBinaryOperators(
		[]l.TokenType{l.TXor},
		p.parseBitAndExpr,
		p.parseBitAndExpr,
	)
}

func (p *Parser) parseBitAndExpr() (Node, error) {
	p.Log("parseBitAndExpr")
	return p.parseBinaryOperators(
		[]l.TokenType{l.TAnd},
		p.parseEqualityExpr,
		p.parseEqualityExpr,
	)
}

func (p *Parser) parseEqualityExpr() (Node, error) {
	p.Log("parseEqualityExpr")
	return p.parseBinaryOperators(
		[]l.TokenType{l.TEqual, l.TNotEqual, l.TStrictEqual, l.TStrictNotEqual},
		p.parseRelationalExpr,
		p.parseRelationalExpr,
	)
}

func (p *Parser) parseRelationalExpr() (Node, error) {
	p.Log("parseRelationalExpr")
	return p.parseBinaryOperators(
		[]l.TokenType{l.TGreaterThan, l.TGreaterThanEqual, l.TLessThan, l.TLessThanEqual, l.TInstanceof, l.TIn},
		p.parseShiftExpr,
		p.parseShiftExpr,
	)
}

func (p *Parser) parseShiftExpr() (Node, error) {
	p.Log("parseShiftExpr")
	return p.parseBinaryOperators(
		[]l.TokenType{l.TLeftShift, l.TRightShift, l.TUnsignedRightShift},
		p.parseAdditiveExpr,
		p.parseAdditiveExpr,
	)
}

func (p *Parser) parseAdditiveExpr() (Node, error) {
	p.Log("parseAdditiveExpr")
	return p.parseBinaryOperators(
		[]l.TokenType{l.TPlus, l.TMinus},
		p.parseMultiplicativeExpr,
		p.parseMultiplicativeExpr,
	)
}

func (p *Parser) parseMultiplicativeExpr() (Node, error) {
	p.Log("parseMultiplicativeExpr")
	return p.parseBinaryOperators(
		[]l.TokenType{l.TStar, l.TSlash, l.TPercent},
		p.parseExponentialExpr,
		p.parseExponentialExpr,
	)
}

func (p *Parser) parseExponentialExpr() (Node, error) {
	p.Log("parseExponentialExpr")
	return p.parseBinaryOperators(
		[]l.TokenType{l.TStarStar},
		p.parseUnaryOperator,
		p.parseUpdateExpr,
	)
}

// UnaryExpression ::=
// | UpdateExpression
// | UnaryOp UnaryExpression
// | AwaitExpression (TODO)
//
// UnaryOp ::= delete | void | typeof | + | - | ~ | !
func (p *Parser) parseUnaryOperator() (Node, error) {
	p.Log("parseUnaryOperator")
	var (
		exprUnary, exprUpdate Node
		err                   error
		unaryOpSet            = newSet(UnaryOperators...)
	)

	// UnaryExpression ::= UpdateExpression
	if exprUpdate, err = p.parseUpdateExpr(); err == nil {
		return exprUpdate, nil
	}

	// UnaryExpression ::= UnaryOp UnaryExpression
	lastCursor := p.cursor
	for {
		token := p.Peek()
		if _, ok := unaryOpSet[token.Type]; !ok {
			break
		}

		p.Next() // consume operator
		exprUnary, err = p.parseUnaryOperator()
		if err != nil {
			return nil, err
		}

		exprUnary = &ExprUnaryOp{
			operator: token,
			operand:  exprUnary,
		}
		p.guardInfiniteLoop(&lastCursor)
	}

	return exprUnary, nil
}

// UpdateExpression ::=
// | LeftHandSideExpression (++ | --)
// | (++ | --) UnaryExpression
func (p *Parser) parseUpdateExpr() (Node, error) {
	p.Log("parseUpdateExpr")
	var (
		exprLeft, exprUnary Node
		err                 error
		unaryOpSet          = newSet(UpdateOperators...)
		match               bool
	)

	// LeftHandSideExpression (++ | --)*
	exprLeft, err = p.parseLeftHandSideExpr()
	if err == nil {
		// (++ | --)
		token := p.Peek()
		if _, ok := unaryOpSet[token.Type]; !ok {
			return exprLeft, nil
		} else {
			return &ExprUnaryOp{ // TODO: make an UpdateExpr
				operand:  exprLeft,
				operator: token,
			}, nil
		}
	}

	// (++ | --) UnaryExpression
	lastCursor := p.cursor
	for {
		token := p.Peek()
		if _, ok := unaryOpSet[token.Type]; !ok {
			break
		}

		p.Next() // consume operator
		operand, err := p.parseUnaryOperator()
		if err != nil {
			return nil, err
		}

		match = true
		exprUnary = &ExprUnaryOp{
			operator: token,
			operand:  operand,
		}
		p.guardInfiniteLoop(&lastCursor)
	}

	if match {
		return exprUnary, nil
	} else {
		return nil, err
	}
}

// parseLeftHandSideExpr parses the following grammar:
//
// LeftHandSideExpression ::=
// | NewExpression
// | CallExpression     (TODO)
// | OptionalExpression (TODO)
func (p *Parser) parseLeftHandSideExpr() (Node, error) {
	p.Log("parseLeftHandSideExpr")
	var (
		expr Node
		err  error
	)
	if expr, err = p.parseNewExpr(); err == nil {
		return expr, nil
	}

	// reset state
	// TODO: parseCallExpr

	// reset state
	// TODO: parseOptionalExpression
	return nil, fmt.Errorf("parseLeftHandSideExpr rejected")
}

// parseLeftHandSideExpr parses the following grammar:
//
// NewExpression ::= MemberExpression | new NewExpression
func (p *Parser) parseNewExpr() (Node, error) {
	p.Log("parseNewExpr")
	var (
		expr Node
		err  error
	)

	token := p.Peek()
	if token.Type == l.TNew {
		// NewExpression ::= new NewExpression
		p.Next() // consumes new
		callee, err := p.parseNewExpr()
		if err != nil {
			return nil, err
		}

		return &ExprNew{
			callee: callee,
		}, nil
	}

	// NewExpression ::= MemberExpression
	if expr, err = p.parseMemberExpr(); err == nil {
		return expr, nil
	}
	return nil, err
}

// parseMemberExpr parses the following grammar:
//
// MemberExpression ::=
// | PrimaryExpression
// | MemberExpression | (. IdentifierName | [ Expression ] | TemplateLiteral | . PrivateIdentifier)
// | SuperProperty                  (TODO)
// | MetaProperty  							    (TODO)
// | new MemberExpression Arguments (TODO)
//
// Transforming to remove the left-recursion:
//
// MemberExpression ::=
// | PrimaryExpression
// | MemberExpression' ([ Expr ] MemberExpression')*
// | MemberExpression' (. IdentifierName MemberExpression')*
// | MemberExpression' (TemplateLiteral MemberExpression')* (TODO)
//
// MemberExpression' ::=
// | PrimaryExpression | SuperProperty | MetaProperty | new MemberExpression Arguments
func (p *Parser) parseMemberExpr() (Node, error) {
	p.Log("parseMemberExpr")
	var (
		exprMember Node
		err        error
	)

	// MemberExpression ::= MemberExpression'
	if exprMember, err = p.parsePrimaryExpr(); err != nil {
		return nil, err
	}

	// MemberExpression ::=
	// | MemberExpression' ([ Expr ] MemberExpression')*
	// | MemberExpression' (. IdentifierName MemberExpression')*
	lastCursor := p.cursor
loop:
	for {
		token := p.Peek()
		switch token.Type {
		case l.TPeriod:
			// MemberExpression ::= (. IdentifierName MemberExpression')*
			p.Next() // consume period
			afterPeriod := p.Peek()
			if afterPeriod.Type == l.TIdentifier {
				// match = true
				exprMember = &ExprMemberAccess{
					object: exprMember,
					property: &ExprIdentifierReference{
						reference: afterPeriod.Lexeme,
					},
				}
				p.Next() // consume identifier
			} else {
				return nil, fmt.Errorf("expected identifier after dot")
			}
		case l.TLeftBracket:
			// MemberExpression ::= ([ Expr ] MemberExpression')*
			p.Next()                                    // consume left bracket
			if expr, err := p.parseExpr(); err == nil { // parseExpr consumes the expression's tokens
				if p.Peek().Type == l.TRightBracket {
					// match = true
					exprMember = &ExprMemberAccess{
						object:   exprMember,
						property: expr,
					}
					p.Next() // consume right bracket
				}
			} else {
				return nil, fmt.Errorf("expected valid expression after left bracket")
			}
		default:
			break loop
		}
		p.guardInfiniteLoop(&lastCursor)
	}

	switch exprMember.(type) {
	case *ExprMemberAccess:
		return exprMember, nil
		// panic(t.S())
	default:
		return exprMember, nil
	}
	// return exprMember, nil
}

// PrimaryExpression ::=
// | this (TODO)
// | IdentifierReference
// | Literal
// | ArrayLiteral (TODO)
// | ObjectLiteral (TODO)
// | FunctionExpression (TODO)
// | ClassExpression (TODO)
// | GeneratorExpression (TODO)
// | AsyncFunctionExpression (TODO)
// | AsyncGeneratorExpression (TODO)
// | RegularExpressionLiteral (TODO)
// | TemplateLiteral (TODO)
// | CoverParenthesizedExpressionAndArrowParameterList (TODO)
func (p *Parser) parsePrimaryExpr() (Node, error) {
	p.Log("parsePrimaryExpr")
	var primaryExpr Node
	token := p.Peek()

	switch token.Type {
	case l.TIdentifier:
		primaryExpr = &ExprIdentifierReference{
			reference: token.Lexeme,
		}
	case l.TNumericLiteral:
		if num, err := strconv.ParseFloat(token.Lexeme, 64); err == nil {
			primaryExpr = &ExprLiteral[float64]{
				tok: l.Token{
					Type:    l.TNumericLiteral,
					Literal: num,
					Lexeme:  token.Lexeme,
				},
			}
		} else {
			return nil, err
		}
	case l.TStringLiteral_SingleQuote:
		primaryExpr = &ExprLiteral[string]{token}
	case l.TStringLiteral_DoubleQuote:
		primaryExpr = &ExprLiteral[string]{token}
	case l.TTrue:
		primaryExpr = ExprLitTrue
	case l.TFalse:
		primaryExpr = ExprLitFalse
	case l.TNull:
		primaryExpr = ExprLitNull
	case l.TUndefined:
		primaryExpr = ExprLitUndefined
	default:
		return nil, fmt.Errorf("primaryExpr rejected")
	}

	p.Next() // consume token
	return primaryExpr, nil
}

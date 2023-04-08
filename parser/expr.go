package parser

import (
	"fmt"
	"strconv"
	"strings"

	l "github.com/ruiconti/gojs/lexer"
)

type ExprType string

const (
	EElision ExprType = "EElision"
)

// /////////////////
// ExprIdentifier //
// /////////////////
const EIdentifier ExprType = "EIdentifier"

type ExprIdentifier struct {
	name string
}

func (e *ExprIdentifier) Type() ExprType {
	return EIdentifier
}

func (e *ExprIdentifier) S() string {
	return e.name
}

// /////////////////
// ExprPrivateIdentifier //
// /////////////////
const EPrivateIdentifierReference ExprType = "EPrivateIdentifierReference"

type ExprPrivateIdentifier struct {
	name string
}

func (e *ExprPrivateIdentifier) Type() ExprType {
	return EPrivateIdentifierReference
}

func (e *ExprPrivateIdentifier) S() string {
	return fmt.Sprintf("#%s", e.name)
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
	ExprLitThis      = MakeLiteralExpr(l.TThis)
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
	operand  Expr
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
	left     Expr
	right    Expr
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
	callee    Expr
	arguments []Expr
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
	object   Expr
	property Expr
	optional bool
}

func (e *ExprMemberAccess) Type() ExprType {
	return EMemberAccess
}

func (e *ExprMemberAccess) S() string {
	if e == nil {
		panic("invalid object: nil")
	}
	var sarg string
	if e.optional {
		sarg = "get?"
	} else {
		sarg = "get"
	}
	return fmt.Sprintf("(%s '%s %s)", sarg, e.property.S(), e.object.S())
}

// ///////////////////
// ExprMetaProperty //
// ///////////////////
const EMetaProperty ExprType = "ExprMetaProperty"

type ExprMetaProperty struct {
	meta     Expr
	property Expr
}

func (e *ExprMetaProperty) Type() ExprType {
	return EMetaProperty
}

func (e *ExprMetaProperty) S() string {
	if e == nil {
		panic("invalid object: nil")
	}
	return fmt.Sprintf("(getmeta %s %s)", e.meta.S(), e.property.S())
}

// ///////////////////
// ExprCall //
// ///////////////////
const ECall ExprType = "ExprCall"

type ExprCall struct {
	callee    Expr
	arguments []Expr
	optional  bool
}

func (e *ExprCall) Type() ExprType {
	return ECall
}

func (e *ExprCall) S() string {
	if e == nil {
		panic("invalid object: nil")
	}
	var args strings.Builder
	for i, arg := range e.arguments {
		args.WriteString(arg.S())
		if i < len(e.arguments)-1 {
			args.WriteString(" ")
		}
	}
	var sarg string
	if e.optional {
		sarg = "λ?"
	} else {
		sarg = "λ"
	}
	return fmt.Sprintf("(%s%s %s)", sarg, e.callee.S(), args.String())
}

// ///////////////////
// SpreadElement //
// ///////////////////
const NSpreadElement ExprType = "SpreadElement"

type SpreadElement struct {
	argument Expr
}

func (e *SpreadElement) Type() ExprType {
	return ECall
}

func (e *SpreadElement) S() string {
	if e == nil {
		panic("invalid object: nil")
	}
	return fmt.Sprintf("(... %s)", e.argument.S())
}

// ///////////////////
// ExprImportCall //
// ///////////////////
const EImportCall ExprType = "ExprImportCall"

type ExprImportCall struct {
	source Expr
}

func (e *ExprImportCall) Type() ExprType {
	return EImportCall
}

func (e *ExprImportCall) S() string {
	if e == nil {
		panic("invalid object: nil")
	}
	return fmt.Sprintf("(import %s)", e.source.S())
}

// /////////////
// ExprAssign //
// /////////////
const EAssign ExprType = "ExprAssign"

type ExprAssign struct {
	operator l.Token
	left     Node
	right    Node
}

func (e *ExprAssign) Type() ExprType {
	return EImportCall
}

func (e *ExprAssign) S() string {
	if e == nil {
		panic("invalid object: nil")
	}
	return fmt.Sprintf("(%s %s <- %s)", e.operator.Type.S(), e.left.S(), e.right.S())
}

// //////////////////////////
// Expressions productions //
// //////////////////////////
func (p *Parser) parseExpr() (Expr, error) {
	return p.parseAssignExpr()
}

// AssignmentExpression :
// | ConditionalExpression
// | [+Yield] YieldExpression
// | ArrowFunction
// | AsyncArrowFunction
// | LeftHandSideExpression '=' AssignmentExpression
// | LeftHandSideExpression AssignmentOperator AssignmentExpression
// | LeftHandSideExpression &&= AssignmentExpression
// | LeftHandSideExpression ||= AssignmentExpression
// | LeftHandSideExpression ??= AssignmentExpression
var assignmentOperators = []l.TokenType{
	l.TAssign,
	l.TPlusAssign,
	l.TMinusAssign,
	l.TSlashAssign,
	l.TStarAssign,
	l.TPercentAssign,
	l.TAndAssign,
	l.TOrAssign,
	l.TLogicalAndAssign,
	l.TLogicalOrAssign,
	l.TXorAssign,
	l.TLeftShiftAssign,
	l.TRightShiftAssign,
	l.TUnsignedRightShiftAssign,
}

func (p *Parser) consumeAssignOp() (*l.Token, error) {
	cur := p.Peek()
	found := false
	for _, op := range assignmentOperators {
		if cur.Type == op {
			p.Next()
			found = true
		}
	}

	if !found {
		return nil, fmt.Errorf("did not find a valid assignment operator")
	}
	return &cur, nil
}

func (p *Parser) parseAssignExpr() (Expr, error) {
	var err error

	cp := p.saveCheckpoint()
	// AssignmentExpression : LeftHandSideExpression '=' AssignmentExpression
	parseLhs := func() (Node, error) {
		lhs, err := p.parseLeftHandSideExpr()
		if err != nil {
			return nil, err
		}

		assignOp, err := p.consumeAssignOp()
		if err != nil {
			return nil, err
		}

		rhs, err := p.parseAssignExpr()
		if err != nil {
			return nil, err
		}
		return &ExprAssign{
			left:     lhs,
			right:    rhs,
			operator: *assignOp,
		}, nil
	}
	if assignExpr, err := parseLhs(); err == nil {
		return assignExpr, nil
	}

	p.restoreCheckpoint(cp)
	// AssignmentExpression : ConditionalExpression
	exprCond, err := p.parseCondExpr()
	if err == nil {
		return exprCond, nil
	}

	return nil, fmt.Errorf("rejected on AssignExpr: %s", err.Error())
}

func (p *Parser) parseCondExpr() (Expr, error) {
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
	higherExprLeft func() (Expr, error),
	higherExprRight func() (Expr, error),
) (Expr, error) {
	var (
		left  Expr
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

func (p *Parser) parseLogOrExpr() (Expr, error) {
	p.Log("parseLogOrExpr")
	return p.parseBinaryOperators(
		[]l.TokenType{l.TLogicalOr},
		p.parseAndExpr,
		p.parseAndExpr,
	)
}

func (p *Parser) parseAndExpr() (Expr, error) {
	p.Log("parseAndExpr")
	return p.parseBinaryOperators(
		[]l.TokenType{l.TLogicalAnd},
		p.parseBitOrExpr,
		p.parseBitOrExpr,
	)
}

func (p *Parser) parseBitOrExpr() (Expr, error) {
	p.Log("parseBitOrExpr")
	return p.parseBinaryOperators(
		[]l.TokenType{l.TOr},
		p.parseBitXorExpr,
		p.parseBitXorExpr,
	)
}

func (p *Parser) parseBitXorExpr() (Expr, error) {
	p.Log("parseBitXorExpr")
	return p.parseBinaryOperators(
		[]l.TokenType{l.TXor},
		p.parseBitAndExpr,
		p.parseBitAndExpr,
	)
}

func (p *Parser) parseBitAndExpr() (Expr, error) {
	p.Log("parseBitAndExpr")
	return p.parseBinaryOperators(
		[]l.TokenType{l.TAnd},
		p.parseEqualityExpr,
		p.parseEqualityExpr,
	)
}

func (p *Parser) parseEqualityExpr() (Expr, error) {
	p.Log("parseEqualityExpr")
	return p.parseBinaryOperators(
		[]l.TokenType{l.TEqual, l.TNotEqual, l.TStrictEqual, l.TStrictNotEqual},
		p.parseRelationalExpr,
		p.parseRelationalExpr,
	)
}

func (p *Parser) parseRelationalExpr() (Expr, error) {
	p.Log("parseRelationalExpr")
	return p.parseBinaryOperators(
		[]l.TokenType{l.TGreaterThan, l.TGreaterThanEqual, l.TLessThan, l.TLessThanEqual, l.TInstanceof, l.TIn},
		p.parseShiftExpr,
		p.parseShiftExpr,
	)
}

func (p *Parser) parseShiftExpr() (Expr, error) {
	p.Log("parseShiftExpr")
	return p.parseBinaryOperators(
		[]l.TokenType{l.TLeftShift, l.TRightShift, l.TUnsignedRightShift},
		p.parseAdditiveExpr,
		p.parseAdditiveExpr,
	)
}

func (p *Parser) parseAdditiveExpr() (Expr, error) {
	p.Log("parseAdditiveExpr")
	return p.parseBinaryOperators(
		[]l.TokenType{l.TPlus, l.TMinus},
		p.parseMultiplicativeExpr,
		p.parseMultiplicativeExpr,
	)
}

func (p *Parser) parseMultiplicativeExpr() (Expr, error) {
	p.Log("parseMultiplicativeExpr")
	return p.parseBinaryOperators(
		[]l.TokenType{l.TStar, l.TSlash, l.TPercent},
		p.parseExponentialExpr,
		p.parseExponentialExpr,
	)
}

func (p *Parser) parseExponentialExpr() (Expr, error) {
	p.Log("parseExponentialExpr")
	return p.parseBinaryOperators(
		[]l.TokenType{l.TStarStar},
		p.parseUnaryOperator,
		p.parseUpdateExpr,
	)
}

// UnaryExpression ::=
// | UnaryOp UnaryExpression
// | UpdateExpression
// | AwaitExpression (TODO)
//
// UnaryOp ::= delete | void | typeof | + | - | ~ | !
func (p *Parser) parseUnaryOperator() (Expr, error) {
	p.Log("parseUnaryOperator")
	var (
		exprUnary, exprUpdate Expr
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

	if err != nil {
		return nil, err
	}
	return exprUnary, nil
}

// UpdateExpression ::=
// | LeftHandSideExpression (++ | --)?
// | (++ | --) UnaryExpression
func (p *Parser) parseUpdateExpr() (Expr, error) {
	p.Log("parseUpdateExpr")
	var (
		exprUpdate Expr
		err        error
		unaryOpSet = newSet(UpdateOperators...)
		match      bool
	)

	// LeftHandSideExpression (++ | --)?
	exprUpdate, err = p.parseLeftHandSideExpr()
	if err == nil {
		token := p.Peek()
		if _, ok := unaryOpSet[token.Type]; ok {
			// UpdateExpression ::= LeftHandSideExpression (++ | --)
			return &ExprUnaryOp{ // TODO: make an UpdateExpr
				operand:  exprUpdate,
				operator: token,
			}, nil
		} else {
			// UpdateExpression ::= LeftHandSideExpression
			return exprUpdate, nil
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
		exprUpdate = &ExprUnaryOp{
			operator: token,
			operand:  operand,
		}
		p.guardInfiniteLoop(&lastCursor)
	}

	// TODO: if exprUpdate != nil {
	if match {
		return exprUpdate, nil
	} else {
		return nil, fmt.Errorf("parseUpdateExpr rejected")
	}
}

// parseLeftHandSideExpr parses the following grammar:
//
// LeftHandSideExpression ::=
// NewExpression
// | CallExpression
// | OptionalExpression (embedded)
func (p *Parser) parseLeftHandSideExpr() (Expr, error) {
	p.Log("parseLeftHandSideExpr")
	var (
		cp   uint32
		expr Expr
		err  error
	)

	cp = p.saveCheckpoint()
	if expr, err = p.parseCallExpr(); err == nil {
		return expr, nil
	}
	p.restoreCheckpoint(cp)

	cp = p.saveCheckpoint()
	if expr, err = p.parseNewExpr(); err == nil {
		return expr, nil
	}
	p.restoreCheckpoint(cp)

	return nil, fmt.Errorf("parseLeftHandSideExpr rejected")
}

// OptionalExpression ::=
// | MemberExpression OptionalChain
// | CallExpression OptionalChain
// | OptionalExpression OptionalChain
//
// simplifying the terms
//
// OptionalExpression : ((MemberExpression | CallExpression) OptionalChain)*
//
// OptionalChain ::=
// '?.' Arguments
// '?.' '[' Expression ']'
// '?.' IdentifierName
// '?.' TemplateLiteral (TODO)
// '?.' PrivateIdentifier
// OptionalChain Arguments
// OptionalChain '[' Expression ']'
// OptionalChain '.' IdentifierName
// OptionalChain TemplateLiteral (TODO)
// OptionalChain '.' PrivateIdentifier
//
// (MemberExpression | CallExpression) OptionalChain OptionalExpressionRest
// OptionalExpressionRest ::= OptionalChain OptionalExpressionRest | ε
func (p *Parser) parseOptionalExpression() (bool, error) {
	if p.Peek().Type == l.TQuestionMark {
		p.Next() // consume '?'
		if p.Peek().Type != l.TPeriod {
			return false, fmt.Errorf("parseOptionalExpression rejected")
		}
		if p.PeekN(1).Type == l.TLeftBracket {
			// only consume '.' if '[' follows, so that we can have a clear and simple
			// separation of the two productions.
			p.Next() // consume '.'
		} else if p.PeekN(1).Type == l.TLeftParen {
			p.Next() // consume '.'
		}
		return true, nil
	}

	return false, nil
}

func (p *Parser) parseImportCall() (Expr, error) {
	if p.Peek().Type == l.TImport {
		p.Next() // consume 'import'
		if p.Peek().Type != l.TLeftParen {
			return nil, fmt.Errorf("parseImportCall rejected")
		}
		p.Next() // consume '('
		expr, err := p.parseAssignExpr()
		if err != nil {
			return nil, err
		}
		if p.Peek().Type != l.TRightParen {
			return nil, fmt.Errorf("parseImportCall rejected")
		}
		p.Next() // consume ')'
		return &ExprImportCall{
			source: expr,
		}, nil
	}
	return nil, fmt.Errorf("parseImportCall rejected")
}

// parses Arguments expression:
// Arguments ::=
// '(' ')'
// '(' ArgumentList ')'
// '(' ArgumentList ',' ')'
//
// ArgumentList ::=
// AssignmentExpression
// '...' AssignmentExpression
// ArgumentList ',' AssignmentExpression
// ArgumentList ',' '...' AssignmentExpression
func (p *Parser) parseArguments() ([]Expr, error) {
	p.Log("parseArguments")
	var (
		err       error
		arguments []Expr
	)

	// simplifying the expression to:
	// Arguments ::= '(' ArgumentList? ')' ExpressionRest
	//
	// expanding ArgumentList:
	// Arguments ::= '(' ('...'? AssignmentExpression ArgumentListRest)* ')' ExpressionRest
	//
	// ArgumentListRest ::= (',' '...'? AssignmentExpression)*
	if p.Peek().Type == l.TLeftParen {
		var exprAssign Expr

		p.Next() // consume '('
		switch p.Peek().Type {
		case l.TEllipsis:
			// fn(...a
			p.Next()                              // consume '...'
			exprAssign, err = p.parseAssignExpr() // consume AssignExpression
			if err != nil {
				return nil, err
			}
			exprAssign = &SpreadElement{exprAssign}
		case l.TRightParen:
			// fn()
			p.Next() // consume ')'
			return []Expr{}, nil

		default:
			// fn(a
			exprAssign, err = p.parseAssignExpr() // consume AssignExpression
			if err != nil {
				return nil, err
			}
		}

		arguments = append(arguments, exprAssign) // populate first arg

	argumentsLoop:
		for {
			// ',' '...'? AssignmentExpression ArgumentListRest
			switch p.Peek().Type {
			case l.TComma:
				p.Next() // consume ','
				if p.Peek().Type == l.TEllipsis {
					p.Next()                                               // consume '...'
					if exprAssign, err = p.parseAssignExpr(); err != nil { // consume AssignExpression
						return nil, err
					}

					exprAssign = &SpreadElement{argument: exprAssign}
				} else if p.Peek().Type == l.TRightParen {
					p.Next() // consume ')'
					break argumentsLoop
				} else {
					if exprAssign, err = p.parseAssignExpr(); err != nil { // consume AssignExpression
						return nil, err
					}
				}

				arguments = append(arguments, exprAssign) // populate arguments
			case l.TRightParen:
				p.Next() // consume ')'
				break argumentsLoop
			}
		}
	}

	return arguments, nil
}

// parseMemberAccess parses the following grammar:
//
// MemberAccess ::=
// | '.' (IdentifierName | PrivateIdentifier)
// | '[' Expr ']'
//
// it has a different behavior; instead of returning the Expr, it modifies the
// Expr passed as argument, as this is a production that often appears within
// recursive productions.
func (p *Parser) parseMemberAccess() (Expr, error) {
	p.Log("parseMemberAccess")

	switch p.Peek().Type {
	case l.TPeriod:
		// MemberExpression ::= ('.' IdentifierName MemberExpression')*
		p.Next() // consume '.'
		afterPeriod := p.Peek()
		switch afterPeriod.Type {
		case l.TNumberSign:
			p.Next() // consume '#'
			afterHash := p.Peek()
			if afterHash.Type == l.TIdentifier {
				p.Next() // consume IdentifierName
				return &ExprPrivateIdentifier{
					name: afterHash.Lexeme,
				}, nil
			} else {
				return nil, fmt.Errorf("expected identifier after '#.'")
			}
		case l.TIdentifier:
			p.Next() // consume IdentifierName
			return &ExprIdentifier{
				name: afterPeriod.Lexeme,
			}, nil
		default:
			return nil, fmt.Errorf("expected identifier after '.'")
		}
	case l.TLeftBracket:
		// MemberExpression ::= '[' Expression ']'
		p.Next()                                    // consume '['
		if expr, err := p.parseExpr(); err == nil { // parseExpr consume Expression
			if p.Peek().Type == l.TRightBracket {
				p.Next() // consume ']'
				return expr, nil
			} else {
				return nil, fmt.Errorf("expected ']' after expression")
			}
		} else {
			return nil, fmt.Errorf("expected valid expression after '['")
		}
	default:
		break
	}

	return nil, fmt.Errorf("expected '.' or '['")
}

// parseCallExpr parses the following grammar:
//
// CallExpression ::=
// | MemberExpression '(' ArgumentList? ')'
// | SuperCall
// | ImportCall
// | CallExpression '(' ArgumentList? ')'
// | CallExpression '[' Expression ']'
// | CallExpression '.' IdentifierName
// | CallExpression TemplateLiteral (TODO)
// | CallExpression '.' PrivateIdentifier
//
// transforming the productions removing left recursion and expanding:
//
// CallExpression ::=
// | MemberExpression CallExpressionRest
// | SuperCall CallExpressionRest
// | ImportCall CallExpressionRest

// CallExpressionRest ::=
// | '(' ArgumentList? ')' CallExpressionRest
// | '[' Expression ']' CallExpressionRest
// | '.' IdentifierName CallExpressionRest
// | TemplateLiteral CallExpressionRest
// | '.' PrivateIdentifier CallExpressionRest
// | ε
func (p *Parser) parseCallExpr() (Expr, error) {
	p.Log("parseCallExpr")
	var (
		exprCall Expr
		err      error
	)

	exprCall, err = p.parseMemberExpr()
	if err != nil {
		// CallExpression : SuperCall CallExpressionRest
		switch p.Peek().Type {
		case l.TSuper:
			exprCall = &ExprCall{
				callee: MakeLiteralExpr(l.TSuper),
			}
		case l.TImport:
			// CallExpression : ImportCall CallExpressionRest
			exprImportCall, err := p.parseImportCall()
			if err == nil {
				exprCall = exprImportCall
			} else {
				return nil, err
			}
		case l.TPeriod:
			panic("damn")
		default:
			return nil, fmt.Errorf("parseCallExpr rejected")
		}
	}

restLoop:
	for {
		optional, err := p.parseOptionalExpression()
		if err != nil {
			return nil, err
		}

		token := p.Peek()
		switch token.Type {
		case l.TLeftParen:
			// T ::= Arguments CallExpressionRest
			if arguments, err := p.parseArguments(); err != nil {
				return nil, err
			} else {
				exprCall = &ExprCall{
					callee:    exprCall,
					arguments: arguments,
					optional:  optional,
				}
			}
		case l.TPeriod, l.TLeftBracket:
			// T ::= MemberAccess CallExpressionRest
			if property, err := p.parseMemberAccess(); err != nil {
				return nil, err
			} else {
				exprCall = &ExprMemberAccess{
					object:   exprCall,
					property: property,
					optional: optional,
				}
			}
		default:
			break restLoop
		}
	}

	return exprCall, nil
}

// parseLeftHandSideExpr parses the following grammar:
//
// NewExpression ::= MemberExpression | 'new' NewExpression
func (p *Parser) parseNewExpr() (Expr, error) {
	p.Log("parseNewExpr")
	var (
		exprNew Expr
		err     error
	)

loop:
	for {
		switch p.Peek().Type {
		case l.TNew:
			// NewExpression ::= 'new' MemberExpression
			p.Next() // consume 'new'
			if p.Peek().Type == l.TNew {
				// NewExpression ::= ('new')+ NewExpression
				// resolves recursive 'new' tokens
				newExprRest, err := p.parseNewExpr()
				if err != nil {
					return nil, err
				}
				exprNew = &ExprNew{
					callee: newExprRest,
				}
				return exprNew, nil
			}

			exprNew, err = p.parseMemberExpr()
			if err != nil {
				return nil, err
			}

			// NewExpression ::= 'new' MemberExpression Arguments?
			arguments, err := p.parseArguments()
			if err != nil {
				return nil, err
			}

			// NewExpression ::= 'new' (Arguments)? MemberExpression
			return &ExprNew{
				callee:    exprNew,
				arguments: arguments,
			}, nil

		default:
			break loop
		}
	}

	return nil, fmt.Errorf("rejected on newExpression")
}

// parseMemberExpr parses the following grammar:
//
// MemberExpression ::=
// PrimaryExpression
// | MemberExpression | ('.' IdentifierName | '[' Expression ']' | TemplateLiteral | '.' PrivateIdentifier)
// | SuperProperty
// | MetaProperty
// | 'new' MemberExpression Arguments
//
// SuperProperty ::= 'super' ('[' Expression ']' | '.' IdentifierName)
// MetaProperty ::= 'new' '.' 'target' | 'import' '.' 'meta'
//
// Transforming to remove the left-recursion:
//
// MemberExpression ::=
// | MemberExpressionRest ('[' Expr ']' MemberExpressionRest)*
// | MemberExpressionRest ('.' IdentifierName MemberExpression')*
// | MemberExpressionRest ('.' PrivateIdentifier)*
// | MemberExpressionRest (TemplateLiteral MemberExpression')* (TODO)
// | 'new' MemberExpression Arguments
//
// MemberExpressionRest ::=
// PrimaryExpression | SuperProperty | MetaProperty | new MemberExpression Arguments
func (p *Parser) parseMemberExpr() (Expr, error) {
	p.Log("parseMemberExpr")
	var (
		exprMember Expr
		err        error
	)

	// MemberExpressionRest ::=
	// PrimaryExpression | SuperProperty | MetaProperty
	exprMember, err = p.parsePrimaryExpr()
	if err != nil {
		switch p.Peek().Type {
		case l.TSuper:
			// SuperProperty :: = 'super' ('[' Expression ']' | '.' IdentifierName)
			p.Next() // consume 'super'
			exprMember = MakeLiteralExpr(l.TSuper)
		case l.TNew:
			// MetaProperty ::= 'new' '.' 'target'
			if p.PeekN(1).Type == l.TPeriod && p.PeekN(2).Lexeme == "target" {
				p.Next() // consume 'new'
				p.Next() // consume '.'
				p.Next() // consume 'target'
				return &ExprMetaProperty{
					meta: MakeLiteralExpr(l.TNew),
					property: &ExprIdentifier{
						name: "target",
					},
				}, nil
			}
			return nil, fmt.Errorf("invalid meta property")
		case l.TImport:
			// MetaProperty ::= 'import' '.' 'meta'
			if p.PeekN(1).Type == l.TPeriod && p.PeekN(2).Lexeme == "meta" {
				p.Next() // consume 'import'
				p.Next() // consume '.'
				p.Next() // consume 'meta'
				return &ExprMetaProperty{
					meta: MakeLiteralExpr(l.TImport),
					property: &ExprIdentifier{
						name: "meta",
					},
				}, nil
			}
			return nil, fmt.Errorf("invalid meta property")
		default:
			return nil, err
		}
	}

	// MemberExpression ::=
	// | ('[' Expr ']' MemberExpressionRest)*
	// | ('.' IdentifierName MemberExpressionRest)*
loop:
	for {
		token := p.Peek()
		switch token.Type {
		case l.TPeriod, l.TLeftBracket:
			if property, err := p.parseMemberAccess(); err != nil {
				return nil, err
			} else {
				exprMember = &ExprMemberAccess{
					object:   exprMember,
					property: property,
				}
			}
		default:
			break loop
		}
	}

	return exprMember, nil
}

// PrimaryExpression ::=
// | this
// | IdentifierReference
// | Literal
// | ArrayLiteral
// | ObjectLiteral (TODO)
// | FunctionExpression (TODO)
// | ClassExpression (TODO)
// | GeneratorExpression (TODO)
// | AsyncFunctionExpression (TODO)
// | AsyncGeneratorExpression (TODO)
// | RegularExpressionLiteral (TODO)
// | TemplateLiteral (TODO)
// | CoverParenthesizedExpressionAndArrowParameterList (TODO)
func (p *Parser) parsePrimaryExpr() (Expr, error) {
	p.Log("parsePrimaryExpr")
	cp := p.saveCheckpoint()
	literal, err := p.parseLiteralAndIdentifier()
	if err == nil {
		return literal, nil
	}

	p.restoreCheckpoint(cp)
	cp = p.saveCheckpoint()
	array, err := p.parseArrayInitializer()
	if err == nil {
		return array, nil
	}
	p.restoreCheckpoint(cp)

	cp = p.saveCheckpoint()
	object, err := p.parseObjectInitializer()
	if err == nil {
		return object, nil
	}
	p.restoreCheckpoint(cp)
	return nil, fmt.Errorf("rejected on primaryExpression")
}

func (p *Parser) parseLiteralAndIdentifier() (Expr, error) {
	p.Log("parseLiteral")
	var primaryExpr Expr
	token := p.Peek()

	switch token.Type {
	case l.TIdentifier:
		primaryExpr = &ExprIdentifier{
			name: token.Lexeme,
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
	case l.TThis:
		primaryExpr = ExprLitThis
	default:
		return nil, fmt.Errorf("primaryExpr rejected")
	}

	p.Next() // consume token
	return primaryExpr, nil
}

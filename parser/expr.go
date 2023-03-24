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
	callee    Node
	arguments []Node
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
	return fmt.Sprintf("(. %s %s)", e.property.S(), e.object.S())
}

// ///////////////////
// ExprSuperProp //
// ///////////////////
const ESuperProp ExprType = "ExprSuperProp"

type ExprSuperProp struct {
	object   Node
	property Node
}

func (e *ExprSuperProp) Type() ExprType {
	return ESuperProp
}

func (e *ExprSuperProp) S() string {
	if e == nil {
		panic("invalid object: nil")
	}
	return fmt.Sprintf("(. %s %s)", e.property.S(), e.object.S())
}

// ///////////////////
// ExprMetaProperty //
// ///////////////////
const EMetaProperty ExprType = "ExprMetaProperty"

type ExprMetaProperty struct {
	meta     Node
	property Node
}

func (e *ExprMetaProperty) Type() ExprType {
	return EMetaProperty
}

func (e *ExprMetaProperty) S() string {
	if e == nil {
		panic("invalid object: nil")
	}
	return fmt.Sprintf("(. %s %s)", e.meta.S(), e.property.S())
}

// ///////////////////
// ExprCall //
// ///////////////////
const ECall ExprType = "ExprCall"

type ExprCall struct {
	callee    Node
	arguments []Node
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
	return fmt.Sprintf("(%s %s)", e.callee.S(), args.String())
}

// ///////////////////
// SpreadElement //
// ///////////////////
const NSpreadElement ExprType = "SpreadElement"

type SpreadElement struct {
	argument Node
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
	source Node
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
// | UnaryOp UnaryExpression
// | UpdateExpression
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
// | LeftHandSideExpression (++ | --)?
// | (++ | --) UnaryExpression
func (p *Parser) parseUpdateExpr() (Node, error) {
	p.Log("parseUpdateExpr")
	var (
		exprUpdate Node
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
// | OptionalExpression (TODO)
func (p *Parser) parseLeftHandSideExpr() (Node, error) {
	p.Log("parseLeftHandSideExpr")
	var (
		expr Node
		err  error
	)

	p.saveCheckpoint()
	if expr, err = p.parseCallExpr(); err == nil {
		return expr, nil
	}
	p.restoreCheckpoint()

	p.saveCheckpoint()
	if expr, err = p.parseNewExpr(); err == nil {
		return expr, nil
	}
	p.restoreCheckpoint()

	// TODO: parseOptionalExpression
	return nil, fmt.Errorf("parseLeftHandSideExpr rejected")
}

func (p *Parser) parseImportCall() (Node, error) {
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

// parseCallExpr parses the following grammar:
//
// CallExpression ::=
// | MemberExpression '(' ArgumentList? ')'
// | SuperCall
// | ImportCall (TODO)
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
// | ImportCall CallExpressionRest (TODO)

// CallExpressionRest ::=
// | '(' ArgumentList? ')' CallExpressionRest
// | '[' Expression ']' CallExpressionRest
// | '.' IdentifierName CallExpressionRest
// | TemplateLiteral CallExpressionRest
// | '.' PrivateIdentifier CallExpressionRest
// | ε
func (p *Parser) parseCallExpr() (Node, error) {
	p.Log("parseCallExpr")
	var (
		exprCall Node
		err      error
	)

	exprCall, err = p.parseMemberExpr()
	if err != nil {
		// CallExpression ::= SuperCall CallExpressionRest
		switch p.Peek().Type {
		case l.TSuper:
			exprCall = &ExprCall{
				callee: MakeLiteralExpr(l.TSuper),
			}
		case l.TImport:
			// CallExpression ::= ImportCall CallExpressionRest
			exprImportCall, err := p.parseImportCall()
			if err == nil {
				exprCall = exprImportCall
			} else {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("parseCallExpr rejected")
		}
	}

restLoop:
	for {
		token := p.Peek()

		switch token.Type {
		case l.TLeftParen:
			// T ::= '(' ArgumentList? ')' CallExpressionRest
			// T ::= '(' ('...'? AssignmentExpression ArgumentListRest)* ')' CallExpressionRest
			//
			// ArgumentListRest ::= (',' '...'? AssignmentExpression)*
			arguments := []Node{}
			p.Next() // consume '('
			var exprAssign Node
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
				exprCall = &ExprCall{
					callee:    exprCall,
					arguments: []Node{},
				}
				continue restLoop
			default:
				// fn(a
				exprAssign, err = p.parseAssignExpr() // consume AssignExpression
				if err != nil {
					return nil, err
				}
			}

			arguments = append(arguments, exprAssign) // populate arguments

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

			exprCall = &ExprCall{
				callee:    exprCall,
				arguments: arguments,
			}

		case l.TPeriod:
			// T ::= '.' PrivateIdentifier CallExpressionRest
			p.Next() // consume '.'
			afterPeriod := p.Peek()
			if afterPeriod.Type == l.TIdentifier {
				exprCall = &ExprMemberAccess{
					object: exprCall,
					property: &ExprIdentifier{
						name: afterPeriod.Lexeme,
					},
				}
				p.Next() // consume identifier
			} else {
				return nil, fmt.Errorf("expected identifier after dot")
			}
		case l.TLeftBracket:
			// T ::= '[' Expression ']' CallExpressionRest
			p.Next()                                    // consume '['
			if expr, err := p.parseExpr(); err == nil { // parseExpr consumes the expression's tokens
				if p.Peek().Type == l.TRightBracket {
					exprCall = &ExprMemberAccess{
						object:   exprCall,
						property: expr,
					}
					p.Next() // consume ']'
				}
			} else {
				return nil, fmt.Errorf("expected valid expression after left bracket")
			}
		default:
			break restLoop
		}
	}

	return exprCall, nil
}

// parseLeftHandSideExpr parses the following grammar:
//
// NewExpression ::= MemberExpression | new NewExpression
func (p *Parser) parseNewExpr() (Node, error) {
	p.Log("parseNewExpr")
	var (
		exprNew Node
		err     error
		match   bool
	)

loop:
	for {
		switch p.Peek().Type {
		case l.TNew:
			// NewExpression ::= 'new' MemberExpression
			p.Next() // consume 'new'
			if p.Peek().Type == l.TNew {
				newExprRest, err := p.parseNewExpr()
				if err != nil {
					return nil, err
				}
				exprNew = &ExprNew{
					callee: newExprRest,
				}
				// this will consume the next 'new' token indefinitely
				// until it ends
				return exprNew, nil
			}

			exprNew, err = p.parseMemberExpr()
			if err != nil {
				return nil, err
			}

			arguments := []Node{}
			if p.Peek().Type == l.TLeftParen {
				// NewExpression ::= 'new' MemberExpression ArgumentsRest
				// ArgumentsRest ::= ('(' ('...'? AssignmentExpression ArgumentListRest)* ')')*
				//
				// ArgumentListRest ::= (',' '...'? AssignmentExpression)*
				p.Next() // consume '('
				var exprAssign Node

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
					exprNew = &ExprNew{
						callee:    exprNew,
						arguments: []Node{},
					}
				default:
					// fn(a
					exprAssign, err = p.parseAssignExpr() // consume AssignExpression
					if err != nil {
						return nil, err
					}
				}

				arguments = append(arguments, exprAssign) // populate arguments

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

			exprNew = &ExprNew{
				callee:    exprNew,
				arguments: arguments,
			}
			match = true
		default:
			break loop
		}
	}

	if match {
		// NewExpression ::= MemberExpression
		return exprNew, nil
	} else {
		return nil, fmt.Errorf("rejected on newExpression")
	}
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
// | MemberExpressionRest ('.' PrivateIdentifier)* (TODO)
// | MemberExpressionRest (TemplateLiteral MemberExpression')* (TODO)
// | 'new' MemberExpression Arguments
//
// MemberExpressionRest ::=
// PrimaryExpression | SuperProperty | MetaProperty | new MemberExpression Arguments
func (p *Parser) parseMemberExpr() (Node, error) {
	p.Log("parseMemberExpr")
	var (
		exprMember Node
		err        error
	)

	// MemberExpressionRest ::=
	// PrimaryExpression | SuperProperty | MetaProperty
	exprMember, err = p.parsePrimaryExpr()
	if err != nil {
		switch p.Peek().Type {
		case l.TSuper:
			// SuperProperty :: = 'super' ('[' Expression ']' | '.' IdentifierName)
			exprMember = MakeLiteralExpr(l.TSuper)
			p.Next() // consume 'super'
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
	lastCursor := p.cursor
loop:
	for {
		token := p.Peek()
		switch token.Type {
		case l.TPeriod:
			// MemberExpression ::= ('.' IdentifierName MemberExpression')*
			p.Next() // consume period
			afterPeriod := p.Peek()
			switch afterPeriod.Type {
			case l.TNumberSign:
				p.Next() // consume '#'
				if p.Peek().Type == l.TIdentifier {
					// match = true
					exprMember = &ExprMemberAccess{
						object: exprMember,
						property: &ExprPrivateIdentifier{
							name: p.Peek().Lexeme,
						},
					}
					p.Next() // consume identifier
				} else {
					return nil, fmt.Errorf("expected identifier after dot")
				}

			case l.TIdentifier:
				exprMember = &ExprMemberAccess{
					object: exprMember,
					property: &ExprIdentifier{
						name: afterPeriod.Lexeme,
					},
				}
				p.Next() // consume identifier
			default:
				return nil, fmt.Errorf("expected identifier after dot")
			}

		case l.TLeftBracket:
			// MemberExpression ::= ('[' Expr ']' MemberExpression')*
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

	return exprMember, nil
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
	default:
		return nil, fmt.Errorf("primaryExpr rejected")
	}

	p.Next() // consume token
	return primaryExpr, nil
}

package parser

import (
	"errors"
	"fmt"

	"github.com/ruiconti/gojs/lex"
)

const EUnaryOp ExprType = "EUnaryOp"

// UnaryExpression[Yield, Await]:
//
//	| UpdateExpression
//	| delete UnaryExpression
//	| void UnaryExpression
//	| typeof UnaryExpression
//	| + UnaryExpression
//	| - UnaryExpression
//	| ~ UnaryExpression
//	| ! UnaryExpression
//
// We'll focus on the 2-4 cases for now
type ExprUnaryOp struct {
	operand  AstNode
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

func isUnaryOperator(t lex.Token) bool {
	for _, op := range UnaryOperators {
		if t.T == op {
			return true
		}
	}
	return false
}

// Parser
var errNotUnaryOperator = errors.New("current token is not an unary operator")

func (p *Parser) parseUnaryOperator(cursor *int) (AstNode, error) {
	// OPERATOR
	// delete 0
	// ˆ
	tok, err := p.peekN(*cursor)
	if err != nil {
		return &ExprUnaryOp{}, err
	}

	// OPERAND
	// delete 0
	//        ˆ
	operator := tok.T
	p.logger.Debug("[%d:%d] parser:unaryOpExpr: %v", p.cursor, *cursor, lex.ResolveName(tok.T))

	if isUnaryOperator(tok) {
		*cursor = *cursor + 1
		opExpr, err := p.parseUnaryOperator(cursor)
		if err == nil {
			expr := &ExprUnaryOp{
				operator: operator,
				operand:  opExpr,
			}
			p.logger.Debug("[%d:%d] parser:unaryOpExpr:unary:acc %v -> %v", p.cursor, *cursor, lex.ResolveName(tok.T), expr.PrettyPrint())
			return expr, nil
		}
	}

	// the other production is UpdateExpression
	// we are, for now, cutting the path:
	// UpdateExpression -> LeftHandSideExpression -> PrimaryExpression -> IdentifierReference
	opExpr, err := p.parsePrimaryExpr(cursor)
	if err == nil {
		p.logger.Debug("[%d:%d] parser:unaryOpExpr:primary:acc %v -> %v", p.cursor, *cursor, lex.ResolveName(tok.T), opExpr.PrettyPrint())
		return opExpr, nil
	}

	p.logger.Debug("[%d:%d] parser:unaryOpExpr:rej %v (%s)", p.cursor, *cursor, tok, err)
	return nil, fmt.Errorf("no productions left in unary operator")
}

var UpdateOperators = []lex.TokenType{
	lex.TMinusMinus,
	lex.TPlusPlus,
}

func isUpdateExpression(t lex.Token) bool {
	for _, op := range UpdateOperators {
		if t.T == op {
			return true
		}
	}
	return false
}

var errNotUpdateOperator = errors.New("current token is not an update operator")

func (p *Parser) parseUpdateExpr(cursor *int) (AstNode, error) {
	// UpdateExpression -> LeftHandSideExpression
	// UpdateExpression -> ++ UnaryExpression
	// UpdateExpression -> -- UnaryExpression

	tok, err := p.peekN(*cursor)
	if err != nil {
		return nil, err
	}

	if !isUpdateExpression(tok) {
		return nil, errNotUpdateOperator
	}

	p.logger.Debug("[%d:%d] parser:updateExpr %v", p.cursor, *cursor, lex.ResolveName(tok.T))

	*cursor = *cursor + 1
	// first tries to resolve it as a primary expr
	opExpr, err := p.parsePrimaryExpr(cursor)
	if err == nil {
		expr := &ExprUnaryOp{
			operator: tok.T,
			operand:  opExpr,
		}
		p.logger.Debug("[%d:%d] parser:updateExpr:primary:acc %v -> %v", p.cursor, *cursor, lex.ResolveName(tok.T), expr.PrettyPrint())
		return expr, nil
	}

	return nil, fmt.Errorf("no productions left in update operator")
}

var errNotPrimaryExpr = errors.New("current token is not a primary expression")

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
		primaryExpr = &ExprNumeric{
			value: token.Lexeme,
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

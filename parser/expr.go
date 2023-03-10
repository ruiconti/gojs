package parser

import (
	"fmt"

	"github.com/ruiconti/gojs/lex"
)

type Expr interface{}

func (p *Parser) parseExpr(c *int) (AstNode, error) {
	return p.parseAssignExpr(c)
}

func (p *Parser) parseAssignExpr(c *int) (AstNode, error) {
	return p.parseCondExpr(c)
}

type ExprBinaryOp struct {
	left     AstNode
	right    AstNode
	operator lex.TokenType
}

const EBinaryOp ExprType = "ExprBinaryOp"

func (e *ExprBinaryOp) Type() ExprType {
	return EBinaryOp
}

func (e *ExprBinaryOp) Source() string {
	return fmt.Sprintf("%s %s %s", e.left.Source(), lex.ResolveName(e.operator), e.right.Source())
}

func (e *ExprBinaryOp) PrettyPrint() string {
	return fmt.Sprintf("(%s %s %s)", lex.ResolveName(e.operator), e.left.PrettyPrint(), e.right.PrettyPrint())
}

func (p *Parser) parseCondExpr(c *int) (AstNode, error) {
	return p.parseLogOrExpr(c)
}

func (p *Parser) parseBinaryExprGeneric(
	cursor *int,
	operator lex.TokenType,
	left func(*int) (AstNode, error),
	right func(*int) (AstNode, error),
) (AstNode, error) {
	str := lex.ResolveName(operator)

	p.logger.Debug("[%d] parseBinExprGeneric(%s)", p.cursor, str)
	c := 0
	expr, err := left(&c)
	if err != nil {
		p.logger.Debug("[%d] bailing", p.cursor)
		return &ExprBinaryOp{}, err
	}

	p.advanceBy(c)
	for p.matchAny(operator) {
		p.logger.Debug("[%d] parseBinExprGeneric(%s): loop", p.cursor, str)
		c = 0
		right, err := right(&c)
		if err != nil {
			// eof or didn't match
			break
		}

		p.advanceBy(c)
		p.logger.Debug("[%d] parseBinExprGeneric(%s): expr (before): %s", p.cursor, str, expr.PrettyPrint())
		expr = &ExprBinaryOp{
			operator: operator,
			left:     expr,
			right:    right,
		}
		p.logger.Debug("[%d] parseBinExprGeneric(%s): expr (after): %s", p.cursor, str, expr.PrettyPrint())
	}
	return expr, nil

}

func (p *Parser) parseLogOrExpr(c *int) (AstNode, error) {
	return p.parseBinaryExprGeneric(
		c,
		lex.TLogicalOr,
		p.parseAndExpr,
		p.parseAndExpr,
	)
}

func (p *Parser) parseAndExpr(c *int) (AstNode, error) {
	return p.parseBinaryExprGeneric(
		c,
		lex.TLogicalAnd,
		p.parseBitOrExpr,
		p.parseBitOrExpr,
	)
}

func (p *Parser) parseBitOrExpr(c *int) (AstNode, error) {
	return p.parseBinaryExprGeneric(
		c,
		lex.TOr,
		p.parseBitXorExpr,
		p.parseBitXorExpr,
	)
}

func (p *Parser) parseBitXorExpr(c *int) (AstNode, error) {
	return p.parseBinaryExprGeneric(
		c,
		lex.TXor,
		p.parsePrimaryExpr,
		p.parsePrimaryExpr,
	)
}

// func (p *Parser) parseBitOrExpr() {}

// func (p *Parser) parseBitXorExpr() {}

// func (p *Parser) parseBitAndExpr() {}

// func (p *Parser) parseEqualityExpr() {}

// func (p *Parser) parseRelationalExpr() {}

// func (p *Parser) parseShiftExpr() {}

// func (p *Parser) parseAdditiveExpr() {}

// func (p *Parser) parseExponentialExpr() {}

// func (p *Parser) parseUnaryExpr() {}

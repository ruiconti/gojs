package parser

import (
	"fmt"
	"strings"

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

	p.logger.Debug("[%d] parseBinExprGeneric(%s) parse-left", p.cursor, str)
	c := 0
	expr, err := left(&c)
	if err != nil {
		p.logger.Debug("[%d] parseBinExprGeneric(%s) parse-left:bailing", p.cursor, str)
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

func (p *Parser) parseBinaryExprGeneric2(
	cursor *int,
	operators []lex.TokenType,
	left func(*int) (AstNode, error),
	right func(*int) (AstNode, error),
) (AstNode, error) {
	var opstr strings.Builder
	for _, op := range operators {
		opstr.Write([]byte(lex.ResolveName(op)))
	}

	p.logger.Debug("[%d] parseBinExprGeneric2(%s) parse-left", p.cursor, opstr.String())
	c := 0
	// try to parse left side
	expr, err := left(&c)
	if err != nil {
		p.logger.Debug("[%d] parseBinExprGeneric2(%s) parse-left:bailing", p.cursor, opstr.String())
		return &ExprBinaryOp{}, err
	}

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

		str := lex.ResolveName(operator.T)
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
			operator: operator.T,
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
		p.parseBitAndExpr,
		p.parseBitAndExpr,
	)
}

func (p *Parser) parseBitAndExpr(c *int) (AstNode, error) {
	return p.parseBinaryExprGeneric(
		c,
		lex.TAnd,
		p.parseEqualityExpr,
		p.parseEqualityExpr,
	)
}

func (p *Parser) parseEqualityExpr(c *int) (AstNode, error) {
	return p.parseBinaryExprGeneric2(
		c,
		[]lex.TokenType{lex.TEqual, lex.TNotEqual, lex.TStrictEqual, lex.TStrictNotEqual},
		p.parseRelationalExpr,
		p.parseRelationalExpr,
	)
}

func (p *Parser) parseRelationalExpr(c *int) (AstNode, error) {
	return p.parseBinaryExprGeneric2(
		c,
		[]lex.TokenType{lex.TGreaterThan, lex.TGreaterThanEqual, lex.TLessThan, lex.TLessThanEqual, lex.TInstanceof, lex.TIn},
		p.parseShiftExpr,
		p.parseShiftExpr,
	)
}

func (p *Parser) parseShiftExpr(c *int) (AstNode, error) {
	return p.parseBinaryExprGeneric2(
		c,
		[]lex.TokenType{lex.TLeftShift, lex.TRightShift, lex.TUnsignedRightShift},
		p.parseAdditiveExpr,
		p.parseAdditiveExpr,
	)
}

func (p *Parser) parseAdditiveExpr(c *int) (AstNode, error) {
	return p.parseBinaryExprGeneric2(
		c,
		[]lex.TokenType{lex.TPlus, lex.TMinus},
		p.parseMultiplicativeExpr,
		p.parseMultiplicativeExpr,
	)
}

func (p *Parser) parseMultiplicativeExpr(c *int) (AstNode, error) {
	return p.parseBinaryExprGeneric2(
		c,
		[]lex.TokenType{lex.TStar, lex.TSlash, lex.TPercent},
		p.parseExponentialExpr,
		p.parseExponentialExpr,
	)
}

func (p *Parser) parseExponentialExpr(c *int) (AstNode, error) {
	return p.parseBinaryExprGeneric2(
		c,
		[]lex.TokenType{lex.TStarStar},
		p.parseUnaryOperator,
		p.parseUnaryOperator,
	)
}

// func (p *Parser) parseUnaryExpr() {}

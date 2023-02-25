package parser

import "github.com/ruiconti/gojs/lex"

const EUnaryOp ExprType = "EUnaryOp"

type ExprUnaryOp struct {
	operand  AstNode
	operator lex.TokenType
}

func (e *ExprUnaryOp) Type() ExprType {
	return EUnaryOp
}

func (e *ExprUnaryOp) Source() string {
	return e.operand.Source()
}

func isUnaryOperator(t lex.Token) bool {
	return t.T == lex.TPlus || t.T == lex.TMinus || t.T == lex.TBang || t.T == lex.TTilde
}

func (p *Parser) parseUnaryOperator(expr *ExprUnaryOp) (*ExprUnaryOp, error) {
	// delete 0
	// ˆ
	// p.cursor: 0
	curr := p.seq[p.cursor]
	p.logger.Debug("parsing unary operator: %d: %v", p.cursor, curr)

	// in unary operators, we first process the operator
	var unaryOpExpr *ExprUnaryOp
	if expr == nil {
		unaryOpExpr = &ExprUnaryOp{
			operator: curr.T,
		}
	} else {
		// expr.operand
	}

	// process the next token
	// delete 0
	//        ˆ
	// p.cursor: 1
	p.cursor++
	next, err := p.LookAhead(0)
	if err != nil {
		return &ExprUnaryOp{}, err
	}
	if isUnaryOperator(next) {
		// recursive descent
		p.parseUnaryOperator(unaryOpExpr)
	} else {
		// we are done parsing the left hand side of the unary operator
		// we now need to recursively parse the operand
		unaryOpExpr.operand = p.parsePrimaryExpr(unaryOpExpr)
	}

	return unaryOpExpr, nil
}

func (p *Parser) parsePrimaryExpr(parent AstNode) AstNode {
	cursorTmp := p.cursor // save cursor to allow backtracking
	reject := false

	// current token position
	// delete s.length
	//        ˆ
	// p.cursor: 1
	curr := p.seq[p.cursor]
	p.logger.Debug("parsing primary expression: %d: %v", cursorTmp, curr)

	// in primary expressions, we first process the operator
	var primaryExpr AstNode
	switch curr.T {
	case lex.TIdentifier:
		primaryExpr = &ExprIdentifierReference{
			reference: curr.Lexeme,
		}
	case lex.TNumericLiteral:
		primaryExpr = &ExprNumeric{
			value: curr.Lexeme,
		}
	case lex.TStringLiteral:
		primaryExpr = &ExprStringLiteral{
			value: curr.Lexeme,
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
		// case lex.TLeftParen:
		// delete 0
		//        ˆ
		// p.cursor: 1
		// p.cursor++
		// primaryExpr = p.parseExpression()
		// delete 0
		//           ˆ
		// p.cursor: 2
		// p.cursor++
	}
	if reject {
		p.cursor = cursorTmp
	} else {
		p.cursor++
	}

	// todo: next token position

	return primaryExpr
}

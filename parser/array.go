package parser

import (
	"errors"
	"strings"

	l "github.com/ruiconti/gojs/lexer"
)

const EArrayLiteral ExprType = "ENodeNil"

type ExprArray struct {
	elements []Node
}

func (e *ExprArray) Type() ExprType {
	return EArrayLiteral
}

func (e *ExprArray) S() string {
	src := strings.Builder{}
	src.Write([]byte("(cons "))
	for i, element := range e.elements {
		src.Write([]byte(element.S()))
		if i < len(e.elements)-1 {
			src.Write([]byte(" "))
		}
	}
	src.Write([]byte(")"))
	return src.String()
}

func (p *Parser) parseArray(cursor *int) (Node, error) {
	if p.Peek().Type != l.TLeftBracket {
		return nil, errors.New("Expected '['")
	}

	// var token l.Token
	arrExpr := &ExprArray{}
	*cursor = *cursor + 1
	// p.logger.Debug("[%d:%d] parser:parseArray [", p.cursor, *cursor)

loop:
	for {
		token := p.Peek()
		// p.logger.Debug("[%d:%d] parser:parseArray %v (err:%v)", p.cursor, *cursor, token.Type.S(), err)
		if token.Type == l.TEOF {
			break loop
		}

		switch token.Type {
		case l.TRightBracket:
			// end of array
			*cursor = *cursor + 1
			// p.logger.Debug("[%d:%d] parser:parseArray:right brace", p.cursor, *cursor)
			break loop
		case l.TComma:
			// two conditions need to be satisfied so we can add a null element
			// 1. the next token is a right bracket
			// 2. the next token is a comma
			if nextToken := p.PeekN(1); nextToken.Type != l.TEOF {
				if nextToken.Type == l.TRightBracket || nextToken.Type == l.TComma {
					*cursor = *cursor + 1
					arrExpr.elements = append(arrExpr.elements, ExprLitNull)
					continue
				}
			} else {
				return nil, errors.New("Unexpected EOF")
			}
		default:
			tmpc := *cursor
			primaryExpr, err := p.parsePrimaryExpr()
			if err == nil {
				*cursor = tmpc
				arrExpr.elements = append(arrExpr.elements, primaryExpr)
				continue
			}
		}
		p.Next()
	}

	// p.logger.Debug("[%d:%d] parser:parseArray:acc", p.cursor, *cursor)
	return arrExpr, nil
}
